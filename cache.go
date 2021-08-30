package gocache

import (
	"github.com/hunterhug/gocache/algorithm"
	"time"
)

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

func New() Cache {
	c := new(cache)
	c.treeMap = algorithm.NewTreeMap()
	c.minHeap = algorithm.NewMinHeap(nil)

	go c.loopCleanExpireItem()
	return c
}
