//TODO: fix cachestats test cases

package cache

import (
	"encoding/json"
	"testing"
)

type cValue struct {
	Size int
}

func (c *cValue) getSize() int {
	return len(string(c.Size))
}

func testInitialState(t *testing.T) {
	cache := LRUCache(10)
	l, s, capacity := cache.CacheStats()
	if l != 0 {
		t.Errorf("length = %v, want 0", l)
	}
	if s != 0 {
		t.Errorf("Size = %v, want 0", s)
	}
	if capacity != 10 {
		t.Errorf("capacity = %v, want 10", capacity)
	}
}

func testInsert(t *testing.T) {
	cache := LRUCache(100)
	data, _ := json.Marshal(&cValue{Size: 0})
	cache.Insert("key", string(data))

	v, ok := cache.Get("key")
	if !ok || v != string(data) {
		t.Errorf("Cache has incorrect value: %v != %v", data, v)
	}
}

func testLRUEviction(t *testing.T) {
	Size := uint64(3)
	cache := LRUCache(Size)

	data, _ := json.Marshal(&cValue{Size: 1})

	cache.Insert("key1", string(data))
	cache.Insert("key2", string(data))
	cache.Insert("key3", string(data))
	// lru: [key3, key2, key1]

	// Look up the elements. This will rearrange the LRU ordering.
	cache.Get("key3")
	cache.Get("key2")
	cache.Get("key1")
	// lru: [key1, key2, key3]

	cache.Insert("key0", string(data))
	// lru: [key0, key1, key2]

	// The least recently used one should have been evicted.
	if _, ok := cache.Get("key3"); ok {
		t.Error("Least recently used element was not evicted.")
	}
}

/*
func testCapacityIsObeyed(t *testing.T) {
	Size := uint64(3)
	cache := LRUCache(Size)
    value,_ := json.Marshal(&cValue{Size: 1})

	// Insert up to the cache's capacity.
	cache.Insert("key1", string(value))
	cache.Insert("key2", string(value))
	cache.Insert("key3", string(value))
	if _, sz, _, _ := cache.CacheStats(); sz != Size {
		t.Errorf("cache.getSize() = %v, expected %v", sz, Size)
	}
	// Insert one more; something should be evicted to make room.
	cache.Insert("key4", string(value))
	if _, sz, _ := cache.CacheStats(); sz != Size {
		t.Errorf("post-evict cache.getSize() = %v, expected %v", sz, Size)
	}
}

func testInsertUpdatesgetSize(t *testing.T) {
	cache := LRUCache(100)
    emptyValue, _ := json.Marshal(&cValue{Size: 0})
	key := "key1"
	cache.Insert(key, string(emptyValue))
	if _, sz, _ := cache.CacheStats(); sz != 0 {
		t.Errorf("cache.getSize() = %v, expected 0", sz)
	}
    someValue, _ := json.Marshal(&cValue{Size: 20})
	key = "key2"
	cache.Insert(key, string(someValue))
	if _, sz, _ := cache.CacheStats(); sz != 20 {
		t.Errorf("cache.getSize() = %v, expected 20", sz)
	}
}
*/

func testInsertWithOldKeyUpdatesValue(t *testing.T) {
	cache := LRUCache(100)
	emptyValue, _ := json.Marshal(&cValue{Size: 0})
	key := "key1"
	cache.Insert(key, string(emptyValue))
	someValue, _ := json.Marshal(&cValue{Size: 20})
	cache.Insert(key, string(someValue))

	v, ok := cache.Get(key)
	if !ok || v != string(someValue) {
		t.Errorf("Cache has incorrect value: %v != %v", someValue, v)
	}
}
