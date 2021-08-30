# Memory Cache auto clean expired data implement by Golang

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/gocache.svg?style=social&label=Forks)](https://github.com/hunterhug/gocache/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/gocache.svg?style=social&label=Stars)](https://github.com/hunterhug/gocache/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/gocache.svg)](https://github.com/hunterhug/gocache)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/gocache.svg)](https://github.com/hunterhug/gocache/issues)

[中文说明](/README_ZH.md)

I use Red-Black Tree Map and Minimum Heap keep the data in memory (the root node of Minimum Heap has the oldest value, so we can fast clean the expired values).

Cache in memory implement very efficient, No pre allocated space required.

## Usage

Simple get it by:

```
go get -v github.com/hunterhug/gocache
```

## Demo

Follow the interface method：

```go
type Cache interface {
    Set(key string, value []byte, expireTime time.Duration)
    SetInterface(key string, value interface{}, expireTime time.Duration)
    SetByExpireUnixNanosecondDateTime(key string, value []byte, expireUnixNanosecondDateTime int64)
    SetInterfaceByExpireUnixNanosecondDateTime(key string, value interface{}, expireUnixNanosecondDateTime int64)
    Delete(key string)
    Get(key string) (value []byte, expireUnixNanosecondDateTime int64, exist bool)
    GetInterface(key string) (value interface{}, expireUnixNanosecondDateTime int64, exist bool)
    GetOldestKey() (key string, expireUnixNanosecondDateTime int64, exist bool)
    Size() int64
    KeyList() []string
    ShutDown()
}
```

You can choose to set cache with expireTime `time.Duration = 1 minute` by call `Set(key string, value []byte, expireTime time.Duration)` or other method.

Example:

```go
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
	cache.SetInterface("c", []byte("c hi"), 2*time.Second)

	fmt.Println(cache.Size())
	fmt.Println(cache.GetOldestKey())
	fmt.Println(cache.KeyList())
	fmt.Println(cache.Get("a"))
	fmt.Println(cache.GetInterface("c"))

	time.Sleep(2 * time.Second)
	fmt.Println(cache.Get("a"))
}
```

# License

```
Copyright [2019-2021] [github.com/hunterhug]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```