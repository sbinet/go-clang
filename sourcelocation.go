package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

// SourceLocation identifies a specific source location within a translation
// unit.
//
// Use clang_getExpansionLocation() or clang_getSpellingLocation()
// to map a source location to a particular file, line, and column.
type SourceLocation struct {
	c C.CXSourceLocation
}

// NewNullLocation creates a NULL (invalid) source location.
func NewNullLocation() SourceLocation {
	return SourceLocation{C.clang_getNullLocation()}
}

// EqualLocations determines whether two source locations, which must refer into
// the same translation unit, refer to exactly the same point in the source
// code.
// Returns non-zero if the source locations refer to the same location, zero
// if they refer to different locations.
func EqualLocations(loc1, loc2 SourceLocation) bool {
	o := C.clang_equalLocations(loc1.c, loc2.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Returns non-zero if the given source location is in a system header.
 */
func (loc SourceLocation) IsInSystemHeader() bool {
	o := C.clang_Location_isInSystemHeader(loc.c)
	if o != 0 {
		return true
	}
	return false
}

/**
 * \brief Returns non-zero if the given source location is in the main file of
 * the corresponding translation unit.
 */
func (loc SourceLocation) IsFromMainFile() bool {
	o := C.clang_Location_isFromMainFile(loc.c)
	if o != 0 {
		return true
	}
	return false
}

// ExpansionLocation returns the file, line, column, and offset represented by
// the given source location.
//
// If the location refers into a macro expansion, retrieves the
// location of the macro expansion.
//
// file: if non-NULL, will be set to the file to which the given
// source location points.
//
// line: if non-NULL, will be set to the line to which the given
// source location points.
//
// column: if non-NULL, will be set to the column to which the given
// source location points.
//
// offset: if non-NULL, will be set to the offset into the
// buffer to which the given source location points.
func (l SourceLocation) ExpansionLocation() (f File, line, column, offset uint) {
	cline := C.uint(0)
	ccol := C.uint(0)
	coff := C.uint(0)
	// FIXME: undefined reference to `clang_getExpansionLocation'
	C.clang_getInstantiationLocation(l.c, &f.c, &cline, &ccol, &coff)
	line = uint(cline)
	column = uint(ccol)
	offset = uint(coff)

	return
}

/**
 * \brief Retrieve the file, line, column, and offset represented by
 * the given source location, as specified in a # line directive.
 *
 * Example: given the following source code in a file somefile.c
 *
 * #123 "dummy.c" 1
 *
 * static int func(void)
 * {
 *     return 0;
 * }
 *
 * the location information returned by this function would be
 *
 * File: dummy.c Line: 124 Column: 12
 *
 * whereas clang_getExpansionLocation would have returned
 *
 * File: somefile.c Line: 3 Column: 12
 *
 * \param location the location within a source file that will be decomposed
 * into its parts.
 *
 * \param filename [out] if non-NULL, will be set to the filename of the
 * source location. Note that filenames returned will be for "virtual" files,
 * which don't necessarily exist on the machine running clang - e.g. when
 * parsing preprocessed output obtained from a different environment. If
 * a non-NULL value is passed in, remember to dispose of the returned value
 * using \c clang_disposeString() once you've finished with it. For an invalid
 * source location, an empty string is returned.
 *
 * \param line [out] if non-NULL, will be set to the line number of the
 * source location. For an invalid source location, zero is returned.
 *
 * \param column [out] if non-NULL, will be set to the column number of the
 * source location. For an invalid source location, zero is returned.
 */
func (l SourceLocation) PresumedLocation() (fname string, line, column uint) {

	cname := cxstring{}
	defer cname.Dispose()
	cline := C.uint(0)
	ccol := C.uint(0)
	C.clang_getPresumedLocation(l.c, &cname.c, &cline, &ccol)
	fname = cname.String()
	line = uint(cline)
	column = uint(ccol)
	return
}

/**
 * \brief Legacy API to retrieve the file, line, column, and offset represented
 * by the given source location.
 *
 * This interface has been replaced by the newer interface
 * \see clang_getExpansionLocation(). See that interface's documentation for
 * details.
 */
func (l SourceLocation) InstantiationLocation() (file File, line, column, offset uint) {

	cline := C.uint(0)
	ccol := C.uint(0)
	coff := C.uint(0)
	C.clang_getInstantiationLocation(l.c,
		&file.c,
		&cline,
		&ccol,
		&coff)
	line = uint(cline)
	column = uint(ccol)
	offset = uint(coff)
	return
}

/**
 * \brief Retrieve the file, line, column, and offset represented by
 * the given source location.
 *
 * If the location refers into a macro instantiation, return where the
 * location was originally spelled in the source file.
 *
 * \param location the location within a source file that will be decomposed
 * into its parts.
 *
 * \param file [out] if non-NULL, will be set to the file to which the given
 * source location points.
 *
 * \param line [out] if non-NULL, will be set to the line to which the given
 * source location points.
 *
 * \param column [out] if non-NULL, will be set to the column to which the given
 * source location points.
 *
 * \param offset [out] if non-NULL, will be set to the offset into the
 * buffer to which the given source location points.
 */
func (l SourceLocation) SpellingLocation() (file File, line, column, offset uint) {

	cline := C.uint(0)
	ccol := C.uint(0)
	coff := C.uint(0)
	C.clang_getSpellingLocation(l.c,
		&file.c,
		&cline,
		&ccol,
		&coff)
	line = uint(cline)
	column = uint(ccol)
	offset = uint(coff)
	return
}

/**
 * \brief Retrieve the file, line, column, and offset represented by
 * the given source location.
 *
 * If the location refers into a macro expansion, return where the macro was
 * expanded or where the macro argument was written, if the location points at
 * a macro argument.
 *
 * \param location the location within a source file that will be decomposed
 * into its parts.
 *
 * \param file [out] if non-NULL, will be set to the file to which the given
 * source location points.
 *
 * \param line [out] if non-NULL, will be set to the line to which the given
 * source location points.
 *
 * \param column [out] if non-NULL, will be set to the column to which the given
 * source location points.
 *
 * \param offset [out] if non-NULL, will be set to the offset into the
 * buffer to which the given source location points.
 */
func (loc SourceLocation) GetFileLocation() (f File, line, column, offset uint) {
	cline := C.uint(0)
	ccol := C.uint(0)
	coff := C.uint(0)
	C.clang_getFileLocation(loc.c,
		&f.c,
		&cline,
		&ccol,
		&coff)
	line = uint(cline)
	column = uint(ccol)
	offset = uint(coff)
	return

}
