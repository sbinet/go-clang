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
