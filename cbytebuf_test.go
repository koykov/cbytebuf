package cbytebuf

import (
	"bytes"
	"github.com/koykov/fastconv"
	"testing"
)

var (
	source = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque euismod ante non arcu " +
		"commodo tempor. Praesent quis nulla sed urna dictum iaculis. Pellentesque malesuada lacinia leo, eu " +
		"hendrerit tellus sodales sit amet. Sed ut finibus purus, ac lacinia metus. Nam tortor nunc, gravida " +
		"hendrerit posuere eu, tristique id elit. Proin id blandit purus. Donec aliquam quam nec erat sodales, eu " +
		"aliquet elit vestibulum. Morbi cursus vehicula semper. Sed dolor lorem, mattis et erat a, elementum " +
		"tincidunt purus. Integer sit amet porta mauris. Curabitur eu est sed augue rutrum tristique et a augue. " +
		"Proin dictum cursus quam vel varius. Duis viverra massa sed lacus gravida, a ullamcorper ipsum iaculis. " +
		"Maecenas interdum congue neque, in ultricies erat ornare id. Suspendisse vitae imperdiet eros.")
	space        = []byte(" ")
	expected     = append(source, space...)
	expectedLong = bytes.Repeat(source, 1000)
	parts        = bytes.Split(source, space)
)

func TestCByteBuf(t *testing.T) {
	buf := NewCByteBuf()
	defer func() {
		buf.Release()
	}()

	for _, part := range parts {
		_, _ = buf.Write(part)
		_ = buf.WriteByte(' ')
	}
	b := buf.Bytes()
	if !bytes.Equal(b, expected) {
		t.Error("not equal")
	}
}

func TestCByteBufLong(t *testing.T) {
	buf := NewCByteBuf()
	defer func() {
		buf.Release()
	}()

	for i := 0; i < 1000; i++ {
		_, _ = buf.Write(source)
	}
	b := buf.Bytes()
	if !bytes.Equal(b, expectedLong) {
		t.Error("not equal")
	}
}

func TestCByteBufAppendBytes(t *testing.T) {
	buf := NewCByteBuf()

	for _, part := range parts {
		_, _ = buf.Write(part)
		_ = buf.WriteByte(' ')
	}
	var b []byte
	b = buf.AppendBytes(b)
	buf.Release()
	if !bytes.Equal(b, expected) {
		t.Error("not equal")
	}
}

func TestCByteBufAppendString(t *testing.T) {
	buf := NewCByteBuf()

	for _, part := range parts {
		_, _ = buf.Write(part)
		_ = buf.WriteByte(' ')
	}
	var s string
	s = buf.AppendString(s)
	buf.Release()
	if !bytes.Equal(fastconv.S2B(s), expected) {
		t.Error("not equal")
	}
}

func BenchmarkCByteBuf_Write(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := NewCByteBuf()
		for _, part := range parts {
			_, _ = buf.Write(part)
			_, _ = buf.Write(space)
		}
		if !bytes.Equal(buf.Bytes(), expected) {
			b.Error("not equal")
		}
		buf.Release()
	}
}

func BenchmarkCByteBuf_WriteLong(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := NewCByteBuf()
		for i := 0; i < 1000; i++ {
			_, _ = buf.Write(source)
		}
		if !bytes.Equal(buf.Bytes(), expectedLong) {
			b.Error("not equal")
		}
		buf.Release()
	}
}

func BenchmarkByteSlice_Append(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0)
		for _, part := range parts {
			buf = append(buf, part...)
			buf = append(buf, ' ')
		}
		if !bytes.Equal(buf, expected) {
			b.Error("not equal")
		}
	}
}

func BenchmarkByteSlice_AppendLong(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0)
		for i := 0; i < 1000; i++ {
			buf = append(buf, source...)
		}
		if !bytes.Equal(buf, expectedLong) {
			b.Error("not equal")
		}
	}
}

func BenchmarkByteBufferNative_Write(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		for _, part := range parts {
			buf.Write(part)
			buf.WriteByte(' ')
		}
		if !bytes.Equal(buf.Bytes(), expected) {
			b.Error("not equal")
		}
		buf.Reset()
	}
}

func BenchmarkByteBufferNative_WriteLong(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		for i := 0; i < 1000; i++ {
			buf.Write(source)
		}
		if !bytes.Equal(buf.Bytes(), expectedLong) {
			b.Error("not equal")
		}
		buf.Reset()
	}
}

func BenchmarkCByteBuf_AppendBytes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := NewCByteBuf()
		_, _ = buf.Write(source)
		var t []byte
		t = buf.AppendBytes(t)
		if !bytes.Equal(t, source) {
			b.Error("not equal")
		}
		buf.Release()
	}
}

func BenchmarkCByteBuf_AppendString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := NewCByteBuf()
		_, _ = buf.Write(source)
		var t string
		t = buf.AppendString(t)
		if !bytes.Equal(fastconv.S2B(t), source) {
			b.Error("not equal")
		}
		buf.Release()
	}
}
