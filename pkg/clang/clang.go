package clang

// #include <stdlib.h>
// #include "clang-c/Index.h"
import "C"
import (
	"time"
	"unsafe"
)

/**
 * \brief Describes the availability of a particular entity, which indicates
 * whether the use of this entity will result in a warning or error due to
 * it being deprecated or unavailable.
 */
type AvailabilityKind uint32
const (

	/**
	 * \brief The entity is available.
	 */
	Available AvailabilityKind = C.CXAvailability_Available

  /**
   * \brief The entity is available, but has been deprecated (and its use is
   * not recommended).
   */
	Deprecated AvailabilityKind = C.CXAvailability_Deprecated
	/**
	 * \brief The entity is not available; any use of it will be an error.
	 */
	NotAvailable AvailabilityKind = C.CXAvailability_NotAvailable
	/**
	 * \brief The entity is available, but not accessible; any use of it will be
	 * an error.
	 */
	NotAccessible AvailabilityKind = C.CXAvailability_NotAccessible
)

// An "index" that consists of a set of translation units that would
// typically be linked together into an executable or library
type Index struct {
	c C.CXIndex
}


// NewIndex provides a shared context for creating
// translation units. It provides two options:
//
// - excludeDeclarationsFromPCH: When non-zero, allows enumeration of "local"
// declarations (when loading any new translation units). A "local" declaration
// is one that belongs in the translation unit itself and not in a precompiled
// header that was used by the translation unit. If zero, all declarations
// will be enumerated.
//
// Here is an example:
//
//   // excludeDeclsFromPCH = 1, displayDiagnostics=1
//   Idx = clang_createIndex(1, 1);
//
//   // IndexTest.pch was produced with the following command:
//   // "clang -x c IndexTest.h -emit-ast -o IndexTest.pch"
//   TU = clang_createTranslationUnit(Idx, "IndexTest.pch");
//
//   // This will load all the symbols from 'IndexTest.pch'
//   clang_visitChildren(clang_getTranslationUnitCursor(TU),
//                       TranslationUnitVisitor, 0);
//   clang_disposeTranslationUnit(TU);
//
//   // This will load all the symbols from 'IndexTest.c', excluding symbols
//   // from 'IndexTest.pch'.
//   char *args[] = { "-Xclang", "-include-pch=IndexTest.pch" };
//   TU = clang_createTranslationUnitFromSourceFile(Idx, "IndexTest.c", 2, args,
//                                                  0, 0);
//   clang_visitChildren(clang_getTranslationUnitCursor(TU),
//                       TranslationUnitVisitor, 0);
//   clang_disposeTranslationUnit(TU);
//
// This process of creating the 'pch', loading it separately, and using it (via
// -include-pch) allows 'excludeDeclsFromPCH' to remove redundant callbacks
// (which gives the indexer the same performance benefit as the compiler).
func NewIndex(excludeDeclarationsFromPCH, displayDiagnostics int) Index {
	idx := C.clang_createIndex(C.int(excludeDeclarationsFromPCH),
		C.int(displayDiagnostics))
	return Index{idx}
}

// Dispose destroys the given index.
//
// The index must not be destroyed until all of the translation units created
// within that index have been destroyed.
func (idx Index) Dispose() {
	C.clang_disposeIndex(idx.c)
}

/**
 * \brief Create a translation unit from an AST file (-emit-ast).
 */
func (idx Index) CreateTranslationUnit(fname string) TranslationUnit {
	cstr := C.CString(fname)
	defer C.free(unsafe.Pointer(cstr))
	o := C.clang_createTranslationUnit(idx.c, cstr)
	return TranslationUnit{o}
}

/**
 * \brief Return the CXTranslationUnit for a given source file and the provided
 * command line arguments one would pass to the compiler.
 *
 * Note: The 'source_filename' argument is optional.  If the caller provides a
 * NULL pointer, the name of the source file is expected to reside in the
 * specified command line arguments.
 *
 * Note: When encountered in 'clang_command_line_args', the following options
 * are ignored:
 *
 *   '-c'
 *   '-emit-ast'
 *   '-fsyntax-only'
 *   '-o <output file>'  (both '-o' and '<output file>' are ignored)
 *
 * \param CIdx The index object with which the translation unit will be
 * associated.
 *
 * \param source_filename - The name of the source file to load, or NULL if the
 * source file is included in \p clang_command_line_args.
 *
 * \param num_clang_command_line_args The number of command-line arguments in
 * \p clang_command_line_args.
 *
 * \param clang_command_line_args The command-line arguments that would be
 * passed to the \c clang executable if it were being invoked out-of-process.
 * These command-line options will be parsed and will affect how the translation
 * unit is parsed. Note that the following options are ignored: '-c',
 * '-emit-ast', '-fsyntex-only' (which is the default), and '-o <output file>'.
 *
 * \param num_unsaved_files the number of unsaved file entries in \p
 * unsaved_files.
 *
 * \param unsaved_files the files that have not yet been saved to disk
 * but may be required for code completion, including the contents of
 * those files.  The contents and name of these files (as specified by
 * CXUnsavedFile) are copied when necessary, so the client only needs to
 * guarantee their validity until the call to this function returns.
 */
//FIXME: handle unsaved-files
func (idx Index) CreateTranslationUnitFromSourceFile(fname string, args []string) TranslationUnit {
	var c_fname *C.char = nil
	if fname != "" {
		c_fname = C.CString(fname)
	}
	defer C.free(unsafe.Pointer(c_fname))

	c_nargs := C.int(len(args))
	c_cmds := make([]*C.char, len(args))
	for i,_ := range args {
		cstr := C.CString(args[i])
		defer C.free(unsafe.Pointer(cstr))
		c_cmds[i] = cstr
	}

	o := C.clang_createTranslationUnitFromSourceFile(
		idx.c,
		c_fname,
		c_nargs, &c_cmds[0],
		C.uint(0), nil)
	return TranslationUnit{o}

}

/**
 * \brief Flags that control the creation of translation units.
 *
 * The enumerators in this enumeration type are meant to be bitwise
 * ORed together to specify which options should be used when
 * constructing the translation unit.
 */
type TranslationUnitFlags uint32
const (
  /**
   * \brief Used to indicate that no special translation-unit options are
   * needed.
   */
	TU_None = C.CXTranslationUnit_None

  /**
   * \brief Used to indicate that the parser should construct a "detailed"
   * preprocessing record, including all macro definitions and instantiations.
   *
   * Constructing a detailed preprocessing record requires more memory
   * and time to parse, since the information contained in the record
   * is usually not retained. However, it can be useful for
   * applications that require more detailed information about the
   * behavior of the preprocessor.
   */
	TU_DetailedPreprocessingRecord = C.CXTranslationUnit_DetailedPreprocessingRecord

  /**
   * \brief Used to indicate that the translation unit is incomplete.
   *
   * When a translation unit is considered "incomplete", semantic
   * analysis that is typically performed at the end of the
   * translation unit will be suppressed. For example, this suppresses
   * the completion of tentative declarations in C and of
   * instantiation of implicitly-instantiation function templates in
   * C++. This option is typically used when parsing a header with the
   * intent of producing a precompiled header.
   */
	TU_Incomplete = C.CXTranslationUnit_Incomplete
  
  /**
   * \brief Used to indicate that the translation unit should be built with an 
   * implicit precompiled header for the preamble.
   *
   * An implicit precompiled header is used as an optimization when a
   * particular translation unit is likely to be reparsed many times
   * when the sources aren't changing that often. In this case, an
   * implicit precompiled header will be built containing all of the
   * initial includes at the top of the main file (what we refer to as
   * the "preamble" of the file). In subsequent parses, if the
   * preamble or the files in it have not changed, \c
   * clang_reparseTranslationUnit() will re-use the implicit
   * precompiled header to improve parsing performance.
   */
	TU_PrecompiledPreamble = C.CXTranslationUnit_PrecompiledPreamble
  
  /**
   * \brief Used to indicate that the translation unit should cache some
   * code-completion results with each reparse of the source file.
   *
   * Caching of code-completion results is a performance optimization that
   * introduces some overhead to reparsing but improves the performance of
   * code-completion operations.
   */
	TU_CacheCompletionResults = C.CXTranslationUnit_CacheCompletionResults

	/**
   * \brief DEPRECATED: Enable precompiled preambles in C++.
   *
   * Note: this is a *temporary* option that is available only while
   * we are testing C++ precompiled preamble support. It is deprecated.
   */
	TU_CXXPrecompiledPreamble = C.CXTranslationUnit_CXXPrecompiledPreamble

  /**
   * \brief DEPRECATED: Enabled chained precompiled preambles in C++.
   *
   * Note: this is a *temporary* option that is available only while
   * we are testing C++ precompiled preamble support. It is deprecated.
   */
	TU_CXXChainedPCH = C.CXTranslationUnit_CXXChainedPCH
)

/**
 * \brief Parse the given source file and the translation unit corresponding
 * to that file.
 *
 * This routine is the main entry point for the Clang C API, providing the
 * ability to parse a source file into a translation unit that can then be
 * queried by other functions in the API. This routine accepts a set of
 * command-line arguments so that the compilation can be configured in the same
 * way that the compiler is configured on the command line.
 *
 * \param CIdx The index object with which the translation unit will be 
 * associated.
 *
 * \param source_filename The name of the source file to load, or NULL if the
 * source file is included in \p command_line_args.
 *
 * \param command_line_args The command-line arguments that would be
 * passed to the \c clang executable if it were being invoked out-of-process.
 * These command-line options will be parsed and will affect how the translation
 * unit is parsed. Note that the following options are ignored: '-c', 
 * '-emit-ast', '-fsyntex-only' (which is the default), and '-o <output file>'.
 *
 * \param num_command_line_args The number of command-line arguments in
 * \p command_line_args.
 *
 * \param unsaved_files the files that have not yet been saved to disk
 * but may be required for parsing, including the contents of
 * those files.  The contents and name of these files (as specified by
 * CXUnsavedFile) are copied when necessary, so the client only needs to
 * guarantee their validity until the call to this function returns.
 *
 * \param num_unsaved_files the number of unsaved file entries in \p
 * unsaved_files.
 *
 * \param options A bitmask of options that affects how the translation unit
 * is managed but not its compilation. This should be a bitwise OR of the
 * CXTranslationUnit_XXX flags.
 *
 * \returns A new translation unit describing the parsed code and containing
 * any diagnostics produced by the compiler. If there is a failure from which
 * the compiler cannot recover, returns NULL.
 */
//FIXME: handle unsaved-files
func (idx Index) Parse(fname string, args []string, options TranslationUnitFlags) TranslationUnit {
	var c_fname *C.char = nil
	if fname != "" {
		c_fname = C.CString(fname)
	}
	defer C.free(unsafe.Pointer(c_fname))

	c_nargs := C.int(len(args))
	c_cmds := make([]*C.char, len(args))
	for i,_ := range args {
		cstr := C.CString(args[i])
		defer C.free(unsafe.Pointer(cstr))
		c_cmds[i] = cstr
	}
	
	var c_args **C.char = nil
	if len(args) > 0 {
		c_args = &c_cmds[0]
	}
	o := C.clang_parseTranslationUnit(
		idx.c,
		c_fname,
		c_args, c_nargs,
		nil, C.uint(0), 
		C.uint(options))
	return TranslationUnit{o}

}

// A single translation unit, which resides in an index
type TranslationUnit struct {
	c C.CXTranslationUnit
}

func (tu TranslationUnit) File(file_name string) File {
	cfname := C.CString(file_name)
	defer C.free(unsafe.Pointer(cfname))
	f := C.clang_getFile(tu.c, cfname)
	return File{f}
}

// IsFileMultipleIncludeGuarded determines whether the given header is guarded
// against multiple inclusions, either with the conventional
// #ifndef/#define/#endif macro guards or with #pragma once.
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
func (tu TranslationUnit) ToCursor() Cursor {
	o := C.clang_getTranslationUnitCursor(tu.c)
	return Cursor{o}
}

/**
 * \brief Get the original translation unit source file name.
 */
func (tu TranslationUnit) Spelling() string {
	cstr := cxstring{C.clang_getTranslationUnitSpelling(tu.c)}
	defer cstr.Dispose()
	return cstr.String()
}

// CursorOf maps a source location to the cursor that describes the entity at that
// location in the source code.
//
// clang_getCursor() maps an arbitrary source location within a translation
// unit down to the most specific cursor that describes the entity at that
// location. For example, given an expression \c x + y, invoking
// clang_getCursor() with a source location pointing to "x" will return the
// cursor for "x"; similarly for "y". If the cursor points anywhere between
// "x" or "y" (e.g., on the + or the whitespace around it), clang_getCursor()
// will return a cursor referring to the "+" expression.
//
// Returns a cursor representing the entity at the given source location, or
// a NULL cursor if no such entity can be found.
func (tu TranslationUnit) Cursor(loc SourceLocation) Cursor {
	o := C.clang_getCursor(tu.c, loc.c)
	return Cursor{o}
}

/**
 * \brief Saves a translation unit into a serialized representation of
 * that translation unit on disk.
 *
 * Any translation unit that was parsed without error can be saved
 * into a file. The translation unit can then be deserialized into a
 * new \c CXTranslationUnit with \c clang_createTranslationUnit() or,
 * if it is an incomplete translation unit that corresponds to a
 * header, used as a precompiled header when parsing other translation
 * units.
 *
 * \param TU The translation unit to save.
 *
 * \param FileName The file to which the translation unit will be saved.
 *
 * \param options A bitmask of options that affects how the translation unit
 * is saved. This should be a bitwise OR of the
 * CXSaveTranslationUnit_XXX flags.
 *
 * \returns A value that will match one of the enumerators of the CXSaveError
 * enumeration. Zero (CXSaveError_None) indicates that the translation unit was 
 * saved successfully, while a non-zero value indicates that a problem occurred.
 */
func (tu TranslationUnit) Save(fname string, options uint) uint32 {
	cstr := C.CString(fname)
	defer C.free(unsafe.Pointer(cstr))
	o := C.clang_saveTranslationUnit(tu.c, cstr, C.uint(options))
	// FIXME: should be a SaveError type...
	return uint32(o)
}

/**
 * \brief Destroy the specified CXTranslationUnit object.
 */
func (tu TranslationUnit) Dispose() {
	C.clang_disposeTranslationUnit(tu.c)
}

// A particular source file that is part of a translation unit.
type File struct {
	c C.CXFile
}

// Name retrieves the complete file and path name of the given file.
func (c File) Name() string {
	cstr := cxstring{C.clang_getFileName(c.c)}
	defer cstr.Dispose()
	return cstr.String()
}

// ModTime retrieves the last modification time of the given file.
func (c File) ModTime() time.Time {
	// time_t is in seconds since epoch
	sec := C.clang_getFileTime(c.c)
	const nsec = 0
	return time.Unix(int64(sec), nsec)
}



// SourceLocation identifies a specific source location within a translation
// unit.
//
// Use clang_getExpansionLocation() or clang_getSpellingLocation()
// to map a source location to a particular file, line, and column.
type SourceLocation struct {
	c C.CXSourceLocation
}

// SourceRange identifies a half-open character range in the source code.
//
// Use clang_getRangeStart() and clang_getRangeEnd() to retrieve the
// starting and end locations from a source range, respectively.
type SourceRange struct {
	c C.CXSourceRange
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

// Location returns the source location associated with a given file/line/column
// in a particular translation unit.
func (tu TranslationUnit) Location(f File, line, column uint) SourceLocation {
	loc := C.clang_getLocation(tu.c, f.c, C.uint(line), C.uint(column))
	return SourceLocation{loc}
}

// LocationForOffset returns the source location associated with a given
// character offset in a particular translation unit.
func (tu TranslationUnit) LocationForOffset(f File, offset uint) SourceLocation {
	loc := C.clang_getLocationForOffset(tu.c, f.c, C.uint(offset))
	return SourceLocation{loc}
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
	ccol  := C.uint(0)
	coff  := C.uint(0)
	// FIXME: undefined reference to `clang_getExpansionLocation'
	C.clang_getInstantiationLocation(l.c, &f.c, &cline, &ccol, &coff)
	line  = uint(cline)
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
	ccol  := C.uint(0)
	C.clang_getPresumedLocation(l.c, &cname.c, &cline, &ccol)
	fname = cname.String()
	line  = uint(cline)
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
	ccol  := C.uint(0)
	coff  := C.uint(0)
	C.clang_getInstantiationLocation(l.c,
                 &file.c,
                 &cline,
                 &ccol,
                 &coff)
	line  = uint(cline)
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
	ccol  := C.uint(0)
	coff  := C.uint(0)
	C.clang_getSpellingLocation(l.c,
                 &file.c,
                 &cline,
                 &ccol,
                 &coff)
	line  = uint(cline)
	column = uint(ccol)
	offset = uint(coff)
	return
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


// EOF
