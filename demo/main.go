package main

import (
	"fmt"
	"github.com/hunterhug/gocache"
	"time"
)

func main() {
	cache := gocache.New()
	defer cache.ShutDown()

	cache.Set("a", []byte("a hi"), 2*time.Second)
	cache.Set("b", []byte("b hi"), 2*time.Second)
	cache.SetInterface("c", []byte("c hi"), 20*time.Second)

	fmt.Println(cache.Size())
	fmt.Println(cache.GetOldestKey())
	fmt.Println(cache.KeyList())
	fmt.Println(cache.Get("a"))
	fmt.Println(cache.GetInterface("c"))

	time.Sleep(2 * time.Second)
	fmt.Println(cache.Get("a"))

	for k := 0; k <= 100; k++ {
		cache.Set(fmt.Sprintf("%v", k), []byte("hi"), 2*time.Second)
	}

	for _, k := range cache.KeyList() {
		v, expireDate, exist := cache.Get(k)
		fmt.Println(k, exist, expireDate, string(v))
	}
}
