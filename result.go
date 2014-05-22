package clang

// #include <stdlib.h>
// #include "clang-c/Index.h"
//
import "C"

// Result is the result of calling a c-clang function
type Result int

const (
	/**
	 * \brief Function returned successfully.
	 */
	Result_Success Result = C.CXResult_Success
	/**
	 * \brief One of the parameters was invalid for the function.
	 */
	Result_Invalid = C.CXResult_Invalid
	/**
	 * \brief The function was terminated by a callback (e.g. it returned
	 * CXVisit_Break)
	 */
	Result_VisitBreak = C.CXResult_VisitBreak
)
