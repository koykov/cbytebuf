package cbytebuf

import (
	"sync"

	"github.com/koykov/lbpool"
)

// Simple byte buffer pool.
type Pool struct {
	p sync.Pool
}

type LBPool struct {
	Size uint
	p    lbpool.Pool
}

// Default instance of the pool for simple cases.
// Just call cbytebuf.P.Get() and cbytebuf.P.Put().
var (
	P   Pool
	LBP = LBPool{Size: 1000}
)

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

// Get old byte buffer from the pool or create a new byte buffer.
func (p *LBPool) Get() *CByteBuf {
	p.p.Size = p.Size
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
func (p *LBPool) Put(b *CByteBuf) {
	if b.h.Data == 0 {
		return
	}
	b.ResetLen()
	p.p.Put(b)
}
