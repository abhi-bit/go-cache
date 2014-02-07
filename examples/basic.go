package main

import (
	"fmt"
	"github.com/abhi-bit/go-lrucache"
)

func main() {
	c := cache.LRUCache(100)
	testKey := "test-key"
	cacheSize := c.Insert(testKey, ("Å i åa ä e "))
	fmt.Println("SET:: CacheSize after insert:", cacheSize)
	doc, _ := c.Get(testKey)
	fmt.Printf("GET::  key: %s value: %s\n", testKey, doc)
	key, doc, size := c.Peek()
	fmt.Println("PEEK:: key:", key, "value:", doc, "size:", size)
	allDocs := c.Keys()
	fmt.Println("Dump of all keys in cache:", allDocs)
}
