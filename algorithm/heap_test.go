package algorithm

import (
	"encoding/json"
	"fmt"
	"testing"
)

func ToJsonString(object interface{}) string {
	js, _ := json.MarshalIndent(object, "", "\t")
	return string(js)
}

func TestNewHeap(t *testing.T) {
	var a []*HeapValue
	a = append(a,
		&HeapValue{
			Value: 9,
			Key:   "91",
			Extra: nil,
		},
		&HeapValue{
			Value: 9,
			Key:   "92",
			Extra: nil,
		},
		&HeapValue{
			Value: 6,
			Key:   "61",
			Extra: nil,
		}, &HeapValue{
			Value: 19,
			Key:   "",
			Extra: nil,
		}, &HeapValue{
			Value: 6,
			Key:   "62",
			Extra: nil,
		}, &HeapValue{
			Value: 60,
			Key:   "",
			Extra: nil,
		}, &HeapValue{
			Value: 6,
			Key:   "63",
			Extra: nil,
		},
	)
	h := NewMinHeap(a)
	fmt.Println(h.Size(), ToJsonString(h.array))

	fmt.Println(h.PopIndex(1))

	v := h.Pop()
	for v != nil {
		fmt.Println(v)
		v = h.Pop()
	}

	min := h.Min()
	fmt.Println(min)

	h.Push(&HeapValue{
		Value: 6,
		Key:   "63",
		Extra: nil,
	})

	h.Pop()
	h.Pop()
	fmt.Println(ToJsonString(h.array))
}

func TestNewHeap2(t *testing.T) {
	h := NewMinHeap(nil)
	h.Push(&HeapValue{
		Value: 63,
		Key:   "63",
		Extra: nil,
	})
	h.Push(&HeapValue{
		Value: 363,
		Key:   "363",
		Extra: nil,
	})
	h.Push(&HeapValue{
		Value: 1,
		Key:   "1",
		Extra: nil,
	})
	h.Push(&HeapValue{
		Value: 13,
		Key:   "13",
		Extra: nil,
	})

	fmt.Println(h.Min())
	fmt.Println(h.Pop())
	fmt.Println(h.Min())
	fmt.Println(h.Pop())
	fmt.Println(h.Min())
	fmt.Println(h.Pop())
	fmt.Println(h.Min())
	fmt.Println(h.Pop())
	fmt.Println(h.Min())
	fmt.Println(h.Pop())

	for i := 10000; i > 1; i-- {
		h.Push(&HeapValue{
			Value: int64(i),
			Key:   "",
			Extra: nil,
		})
	}

	fmt.Println(cap(h.array))

	fmt.Println(h.Min())
	//fmt.Println(h.Size())
	fmt.Println(h.Pop())
	//fmt.Println(h.Size())

	for i := 10000; i > 1; i-- {
		h.Pop()
		fmt.Println(cap(h.array))
	}

	fmt.Println(cap(h.array))
}
