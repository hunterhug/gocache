package gocache

import (
	"github.com/hunterhug/gocache/algorithm"
	"sync"
	"time"
)

type cache struct {
	minHeap *algorithm.Heap
	treeMap algorithm.TreeMap
	close   bool
	locker  sync.Mutex
}

type cacheItem struct {
	RawByte                      []byte
	Raw                          interface{}
	expireUnixNanosecondDateTime int64
}

func (i *cacheItem) GetExpireUnixNanosecondDateTime() int64 {
	return i.expireUnixNanosecondDateTime
}

func (i *cacheItem) IsExpire() bool {
	return i.expireUnixNanosecondDateTime <= time.Now().UnixNano()
}

func (c *cache) loopCleanExpireItem() {
	timer := time.NewTimer(time.Second)
	for {
		if c.close {
			timer.Stop()
			return
		}
		select {
		case <-timer.C:
			c.cleanOlder()
			timer.Reset(time.Second)
		}
	}
}

func (c *cache) cleanOlder() {
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.close {
		c.minHeap = nil
		return
	}

	i := 0
	for i < 30 {
		min := c.minHeap.Min()
		if min == nil {
			return
		}

		if min.Value > time.Now().UnixNano() {
			return
		}

		c.treeMap.Delete(min.Key)
		c.minHeap.Pop()
		i++
	}
}

func (c *cache) ShutDown() {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.close = true
}

func (c *cache) Set(key string, value []byte, expireTime time.Duration) {
	item := cacheItem{
		RawByte: value,
	}

	c.set(key, item, expireTime)
}

func (c *cache) SetInterface(key string, value interface{}, expireTime time.Duration) {
	item := cacheItem{
		Raw: value,
	}

	c.set(key, item, expireTime)
}

func (c *cache) SetByExpireUnixNanosecondDateTime(key string, value []byte, expireUnixNanosecondDateTime int64) {
	item := cacheItem{
		RawByte: value,
	}

	c.setByExpireDateTime(key, item, expireUnixNanosecondDateTime)
}

func (c *cache) SetInterfaceByExpireUnixNanosecondDateTime(key string, value interface{}, expireUnixNanosecondDateTime int64) {
	item := cacheItem{
		Raw: value,
	}

	c.setByExpireDateTime(key, item, expireUnixNanosecondDateTime)
}

func (c *cache) setByExpireDateTime(key string, value cacheItem, expireUnixNanosecondDateTime int64) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.close {
		return
	}

	value.expireUnixNanosecondDateTime = expireUnixNanosecondDateTime

	oldTreeMapValue, exist := c.treeMap.Get(key)
	if !exist {
		innerValue := &algorithm.HeapValue{
			Value: expireUnixNanosecondDateTime,
			Key:   key,
			Extra: value,
		}
		c.treeMap.Put(key, innerValue)
		c.minHeap.Push(innerValue)
		return
	}

	oldTreeMapValueReal := oldTreeMapValue.(*algorithm.HeapValue)
	oldHeapValue := c.minHeap.PopIndex(oldTreeMapValueReal.Index)
	oldHeapValue.Value = expireUnixNanosecondDateTime
	oldHeapValue.Extra = value
	c.minHeap.Push(oldHeapValue)
}

func (c *cache) set(key string, value cacheItem, expireTime time.Duration) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.close {
		return
	}

	expireUnixNanosecondDateTime := time.Now().UnixNano() + int64(expireTime/time.Nanosecond)
	value.expireUnixNanosecondDateTime = expireUnixNanosecondDateTime

	oldTreeMapValue, exist := c.treeMap.Get(key)
	if !exist {
		innerValue := &algorithm.HeapValue{
			Value: expireUnixNanosecondDateTime,
			Key:   key,
			Extra: value,
		}
		c.treeMap.Put(key, innerValue)
		c.minHeap.Push(innerValue)
		return
	}

	oldTreeMapValueReal := oldTreeMapValue.(*algorithm.HeapValue)
	oldHeapValue := c.minHeap.PopIndex(oldTreeMapValueReal.Index)
	oldHeapValue.Value = expireUnixNanosecondDateTime
	oldHeapValue.Extra = value
	c.minHeap.Push(oldHeapValue)
}

func (c *cache) Delete(key string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.close {
		return
	}

	treeMapValue, exist := c.treeMap.Get(key)
	if !exist {
		return
	}

	treeMapValueReal := treeMapValue.(*algorithm.HeapValue)
	c.minHeap.PopIndex(treeMapValueReal.Index)
	c.treeMap.Delete(key)
}

func (c *cache) get(key string) (value *cacheItem, exist bool) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.close {
		return
	}

	treeMapValue, exist := c.treeMap.Get(key)
	if !exist {
		return nil, false
	}

	treeMapValueReal := treeMapValue.(*algorithm.HeapValue)
	item := treeMapValueReal.Extra.(cacheItem)
	if item.IsExpire() {
		c.minHeap.PopIndex(treeMapValueReal.Index)
		c.treeMap.Delete(key)
		return nil, false
	}

	return &item, true
}

func (c *cache) Get(key string) (value []byte, expireUnixNanosecondDateTime int64, exist bool) {
	result, has := c.get(key)
	if has {
		return result.RawByte, result.expireUnixNanosecondDateTime, true
	}

	return
}

func (c *cache) GetInterface(key string) (value interface{}, expireUnixNanosecondDateTime int64, exist bool) {
	result, has := c.get(key)
	if has {
		return result.Raw, result.expireUnixNanosecondDateTime, true
	}

	return
}

func (c *cache) Size() int64 {
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.close {
		return 0
	}

	return int64(c.minHeap.Size())
}

func (c *cache) GetOldestKey() (key string, expireUnixNanosecondDateTime int64, exist bool) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.close {
		return
	}

	min := c.minHeap.Min()
	if min != nil {
		return min.Key, min.Value, true
	}

	return
}

func (c *cache) KeyList() []string {
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.close {
		return nil
	}

	return c.treeMap.KeyList()
}
