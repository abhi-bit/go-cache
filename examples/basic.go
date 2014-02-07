package main

import (
    "github.com/abhi-bit/go-lrucache"
    "fmt"
)

func main() {
    c := cache.New(100)
    testKey := "test-key"
    cacheSize := c.Insert(testKey, []rune("Å i åa ä e "))
    fmt.Println("SET:: CacheSize after insert:", cacheSize)
    doc, _ := c.Get(testKey)
    fmt.Printf("GET::  key: %s value: %s\n", testKey, string(doc))
    key, doc, size := c.Peek()
    fmt.Println("PEEK:: key:",key, "value:", string(doc), "size:",size)
    allDocs := c.Keys()
    fmt.Println("Dump of all keys in cache:", allDocs)
}
