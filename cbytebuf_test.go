package cbytebuf

import (
	"bytes"
	"math"
	"testing"

	"github.com/koykov/fastconv"
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
	defer buf.Release()

	for _, part := range parts {
		_, _ = buf.Write(part)
		_ = buf.WriteByte(' ')
	}
	b := buf.Bytes()
	if !bytes.Equal(b, expected) {
		t.Error("not equal")
	}
}

func TestCByteBuf_Long(t *testing.T) {
	buf := NewCByteBuf()
	defer buf.Release()

	for i := 0; i < 1000; i++ {
		_, _ = buf.Write(source)
	}
	b := buf.Bytes()
	if !bytes.Equal(b, expectedLong) {
		t.Error("not equal")
	}
}

func TestCByteBuf_AppendBytes(t *testing.T) {
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

func TestCByteBuf_AppendString(t *testing.T) {
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

func TestCByteBuf_WriteInt(t *testing.T) {
	buf := NewCByteBuf()
	_, err := buf.WriteInt(math.MaxInt64)
	if err != nil {
		t.Error(err)
	}
	if buf.String() != "9223372036854775807" {
		t.Error("not equal")
	}
	buf.Release()
}

func TestCByteBuf_WriteUint(t *testing.T) {
	buf := NewCByteBuf()
	_, err := buf.WriteUint(math.MaxUint64)
	if err != nil {
		t.Error(err)
	}
	if buf.String() != "18446744073709551615" {
		t.Error("not equal")
	}
	buf.Release()
}

func TestCByteBuf_WriteFloat(t *testing.T) {
	buf := NewCByteBuf()
	_, err := buf.WriteFloat(math.MaxFloat64, 6)
	if err != nil {
		t.Error(err)
	}
	if buf.String() != "179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000" {
		t.Error("not equal")
	}
	buf.Release()
}

func TestCByteBuf_WriteBool(t *testing.T) {
	buf := NewCByteBuf()
	_, err := buf.WriteBool(false)
	if err != nil {
		t.Error(err)
	}
	if buf.String() != "false" {
		t.Error("not equal")
	}
	buf.Release()
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
