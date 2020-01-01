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

const (
	// Limit to switch to loop rolling write.
	shortInputThreshold = 256
	// Buffer size limit to use malloc to grow.
	mallocGrowThreshold = 1024
)

// Variable-size alloc-free buffer.
// Also no escapes to the heap since buffer doesn't contain any pointer.
type CByteBuf struct {
	// Header to manipulate buffer size and fast slice construction.
	h reflect.SliceHeader
	// Helper temporary variable.
	t int
}

// MarshalerTo interface to write struct like Protobuf.
type MarshalerTo interface {
	Size() int
	MarshalTo(data []byte) (int, error)
}

var (
	// Error constants.
	ErrOk           error = nil
	ErrBadAlloc           = errors.New("bad alloc on buffer init or grow")
	ErrNegativeCap        = errors.New("negative cap on the grow")
	ErrNegativeRead       = errors.New("reader returned negative count from Read")
)

// Shorthand buffer make func.
func NewCByteBuf() *CByteBuf {
	b := CByteBuf{}
	return &b
}

// Get length of the buffer.
func (b *CByteBuf) Len() int {
	return b.h.Len
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
		if b.t == b.h.Cap {
			if err = b.Grow(b.h.Cap * 2); err != nil {
				return 0, err
			}
		}
		b.t, err = r.Read(b.Bytes()[b.t:])
		if b.t < 0 {
			return n, ErrNegativeRead
		}
		b.h.Len += b.t
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
	n, err := w.Write(b.Bytes())
	return int64(n), err
}

// Implement io.Writer.
func (b *CByteBuf) Write(data []byte) (int, error) {
	b.t = len(data)
	if b.h.Data == 0 {
		// First write, need to create internal byte slice.
		// Check write after reset.
		if b.h.Cap == 0 {
			b.h.Cap = b.t * 2
		}
		// Create underlying byte array in the C memory, outside of GC's eyes.
		b.h.Data = uintptr(C.cbb_init(C.int(b.h.Cap)))
		if b.h.Data == 0 {
			return 0, ErrBadAlloc
		}
	}

	if b.h.Len+b.t > b.h.Cap {
		// Increase capacity of the byte array due to not enough space in it.
		err := b.Grow((b.h.Len + b.t) * 2)
		if err != nil {
			return 0, err
		}
	}

	// Add data to the slice.
	if b.t > shortInputThreshold {
		// Write long data using loop rolling.
		for b.t >= 8 {
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len))) = data[0]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+1))) = data[1]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+2))) = data[2]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+3))) = data[3]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+4))) = data[4]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+5))) = data[5]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+6))) = data[6]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+7))) = data[7]
			b.h.Len += 8
			b.t -= 8
			data = data[8:]
		}
		for b.t >= 4 {
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len))) = data[0]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+1))) = data[1]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+2))) = data[2]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+3))) = data[3]
			b.h.Len += 4
			b.t -= 4
			data = data[4:]
		}
		for b.t >= 2 {
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len))) = data[0]
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+1))) = data[1]
			b.h.Len += 2
			b.t -= 2
			data = data[2:]
		}
		if b.t > 0 {
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len))) = data[0]
			b.h.Len++
			b.t--
		}
	} else {
		for i := 0; i < b.t; i++ {
			*(*byte)(unsafe.Pointer(b.h.Data + uintptr(b.h.Len+i))) = data[i]
		}
		b.h.Len += b.t
	}

	return b.t, ErrOk
}

// Marshal data of struct implemented MarshalerTo interface.
func (b *CByteBuf) WriteMarshalerTo(m MarshalerTo) (int, error) {
	err := b.Grow(m.Size())
	if err != nil {
		return 0, err
	}
	b.h.Len = b.h.Cap
	return m.MarshalTo(b.Bytes())
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
	if b.h.Data != 0 && b.h.Cap >= cap {
		// Buffer has necessary space actually.
		return ErrOk
	}
	// Save new capacity.
	b.h.Cap = cap
	if b.h.Len > b.h.Cap {
		// Trim length to new capacity.
		b.h.Len = b.h.Cap
	}

	// Allocate memory.
	if b.h.Data == 0 {
		// Grow after reset detected.
		// Allocate underlying byte array in C memory.
		b.h.Data = uintptr(C.cbb_init(C.int(b.h.Cap)))
	} else {
		// Reallocate underlying byte array in C memory.
		// All necessary copying/free will perform implicitly, don't worry about this.
		// Using combination of malloc()+memcpy()+free() to grow for short buffers is more efficient than simple using
		// of realloc().
		if b.h.Len > mallocGrowThreshold {
			b.h.Data = uintptr(C.cbb_grow_r(C.ulong(b.h.Data), C.int(b.h.Cap)))
		} else {
			b.h.Data = uintptr(C.cbb_grow_m(C.ulong(b.h.Data), C.int(b.h.Len), C.int(b.h.Cap)))
		}
	}
	if b.h.Data == 0 {
		return ErrBadAlloc
	}
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
	return *(*[]byte)(unsafe.Pointer(&b.h))
}

// Append buffer value to destination and return it.
func (b *CByteBuf) AppendBytes(dst []byte) []byte {
	return append(dst[:0], b.Bytes()...)
}

// Get the contents of the buffer as string.
func (b *CByteBuf) String() string {
	return fastconv.B2S(b.Bytes())
}

// Append buffer value to destination string and return it.
func (b *CByteBuf) AppendString(dst string) string {
	return fastconv.B2S(append(fastconv.S2B(dst)[:0], b.Bytes()...))
}

// Reset buffer length without releasing memory.
func (b *CByteBuf) ResetLen() {
	b.h.Len = 0
}

// Reset all data accumulated in buffer.
//
// Using the buffer data after call this func may crash your app.
// Buffer capacity keeps to reduce amount of further CGO calls.
func (b *CByteBuf) Reset() {
	b.release()
	b.h.Len = 0
}

// Manually release of the underlying byte array.
//
// Using the buffer data after call this func may crash your app.
// This method truncates buffer's capacity.
func (b *CByteBuf) Release() {
	b.release()
	b.h.Len = 0
	b.h.Cap = 0
}

// Internal release method.
func (b *CByteBuf) release() {
	if b.h.Data == 0 {
		return
	}
	// Free memory and reset pointer.
	C.cbb_release(C.ulong(b.h.Data))
	b.h.Data = 0
}
