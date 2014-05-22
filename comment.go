package clang

// #include <stdlib.h>
// #include "clang-c/Index.h"
//
import "C"

/**
 * \brief A comment AST node.
 */
type Comment struct {
	c C.CXComment
}

/**
 * \param Comment AST node of any kind.
 *
 * \returns the type of the AST node.
 */
func (c Comment) Kind() CommentKind {
	return CommentKind(C.clang_Comment_getKind(c.c))
}

/**
 * \param Comment AST node of any kind.
 *
 * \returns number of children of the AST node.
 */
func (c Comment) NumChildren() int {
	return int(C.clang_Comment_getNumChildren(c.c))
}

/**
 * \param Comment AST node of any kind.
 *
 * \param ChildIdx child index (zero-based).
 *
 * \returns the specified child of the AST node.
 */
func (c Comment) Child(idx int) Comment {
	return Comment{C.clang_Comment_getChild(c.c, C.unsigned(idx))}
}

/**
 * \brief A \c CXComment_Paragraph node is considered whitespace if it contains
 * only \c CXComment_Text nodes that are empty or whitespace.
 *
 * Other AST nodes (except \c CXComment_Paragraph and \c CXComment_Text) are
 * never considered whitespace.
 *
 * \returns non-zero if \c Comment is whitespace.
 */
func (c Comment) IsWhitespace() bool {
	o := C.clang_Comment_isWhitespace(c.c)
	if o != 0 {
		return true
	}
	return false
}

/**
 * \returns non-zero if \c Comment is inline content and has a newline
 * immediately following it in the comment text.  Newlines between paragraphs
 * do not count.
 */
func (c Comment) HasTrailingNewline() bool {
	o := C.clang_InlineContentComment_hasTrailingNewline(c.c)
	if 0 != o {
		return true
	}
	return false
}

/**
 * \param Comment a \c CXComment_Text AST node.
 *
 * \returns text contained in the AST node.
 */
func (c Comment) TextComment() string {
	o := cxstring{C.clang_TextComment_getText(c.c)}
	defer o.Dispose()
	return o.String()
}

// TODO: implement more of Comment API

// /**
//  * \param Comment a \c CXComment_InlineCommand AST node.
//  *
//  * \returns name of the inline command.
//  */
// CINDEX_LINKAGE
// CXString clang_InlineCommandComment_getCommandName(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_InlineCommand AST node.
//  *
//  * \returns the most appropriate rendering mode, chosen on command
//  * semantics in Doxygen.
//  */
// CINDEX_LINKAGE enum CXCommentInlineCommandRenderKind
// clang_InlineCommandComment_getRenderKind(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_InlineCommand AST node.
//  *
//  * \returns number of command arguments.
//  */
// CINDEX_LINKAGE
// unsigned clang_InlineCommandComment_getNumArgs(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_InlineCommand AST node.
//  *
//  * \param ArgIdx argument index (zero-based).
//  *
//  * \returns text of the specified argument.
//  */
// CINDEX_LINKAGE
// CXString clang_InlineCommandComment_getArgText(CXComment Comment,
//                                                unsigned ArgIdx);

// /**
//  * \param Comment a \c CXComment_HTMLStartTag or \c CXComment_HTMLEndTag AST
//  * node.
//  *
//  * \returns HTML tag name.
//  */
// CINDEX_LINKAGE CXString clang_HTMLTagComment_getTagName(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_HTMLStartTag AST node.
//  *
//  * \returns non-zero if tag is self-closing (for example, &lt;br /&gt;).
//  */
// CINDEX_LINKAGE
// unsigned clang_HTMLStartTagComment_isSelfClosing(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_HTMLStartTag AST node.
//  *
//  * \returns number of attributes (name-value pairs) attached to the start tag.
//  */
// CINDEX_LINKAGE unsigned clang_HTMLStartTag_getNumAttrs(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_HTMLStartTag AST node.
//  *
//  * \param AttrIdx attribute index (zero-based).
//  *
//  * \returns name of the specified attribute.
//  */
// CINDEX_LINKAGE
// CXString clang_HTMLStartTag_getAttrName(CXComment Comment, unsigned AttrIdx);

// /**
//  * \param Comment a \c CXComment_HTMLStartTag AST node.
//  *
//  * \param AttrIdx attribute index (zero-based).
//  *
//  * \returns value of the specified attribute.
//  */
// CINDEX_LINKAGE
// CXString clang_HTMLStartTag_getAttrValue(CXComment Comment, unsigned AttrIdx);

// /**
//  * \param Comment a \c CXComment_BlockCommand AST node.
//  *
//  * \returns name of the block command.
//  */
// CINDEX_LINKAGE
// CXString clang_BlockCommandComment_getCommandName(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_BlockCommand AST node.
//  *
//  * \returns number of word-like arguments.
//  */
// CINDEX_LINKAGE
// unsigned clang_BlockCommandComment_getNumArgs(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_BlockCommand AST node.
//  *
//  * \param ArgIdx argument index (zero-based).
//  *
//  * \returns text of the specified word-like argument.
//  */
// CINDEX_LINKAGE
// CXString clang_BlockCommandComment_getArgText(CXComment Comment,
//                                               unsigned ArgIdx);

// /**
//  * \param Comment a \c CXComment_BlockCommand or
//  * \c CXComment_VerbatimBlockCommand AST node.
//  *
//  * \returns paragraph argument of the block command.
//  */
// CINDEX_LINKAGE
// CXComment clang_BlockCommandComment_getParagraph(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_ParamCommand AST node.
//  *
//  * \returns parameter name.
//  */
// CINDEX_LINKAGE
// CXString clang_ParamCommandComment_getParamName(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_ParamCommand AST node.
//  *
//  * \returns non-zero if the parameter that this AST node represents was found
//  * in the function prototype and \c clang_ParamCommandComment_getParamIndex
//  * function will return a meaningful value.
//  */
// CINDEX_LINKAGE
// unsigned clang_ParamCommandComment_isParamIndexValid(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_ParamCommand AST node.
//  *
//  * \returns zero-based parameter index in function prototype.
//  */
// CINDEX_LINKAGE
// unsigned clang_ParamCommandComment_getParamIndex(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_ParamCommand AST node.
//  *
//  * \returns non-zero if parameter passing direction was specified explicitly in
//  * the comment.
//  */
// CINDEX_LINKAGE
// unsigned clang_ParamCommandComment_isDirectionExplicit(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_ParamCommand AST node.
//  *
//  * \returns parameter passing direction.
//  */
// CINDEX_LINKAGE
// enum CXCommentParamPassDirection clang_ParamCommandComment_getDirection(
//                                                             CXComment Comment);

// /**
//  * \param Comment a \c CXComment_TParamCommand AST node.
//  *
//  * \returns template parameter name.
//  */
// CINDEX_LINKAGE
// CXString clang_TParamCommandComment_getParamName(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_TParamCommand AST node.
//  *
//  * \returns non-zero if the parameter that this AST node represents was found
//  * in the template parameter list and
//  * \c clang_TParamCommandComment_getDepth and
//  * \c clang_TParamCommandComment_getIndex functions will return a meaningful
//  * value.
//  */
// CINDEX_LINKAGE
// unsigned clang_TParamCommandComment_isParamPositionValid(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_TParamCommand AST node.
//  *
//  * \returns zero-based nesting depth of this parameter in the template parameter list.
//  *
//  * For example,
//  * \verbatim
//  *     template<typename C, template<typename T> class TT>
//  *     void test(TT<int> aaa);
//  * \endverbatim
//  * for C and TT nesting depth is 0,
//  * for T nesting depth is 1.
//  */
// CINDEX_LINKAGE
// unsigned clang_TParamCommandComment_getDepth(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_TParamCommand AST node.
//  *
//  * \returns zero-based parameter index in the template parameter list at a
//  * given nesting depth.
//  *
//  * For example,
//  * \verbatim
//  *     template<typename C, template<typename T> class TT>
//  *     void test(TT<int> aaa);
//  * \endverbatim
//  * for C and TT nesting depth is 0, so we can ask for index at depth 0:
//  * at depth 0 C's index is 0, TT's index is 1.
//  *
//  * For T nesting depth is 1, so we can ask for index at depth 0 and 1:
//  * at depth 0 T's index is 1 (same as TT's),
//  * at depth 1 T's index is 0.
//  */
// CINDEX_LINKAGE
// unsigned clang_TParamCommandComment_getIndex(CXComment Comment, unsigned Depth);

// /**
//  * \param Comment a \c CXComment_VerbatimBlockLine AST node.
//  *
//  * \returns text contained in the AST node.
//  */
// CINDEX_LINKAGE
// CXString clang_VerbatimBlockLineComment_getText(CXComment Comment);

// /**
//  * \param Comment a \c CXComment_VerbatimLine AST node.
//  *
//  * \returns text contained in the AST node.
//  */
// CINDEX_LINKAGE CXString clang_VerbatimLineComment_getText(CXComment Comment);

// /**
//  * \brief Convert an HTML tag AST node to string.
//  *
//  * \param Comment a \c CXComment_HTMLStartTag or \c CXComment_HTMLEndTag AST
//  * node.
//  *
//  * \returns string containing an HTML tag.
//  */
// CINDEX_LINKAGE CXString clang_HTMLTagComment_getAsString(CXComment Comment);

// /**
//  * \brief Convert a given full parsed comment to an HTML fragment.
//  *
//  * Specific details of HTML layout are subject to change.  Don't try to parse
//  * this HTML back into an AST, use other APIs instead.
//  *
//  * Currently the following CSS classes are used:
//  * \li "para-brief" for \\brief paragraph and equivalent commands;
//  * \li "para-returns" for \\returns paragraph and equivalent commands;
//  * \li "word-returns" for the "Returns" word in \\returns paragraph.
//  *
//  * Function argument documentation is rendered as a \<dl\> list with arguments
//  * sorted in function prototype order.  CSS classes used:
//  * \li "param-name-index-NUMBER" for parameter name (\<dt\>);
//  * \li "param-descr-index-NUMBER" for parameter description (\<dd\>);
//  * \li "param-name-index-invalid" and "param-descr-index-invalid" are used if
//  * parameter index is invalid.
//  *
//  * Template parameter documentation is rendered as a \<dl\> list with
//  * parameters sorted in template parameter list order.  CSS classes used:
//  * \li "tparam-name-index-NUMBER" for parameter name (\<dt\>);
//  * \li "tparam-descr-index-NUMBER" for parameter description (\<dd\>);
//  * \li "tparam-name-index-other" and "tparam-descr-index-other" are used for
//  * names inside template template parameters;
//  * \li "tparam-name-index-invalid" and "tparam-descr-index-invalid" are used if
//  * parameter position is invalid.
//  *
//  * \param Comment a \c CXComment_FullComment AST node.
//  *
//  * \returns string containing an HTML fragment.
//  */
// CINDEX_LINKAGE CXString clang_FullComment_getAsHTML(CXComment Comment);

// /**
//  * \brief Convert a given full parsed comment to an XML document.
//  *
//  * A Relax NG schema for the XML can be found in comment-xml-schema.rng file
//  * inside clang source tree.
//  *
//  * \param Comment a \c CXComment_FullComment AST node.
//  *
//  * \returns string containing an XML document.
//  */
// CINDEX_LINKAGE CXString clang_FullComment_getAsXML(CXComment Comment);
