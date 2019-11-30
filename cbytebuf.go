package cbytebuf

/*
#include "stdlib.h"
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
	errOk          = 0
	errBadAlloc    = 1
	errNegativeCap = 2
)

type CByteBuf struct {
	sh reflect.SliceHeader
	b  []byte
	l  int
	e  uint
}

var errs = []error{
	errOk:          nil,
	errBadAlloc:    errors.New("bad alloc on buffer init or grow"),
	errNegativeCap: errors.New("negative cap on the grow"),
}

func NewCByteBuf() (*CByteBuf, error) {
	b := CByteBuf{}
	return &b, nil
}

// Get length of the buffer.
func (b *CByteBuf) Len() int {
	return b.sh.Len
}

// Implement io.WriterTo.
func (b *CByteBuf) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.b)
	return int64(n), err
}

// Implement io.Writer.
func (b *CByteBuf) Write(data []byte) (int, error) {
	b.l = len(data)
	if b.b == nil {
		// First write, need to create internal byte slice.
		b.sh.Cap = b.l * 2
		ptrErr := (*C.uint)(unsafe.Pointer(&b.e))
		ptrAddr := (*C.uintptr)(unsafe.Pointer(&b.sh.Data))
		ptrCap := (*C.int)(unsafe.Pointer(&b.sh.Cap))
		// Create underlying byte array in the C memory, outside of GC's eyes.
		C.cbb_init(ptrErr, ptrAddr, ptrCap)
		// Manually create the byte slice.
		b.b = *(*[]byte)(unsafe.Pointer(&b.sh))
	}

	if b.sh.Len+b.l > b.sh.Cap {
		// Increase capacity of the byte array due to not enough space in it.
		_ = b.Grow((b.sh.Len + b.l) * 2)
	}

	// Add data to the slice.
	b.b = append(b.b, data...)
	// Increase internal len for further grows.
	b.sh.Len += b.l

	return b.l, errs[b.e]
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
	return b.Write(fastconv.StringToBytes(s))
}

// Increase or decrease capacity of the buffer.
func (b *CByteBuf) Grow(cap int) error {
	if cap < 0 {
		return errs[errNegativeCap]
	}
	// Save new capacity.
	b.sh.Cap = cap
	ptrErr := (*C.uint)(unsafe.Pointer(&b.e))
	ptrAddr := (*C.uintptr)(unsafe.Pointer(&b.sh.Data))
	ptrCap := (*C.int)(unsafe.Pointer(&b.sh.Cap))
	// Reallocate underlying byte array in C memory.
	// New array may overlap with the previous if it's possible to resize it (there is free space at the right side).
	// All necessary copying/free will perform implicitly, don't worry about this.
	C.cbb_grow(ptrErr, ptrAddr, ptrCap)
	// Recreate the slice (old accumulated data keeps).
	b.b = *(*[]byte)(unsafe.Pointer(&b.sh))
	return errs[b.e]
}

// Get the contents of the buffer.
func (b *CByteBuf) Bytes() []byte {
	return b.b
}

// Get the contents of the buffer as string.
func (b *CByteBuf) String() string {
	return fastconv.BytesToString(b.b)
}

// Reset all data accumulated in buffer.
func (b *CByteBuf) Reset() {
	b.b = b.b[:0]
}

// Manually release of the underlying byte array.
//
// Using the buffer data after call this func may crash your app.
func (b *CByteBuf) Release() error {
	ptrErr := (*C.uint)(unsafe.Pointer(&b.e))
	ptrAddr := (*C.uintptr)(unsafe.Pointer(&b.sh.Data))
	// Free memory.
	C.cbb_release(ptrErr, ptrAddr)
	// Truncate length and capacity.
	b.sh.Len, b.sh.Cap = 0, 0
	// Slice is broken here, therefore kill it.
	b.b = nil
	return errs[b.e]
}
