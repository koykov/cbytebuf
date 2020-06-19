package cbytebuf

import "testing"

func TestLBPool(t *testing.T) {
	var p = LBPool{Size: 10}
	b := p.Get()
	_, _ = b.WriteString("foobar")
	p.Put(b)
}

func BenchmarkLBPool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b := LBAcquire()
		_, _ = b.WriteString("foobar")
		LBRelease(b)
	}
}
