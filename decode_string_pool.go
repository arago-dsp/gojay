package gojay

import (
	"sync"
)

var strPool = sync.Pool{
	New: func() any { return new(string) },
}

func init() {
	for range 32 {
		strPool.Put(new(string))
	}
}

// AcquireString gets a string pointer from the pool.
func AcquireString() *string {
	//nolint:forcetypeassert
	return strPool.Get().(*string)
}

// ReleaseString sends back a string to the pool.
func ReleaseString(s *string) {
	*s = (*s)[:0]
	strPool.Put(s)
}
