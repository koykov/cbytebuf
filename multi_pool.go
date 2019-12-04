package cbytebuf

import "math/rand"

// Default multi-pool size.
const defaultMultiPollSize uint = 16

// Special pool for highly load systems to reduce runtime.futex overwork.
//
// It's just a collection of simple pools with random choose between them.
type MultiPool struct {
	sz uint
	p  []Pool
}

// Default instance of the multi-pool.
// Just call cbytebuf.MP.Get() and cbytebuf.MP.Put().
var MP MultiPool

// Check and initialize multi-pool if needed.
func (m *MultiPool) initMP() {
	if m.sz == 0 {
		m.sz = defaultMultiPollSize
	}
	if m.p == nil {
		m.p = make([]Pool, m.sz, m.sz)
	}
}

// Get byte buffer from one of pools or create a new one.
func (m *MultiPool) Get() *CByteBuf {
	m.initMP()
	if b := m.p[rand.Intn(int(m.sz))].Get(); b != nil {
		return b
	}
	return &CByteBuf{}
}

// Put byte buffer back to one of pools.
func (m *MultiPool) Put(b *CByteBuf) {
	m.initMP()
	b.Reset()
	m.p[rand.Intn(int(m.sz))].Put(b)
}
