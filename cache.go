package cache

import (
    "container/list"
    "sync"
)

type cacheValue struct {
    key string
    bytes []byte
}

func (v *cacheValue) size() uint64{
    return uint64(len([]byte(v.key)) + len(v.bytes))
}

type Cache struct {
    sync.Mutex

    Size uint64

    capacity uint64
    list *list.List
    table map[string]*list.Element
}

func New(capacity uint64) *Cache {
    return &Cache{
        capacity: capacity,
        list: list.New(),
        table: make(map[string]*list.Element),
    }
}

func (c *Cache) trim() {
    for c.Size > c.capacity {
        elt := c.list.Back()
        if elt == nil {
            return
        }
        v := c.list.Remove(elt).(*cacheValue)
        delete(c.table, v.key)
        c.Size -= v.size()
    }
}

func (c *Cache) Insert(key string, document []byte) (cacheSize uint64) {
    c.Lock()
    defer c.Unlock()

    _, ok := c.table[key]
    if ok {
        return
    }

    v := &cacheValue{key, document}
    elt := c.list.PushFront(v)
    c.table[key] = elt
    c.Size += v.size()
    c.trim()
    return c.Size
}

func (c *Cache) Get(key string) (document []byte, ok bool) {
    c.Lock()
    defer c.Unlock()

    elt, ok := c.table[key]
    if !ok {
        return nil, false
    }
    c.list.MoveToFront(elt)
    return elt.Value.(*cacheValue).bytes, true
}

func (c *Cache) Update(key string) {
    c.Lock()
    defer c.Unlock()

    elt, ok := c.table[key]
    if !ok {
        return
    }
    c.list.MoveToFront(elt)
}

func (c *Cache) Delete(key string) {
    c.Lock()
    defer c.Unlock()

    elt, ok := c.table[key]
    if !ok {
        return
    }
    v := c.list.Remove(elt).(*cacheValue)
    delete(c.table, key)
    c.Size -= v.size()
}

func (c *Cache) Peek() (key string, document []byte, size uint64) {
    elt := c.list.Back()
    if elt == nil {
        return
    }
    v := elt.Value.(*cacheValue)
   return v.key, v.bytes, v.size()
}
