package clang

// #include <stdlib.h>
// #include "clang-c/Index.h"
import "C"

import (
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

func (ak AvailabilityKind) String() string {
	switch ak {
	case Available:
		return "Available"
	case Deprecated:
		return "Deprecated"
	case NotAvailable:
		return "NotAvailable"
	case NotAccessible:
		return "NotAccessible"
	}
	return ""
}

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
func (idx Index) CreateTranslationUnitFromSourceFile(fname string, args []string, us UnsavedFiles) TranslationUnit {
	var (
		c_fname *C.char = nil
		c_us            = us.to_c()
	)
	defer c_us.Dispose()
	if fname != "" {
		c_fname = C.CString(fname)
	}
	defer C.free(unsafe.Pointer(c_fname))

	c_nargs := C.int(len(args))
	c_cmds := make([]*C.char, len(args))
	for i, _ := range args {
		cstr := C.CString(args[i])
		defer C.free(unsafe.Pointer(cstr))
		c_cmds[i] = cstr
	}

	var c_argv **C.char = nil
	if c_nargs > 0 {
		c_argv = &c_cmds[0]
	}

	o := C.clang_createTranslationUnitFromSourceFile(
		idx.c,
		c_fname,
		c_nargs, c_argv,
		C.uint(len(c_us)), c_us.ptr())
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
func (idx Index) Parse(fname string, args []string, us UnsavedFiles, options TranslationUnitFlags) TranslationUnit {
	var (
		c_fname *C.char = nil
		c_us            = us.to_c()
	)
	defer c_us.Dispose()
	if fname != "" {
		c_fname = C.CString(fname)
	}
	defer C.free(unsafe.Pointer(c_fname))

	c_nargs := C.int(len(args))
	c_cmds := make([]*C.char, len(args))
	for i, _ := range args {
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
		c_us.ptr(), C.uint(len(c_us)),
		C.uint(options))
	return TranslationUnit{o}

}

// EOF
