package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// Type represents the type of an element in the abstract syntax tree.
type Type struct {
	c C.CXType
}

func (c Type) Kind() TypeKind {
	return TypeKind(c.c.kind)
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
 * \brief Pretty-print the underlying type using the rules of the
 * language of the translation unit from which it came.
 *
 * If the type is invalid, an empty string is returned.
 */
func (t Type) TypeSpelling() string {
	o := cxstring{C.clang_getTypeSpelling(t.c)}
	defer o.Dispose()
	return o.String()
}

// CanonicalType returns the canonical type for a Type.
//
// Clang's type system explicitly models typedefs and all the ways
// a specific type can be represented.  The canonical type is the underlying
// type with all the "sugar" removed.  For example, if 'T' is a typedef
// for 'int', the canonical type for 'T' would be 'int'.
func (t Type) CanonicalType() Type {
	o := C.clang_getCanonicalType(t.c)
	return Type{o}
}

// IsConstQualified determines whether a Type has the "const" qualifier set,
// without looking through typedefs that may have added "const" at a
// different level.
func (t Type) IsConstQualified() bool {
	o := C.clang_isConstQualifiedType(t.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

// IsVolatileQualified determines whether a Type has the "volatile" qualifier
// set, without looking through typedefs that may have added "volatile" at a
// different level.
func (t Type) IsVolatileQualified() bool {
	o := C.clang_isVolatileQualifiedType(t.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

// IsRestrictQualified determines whether a Type has the "restrict" qualifier
// set, without looking through typedefs that may have added "restrict" at a
// different level.
func (t Type) IsRestrictQualified() bool {
	o := C.clang_isRestrictQualifiedType(t.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

// PointeeType (for pointer types), returns the type of the pointee.
func (t Type) PointeeType() Type {
	o := C.clang_getPointeeType(t.c)
	return Type{o}
}

// Declaration returns the cursor for the declaration of the given type.
func (t Type) Declaration() Cursor {
	o := C.clang_getTypeDeclaration(t.c)
	return Cursor{o}
}

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

/**
 * \brief Return the alignment of a type in bytes as per C++[expr.alignof]
 *   standard.
 *
 * If the type declaration is invalid, CXTypeLayoutError_Invalid is returned.
 * If the type declaration is an incomplete type, CXTypeLayoutError_Incomplete
 *   is returned.
 * If the type declaration is a dependent type, CXTypeLayoutError_Dependent is
 *   returned.
 * If the type declaration is not a constant size type,
 *   CXTypeLayoutError_NotConstantSize is returned.
 */
func (t Type) AlignOf() (int, error) {
	o := C.clang_Type_getAlignOf(t.c)
	if o < 0 {
		return int(o), TypeLayoutError(o)
	}
	return int(o), nil
}

/**
 * \brief Return the class type of an member pointer type.
 *
 * If a non-member-pointer type is passed in, an invalid type is returned.
 */
func (t Type) ClassType() Type {
	return Type{C.clang_Type_getClassType(t.c)}
}

/**
 * \brief Return the size of a type in bytes as per C++[expr.sizeof] standard.
 *
 * If the type declaration is invalid, CXTypeLayoutError_Invalid is returned.
 * If the type declaration is an incomplete type, CXTypeLayoutError_Incomplete
 *   is returned.
 * If the type declaration is a dependent type, CXTypeLayoutError_Dependent is
 *   returned.
 */
func (t Type) SizeOf() (int, error) {
	o := C.clang_Type_getSizeOf(t.c)
	if o < 0 {
		return int(o), TypeLayoutError(o)
	}
	return int(o), nil
}

/**
 * \brief Return the offset of a field named S in a record of type T in bits
 *   as it would be returned by __offsetof__ as per C++11[18.2p4]
 *
 * If the cursor is not a record field declaration, CXTypeLayoutError_Invalid
 *   is returned.
 * If the field's type declaration is an incomplete type,
 *   CXTypeLayoutError_Incomplete is returned.
 * If the field's type declaration is a dependent type,
 *   CXTypeLayoutError_Dependent is returned.
 * If the field's name S is not found,
 *   CXTypeLayoutError_InvalidFieldName is returned.
 */
func (t Type) OffsetOf(s string) (int, error) {
	c_str := C.CString(s)
	defer C.free(unsafe.Pointer(c_str))
	o := C.clang_Type_getOffsetOf(t.c, c_str)
	if o < 0 {
		return int(o), TypeLayoutError(o)
	}
	return int(o), nil
}

/**
 * \brief Retrieve the ref-qualifier kind of a function or method.
 *
 * The ref-qualifier is returned for C++ functions or methods. For other types
 * or non-C++ declarations, CXRefQualifier_None is returned.
 */
func (t Type) CXXRefQualifier() RefQualifierKind {
	return RefQualifierKind(C.clang_Type_getCXXRefQualifier(t.c))
}

// EOF
