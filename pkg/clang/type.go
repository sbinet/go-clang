package clang

// #include <stdlib.h>
// #include "clang-c/Index.h"
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

	TK_Void       = C.CXType_Void
	TK_Bool       = C.CXType_Bool
	TK_Char_U     = C.CXType_Char_U
	TK_UChar      = C.CXType_UChar
	TK_Char16     = C.CXType_Char16
	TK_Char32     = C.CXType_Char32
	TK_UShort     = C.CXType_UShort
	TK_UInt       = C.CXType_UInt
	TK_ULong      = C.CXType_ULong
	TK_ULongLong  = C.CXType_ULongLong
	TK_UInt128    = C.CXType_UInt128
	TK_Char_S     = C.CXType_Char_S
	TK_SChar      = C.CXType_SChar
	TK_WChar      = C.CXType_WChar
	TK_Short      = C.CXType_Short
	TK_Int        = C.CXType_Int
	TK_Long       = C.CXType_Long
	TK_LongLong   = C.CXType_LongLong
	TK_Int128     = C.CXType_Int128
	TK_Float      = C.CXType_Float
	TK_Double     = C.CXType_Double
	TK_LongDouble = C.CXType_LongDouble
	TK_NullPtr    = C.CXType_NullPtr
	TK_Overload   = C.CXType_Overload
	TK_Dependent  = C.CXType_Dependent
	TK_ObjCId     = C.CXType_ObjCId
	TK_ObjCClass  = C.CXType_ObjCClass
	TK_ObjCSel    = C.CXType_ObjCSel

	TK_FirstBuiltin = C.CXType_FirstBuiltin
	TK_LastBuiltin  = C.CXType_LastBuiltin

	TK_Complex           = C.CXType_Complex
	TK_Pointer           = C.CXType_Pointer
	TK_BlockPointer      = C.CXType_BlockPointer
	TK_LValueReference   = C.CXType_LValueReference
	TK_RValueReference   = C.CXType_RValueReference
	TK_Record            = C.CXType_Record
	TK_Enum              = C.CXType_Enum
	TK_Typedef           = C.CXType_Typedef
	TK_ObjCInterface     = C.CXType_ObjCInterface
	TK_ObjCObjectPointer = C.CXType_ObjCObjectPointer
	TK_FunctionNoProto   = C.CXType_FunctionNoProto
	TK_FunctionProto     = C.CXType_FunctionProto
	TK_ConstantArray     = C.CXType_ConstantArray
)

func (t TypeKind) to_c() uint32 {
	return uint32(t)
}

// Type represents the type of an element in the abstract syntax tree.
type Type struct {
	c C.CXType
}

// Spelling returns the spelling of a given TypeKind.
func (t TypeKind) Spelling() string {
	cstr := cxstring{C.clang_getTypeKindSpelling(t.to_c())}
	defer cstr.Dispose()
	return cstr.String()
}

// EqualTypes determines whether two Types represent the same type.
func EqualTypes(t1, t2 Type) bool {
	o := C.clang_equalTypes(t1.c, t2.c)
	if o != C.uint(0) {
		return true
	}
	return false
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

// EOF
