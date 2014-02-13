package main

import (
	"crypto/rand"
	"fmt"
	"github.com/abhi-bit/go-lrucache"
	"time"
)

type CacheValue struct {
	size string
}

type MyValue []byte

func (cv *CacheValue) Size() int {
	return len(cv.size)
}

func (mval MyValue) Size() int {
	return cap(mval)
}

func main() {
	c := lrucache.NewLRUCache(64 * 1024 * 1024)
	testKey := "test-key"

	//Benchmarking Get
	value := make(MyValue, 1000)
	c.Set(testKey, value)
	counter := 0
	now := time.Now()

	for time.Now().Sub(now).Seconds() < 2 {
		_, ok := c.Get(testKey)
		counter += 1
		if !ok {
			panic("error")
		}
	}
	fmt.Println("GET:: Requests per sec:", counter/2)

	//Benchmarking Set
	counter = 0
	now = time.Now()
	for time.Now().Sub(now).Seconds() < 2 {
		//Static value
		//c.Set(testKey, value)

		//Random values of key size 5 and val size 10
		c.Set(randString(5), &CacheValue{randString(10)})
		counter += 1
	}

	fmt.Println("SET:: Requests per sec:", counter/2)
	stats := c.StatsJSON()
	fmt.Println("Cache Stats:", stats)
}

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
