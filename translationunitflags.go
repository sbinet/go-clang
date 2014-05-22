package clang

// #include <stdlib.h>
// #include "clang-c/Index.h"
import "C"

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
	 * \brief Used to indicate that the translation unit will be serialized with
	 * \c clang_saveTranslationUnit.
	 *
	 * This option is typically used when parsing a header with the intent of
	 * producing a precompiled header.
	 */
	TU_ForSerialization = C.CXTranslationUnit_ForSerialization

	/**
	 * \brief DEPRECATED: Enabled chained precompiled preambles in C++.
	 *
	 * Note: this is a *temporary* option that is available only while
	 * we are testing C++ precompiled preamble support. It is deprecated.
	 */
	TU_CXXChainedPCH = C.CXTranslationUnit_CXXChainedPCH

	/**
	 * \brief Used to indicate that function/method bodies should be skipped while
	 * parsing.
	 *
	 * This option can be used to search for declarations/definitions while
	 * ignoring the usages.
	 */
	TU_SkipFunctionBodies = C.CXTranslationUnit_SkipFunctionBodies

	/**
	 * \brief Used to indicate that brief documentation comments should be
	 * included into the set of code completions returned from this translation
	 * unit.
	 */
	TU_IncludeBriefCommentsInCodeCompletion = C.CXTranslationUnit_IncludeBriefCommentsInCodeCompletion
)
