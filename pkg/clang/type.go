package clang

// #include <stdlib.h>
// #cgo LDFLAGS: -L/opt/local/libexec/llvm-3.0/lib -lclang
// #include "/opt/local/libexec/llvm-3.0/include/clang-c/Index.h"
// inline static
// CXCursor _go_clang_ocursor_at(CXCursor *c, int idx) {
//   return c[idx];
// }
//
import "C"
import (
	//"unsafe"
)

// TypeKind describes the kind of a type
type TypeKind uint32

const (
	// Represents an invalid type (e.g., where no type is available).
	TK_Invalid TypeKind = C.CXType_Invalid

	// A type whose specific kind is not exposed via this interface.
	TK_Unexposed = C.CXType_Unexposed
	
	//FIXME
)

// Type represents the type of an element in the abstract syntax tree.
type Type struct {
	c C.CXType
}

// EqualTypes determines whether two Types represent the same type.
func EqualTypes(t1, t2 Type) bool {
	o := C.clang_equalTypes(t1.c, t2.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

// EOF
