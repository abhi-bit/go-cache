package cache

import (
	"encoding/json"
	"testing"
)

type cValue struct {
	size int
}

func (c *cValue) Size() int {
	return len(string(c.size))
}

func testInitialState(t *testing.T) {
	cache := LRUCache(10)
	l, s, capacity := cache.CacheStats()
	if l != 0 {
		t.Errorf("length = %v, want 0", l)
	}
	if s != 0 {
		t.Errorf("size = %v, want 0", s)
	}
	if capacity != 10 {
		t.Errorf("capacity = %v, want 10", capacity)
	}
}

func testInsert(t *testing.T) {
	cache := LRUCache(100)
	data := json.Marshal(&cValue{0})
	cache.Insert("key", string(data))

	v, ok := cache.Get("key")
	if !ok || v(*cValue) != data {
		t.Errorf("Cache has incorrect value: %v != %v", data, v)
	}
}

func testLRUEviction(t *testing.T) {
	size := uint64(3)
	cache := LRUCache(size)

	cache.Insert("key1", &cValue{1})
	cache.Insert("key2", &cValue{1})
	cache.Insert("key3", &cValue{1})
	// lru: [key3, key2, key1]

	// Look up the elements. This will rearrange the LRU ordering.
	cache.Get("key3")
	cache.Get("key2")
	cache.Get("key1")
	// lru: [key1, key2, key3]

	cache.Insert("key0", &cValue{1})
	// lru: [key0, key1, key2]

	// The least recently used one should have been evicted.
	if _, ok := cache.Get("key3"); ok {
		t.Error("Least recently used element was not evicted.")
	}
}

func testCapacityIsObeyed(t *testing.T) {
	size := uint64(3)
	cache := LRUCache(size)
	value := &cValue{1}

	// Insert up to the cache's capacity.
	cache.Insert("key1", value)
	cache.Insert("key2", value)
	cache.Insert("key3", value)
	if _, sz, _, _ := cache.CacheStats(); sz != size {
		t.Errorf("cache.Size() = %v, expected %v", sz, size)
	}
	// Insert one more; something should be evicted to make room.
	cache.Insert("key4", value)
	if _, sz, _ := cache.CacheStats(); sz != size {
		t.Errorf("post-evict cache.Size() = %v, expected %v", sz, size)
	}
}

func testInsertUpdatesSize(t *testing.T) {
	cache := LRUCache(100)
	emptyValue := &cValue{0}
	key := "key1"
	cache.Insert(key, emptyValue)
	if _, sz, _ := cache.CacheStats(); sz != 0 {
		t.Errorf("cache.Size() = %v, expected 0", sz)
	}
	someValue := &cValue{20}
	key = "key2"
	cache.Insert(key, someValue)
	if _, sz, _ := cache.CacheStats(); sz != 20 {
		t.Errorf("cache.Size() = %v, expected 20", sz)
	}
}

func testInsertWithOldKeyUpdatesValue(t *testing.T) {
	cache := LRUCache(100)
	emptyValue := &cValue{0}
	key := "key1"
	cache.Insert(key, emptyValue)
	someValue := &cValue{20}
	cache.Insert(key, someValue)

	v, ok := cache.Get(key)
	if !ok || v.(*cValue) != someValue {
		t.Errorf("Cache has incorrect value: %v != %v", someValue, v)
	}
}
