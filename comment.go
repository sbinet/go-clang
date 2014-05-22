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
