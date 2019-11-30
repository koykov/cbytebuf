package cbytebuf

/*
#include "stdlib.h"
#include "cbytebuf.h"
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CByteBuf struct {
	sh reflect.SliceHeader
	b  []byte
	l  int
}

func NewCByteBuf() (*CByteBuf, error) {
	b := CByteBuf{}
	return &b, nil
}

func (b *CByteBuf) Write(data []byte) (int, error) {
	b.l = len(data)
	if b.b == nil {
		b.sh.Cap = b.l * 2
		ptrAddr := (*C.uintptr)(unsafe.Pointer(&b.sh.Data))
		ptrCap := (*C.int)(unsafe.Pointer(&b.sh.Cap))
		C.cbb_init(ptrAddr, ptrCap)
		b.b = *(*[]byte)(unsafe.Pointer(&b.sh))
	}

	if b.sh.Len+b.l > b.sh.Cap {
		_ = b.Grow((b.sh.Len + b.l) * 2)
	}

	b.b = append(b.b, data...)
	b.sh.Len += b.l

	return b.l, nil
}

func (b *CByteBuf) Grow(cap int) error {
	b.sh.Cap = cap
	ptrAddr := (*C.uintptr)(unsafe.Pointer(&b.sh.Data))
	ptrCap := (*C.int)(unsafe.Pointer(&b.sh.Cap))
	C.cbb_grow(ptrAddr, ptrCap)
	b.b = *(*[]byte)(unsafe.Pointer(&b.sh))
	return nil
}

func (b *CByteBuf) Bytes() []byte {
	return b.b
}

func (b *CByteBuf) Reset() {
	b.b = b.b[:0]
}

func (b *CByteBuf) Release() {
	ptrAddr := (*C.uintptr)(unsafe.Pointer(&b.sh.Data))
	C.cbb_release(ptrAddr)
	b.sh.Len, b.sh.Cap = 0, 0
	b.b = nil
}
