package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"
import (
	"fmt"
)

/**
 * \brief List the possible error codes for \c clang_Type_getSizeOf,
 *   \c clang_Type_getAlignOf, \c clang_Type_getOffsetOf and
 *   \c clang_Cursor_getOffsetOf.
 *
 * A value of this enumeration type can be returned if the target type is not
 * a valid argument to sizeof, alignof or offsetof.
 */
type TypeLayoutError int

const (
	/**
	 * \brief Type is of kind CXType_Invalid.
	 */
	TLE_Invalid TypeLayoutError = C.CXTypeLayoutError_Invalid

	/**
	 * \brief The type is an incomplete Type.
	 */
	TLE_Incomplete = C.CXTypeLayoutError_Incomplete

	/**
	 * \brief The type is a dependent Type.
	 */
	TLE_Dependent = C.CXTypeLayoutError_Dependent

	/**
	 * \brief The type is not a constant size type.
	 */
	TLE_NotConstantSize = C.CXTypeLayoutError_NotConstantSize

	/**
	 * \brief The Field name is not valid for this record.
	 */
	TLE_InvalidFieldName = C.CXTypeLayoutError_InvalidFieldName
)

func (tle TypeLayoutError) Error() string {
	switch tle {
	case TLE_Invalid:
		return "TypeLayout=Invalid"

	case TLE_Incomplete:
		return "TypeLayout=Incomplete"

	case TLE_Dependent:
		return "TypeLayout=Dependent"

	case TLE_NotConstantSize:
		return "TypeLayout=NotConstantSize"

	case TLE_InvalidFieldName:
		return "TypeLayout=InvalidFieldName"
	}

	return fmt.Sprintf("TypeLayout=unknown (%d)", int(tle))
}
