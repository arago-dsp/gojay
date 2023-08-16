package gojay

import (
	"sync"
)

var strPool = sync.Pool{
	New: func() interface{} { return new(string) },
}

func init() {
	for i := 0; i < 32; i++ {
		strPool.Put(new(string))
	}
}

// AcquireString gets a string pointer from the pool.
func AcquireString() *string {
	return strPool.Get().(*string)
}

// ReleaseString sends back a string to the pool.
func ReleaseString(s *string) {
	*s = (*s)[:0]
	strPool.Put(s)
}
