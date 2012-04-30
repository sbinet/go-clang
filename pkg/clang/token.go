package clang

// #include <stdlib.h>
// #cgo LDFLAGS: -L/opt/local/libexec/llvm-3.0/lib -L/usr/lib/llvm -lclang
// #cgo CFLAGS: -I/opt/local/libexec/llvm-3.0/include -I.
// #include "clang-c/Index.h"
// #include "go-clang.h"
// inline static
// CXCursor _go_clang_ocursor_at(CXCursor *c, int idx) {
//   return c[idx];
// }
//
import "C"

// TokenKind describes a kind of token
type TokenKind uint32
const (
  /**
   * \brief A token that contains some kind of punctuation.
   */
	TK_Punctuation = C.CXToken_Punctuation

  /**
   * \brief A language keyword.
   */
	TK_Keyword = C.CXToken_Keyword

	/**
	 * \brief An identifier (that is not a keyword).
	 */
	TK_Identifier = C.CXToken_Identifier

	/**
	 * \brief A numeric, string, or character literal.
	 */
	TK_Literal = C.CXToken_Literal

	/**
	 * \brief A comment.
	 */
	TK_Comment = C.CXToken_Comment
)