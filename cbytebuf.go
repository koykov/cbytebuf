package cbytebuf

/*
#include "stdlib.h"
#include "export.h"
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type CByteBuf struct {
	cp *C.CByteBuf
	sh reflect.SliceHeader
	a  int
	l  int
	c  int
}

// Init new instance of ByteBuf.
func NewCByteBuf() (*CByteBuf, error) {
	cbb := C.cbb_new()
	bb := CByteBuf{
		cp: (*C.CByteBuf)(unsafe.Pointer(cbb)),
	}
	return &bb, nil
}

func (b *CByteBuf) Write(data []byte) (int, error) {
	ptrData := (*C.uchar)(unsafe.Pointer(&data[0]))
	b.l = len(data)
	ptrLen := (*C.int)(unsafe.Pointer(&b.l))
	C.cbb_write(b.cp, ptrData, ptrLen)
	return len(data), nil
}

func (b *CByteBuf) Bytes() []byte {
	ptrAddr := (*C.uintptr)(unsafe.Pointer(&b.a))
	ptrLen := (*C.int)(unsafe.Pointer(&b.l))
	ptrCap := (*C.int)(unsafe.Pointer(&b.c))
	C.cbb_bytes(b.cp, ptrAddr, ptrLen, ptrCap)

	b.sh.Data = uintptr(b.a)
	b.sh.Len = b.l
	b.sh.Cap = b.c
	return *(*[]byte)(unsafe.Pointer(&b.sh))
}

func (b *CByteBuf) Reset() {
	C.cbb_release(b.cp)
}
