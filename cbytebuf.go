package cbytebuf

/*
#include "cbytebuf.h"
*/
import "C"
import (
	"errors"
	"github.com/koykov/fastconv"
	"io"
	"reflect"
	"unsafe"
)

const shortInputThreshold = 256

// Variable-size alloc-free buffer.
type CByteBuf struct {
	h reflect.SliceHeader // header to fast slice construct
	b []byte              // buffer slice
	l int                 // actual length
	t int                 // temporary int var
}

var (
	ErrOk           error = nil
	ErrBadAlloc           = errors.New("bad alloc on buffer init or grow")
	ErrNegativeCap        = errors.New("negative cap on the grow")
	ErrNegativeRead       = errors.New("reader returned negative count from Read")
)

func NewCByteBuf() (*CByteBuf, error) {
	b := CByteBuf{}
	return &b, nil
}

// Get length of the buffer.
func (b *CByteBuf) Len() int {
	return b.l
}

// Get capacity of the buffer.
func (b *CByteBuf) Cap() int {
	return b.h.Cap
}

// Implement io.ReaderFrom.
func (b *CByteBuf) ReadFrom(r io.Reader) (n int64, err error) {
	if b.h.Cap == 0 {
		if err = b.Grow(64); err != nil {
			return 0, err
		}
	}
	for {
		if b.l == b.h.Cap {
			if err = b.Grow(b.h.Cap * 2); err != nil {
				return 0, err
			}
		}
		b.t, err = r.Read(b.b[b.l:])
		if b.t < 0 {
			return n, ErrNegativeRead
		}
		b.l += b.t
		n += int64(b.t)
		if err == io.EOF {
			return n, nil
		}
		if err != nil {
			return n, err
		}
	}
}

// Implement io.WriterTo.
func (b *CByteBuf) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.b)
	return int64(n), err
}

// Implement io.Writer.
func (b *CByteBuf) Write(data []byte) (int, error) {
	b.t = len(data)
	if b.b == nil {
		// First write, need to create internal byte slice.
		b.h.Cap = b.t * 2
		b.h.Len = b.h.Cap
		// Create underlying byte array in the C memory, outside of GC's eyes.
		//C.cbb_init((*C.uint)(unsafe.Pointer(&b.e)), (*C.uintptr)(unsafe.Pointer(&b.h.Data)), (*C.int)(unsafe.Pointer(&b.h.Cap)))
		b.h.Data = uintptr(C.cbb_init_np(C.int(b.h.Cap)))
		if b.h.Data == 0 {
			return 0, ErrBadAlloc
		}
		// Manually create the byte slice.
		b.b = *(*[]byte)(unsafe.Pointer(&b.h))
	}

	if b.l+b.t > b.h.Cap {
		// Increase capacity of the byte array due to not enough space in it.
		err := b.Grow((b.l + b.t) * 2)
		if err != nil {
			return 0, err
		}
	}

	// Add data to the slice.
	//b.b = append(b.b, data...)
	if b.t > shortInputThreshold {
		for len(data) >= 8 {
			b.b[b.l], b.b[b.l+1], b.b[b.l+2], b.b[b.l+3], b.b[b.l+4], b.b[b.l+5], b.b[b.l+6], b.b[b.l+7] =
				data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7]
			b.l += 8
			b.t -= 8
			data = data[8:]
		}
		for len(data) >= 4 {
			b.b[b.l], b.b[b.l+1], b.b[b.l+2], b.b[b.l+3] = data[0], data[1], data[2], data[3]
			b.l += 4
			b.t -= 4
			data = data[4:]
		}
		for len(data) >= 2 {
			b.b[b.l], b.b[b.l+1] = data[0], data[1]
			b.l += 2
			b.t -= 2
			data = data[2:]
		}
		if b.t > 0 {
			b.b[b.l] = data[0]
			b.l++
			b.t--
		}
	} else {
		for i := 0; i < b.t; i++ {
			b.b[b.l+i] = data[i]
		}
	}
	// Increase internal len for further grows.
	b.l += b.t

	return b.t, ErrOk
}

// Write single byte in the buffer.
//
// Implement io.ByteWriter.
func (b *CByteBuf) WriteByte(c byte) error {
	_, err := b.Write([]byte{c})
	return err
}

// Write string in the buffer.
//
// String will be convert to byte slice on the fly.
func (b *CByteBuf) WriteString(s string) (int, error) {
	return b.Write(fastconv.S2B(s))
}

// Increase or decrease capacity of the buffer.
func (b *CByteBuf) Grow(cap int) error {
	if cap < 0 {
		return ErrNegativeCap
	}
	if b.b != nil && b.h.Cap >= cap {
		return ErrOk
	}
	// Save new capacity.
	b.h.Cap = cap
	b.h.Len = cap
	// Reallocate underlying byte array in C memory.
	// New array may overlap with the previous if it's possible to resize it (there is free space at the right side).
	// All necessary copying/free will perform implicitly, don't worry about this.
	//C.cbb_grow((*C.uint)(unsafe.Pointer(&b.e)), (*C.uintptr)(unsafe.Pointer(&b.h.Data)), (*C.int)(unsafe.Pointer(&b.h.Cap)))

	//b.h.Data = uintptr(C.cbb_grow_np(C.ulong(b.h.Data), C.int(b.h.Cap)))

	b.h.Data = uintptr(C.cbb_grow_np1(C.ulong(b.h.Data), C.int(b.l), C.int(b.h.Cap)))

	if b.h.Data == 0 {
		return ErrBadAlloc
	}
	// Recreate the slice (old accumulated data keeps).
	b.b = *(*[]byte)(unsafe.Pointer(&b.h))
	return ErrOk
}

// Increase or decrease capacity of the buffer using delta value.
//
// Delta may be negative, but if delta will less than -capacity, the error will be triggered.
func (b *CByteBuf) GrowDelta(delta int) error {
	return b.Grow(b.h.Cap + delta)
}

// Get the contents of the buffer.
func (b *CByteBuf) Bytes() []byte {
	return b.b[:b.l]
}

// Get the contents of the buffer as string.
func (b *CByteBuf) String() string {
	return fastconv.B2S(b.b)
}

// Reset all data accumulated in buffer.
func (b *CByteBuf) Reset() {
	b.l = 0
}

// Manually release of the underlying byte array.
//
// Using the buffer data after call this func may crash your app.
func (b *CByteBuf) Release() error {
	// Free memory.
	//C.cbb_release((*C.uint)(unsafe.Pointer(&b.e)), (*C.uintptr)(unsafe.Pointer(&b.h.Data)))
	C.cbb_release_np(C.ulong(b.h.Data))
	// Truncate length and capacity.
	b.h.Data = 0
	b.h.Len, b.h.Cap, b.l = 0, 0, 0
	// Slice is broken here, therefore kill it.
	b.b = nil
	return ErrOk
}
