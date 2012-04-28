package clang

// #include <stdlib.h>
// #cgo LDFLAGS: -L/opt/local/libexec/llvm-3.0/lib -lclang
// #include "/opt/local/libexec/llvm-3.0/include/clang-c/Index.h"
import "C"
import (
	//"unsafe"
)

/**
 * \brief Describes the kind of entity that a cursor refers to.
 */
type CursorKind int
const (
  /**
   * \brief A declaration whose specific kind is not exposed via this
   * interface.
   *
   * Unexposed declarations have the same operations as any other kind
   * of declaration; one can extract their location information,
   * spelling, find their definitions, etc. However, the specific kind
   * of the declaration is not reported.
   */
	CK_UnexposedDecl CursorKind = C.CXCursor_UnexposedDecl

	// A C or C++ struct.
	CK_StructDecl = C.CXCursor_StructDecl
)

/**
 * \brief A cursor representing some element in the abstract syntax tree for
 * a translation unit.
 *
 * The cursor abstraction unifies the different kinds of entities in a
 * program--declaration, statements, expressions, references to declarations,
 * etc.--under a single "cursor" abstraction with a common set of operations.
 * Common operation for a cursor include: getting the physical location in
 * a source file where the cursor points, getting the name associated with a
 * cursor, and retrieving cursors for any child nodes of a particular cursor.
 *
 * Cursors can be produced in two specific ways.
 * clang_getTranslationUnitCursor() produces a cursor for a translation unit,
 * from which one can use clang_visitChildren() to explore the rest of the
 * translation unit. clang_getCursor() maps from a physical source location
 * to the entity that resides at that location, allowing one to map from the
 * source code into the AST.
 */
type Cursor struct {
	c C.CXCursor
}

// Retrieve the NULL cursor, which represents no entity.
func GetNullCursor() Cursor {
	return Cursor{C.clang_getNullCursor()}
}

// EOF
