package clang

// #include "go-clang.h"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

/**
 * \brief A semantic string that describes a code-completion result.
 *
 * A semantic string that describes the formatting of a code-completion
 * result as a single "template" of text that should be inserted into the
 * source buffer when a particular code-completion result is selected.
 * Each semantic string is made up of some number of "chunks", each of which
 * contains some text along with a description of what that text means, e.g.,
 * the name of the entity being referenced, whether the text chunk is part of
 * the template, or whether it is a "placeholder" that the user should replace
 * with actual code,of a specific kind. See \c CXCompletionChunkKind for a
 * description of the different kinds of chunks.
 */
type CompletionString struct {
	c C.CXCompletionString
}

/**
 * \brief Determine the priority of this code completion.
 *
 * The priority of a code completion indicates how likely it is that this
 * particular completion is the completion that the user will select. The
 * priority is selected by various internal heuristics.
 *
 * \param completion_string The completion string to query.
 *
 * \returns The priority of this completion string. Smaller values indicate
 * higher-priority (more likely) completions.
 */
func (cs CompletionString) Priority() int {
	return int(C.clang_getCompletionPriority(cs.c))
}

/**
 * \brief Determine the availability of the entity that this code-completion
 * string refers to.
 *
 * \param completion_string The completion string to query.
 *
 * \returns The availability of the completion string.
 */
func (cs CompletionString) Availability() AvailabilityKind {
	return AvailabilityKind(C.clang_getCompletionAvailability(cs.c))
}

/**
 * \brief Retrieve the number of annotations associated with the given
 * completion string.
 *
 * \param completion_string the completion string to query.
 *
 * \returns the number of annotations associated with the given completion
 * string.
 */
func (cs CompletionString) NumAnnotations() int {
	return int(C.clang_getCompletionNumAnnotations(cs.c))
}

/**
 * \brief Retrieve the annotation associated with the given completion string.
 *
 * \param completion_string the completion string to query.
 *
 * \param annotation_number the 0-based index of the annotation of the
 * completion string.
 *
 * \returns annotation string associated with the completion at index
 * \c annotation_number, or a NULL string if that annotation is not available.
 */
func (cs CompletionString) Annotation(i int) string {
	cx := cxstring{C.clang_getCompletionAnnotation(cs.c, C.uint(i))}
	defer cx.Dispose()
	return cx.String()
}

/**
 * \brief Retrieve the parent context of the given completion string.
 *
 * The parent context of a completion string is the semantic parent of
 * the declaration (if any) that the code completion represents. For example,
 * a code completion for an Objective-C method would have the method's class
 * or protocol as its context.
 *
 * \param completion_string The code completion string whose parent is
 * being queried.
 *
 * \param kind DEPRECATED: always set to CXCursor_NotImplemented if non-NULL.
 *
 * \returns The name of the completion parent, e.g., "NSObject" if
 * the completion string represents a method in the NSObject class.
 */
func (cs CompletionString) CompletionParent() string {
	o := cxstring{C.clang_getCompletionParent(cs.c, nil)}
	defer o.Dispose()
	return o.String()
}

/**
 * \brief Retrieve the brief documentation comment attached to the declaration
 * that corresponds to the given completion string.
 */
func (cs CompletionString) CompletionBriefComment() string {
	o := cxstring{C.clang_getCompletionBriefComment(cs.c)}
	defer o.Dispose()
	return o.String()
}

/**
 * \brief Retrieve the annotation associated with the given completion string.
 *
 * \param completion_string the completion string to query.
 *
 * \param annotation_number the 0-based index of the annotation of the
 * completion string.
 *
 * \returns annotation string associated with the completion at index
 * \c annotation_number, or a NULL string if that annotation is not available.
 */

func (cs CompletionString) Chunks() (ret []CompletionChunk) {
	ret = make([]CompletionChunk, C.clang_getNumCompletionChunks(cs.c))
	for i := range ret {
		ret[i].cs = cs.c
		ret[i].number = C.uint(i)
	}
	return
}

type CompletionChunk struct {
	cs     C.CXCompletionString
	number C.uint
}

func (cc CompletionChunk) String() string {
	return fmt.Sprintf("%s %s", cc.Kind(), cc.Text())
}

/**
 * \brief Retrieve the text associated with a particular chunk within a
 * completion string.
 *
 * \param completion_string the completion string to query.
 *
 * \param chunk_number the 0-based index of the chunk in the completion string.
 *
 * \returns the text associated with the chunk at index \c chunk_number.
 */
func (cc CompletionChunk) Text() string {
	cx := cxstring{C.clang_getCompletionChunkText(cc.cs, cc.number)}
	defer cx.Dispose()
	return cx.String()
}

/**
 * \brief Determine the kind of a particular chunk within a completion string.
 *
 * \param completion_string the completion string to query.
 *
 * \param chunk_number the 0-based index of the chunk in the completion string.
 *
 * \returns the kind of the chunk at the index \c chunk_number.
 */
func (cs CompletionChunk) Kind() CompletionChunkKind {
	return CompletionChunkKind(C.clang_getCompletionChunkKind(cs.cs, cs.number))
}

/**
 * \brief A single result of code completion.
 */
type CompletionResult struct {
	/**
	 * \brief The kind of entity that this completion refers to.
	 *
	 * The cursor kind will be a macro, keyword, or a declaration (one of the
	 * *Decl cursor kinds), describing the entity that the completion is
	 * referring to.
	 *
	 * \todo In the future, we would like to provide a full cursor, to allow
	 * the client to extract additional information from declaration.
	 */
	CursorKind CursorKind
	/**
	 * \brief The code-completion string that describes how to insert this
	 * code-completion result into the editing buffer.
	 */
	CompletionString CompletionString
}

/**
 * \brief Describes a single piece of text within a code-completion string.
 *
 * Each "chunk" within a code-completion string (\c CXCompletionString) is
 * either a piece of text with a specific "kind" that describes how that text
 * should be interpreted by the client or is another completion string.
 */
type CompletionChunkKind int

const (
	/**
	 * \brief A code-completion string that describes "optional" text that
	 * could be a part of the template (but is not required).
	 *
	 * The Optional chunk is the only kind of chunk that has a code-completion
	 * string for its representation, which is accessible via
	 * \c clang_getCompletionChunkCompletionString(). The code-completion string
	 * describes an additional part of the template that is completely optional.
	 * For example, optional chunks can be used to describe the placeholders for
	 * arguments that match up with defaulted function parameters, e.g. given:
	 *
	 * \code
	 * void f(int x, float y = 3.14, double z = 2.71828);
	 * \endcode
	 *
	 * The code-completion string for this function would contain:
	 *   - a TypedText chunk for "f".
	 *   - a LeftParen chunk for "(".
	 *   - a Placeholder chunk for "int x"
	 *   - an Optional chunk containing the remaining defaulted arguments, e.g.,
	 *       - a Comma chunk for ","
	 *       - a Placeholder chunk for "float y"
	 *       - an Optional chunk containing the last defaulted argument:
	 *           - a Comma chunk for ","
	 *           - a Placeholder chunk for "double z"
	 *   - a RightParen chunk for ")"
	 *
	 * There are many ways to handle Optional chunks. Two simple approaches are:
	 *   - Completely ignore optional chunks, in which case the template for the
	 *     function "f" would only include the first parameter ("int x").
	 *   - Fully expand all optional chunks, in which case the template for the
	 *     function "f" would have all of the parameters.
	 */
	CompletionChunk_Optional CompletionChunkKind = C.CXCompletionChunk_Optional
	/**
	 * \brief Text that a user would be expected to type to get this
	 * code-completion result.
	 *
	 * There will be exactly one "typed text" chunk in a semantic string, which
	 * will typically provide the spelling of a keyword or the name of a
	 * declaration that could be used at the current code point. Clients are
	 * expected to filter the code-completion results based on the text in this
	 * chunk.
	 */
	CompletionChunk_TypedText CompletionChunkKind = C.CXCompletionChunk_TypedText
	/**
	 * \brief Text that should be inserted as part of a code-completion result.
	 *
	 * A "text" chunk represents text that is part of the template to be
	 * inserted into user code should this particular code-completion result
	 * be selected.
	 */
	CompletionChunk_Text CompletionChunkKind = C.CXCompletionChunk_Text
	/**
	 * \brief Placeholder text that should be replaced by the user.
	 *
	 * A "placeholder" chunk marks a place where the user should insert text
	 * into the code-completion template. For example, placeholders might mark
	 * the function parameters for a function declaration, to indicate that the
	 * user should provide arguments for each of those parameters. The actual
	 * text in a placeholder is a suggestion for the text to display before
	 * the user replaces the placeholder with real code.
	 */
	CompletionChunk_Placeholder CompletionChunkKind = C.CXCompletionChunk_Placeholder
	/**
	 * \brief Informative text that should be displayed but never inserted as
	 * part of the template.
	 *
	 * An "informative" chunk contains annotations that can be displayed to
	 * help the user decide whether a particular code-completion result is the
	 * right option, but which is not part of the actual template to be inserted
	 * by code completion.
	 */
	CompletionChunk_Informative CompletionChunkKind = C.CXCompletionChunk_Informative
	/**
	 * \brief Text that describes the current parameter when code-completion is
	 * referring to function call, message send, or template specialization.
	 *
	 * A "current parameter" chunk occurs when code-completion is providing
	 * information about a parameter corresponding to the argument at the
	 * code-completion point. For example, given a function
	 *
	 * \code
	 * int add(int x, int y);
	 * \endcode
	 *
	 * and the source code \c add(, where the code-completion point is after the
	 * "(", the code-completion string will contain a "current parameter" chunk
	 * for "int x", indicating that the current argument will initialize that
	 * parameter. After typing further, to \c add(17, (where the code-completion
	 * point is after the ","), the code-completion string will contain a
	 * "current paremeter" chunk to "int y".
	 */
	CompletionChunk_CurrentParameter CompletionChunkKind = C.CXCompletionChunk_CurrentParameter
	/**
	 * \brief A left parenthesis ('('), used to initiate a function call or
	 * signal the beginning of a function parameter list.
	 */
	CompletionChunk_LeftParen CompletionChunkKind = C.CXCompletionChunk_LeftParen
	/**
	 * \brief A right parenthesis (')'), used to finish a function call or
	 * signal the end of a function parameter list.
	 */
	CompletionChunk_RightParen CompletionChunkKind = C.CXCompletionChunk_RightParen
	/**
	 * \brief A left bracket ('[').
	 */
	CompletionChunk_LeftBracket CompletionChunkKind = C.CXCompletionChunk_LeftBracket
	/**
	 * \brief A right bracket (']').
	 */
	CompletionChunk_RightBracket CompletionChunkKind = C.CXCompletionChunk_RightBracket
	/**
	 * \brief A left brace ('{').
	 */
	CompletionChunk_LeftBrace CompletionChunkKind = C.CXCompletionChunk_LeftBrace
	/**
	 * \brief A right brace ('}').
	 */
	CompletionChunk_RightBrace CompletionChunkKind = C.CXCompletionChunk_RightBrace
	/**
	 * \brief A left angle bracket ('<').
	 */
	CompletionChunk_LeftAngle CompletionChunkKind = C.CXCompletionChunk_LeftAngle
	/**
	 * \brief A right angle bracket ('>').
	 */
	CompletionChunk_RightAngle CompletionChunkKind = C.CXCompletionChunk_RightAngle
	/**
	 * \brief A comma separator (',').
	 */
	CompletionChunk_Comma CompletionChunkKind = C.CXCompletionChunk_Comma
	/**
	 * \brief Text that specifies the result type of a given result.
	 *
	 * This special kind of informative chunk is not meant to be inserted into
	 * the text buffer. Rather, it is meant to illustrate the type that an
	 * expression using the given completion string would have.
	 */
	CompletionChunk_ResultType CompletionChunkKind = C.CXCompletionChunk_ResultType
	/**
	 * \brief A colon (':').
	 */
	CompletionChunk_Colon CompletionChunkKind = C.CXCompletionChunk_Colon
	/**
	 * \brief A semicolon (';').
	 */
	CompletionChunk_SemiColon CompletionChunkKind = C.CXCompletionChunk_SemiColon
	/**
	 * \brief An '=' sign.
	 */
	CompletionChunk_Equal CompletionChunkKind = C.CXCompletionChunk_Equal
	/**
	 * Horizontal space (' ').
	 */
	CompletionChunk_HorizontalSpace CompletionChunkKind = C.CXCompletionChunk_HorizontalSpace
	/**
	 * Vertical space ('\n'), after which it is generally a good idea to
	 * perform indentation.
	 */
	CompletionChunk_VerticalSpace CompletionChunkKind = C.CXCompletionChunk_VerticalSpace
)

func (cck CompletionChunkKind) String() string {
	switch cck {
	case CompletionChunk_Optional:
		return "Optional"
	case CompletionChunk_TypedText:
		return "TypedText"
	case CompletionChunk_Text:
		return "Text"
	case CompletionChunk_Placeholder:
		return "Placeholder"
	case CompletionChunk_Informative:
		return "Informative"
	case CompletionChunk_CurrentParameter:
		return "CurrentParameter"
	case CompletionChunk_LeftParen:
		return "LeftParen"
	case CompletionChunk_RightParen:
		return "RightParen"
	case CompletionChunk_LeftBracket:
		return "LeftBracket"
	case CompletionChunk_RightBracket:
		return "RightBracket"
	case CompletionChunk_LeftBrace:
		return "LeftBrace"
	case CompletionChunk_RightBrace:
		return "RightBrace"
	case CompletionChunk_LeftAngle:
		return "LeftAngle"
	case CompletionChunk_RightAngle:
		return "RightAngle"
	case CompletionChunk_Comma:
		return "Comma"
	case CompletionChunk_ResultType:
		return "ResultType"
	case CompletionChunk_Colon:
		return "Colon"
	case CompletionChunk_SemiColon:
		return "SemiColon"
	case CompletionChunk_Equal:
		return "Equal"
	case CompletionChunk_HorizontalSpace:
		return "HorizontalSpace"
	case CompletionChunk_VerticalSpace:
		return "VerticalSpace"
	default:
		return "Invalid"
	}
}

/**
 * \brief Contains the results of code-completion.
 *
 * This data structure contains the results of code completion, as
 * produced by \c clang_codeCompleteAt(). Its contents must be freed by
 * \c clang_disposeCodeCompleteResults.
 */
type CodeCompleteResults struct {
	c *C.CXCodeCompleteResults
}

func (ccr CodeCompleteResults) IsValid() bool {
	return ccr.c != nil
}

// TODO(): is there a better way to handle this?
func (ccr CodeCompleteResults) Results() (ret []CompletionResult) {
	header := (*reflect.SliceHeader)((unsafe.Pointer(&ret)))
	header.Cap = int(ccr.c.NumResults)
	header.Len = int(ccr.c.NumResults)
	header.Data = uintptr(unsafe.Pointer(ccr.c.Results))
	return
}

/**
 * \brief Sort the code-completion results in case-insensitive alphabetical
 * order.
 *
 * \param Results The set of results to sort.
 * \param NumResults The number of results in \p Results.
 */
func (ccr CodeCompleteResults) Sort() {
	C.clang_sortCodeCompletionResults(ccr.c.Results, ccr.c.NumResults)
}

/**
 * \brief Free the given set of code-completion results.
 */
func (ccr CodeCompleteResults) Dispose() {
	C.clang_disposeCodeCompleteResults(ccr.c)
}

/**
 * \brief Retrieve a diagnostic associated with the given code completion.
 *
 * \param Results the code completion results to query.
 * \param Index the zero-based diagnostic number to retrieve.
 *
 * \returns the requested diagnostic. This diagnostic must be freed
 * via a call to \c clang_disposeDiagnostic().
 */
func (ccr CodeCompleteResults) Diagnostics() (ret Diagnostics) {
	ret = make(Diagnostics, C.clang_codeCompleteGetNumDiagnostics(ccr.c))
	for i := range ret {
		ret[i].c = C.clang_codeCompleteGetDiagnostic(ccr.c, C.uint(i))
	}
	return
}

/**
 * \brief Flags that can be passed to \c clang_codeCompleteAt() to
 * modify its behavior.
 *
 * The enumerators in this enumeration can be bitwise-OR'd together to
 * provide multiple options to \c clang_codeCompleteAt().
 */
type CodeCompleteFlags int

const (
	/**
	 * \brief Whether to include macros within the set of code
	 * completions returned.
	 */
	CodeCompleteFlags_IncludeMacros CodeCompleteFlags = C.CXCodeComplete_IncludeMacros

	/**
	 * \brief Whether to include code patterns for language constructs
	 * within the set of code completions, e.g., for loops.
	 */
	CodeCompleteFlags_IncludeCodePatterns = C.CXCodeComplete_IncludeCodePatterns

	/**
	 * \brief Whether to include brief documentation within the set of code
	 * completions returned.
	 */
	CodeCompleteFlags_IncludeBriefComments = C.CXCodeComplete_IncludeBriefComments
)

/**
 * \brief Bits that represent the context under which completion is occurring.
 *
 * The enumerators in this enumeration may be bitwise-OR'd together if multiple
 * contexts are occurring simultaneously.
 */
type CompletionContext int

const (
	/**
	 * \brief The context for completions is unexposed, as only Clang results
	 * should be included. (This is equivalent to having no context bits set.)
	 */
	CompletionContext_Unexposed CompletionContext = C.CXCompletionContext_Unexposed

	/**
	 * \brief Completions for any possible type should be included in the results.
	 */
	CompletionContext_AnyType CompletionContext = C.CXCompletionContext_AnyType

	/**
	 * \brief Completions for any possible value (variables, function calls, etc.)
	 * should be included in the results.
	 */
	CompletionContext_AnyValue CompletionContext = C.CXCompletionContext_AnyValue
	/**
	 * \brief Completions for values that resolve to an Objective-C object should
	 * be included in the results.
	 */
	CompletionContext_ObjCObjectValue CompletionContext = C.CXCompletionContext_ObjCObjectValue
	/**
	 * \brief Completions for values that resolve to an Objective-C selector
	 * should be included in the results.
	 */
	CompletionContext_ObjCSelectorValue CompletionContext = C.CXCompletionContext_ObjCSelectorValue
	/**
	 * \brief Completions for values that resolve to a C++ class type should be
	 * included in the results.
	 */
	CompletionContext_CXXClassTypeValue CompletionContext = C.CXCompletionContext_CXXClassTypeValue

	/**
	 * \brief Completions for fields of the member being accessed using the dot
	 * operator should be included in the results.
	 */
	CompletionContext_DotMemberAccess CompletionContext = C.CXCompletionContext_DotMemberAccess
	/**
	 * \brief Completions for fields of the member being accessed using the arrow
	 * operator should be included in the results.
	 */
	CompletionContext_ArrowMemberAccess CompletionContext = C.CXCompletionContext_ArrowMemberAccess
	/**
	 * \brief Completions for properties of the Objective-C object being accessed
	 * using the dot operator should be included in the results.
	 */
	CompletionContext_ObjCPropertyAccess CompletionContext = C.CXCompletionContext_ObjCPropertyAccess

	/**
	 * \brief Completions for enum tags should be included in the results.
	 */
	CompletionContext_EnumTag CompletionContext = C.CXCompletionContext_EnumTag
	/**
	 * \brief Completions for union tags should be included in the results.
	 */
	CompletionContext_UnionTag CompletionContext = C.CXCompletionContext_UnionTag
	/**
	 * \brief Completions for struct tags should be included in the results.
	 */
	CompletionContext_StructTag CompletionContext = C.CXCompletionContext_StructTag

	/**
	 * \brief Completions for C++ class names should be included in the results.
	 */
	CompletionContext_ClassTag CompletionContext = C.CXCompletionContext_ClassTag
	/**
	 * \brief Completions for C++ namespaces and namespace aliases should be
	 * included in the results.
	 */
	CompletionContext_Namespace CompletionContext = C.CXCompletionContext_Namespace
	/**
	 * \brief Completions for C++ nested name specifiers should be included in
	 * the results.
	 */
	CompletionContext_NestedNameSpecifier CompletionContext = C.CXCompletionContext_NestedNameSpecifier

	/**
	 * \brief Completions for Objective-C interfaces (classes) should be included
	 * in the results.
	 */
	CompletionContext_ObjCInterface CompletionContext = C.CXCompletionContext_ObjCInterface
	/**
	 * \brief Completions for Objective-C protocols should be included in
	 * the results.
	 */
	CompletionContext_ObjCProtocol CompletionContext = C.CXCompletionContext_ObjCProtocol
	/**
	 * \brief Completions for Objective-C categories should be included in
	 * the results.
	 */
	CompletionContext_ObjCCategory CompletionContext = C.CXCompletionContext_ObjCCategory
	/**
	 * \brief Completions for Objective-C instance messages should be included
	 * in the results.
	 */
	CompletionContext_ObjCInstanceMessage CompletionContext = C.CXCompletionContext_ObjCInstanceMessage
	/**
	 * \brief Completions for Objective-C class messages should be included in
	 * the results.
	 */
	CompletionContext_ObjCClassMessage CompletionContext = C.CXCompletionContext_ObjCClassMessage
	/**
	 * \brief Completions for Objective-C selector names should be included in
	 * the results.
	 */
	CompletionContext_ObjCSelectorName CompletionContext = C.CXCompletionContext_ObjCSelectorName

	/**
	 * \brief Completions for preprocessor macro names should be included in
	 * the results.
	 */
	CompletionContext_MacroName CompletionContext = C.CXCompletionContext_MacroName

	/**
	 * \brief Natural language completions should be included in the results.
	 */
	CompletionContext_NaturalLanguage CompletionContext = C.CXCompletionContext_NaturalLanguage

	/**
	 * \brief The current context is unknown, so set all contexts.
	 */
	CompletionContext_Unknown CompletionContext = C.CXCompletionContext_Unknown
)
