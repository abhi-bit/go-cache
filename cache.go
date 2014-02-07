//Package cache: Go-routine safe, simple LRU cache for storing documents([] bytes)
package cache

import (
	"container/list"
	"sync"
)

//Key - value pairs inside cache
type cacheValue struct {
	key   string
	bytes []byte
}

//Size of key - value pair. Not counting their metadata
func (v *cacheValue) size() uint64 {
	return uint64(len([]byte(v.key)) + len(v.bytes))
}

//Base struct for LRU cache
//Need to grab a Lock before any goroutine can edit
//TODO - make it non-blocking using channels
type Cache struct {
	sync.Mutex

	Size uint64

	//Capacity of LRU Cache
	capacity uint64

	//List handling the eviction based on LRU
	list *list.List

	//Hash table to make data retrival quick
	table map[string]*list.Element
}

//Cache with max size of capacity bytes
func New(capacity uint64) *Cache {
	return &Cache{
		capacity: capacity,
		list:     list.New(),
		table:    make(map[string]*list.Element),
	}
}

//Prunes the Cache on LRU
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

//Inserts key and doc([] byte). Doesn't overwrite if key exists
//Returns LRU cache size at that point
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

//Retrives a key if its present
//Returns doc and boolean flag to tell if key was there or not
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

//Updates the LRU timestamp of input key
func (c *Cache) Update(key string) {
	c.Lock()
	defer c.Unlock()

	elt, ok := c.table[key]
	if !ok {
		return
	}
	c.list.MoveToFront(elt)
}

//Deletes input key
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

//Peeks into LRU to spit out key that will get evicted from cache first
func (c *Cache) Peek() (key string, document []byte, size uint64) {
	elt := c.list.Back()
	if elt == nil {
		return
	}
	v := elt.Value.(*cacheValue)
	return v.key, v.bytes, v.size()
}
