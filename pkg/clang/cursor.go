package clang

// #include <stdlib.h>
// #cgo LDFLAGS: -L/opt/local/libexec/llvm-3.0/lib -lclang
// #include "/opt/local/libexec/llvm-3.0/include/clang-c/Index.h"
// inline static
// CXCursor _go_clang_ocursor_at(CXCursor *c, int idx) {
//   return c[idx];
// }
//
import "C"
import (
	//"unsafe"
)

/**
 * \brief Describes the kind of entity that a cursor refers to.
 */
type CursorKind uint32 //FIXME: use uint32? int64?
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

// Determine whether two cursors are equivalent
func EqualCursors(c1, c2 Cursor) bool {
	o := C.clang_equalCursors(c1.c, c2.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

// IsNull returns true if the underlying Cursor is null
func (c Cursor) IsNull() bool {
	o := C.clang_Cursor_isNull(c.c)
	if o != C.int(0) {
		return true
	}
	return false
}

// Hash computes a hash value for the cursor
func (c Cursor) Hash() uint {
	o := C.clang_hashCursor(c.c)
	return uint(o)
}

// Kind returns the cursor's kind.
func (c Cursor) Kind() CursorKind {
	o := C.clang_getCursorKind(c.c)
	return CursorKind(o)
}

// IsDeclaration determines whether the cursor kind represents a declaration
func (ck CursorKind) IsDeclaration() bool {
	o := C.clang_isDeclaration(uint32(ck))
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * IsReference determines whether the given cursor kind represents a simple
 * reference.
 *
 * Note that other kinds of cursors (such as expressions) can also refer to
 * other cursors. Use clang_getCursorReferenced() to determine whether a
 * particular cursor refers to another entity.
 */
func (ck CursorKind) IsReference() bool {
	o := C.clang_isReference(uint32(ck))
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Determine whether the given cursor kind represents an expression.
 */
func (ck CursorKind) IsExpression() bool {
	o := C.clang_isExpression(uint32(ck))
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Determine whether the given cursor kind represents a statement.
 */
func (ck CursorKind) IsStatement() bool {
	o := C.clang_isStatement(uint32(ck))
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Determine whether the given cursor kind represents an attribute.
 */
func (ck CursorKind) IsAttribute() bool {
	o := C.clang_isAttribute(uint32(ck))
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Determine whether the given cursor kind represents an invalid
 * cursor.
 */
func (ck CursorKind) IsInvalid() bool {
	o := C.clang_isInvalid(uint32(ck))
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Determine whether the given cursor kind represents a translation
 * unit.
 */
func (ck CursorKind) IsTranslationUnit() bool {
	o := C.clang_isTranslationUnit(uint32(ck))
	if o != C.uint(0) {
		return true
	}
	return false
}

/***
 * \brief Determine whether the given cursor represents a preprocessing
 * element, such as a preprocessor directive or macro instantiation.
 */
func (ck CursorKind) IsPreprocessing() bool {
	o := C.clang_isPreprocessing(uint32(ck))
	if o != C.uint(0) {
		return true
	}
	return false
}
  
/***
 * \brief Determine whether the given cursor represents a currently
 *  unexposed piece of the AST (e.g., CXCursor_UnexposedStmt).
 */
func (ck CursorKind) IsUnexposed() bool {
	o := C.clang_isUnexposed(uint32(ck))
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Describe the linkage of the entity referred to by a cursor.
 */
type LinkageKind uint32
const (
	/** \brief This value indicates that no linkage information is available
	 * for a provided CXCursor. */
	LK_NoLinkage LinkageKind = C.CXLinkage_NoLinkage

	// This is the linkage for static variables and static functions.
	LK_Internal = C.CXLinkage_Internal

	// This is the linkage for entities with external linkage that live
	// in C++ anonymous namespaces.
	LK_UniqueExternal = C.CXLinkage_UniqueExternal

	// This is the linkage for entities with true, external linkage
	LK_External = C.CXLinkage_External
)

// Linkage returns the linkage of the entity referred to by a cursor
func (c Cursor) Linkage() LinkageKind {
	o := C.clang_getCursorLinkage(c.c)
	return LinkageKind(o)
}

//
type AvailabilityKind uint32

// Availability returns the availability of the entity that this cursor refers to
func (c Cursor) Availability() AvailabilityKind {
	o := C.clang_getCursorAvailability(c.c)
	return AvailabilityKind(o)
}

// LanguageKind describes the "language" of the entity referred to by a cursor.
type LanguageKind uint32
const (
	LanguageInvalid LanguageKind = C.CXLanguage_Invalid
	LanguageC  = C.CXLanguage_C
	LanguageObjC  = C.CXLanguage_ObjC
	LanguageCPlusPlus = C.CXLanguage_CPlusPlus
)

// Language returns the "language" of the entity referred to by a cursor.
func (c Cursor) Language() LanguageKind {
	o := C.clang_getCursorLanguage(c.c)
	return LanguageKind(o)
}

// TranslationUnit returns the translation unit that a cursor originated from
func (c Cursor) TranslationUnit() TranslationUnit {
	o := C.clang_Cursor_getTranslationUnit(c.c)
	return TranslationUnit{o}
}

// CursorSet is a fast container representing a set of Cursors.
type CursorSet struct {
	c C.CXCursorSet
}

// NewCursorSet creates an empty CursorSet
func NewCursorSet() CursorSet {
	return CursorSet{C.clang_createCXCursorSet()}
}

// Dispose releases the memory associated with a CursorSet
func (c CursorSet) Dispose() {
	C.clang_disposeCXCursorSet(c.c)
}

// Contains queries a CursorSet to see if it contains a specific Cursor
func (c CursorSet) Contains(cursor Cursor) bool {
	o := C.clang_CXCursorSet_contains(c.c, cursor.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

// Insert inserts a Cursor into the set and returns false if the cursor was
// already in that set.
func (c CursorSet) Insert(cursor Cursor) bool {
	o := C.clang_CXCursorSet_insert(c.c, cursor.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Determine the semantic parent of the given cursor.
 *
 * The semantic parent of a cursor is the cursor that semantically contains
 * the given \p cursor. For many declarations, the lexical and semantic parents
 * are equivalent (the lexical parent is returned by 
 * \c clang_getCursorLexicalParent()). They diverge when declarations or
 * definitions are provided out-of-line. For example:
 *
 * \code
 * class C {
 *  void f();
 * };
 *
 * void C::f() { }
 * \endcode
 *
 * In the out-of-line definition of \c C::f, the semantic parent is the 
 * the class \c C, of which this function is a member. The lexical parent is
 * the place where the declaration actually occurs in the source code; in this
 * case, the definition occurs in the translation unit. In general, the 
 * lexical parent for a given entity can change without affecting the semantics
 * of the program, and the lexical parent of different declarations of the
 * same entity may be different. Changing the semantic parent of a declaration,
 * on the other hand, can have a major impact on semantics, and redeclarations
 * of a particular entity should all have the same semantic context.
 *
 * In the example above, both declarations of \c C::f have \c C as their
 * semantic context, while the lexical context of the first \c C::f is \c C
 * and the lexical context of the second \c C::f is the translation unit.
 *
 * For global declarations, the semantic parent is the translation unit.
 */
func (c Cursor) SemanticParent() Cursor {
	o := C.clang_getCursorSemanticParent(c.c)
	return Cursor{o}
}

/**
 * \brief Determine the lexical parent of the given cursor.
 *
 * The lexical parent of a cursor is the cursor in which the given \p cursor
 * was actually written. For many declarations, the lexical and semantic parents
 * are equivalent (the semantic parent is returned by 
 * \c clang_getCursorSemanticParent()). They diverge when declarations or
 * definitions are provided out-of-line. For example:
 *
 * \code
 * class C {
 *  void f();
 * };
 *
 * void C::f() { }
 * \endcode
 *
 * In the out-of-line definition of \c C::f, the semantic parent is the 
 * the class \c C, of which this function is a member. The lexical parent is
 * the place where the declaration actually occurs in the source code; in this
 * case, the definition occurs in the translation unit. In general, the 
 * lexical parent for a given entity can change without affecting the semantics
 * of the program, and the lexical parent of different declarations of the
 * same entity may be different. Changing the semantic parent of a declaration,
 * on the other hand, can have a major impact on semantics, and redeclarations
 * of a particular entity should all have the same semantic context.
 *
 * In the example above, both declarations of \c C::f have \c C as their
 * semantic context, while the lexical context of the first \c C::f is \c C
 * and the lexical context of the second \c C::f is the translation unit.
 *
 * For declarations written in the global scope, the lexical parent is
 * the translation unit.
 */
func (c Cursor) LexicalParent() Cursor {
	o := C.clang_getCursorLexicalParent(c.c)
	return Cursor{o}
}

/**
 * \brief Determine the set of methods that are overridden by the given
 * method.
 *
 * In both Objective-C and C++, a method (aka virtual member function,
 * in C++) can override a virtual method in a base class. For
 * Objective-C, a method is said to override any method in the class's
 * interface (if we're coming from an implementation), its protocols,
 * or its categories, that has the same selector and is of the same
 * kind (class or instance). If no such method exists, the search
 * continues to the class's superclass, its protocols, and its
 * categories, and so on.
 *
 * For C++, a virtual member function overrides any virtual member
 * function with the same signature that occurs in its base
 * classes. With multiple inheritance, a virtual member function can
 * override several virtual member functions coming from different
 * base classes.
 *
 * In all cases, this function determines the immediate overridden
 * method, rather than all of the overridden methods. For example, if
 * a method is originally declared in a class A, then overridden in B
 * (which in inherits from A) and also in C (which inherited from B),
 * then the only overridden method returned from this function when
 * invoked on C's method will be B's method. The client may then
 * invoke this function again, given the previously-found overridden
 * methods, to map out the complete method-override set.
 *
 * \param cursor A cursor representing an Objective-C or C++
 * method. This routine will compute the set of methods that this
 * method overrides.
 * 
 * \param overridden A pointer whose pointee will be replaced with a
 * pointer to an array of cursors, representing the set of overridden
 * methods. If there are no overridden methods, the pointee will be
 * set to NULL. The pointee must be freed via a call to 
 * \c clang_disposeOverriddenCursors().
 *
 * \param num_overridden A pointer to the number of overridden
 * functions, will be set to the number of overridden functions in the
 * array pointed to by \p overridden.
 */
func (c Cursor) OverriddenCursors() (o OverriddenCursors) {
	C.clang_getOverriddenCursors(c.c, &o.c, &o.n)

	return o
}

type OverriddenCursors struct {
	c *C.CXCursor
	n C.uint
}

// Dispose frees the set of overridden cursors
func (c OverriddenCursors) Dispose() {
	C.clang_disposeOverriddenCursors(c.c)
}

func (c OverriddenCursors) Len() int {
	return int(c.n)
}

func (c OverriddenCursors) At(i int) Cursor {
	if i >= int(c.n) {
		panic("clang: index out of range")
	}
	return Cursor{C._go_clang_ocursor_at(c.c, C.int(i))}
}

// IncludedFile returns the file that is included by the given inclusion directive
func (c Cursor) IncludedFile() File {
	o := C.clang_getIncludedFile(c.c)
	return File{o}
}

/**
 * \brief Retrieve the physical location of the source constructor referenced
 * by the given cursor.
 *
 * The location of a declaration is typically the location of the name of that
 * declaration, where the name of that declaration would occur if it is
 * unnamed, or some keyword that introduces that particular declaration.
 * The location of a reference is where that reference occurs within the
 * source code.
 */
func (c Cursor) Location() SourceLocation {
	o := C.clang_getCursorLocation(c.c)
	return SourceLocation{o}
}

/**
 * \brief Retrieve the physical extent of the source construct referenced by
 * the given cursor.
 *
 * The extent of a cursor starts with the file/line/column pointing at the
 * first character within the source construct that the cursor refers to and
 * ends with the last character withinin that source construct. For a
 * declaration, the extent covers the declaration itself. For a reference,
 * the extent covers the location of the reference (e.g., where the referenced
 * entity was actually used).
 */
func (c Cursor) Extent() SourceRange {
	o := C.clang_getCursorExtent(c.c)
	return SourceRange{o}
}

// Type retrieves the type of a cursor (if any).
func (c Cursor) Type() Type {
	o := C.clang_getCursorType(c.c)
	return Type{o}
}

/**
 * \brief Retrieve the result type associated with a given cursor.  This only
 *  returns a valid type of the cursor refers to a function or method.
 */
func (c Cursor) ResultType() Type {
	o := C.clang_getCursorResultType(c.c)
	return Type{o}
}


// EOF
