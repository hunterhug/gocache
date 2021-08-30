package gocache

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c := New()
	defer c.ShutDown()

	expireTime := 10 * time.Second
	key := "a"
	value := []byte("value")

	c.Set(key, value, expireTime)

	v, expireMillDateTime, exist := c.Get(key)
	if exist {
		fmt.Println(key, "is", string(v), expireMillDateTime, exist)
	} else {
		fmt.Println(key, " not exist")
	}

	for i := 0; i < 5; i++ {
		c.Set(fmt.Sprintf("%s-%d", key, i), []byte(fmt.Sprintf("hi:%d", i)), 5*time.Second)
	}

	oldK, oldExpireMillDateTime, oldExist := c.GetOldestKey()
	fmt.Println("oldest:", oldK, oldExpireMillDateTime, oldExist)

	for _, k := range c.KeyList() {
		v, expireMillDateTime, exist = c.Get(k)
		if exist {
			fmt.Println(k, "is", string(v), expireMillDateTime, exist)
		} else {
			fmt.Println(k, "not exist")
		}
	}

	time.Sleep(10 * time.Second)

	v, expireMillDateTime, exist = c.Get(key)
	if exist {
		fmt.Println(key, "is", string(v), expireMillDateTime, exist)
	} else {
		fmt.Println(key, "not exist")
	}
}
