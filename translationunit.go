package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"
import (
	"unsafe"
)

// A single translation unit, which resides in an index
type TranslationUnit struct {
	c C.CXTranslationUnit
}

/**
 * \brief Perform code completion at a given location in a translation unit.
 *
 * This function performs code completion at a particular file, line, and
 * column within source code, providing results that suggest potential
 * code snippets based on the context of the completion. The basic model
 * for code completion is that Clang will parse a complete source file,
 * performing syntax checking up to the location where code-completion has
 * been requested. At that point, a special code-completion token is passed
 * to the parser, which recognizes this token and determines, based on the
 * current location in the C/Objective-C/C++ grammar and the state of
 * semantic analysis, what completions to provide. These completions are
 * returned via a new \c CXCodeCompleteResults structure.
 *
 * Code completion itself is meant to be triggered by the client when the
 * user types punctuation characters or whitespace, at which point the
 * code-completion location will coincide with the cursor. For example, if \c p
 * is a pointer, code-completion might be triggered after the "-" and then
 * after the ">" in \c p->. When the code-completion location is afer the ">",
 * the completion results will provide, e.g., the members of the struct that
 * "p" points to. The client is responsible for placing the cursor at the
 * beginning of the token currently being typed, then filtering the results
 * based on the contents of the token. For example, when code-completing for
 * the expression \c p->get, the client should provide the location just after
 * the ">" (e.g., pointing at the "g") to this code-completion hook. Then, the
 * client can filter the results based on the current token text ("get"), only
 * showing those results that start with "get". The intent of this interface
 * is to separate the relatively high-latency acquisition of code-completion
 * results from the filtering of results on a per-character basis, which must
 * have a lower latency.
 *
 * \param TU The translation unit in which code-completion should
 * occur. The source files for this translation unit need not be
 * completely up-to-date (and the contents of those source files may
 * be overridden via \p unsaved_files). Cursors referring into the
 * translation unit may be invalidated by this invocation.
 *
 * \param complete_filename The name of the source file where code
 * completion should be performed. This filename may be any file
 * included in the translation unit.
 *
 * \param complete_line The line at which code-completion should occur.
 *
 * \param complete_column The column at which code-completion should occur.
 * Note that the column should point just after the syntactic construct that
 * initiated code completion, and not in the middle of a lexical token.
 *
 * \param unsaved_files the Tiles that have not yet been saved to disk
 * but may be required for parsing or code completion, including the
 * contents of those files.  The contents and name of these files (as
 * specified by CXUnsavedFile) are copied when necessary, so the
 * client only needs to guarantee their validity until the call to
 * this function returns.
 *
 * \param num_unsaved_files The number of unsaved file entries in \p
 * unsaved_files.
 *
 * \param options Extra options that control the behavior of code
 * completion, expressed as a bitwise OR of the enumerators of the
 * CXCodeComplete_Flags enumeration. The
 * \c clang_defaultCodeCompleteOptions() function returns a default set
 * of code-completion options.
 *
 * \returns If successful, a new \c CXCodeCompleteResults structure
 * containing code-completion results, which should eventually be
 * freed with \c clang_disposeCodeCompleteResults(). If code
 * completion fails, returns NULL.
 */
func (tu TranslationUnit) CompleteAt(complete_filename string, complete_line, complete_column int, us UnsavedFiles, options CodeCompleteFlags) CodeCompleteResults {
	cfname := C.CString(complete_filename)
	defer C.free(unsafe.Pointer(cfname))
	c_us := us.to_c()
	defer c_us.Dispose()
	cr := C.clang_codeCompleteAt(tu.c, cfname, C.uint(complete_line), C.uint(complete_column), c_us.ptr(), C.uint(len(c_us)), C.uint(options))
	return CodeCompleteResults{cr}
}

func (tu TranslationUnit) IsValid() bool {
	return tu.c != nil
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
 * \brief Reparse the source files that produced this translation unit.
 *
 * This routine can be used to re-parse the source files that originally
 * created the given translation unit, for example because those source files
 * have changed (either on disk or as passed via \p unsaved_files). The
 * source code will be reparsed with the same command-line options as it
 * was originally parsed.
 *
 * Reparsing a translation unit invalidates all cursors and source locations
 * that refer into that translation unit. This makes reparsing a translation
 * unit semantically equivalent to destroying the translation unit and then
 * creating a new translation unit with the same command-line arguments.
 * However, it may be more efficient to reparse a translation
 * unit using this routine.
 *
 * \param TU The translation unit whose contents will be re-parsed. The
 * translation unit must originally have been built with
 * \c clang_createTranslationUnitFromSourceFile().
 *
 * \param num_unsaved_files The number of unsaved file entries in \p
 * unsaved_files.
 *
 * \param unsaved_files The files that have not yet been saved to disk
 * but may be required for parsing, including the contents of
 * those files.  The contents and name of these files (as specified by
 * CXUnsavedFile) are copied when necessary, so the client only needs to
 * guarantee their validity until the call to this function returns.
 *
 * \param options A bitset of options composed of the flags in CXReparse_Flags.
 * The function \c clang_defaultReparseOptions() produces a default set of
 * options recommended for most uses, based on the translation unit.
 *
 * \returns 0 if the sources could be reparsed. A non-zero value will be
 * returned if reparsing was impossible, such that the translation unit is
 * invalid. In such cases, the only valid call for \p TU is
 * \c clang_disposeTranslationUnit(TU).
 */
func (tu TranslationUnit) Reparse(us UnsavedFiles, options TranslationUnitFlags) int {
	c_us := us.to_c()
	defer c_us.Dispose()
	return int(C.clang_reparseTranslationUnit(tu.c, C.uint(len(c_us)), c_us.ptr(), C.uint(options)))
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
 * \brief Retrieve a diagnostic associated with the given translation unit.
 *
 * \param Unit the translation unit to query.
 * \param Index the zero-based diagnostic number to retrieve.
 *
 * \returns the requested diagnostic. This diagnostic must be freed
 * via a call to \c clang_disposeDiagnostic().
 */
func (tu TranslationUnit) Diagnostics() (ret Diagnostics) {
	ret = make(Diagnostics, C.clang_getNumDiagnostics(tu.c))
	for i := range ret {
		ret[i].c = C.clang_getDiagnostic(tu.c, C.uint(i))
	}
	return
}

/**
 * \brief Destroy the specified CXTranslationUnit object.
 */
func (tu TranslationUnit) Dispose() {
	C.clang_disposeTranslationUnit(tu.c)
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

/**
 * \param Module a module object.
 *
 * \returns the number of top level headers associated with this module.
 */
func (tu TranslationUnit) NumTopLevelHeaders(m Module) int {
	return int(C.clang_Module_getNumTopLevelHeaders(tu.c, m.c))
}

/**
 * \param Module a module object.
 *
 * \param Index top level header index (zero-based).
 *
 * \returns the specified top level header associated with the module.
 */
func (tu TranslationUnit) TopLevelHeader(m Module, i int) File {
	return File{C.clang_Module_getTopLevelHeader(tu.c, m.c, C.unsigned(i))}
}

// TODO
//
// /**
//  * \brief Find #import/#include directives in a specific file.
//  *
//  * \param TU translation unit containing the file to query.
//  *
//  * \param file to search for #import/#include directives.
//  *
//  * \param visitor callback that will receive pairs of CXCursor/CXSourceRange for
//  * each directive found.
//  *
//  * \returns one of the CXResult enumerators.
//  */
// CINDEX_LINKAGE CXResult clang_findIncludesInFile(CXTranslationUnit TU,
//                                                  CXFile file,
//                                               CXCursorAndRangeVisitor visitor);

// TODO
//
// CINDEX_LINKAGE
// CXResult clang_findIncludesInFileWithBlock(CXTranslationUnit, CXFile,
//                                            CXCursorAndRangeVisitorBlock);
