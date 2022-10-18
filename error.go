package cbytebuf

import "errors"

var (
	ErrOk            error = nil
	ErrBadAlloc            = errors.New("bad alloc on buffer init or grow")
	ErrNegativeCap         = errors.New("negative cap on the grow")
	ErrNegativeRead        = errors.New("reader returned negative count from Read")
	ErrNilMarshaller       = errors.New("marshaller object is nil")
)
