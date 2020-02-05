package cbytebuf

import "testing"

func TestPool(t *testing.T) {
	var p Pool
	b := p.Get()
	_, _ = b.WriteString("foobar")
	p.Put(b)
}

func BenchmarkPool(b *testing.B) {
	var p Pool
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b := p.Get()
		_, _ = b.WriteString("foobar")
		p.Put(b)
	}
}
