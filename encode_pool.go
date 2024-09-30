package gojay

import (
	"io"
	"sync"
)

var encPool = sync.Pool{
	New: func() any {
		return NewEncoder(nil)
	},
}

func init() {
	for range 32 {
		encPool.Put(NewEncoder(nil))
	}
}

// NewEncoder returns a new encoder or borrows one from the pool.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// BorrowEncoder borrows an Encoder from the pool.
func BorrowEncoder(w io.Writer) *Encoder {
	//nolint:forcetypeassert
	enc := encPool.Get().(*Encoder)
	enc.w = w
	enc.buf = enc.buf[:0]
	enc.isPooled = 0
	enc.err = nil
	enc.hasKeys = false
	enc.keys = nil
	return enc
}

// Release sends back a Encoder to the pool.
func (enc *Encoder) Release() {
	enc.isPooled = 1
	encPool.Put(enc)
}
