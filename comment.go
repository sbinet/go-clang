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
