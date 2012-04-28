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

/**
 * \brief Return the canonical type for a CXType.
 *
 * Clang's type system explicitly models typedefs and all the ways
 * a specific type can be represented.  The canonical type is the underlying
 * type with all the "sugar" removed.  For example, if 'T' is a typedef
 * for 'int', the canonical type for 'T' would be 'int'.
 */
func (t Type) CanonicalType() Type {
	o := C.clang_getCanonicalType(t.c)
	return Type{o}
}

/**
 *  \determine Determine whether a CXType has the "const" qualifier set, 
 *  without looking through typedefs that may have added "const" at a different level.
 */
func (t Type) IsConstQualified() bool {
	o := C.clang_isConstQualifiedType(t.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 *  \determine Determine whether a CXType has the "volatile" qualifier set,
 *  without looking through typedefs that may have added "volatile" at a different level.
 */
func (t Type) IsVolatileQualified() bool {
	o := C.clang_isVolatileQualifiedType(t.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 *  \determine Determine whether a CXType has the "restrict" qualifier set,
 *  without looking through typedefs that may have added "restrict" at a different level.
 */
func (t Type) IsRestrictQualified() bool {
	o := C.clang_isRestrictQualifiedType(t.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief For pointer types, returns the type of the pointee.
 *
 */
func (t Type) PointeeType() Type {
	o := C.clang_getPointeeType(t.c)
	return Type{o}
}

/**
 * \brief Return the cursor for the declaration of the given type.
 */
func (t Type) Declaration() Cursor {
	o := C.clang_getTypeDeclaration(t.c)
	return Cursor{o}
}

/**
 * Returns the Objective-C type encoding for the specified declaration.
 */
//FIXME
//CINDEX_LINKAGE CXString clang_getDeclObjCTypeEncoding(CXCursor C);

/**
 * \brief Retrieve the spelling of a given CXTypeKind.
 */
//FIXME
//CINDEX_LINKAGE CXString clang_getTypeKindSpelling(enum CXTypeKind K);

/**
 * \brief Retrieve the result type associated with a function type.
 */
func (t Type) ResultType() Type {
	o := C.clang_getResultType(t.c)
	return Type{o}
}

/**
 * \brief Return 1 if the CXType is a POD (plain old data) type, and 0
 *  otherwise.
 */
func (t Type) IsPOD() bool {
	o := C.clang_isPODType(t.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Return the element type of an array type.
 *
 * If a non-array type is passed in, an invalid type is returned.
 */
func (t Type) ArrayElementType() Type {
	o := C.clang_getArrayElementType(t.c)
	return Type{o}
}

/**
 * \brief Return the the array size of a constant array.
 *
 * If a non-array type is passed in, -1 is returned.
 */
func (t Type) ArraySize() int64 {
	o := C.clang_getArraySize(t.c)
	return int64(o)
}




// EOF
