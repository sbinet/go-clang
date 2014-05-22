package clang

// #include <stdlib.h>
// #include "clang-c/Index.h"
//
import "C"

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

	TK_Complex             = C.CXType_Complex
	TK_Pointer             = C.CXType_Pointer
	TK_BlockPointer        = C.CXType_BlockPointer
	TK_LValueReference     = C.CXType_LValueReference
	TK_RValueReference     = C.CXType_RValueReference
	TK_Record              = C.CXType_Record
	TK_Enum                = C.CXType_Enum
	TK_Typedef             = C.CXType_Typedef
	TK_ObjCInterface       = C.CXType_ObjCInterface
	TK_ObjCObjectPointer   = C.CXType_ObjCObjectPointer
	TK_FunctionNoProto     = C.CXType_FunctionNoProto
	TK_FunctionProto       = C.CXType_FunctionProto
	TK_ConstantArray       = C.CXType_ConstantArray
	TK_Vector              = C.CXType_Vector
	TK_IncompleteArray     = C.CXType_IncompleteArray
	TK_VariableArray       = C.CXType_VariableArray
	TK_DependentSizedArray = C.CXType_DependentSizedArray
	TK_MemberPointer       = C.CXType_MemberPointer
)

func (t TypeKind) to_c() uint32 {
	return uint32(t)
}

// Spelling returns the spelling of a given TypeKind.
func (t TypeKind) Spelling() string {
	cstr := cxstring{C.clang_getTypeKindSpelling(t.to_c())}
	defer cstr.Dispose()
	return cstr.String()
}
