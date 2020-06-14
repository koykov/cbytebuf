package cbytebuf

import "github.com/koykov/lbpool"

type LBPool struct {
	Size uint
	p    lbpool.Pool
}

// Default instance of the LB pool for simple cases.
// Just call cbytebuf.LBAcquire() and cbytebuf.LBRelease().
var (
	LBP = LBPool{Size: 1000}

	_, _ = LBAcquire, LBRelease
)

// Get old byte buffer from the LB pool or create a new byte buffer.
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

// Put byte buffer back to the LB pool.
//
// Using data returned from the buffer after putting is unsafe.
func (p *LBPool) Put(b *CByteBuf) {
	if b.h.Data == 0 {
		return
	}
	b.ResetLen()
	p.p.Put(b)
}

// Get byte buffer from default LB pool instance.
func LBAcquire() *CByteBuf {
	return LBP.Get()
}

// Put byte buffer back to default LB pool instance.
func LBRelease(b *CByteBuf) {
	LBP.Put(b)
}
