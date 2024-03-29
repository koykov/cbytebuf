package cbytebuf

import (
	"sync"
)

// Pool is a pool implementation based on sync.Pool.
type Pool struct {
	p sync.Pool
}

// P is a default instance of native pool for simple cases.
// Just call cbytebuf.Acquire() and cbytebuf.Release().
var P Pool

// Get old byte buffer from the pool or create a new byte buffer.
func (p *Pool) Get() *CByteBuf {
	v := p.p.Get()
	if v != nil {
		if b, ok := v.(*CByteBuf); ok {
			return b
		}
	}
	return &CByteBuf{}
}

// Put byte buffer back to the pool.
//
// Using data returned from the buffer after putting is unsafe.
func (p *Pool) Put(b *CByteBuf) {
	if b.h.Data == 0 {
		return
	}
	b.Reset()
	p.p.Put(b)
}

// Acquire gets byte buffer from default pool instance.
func Acquire() *CByteBuf {
	return P.Get()
}

// Release puts byte buffer back to default pool instance.
func Release(b *CByteBuf) {
	P.Put(b)
}
