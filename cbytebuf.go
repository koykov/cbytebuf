package cbytebuf

/*
#include "stdlib.h"
#include "cbytebuf.h"
*/
import "C"
import (
	"errors"
	"reflect"
	"unsafe"
)

const (
	errOk       = 0
	errBadAlloc = 1
)

type CByteBuf struct {
	sh reflect.SliceHeader
	b  []byte
	l  int
	e  uint
}

var errs = []error{
	errOk:       nil,
	errBadAlloc: errors.New("bad alloc on buffer init or grow"),
}

func NewCByteBuf() (*CByteBuf, error) {
	b := CByteBuf{}
	return &b, nil
}

func (b *CByteBuf) Write(data []byte) (int, error) {
	b.l = len(data)
	if b.b == nil {
		b.sh.Cap = b.l * 2
		ptrErr := (*C.uint)(unsafe.Pointer(&b.e))
		ptrAddr := (*C.uintptr)(unsafe.Pointer(&b.sh.Data))
		ptrCap := (*C.int)(unsafe.Pointer(&b.sh.Cap))
		C.cbb_init(ptrErr, ptrAddr, ptrCap)
		b.b = *(*[]byte)(unsafe.Pointer(&b.sh))
	}

	if b.sh.Len+b.l > b.sh.Cap {
		_ = b.Grow((b.sh.Len + b.l) * 2)
	}

	b.b = append(b.b, data...)
	b.sh.Len += b.l

	return b.l, errs[b.e]
}

func (b *CByteBuf) Grow(cap int) error {
	b.sh.Cap = cap
	ptrErr := (*C.uint)(unsafe.Pointer(&b.e))
	ptrAddr := (*C.uintptr)(unsafe.Pointer(&b.sh.Data))
	ptrCap := (*C.int)(unsafe.Pointer(&b.sh.Cap))
	C.cbb_grow(ptrErr, ptrAddr, ptrCap)
	b.b = *(*[]byte)(unsafe.Pointer(&b.sh))
	return errs[b.e]
}

func (b *CByteBuf) Bytes() []byte {
	return b.b
}

func (b *CByteBuf) Reset() {
	b.b = b.b[:0]
}

func (b *CByteBuf) Release() error {
	ptrErr := (*C.uint)(unsafe.Pointer(&b.e))
	ptrAddr := (*C.uintptr)(unsafe.Pointer(&b.sh.Data))
	C.cbb_release(ptrErr, ptrAddr)
	b.sh.Len, b.sh.Cap = 0, 0
	b.b = nil
	return errs[b.e]
}
