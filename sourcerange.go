package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

// SourceRange identifies a half-open character range in the source code.
//
// Use clang_getRangeStart() and clang_getRangeEnd() to retrieve the
// starting and end locations from a source range, respectively.
type SourceRange struct {
	c C.CXSourceRange
}

// NewNullRange creates a NULL (invalid) source range.
func NewNullRange() SourceRange {
	return SourceRange{C.clang_getNullRange()}
}

// NewRange creates a source range given the beginning and ending source
// locations.
func NewRange(beg, end SourceLocation) SourceRange {
	o := C.clang_getRange(beg.c, end.c)
	return SourceRange{o}
}

// EqualRanges determines whether two ranges are equivalent.
func EqualRanges(r1, r2 SourceRange) bool {
	o := C.clang_equalRanges(r1.c, r2.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

// IsNull checks if the underlying source range is null.
func (r SourceRange) IsNull() bool {
	o := C.clang_Range_isNull(r.c)
	if o != C.int(0) {
		return true
	}
	return false
}

/**
 * \brief Retrieve a source location representing the first character within a
 * source range.
 */
func (s SourceRange) Start() SourceLocation {
	o := C.clang_getRangeStart(s.c)
	return SourceLocation{o}
}

/**
 * \brief Retrieve a source location representing the last character within a
 * source range.
 */
func (s SourceRange) End() SourceLocation {
	o := C.clang_getRangeEnd(s.c)
	return SourceLocation{o}
}
