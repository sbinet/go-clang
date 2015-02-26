package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

// RefQualifierKind describes the kind of reference a Type is decorated with
type RefQualifierKind int

const (

	/** \brief No ref-qualifier was provided. */
	RQK_None RefQualifierKind = C.CXRefQualifier_None

	/** \brief An lvalue ref-qualifier was provided (\c &). */
	RQK_LValue = C.CXRefQualifier_LValue
	/** \brief An rvalue ref-qualifier was provided (\c &&). */
	RQK_RValue = C.CXRefQualifier_RValue
)
