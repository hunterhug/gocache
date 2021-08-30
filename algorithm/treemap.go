package algorithm

import "strings"

type comparator func(key1, key2 string) int64

// TreeMap method
// design to be concurrent safe
// should support int key?
type TreeMap interface {
	Put(key string, value interface{})                            // put key pairs
	Delete(key string)                                            // delete a key
	Get(key string) (value interface{}, exist bool)               // get value from key
	GetInt(key string) (value int, exist bool, err error)         // get value auto change to Int
	GetInt64(key string) (value int64, exist bool, err error)     // get value auto change to Int64
	GetString(key string) (value string, exist bool, err error)   // get value auto change to string
	GetFloat64(key string) (value float64, exist bool, err error) // get value auto change to string
	GetBytes(key string) (value []byte, exist bool, err error)    // get value auto change to []byte
	Contains(key string) (exist bool)                             // map contains key?
	Len() int64                                                   // map key pairs num
	KeyList() []string                                            // map key out to list from top to bottom which is layer order
	KeySortedList() []string                                      // map key out to list sorted
	Iterator() TreeMapIterator                                    // map iterator, iterator from top to bottom which is layer order
	MaxKey() (key string, value interface{}, exist bool)          // find max key pairs
	MinKey() (key string, value interface{}, exist bool)          // find min key pairs
	SetComparator(comparator) TreeMap                             // set compare func to control key compare
	Check() bool                                                  // just help
	Height() int64                                                // just help
}

// TreeMapIterator Iterator concurrent not safe
// you should deal by yourself
type TreeMapIterator interface {
	HasNext() bool
	Next() (key string, value interface{})
}

// NewTreeMap default map is rbt implement
func NewTreeMap() TreeMap {
	t := new(rbTree)
	t.c = comparatorDefault
	return t
}

// compare two key
func comparatorDefault(key1, key2 string) int64 {
	return int64(strings.Compare(key1, key2))
}
