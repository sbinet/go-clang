package clang

// #include <stdlib.h>
// #cgo LDFLAGS: -L/opt/local/libexec/llvm-3.0/lib -L/usr/lib/llvm -lclang
// #cgo CFLAGS: -I/opt/local/libexec/llvm-3.0/include
// #include "clang-c/Index.h"
import "C"
import (
	//"unsafe"
)

// cxstring is a character string.
//
// The cxstring type is used to return strings from the interface when the
// ownership of that string might different from one call to the next.
type cxstring struct {
	c C.CXString
}

func (c cxstring) String() string {
	cstr := C.clang_getCString(c.c)
	return C.GoString(cstr)
}

func (c cxstring) Dispose() {
	C.clang_disposeString(c.c)
}

// EOF
