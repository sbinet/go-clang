package clang

// #include <stdlib.h>
// #cgo LDFLAGS: -L/opt/local/libexec/llvm-3.0/lib -lclang
// #include "/opt/local/libexec/llvm-3.0/include/clang-c/Index.h"
import "C"
import (
	"unsafe"
)

// An "index" that consists of a set of translation units that would
// typically be linked together into an executable or library
type Index struct {
	c C.CXIndex
}

/**
 * CreateIndex provides a shared context for creating
 * translation units. It provides two options:
 *
 * - excludeDeclarationsFromPCH: When non-zero, allows enumeration of "local"
 * declarations (when loading any new translation units). A "local" declaration
 * is one that belongs in the translation unit itself and not in a precompiled
 * header that was used by the translation unit. If zero, all declarations
 * will be enumerated.
 *
 * Here is an example:
 *
 *   // excludeDeclsFromPCH = 1, displayDiagnostics=1
 *   Idx = clang_createIndex(1, 1);
 *
 *   // IndexTest.pch was produced with the following command:
 *   // "clang -x c IndexTest.h -emit-ast -o IndexTest.pch"
 *   TU = clang_createTranslationUnit(Idx, "IndexTest.pch");
 *
 *   // This will load all the symbols from 'IndexTest.pch'
 *   clang_visitChildren(clang_getTranslationUnitCursor(TU),
 *                       TranslationUnitVisitor, 0);
 *   clang_disposeTranslationUnit(TU);
 *
 *   // This will load all the symbols from 'IndexTest.c', excluding symbols
 *   // from 'IndexTest.pch'.
 *   char *args[] = { "-Xclang", "-include-pch=IndexTest.pch" };
 *   TU = clang_createTranslationUnitFromSourceFile(Idx, "IndexTest.c", 2, args,
 *                                                  0, 0);
 *   clang_visitChildren(clang_getTranslationUnitCursor(TU),
 *                       TranslationUnitVisitor, 0);
 *   clang_disposeTranslationUnit(TU);
 *
 * This process of creating the 'pch', loading it separately, and using it (via
 * -include-pch) allows 'excludeDeclsFromPCH' to remove redundant callbacks
 * (which gives the indexer the same performance benefit as the compiler).
 */
func CreateIndex(excludeDeclarationsFromPCH, displayDiagnostics int) Index {
	idx := C.clang_createIndex(C.int(excludeDeclarationsFromPCH),
		C.int(displayDiagnostics))
	return Index{idx}
}

/**
 * Destroy the given index.
 *
 * The index must not be destroyed until all of the translation units created
 * within that index have been destroyed.
 */
func (idx Index) Dispose() {
	C.clang_disposeIndex(idx.c)
}

// A single translation unit, which resides in an index
type TranslationUnit struct {
	c C.CXTranslationUnit
}

func (tu TranslationUnit) GetFile(file_name string) File {
	cfname := C.CString(file_name)
	defer C.free(unsafe.Pointer(cfname))
	f := C.clang_getFile(tu.c, cfname)
	return File{f}
}

/**
 * \brief Determine whether the given header is guarded against
 * multiple inclusions, either with the conventional
 * #ifndef/#define/#endif macro guards or with #pragma once.
 */
func (tu TranslationUnit) IsFileMultipleIncludeGuarded(file File) bool {
	o := C.clang_isFileMultipleIncludeGuarded(tu.c, file.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Retrieve the cursor that represents the given translation unit.
 *
 * The translation unit cursor can be used to start traversing the
 * various declarations within the given translation unit.
 */
func (tu TranslationUnit) GetCursor() Cursor {
	o := C.clang_getTranslationUnitCursor(tu.c)
	return Cursor{o}
}

/**
 * \brief Map a source location to the cursor that describes the entity at that
 * location in the source code.
 *
 * clang_getCursor() maps an arbitrary source location within a translation
 * unit down to the most specific cursor that describes the entity at that
 * location. For example, given an expression \c x + y, invoking
 * clang_getCursor() with a source location pointing to "x" will return the
 * cursor for "x"; similarly for "y". If the cursor points anywhere between
 * "x" or "y" (e.g., on the + or the whitespace around it), clang_getCursor()
 * will return a cursor referring to the "+" expression.
 *
 * \returns a cursor representing the entity at the given source location, or
 * a NULL cursor if no such entity can be found.
 */
func (tu TranslationUnit) Cursor(loc SourceLocation) Cursor {
	o := C.clang_getCursor(tu.c, loc.c)
	return Cursor{o}
}

// A particular source file that is part of a translation unit.
type File struct {
	c C.CXFile
}

// FIXME
// /**
//  * \brief Retrieve the complete file and path name of the given file.
//  */
// CINDEX_LINKAGE CXString clang_getFileName(CXFile SFile);

// /**
//  * \brief Retrieve the last modification time of the given file.
//  */
// CINDEX_LINKAGE time_t clang_getFileTime(CXFile SFile);



/**
 * \brief Identifies a specific source location within a translation
 * unit.
 *
 * Use clang_getExpansionLocation() or clang_getSpellingLocation()
 * to map a source location to a particular file, line, and column.
 */
type SourceLocation struct {
	c C.CXSourceLocation
}

/**
 * \brief Identifies a half-open character range in the source code.
 *
 * Use clang_getRangeStart() and clang_getRangeEnd() to retrieve the
 * starting and end locations from a source range, respectively.
 */
type SourceRange struct {
	c C.CXSourceRange
}

/**
 * \brief Retrieve a NULL (invalid) source location.
 */
func GetNullLocation() SourceLocation {
	return SourceLocation{C.clang_getNullLocation()}
}

/**
 * \determine Determine whether two source locations, which must refer into
 * the same translation unit, refer to exactly the same point in the source
 * code.
 *
 * \returns non-zero if the source locations refer to the same location, zero
 * if they refer to different locations.
 */
func EqualLocations(loc1, loc2 SourceLocation) bool {
	o := C.clang_equalLocations(loc1.c, loc2.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Retrieves the source location associated with a given file/line/column
 * in a particular translation unit.
 */
func (tu TranslationUnit) GetLocation(f File, line, column uint) SourceLocation {
	loc := C.clang_getLocation(tu.c, f.c, C.uint(line), C.uint(column))
	return SourceLocation{loc}
}

/**
 * \brief Retrieves the source location associated with a given character offset
 * in a particular translation unit.
 */
func (tu TranslationUnit) GetLocationForOffset(f File, offset uint) SourceLocation {
	loc := C.clang_getLocationForOffset(tu.c, f.c, C.uint(offset))
	return SourceLocation{loc}
}

/**
 * \brief Retrieve a NULL (invalid) source range.
 */
func GetNullRange() SourceRange {
	return SourceRange{C.clang_getNullRange()}
}

/**
 * \brief Retrieve a source range given the beginning and ending source
 * locations.
 */
func GetRange(beg, end SourceLocation) SourceRange {
	o := C.clang_getRange(beg.c, end.c)
	return SourceRange{o}
}

/**
 * \brief Determine whether two ranges are equivalent.
 *
 * \returns non-zero if the ranges are the same, zero if they differ.
 */
func EqualRanges(r1, r2 SourceRange) bool {
	o := C.clang_equalRanges(r1.c, r2.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Returns non-zero if \arg range is null.
 */
func (r SourceRange) IsNull() bool {
	o := C.clang_Range_isNull(r.c)
	if o != C.int(0) {
		return true
	}
	return false
}

/**
 * \brief Retrieve the file, line, column, and offset represented by
 * the given source location.
 *
 * If the location refers into a macro expansion, retrieves the
 * location of the macro expansion.
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
func (l SourceLocation) GetExpansionLocation() (f File, line, column, offset uint) {
	cline := C.uint(0)
	ccol  := C.uint(0)
	coff  := C.uint(0)
	// FIXME
	// C.clang_getExpansionLocation(l.c,
        //         unsafe.Pointer(&f.c),
        //         unsafe.Pointer(&cline),
        //         unsafe.Pointer(&ccol),
        //         unsafe.Pointer(&coff))
	line  = uint(cline)
	column = uint(ccol)
	offset = uint(coff)

	return
}

// EOF
