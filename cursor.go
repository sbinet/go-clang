package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"
import (
	"unsafe"
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
func NewNullCursor() Cursor {
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

// Spelling returns the name of the entity referenced by this cursor.
func (c Cursor) Spelling() string {
	cstr := cxstring{C.clang_getCursorSpelling(c.c)}
	defer cstr.Dispose()
	return cstr.String()
}

/**
 * \brief Retrieve the display name for the entity referenced by this cursor.
 *
 * The display name contains extra information that helps identify the cursor,
 * such as the parameters of a function or template or the arguments of a
 * class template specialization.
 */
func (c Cursor) DisplayName() string {
	cstr := cxstring{C.clang_getCursorDisplayName(c.c)}
	defer cstr.Dispose()
	return cstr.String()
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

// Availability returns the availability of the entity that this cursor refers to
func (c Cursor) Availability() AvailabilityKind {
	o := C.clang_getCursorAvailability(c.c)
	return AvailabilityKind(o)
}

/**
 * \brief Determine the availability of the entity that this cursor refers to
 * on any platforms for which availability information is known.
 *
 * \param cursor The cursor to query.
 *
 * \param always_deprecated If non-NULL, will be set to indicate whether the
 * entity is deprecated on all platforms.
 *
 * \param deprecated_message If non-NULL, will be set to the message text
 * provided along with the unconditional deprecation of this entity. The client
 * is responsible for deallocating this string.
 *
 * \param always_unavailable If non-NULL, will be set to indicate whether the
 * entity is unavailable on all platforms.
 *
 * \param unavailable_message If non-NULL, will be set to the message text
 * provided along with the unconditional unavailability of this entity. The
 * client is responsible for deallocating this string.
 *
 * \param availability If non-NULL, an array of CXPlatformAvailability instances
 * that will be populated with platform availability information, up to either
 * the number of platforms for which availability information is available (as
 * returned by this function) or \c availability_size, whichever is smaller.
 *
 * \param availability_size The number of elements available in the
 * \c availability array.
 *
 * \returns The number of platforms (N) for which availability information is
 * available (which is unrelated to \c availability_size).
 *
 * Note that the client is responsible for calling
 * \c clang_disposeCXPlatformAvailability to free each of the
 * platform-availability structures returned. There are
 * \c min(N, availability_size) such structures.
 */
func (c Cursor) PlatformAvailability(availability []PlatformAvailability) (always_deprecated bool, deprecated_msg string, always_unavailable bool, unavailable_msg string) {
	var c_always_deprecated C.int
	var c_deprecated_msg cxstring
	defer c_deprecated_msg.Dispose()
	var c_always_unavailable C.int
	var c_unavailable_msg cxstring
	defer c_unavailable_msg.Dispose()
	c_platforms := make([]C.CXPlatformAvailability, len(availability))

	nn := int(C.clang_getCursorPlatformAvailability(
		c.c,
		&c_always_deprecated,
		&c_deprecated_msg.c,
		&c_always_unavailable,
		&c_unavailable_msg.c,
		&c_platforms[0],
		C.int(len(c_platforms)),
	))

	if c_always_deprecated != 0 {
		always_deprecated = true
	}
	deprecated_msg = c_deprecated_msg.String()

	if c_always_unavailable != 0 {
		always_unavailable = true
	}
	unavailable_msg = c_unavailable_msg.String()

	if nn > len(availability) {
		nn = len(availability)
	}

	availability = make([]PlatformAvailability, nn)
	for i := 0; i < nn; i++ {
		availability[i] = PlatformAvailability{C._goclang_get_platform_availability_at(&c_platforms[0], C.int(i))}
	}

	return
}

// LanguageKind describes the "language" of the entity referred to by a cursor.
type LanguageKind uint32

const (
	LanguageInvalid   LanguageKind = C.CXLanguage_Invalid
	LanguageC                      = C.CXLanguage_C
	LanguageObjC                   = C.CXLanguage_ObjC
	LanguageCPlusPlus              = C.CXLanguage_CPlusPlus
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

// DeclObjCTypeEncoding returns the Objective-C type encoding for the
// specified declaration.
func (c Cursor) DeclObjCTypeEncoding() string {
	o := C.clang_getDeclObjCTypeEncoding(c.c)
	cstr := cxstring{o}
	defer cstr.Dispose()
	return cstr.String()
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
 * \brief Retrieve the underlying type of a typedef declaration.
 *
 * If the cursor does not reference a typedef declaration, an invalid type is
 * returned.
 */
func (c Cursor) TypedefDeclUnderlyingType() Type {
	o := C.clang_getTypedefDeclUnderlyingType(c.c)
	return Type{o}
}

/**
 * \brief Retrieve the integer type of an enum declaration.
 *
 * If the cursor does not reference an enum declaration, an invalid type is
 * returned.
 */
func (c Cursor) EnumDeclIntegerType() Type {
	o := C.clang_getEnumDeclIntegerType(c.c)
	return Type{o}
}

/**
 * \brief Retrieve the integer value of an enum constant declaration as a signed
 *  long long.
 *
 * If the cursor does not reference an enum constant declaration, LLONG_MIN is returned.
 * Since this is also potentially a valid constant value, the kind of the cursor
 * must be verified before calling this function.
 */
func (c Cursor) EnumConstantDeclValue() int64 {
	return int64(C.clang_getEnumConstantDeclValue(c.c))
}

/**
 * \brief Retrieve the integer value of an enum constant declaration as an unsigned
 *  long long.
 *
 * If the cursor does not reference an enum constant declaration, ULLONG_MAX is returned.
 * Since this is also potentially a valid constant value, the kind of the cursor
 * must be verified before calling this function.
 */
func (c Cursor) EnumConstantDeclUnsignedValue() uint64 {
	return uint64(C.clang_getEnumConstantDeclUnsignedValue(c.c))
}

/**
 * \brief Retrieve the bit width of a bit field declaration as an integer.
 *
 * If a cursor that is not a bit field declaration is passed in, -1 is returned.
 */
func (c Cursor) FieldDeclBitWidth() int {
	return int(C.clang_getFieldDeclBitWidth(c.c))
}

/**
 * \brief Retrieve the number of non-variadic arguments associated with a given
 * cursor.
 *
 * If a cursor that is not a function or method is passed in, -1 is returned.
 */
// CINDEX_LINKAGE int clang_Cursor_getNumArguments(CXCursor C);
func (c Cursor) NumArguments() int {
	n := C.clang_Cursor_getNumArguments(c.c)
	return int(n)
}

/**
 * \brief Retrieve the argument cursor of a function or method.
 *
 * If a cursor that is not a function or method is passed in or the index
 * exceeds the number of arguments, an invalid cursor is returned.
 */
// CINDEX_LINKAGE CXCursor clang_Cursor_getArgument(CXCursor C, unsigned i);
func (c Cursor) Argument(i uint) Cursor {
	o := C.clang_Cursor_getArgument(c.c, C.uint(i))
	return Cursor{o}
}

/**
 * \brief Retrieve the result type associated with a given cursor.  This only
 *  returns a valid type of the cursor refers to a function or method.
 */
func (c Cursor) ResultType() Type {
	o := C.clang_getCursorResultType(c.c)
	return Type{o}
}

/**
 * \brief Returns non-zero if the cursor specifies a Record member that is a
 *   bitfield.
 */
func (c Cursor) IsBitField() bool {
	o := C.clang_Cursor_isBitField(c.c)
	if o != 0 {
		return true
	}
	return false
}

/**
 * \brief Returns 1 if the base class specified by the cursor with kind
 *   CX_CXXBaseSpecifier is virtual.
 */
func (c Cursor) IsVirtualBase() bool {
	o := C.clang_isVirtualBase(c.c)
	return o == C.uint(1)
}

/**
 * \brief Represents the C++ access control level to a base class for a
 * cursor with kind CX_CXXBaseSpecifier.
 */
type AccessSpecifier uint32

const (
	AS_Invalid   AccessSpecifier = C.CX_CXXInvalidAccessSpecifier
	AS_Public                    = C.CX_CXXPublic
	AS_Protected                 = C.CX_CXXProtected
	AS_Private                   = C.CX_CXXPrivate
)

/**
 * \brief Returns the access control level for the C++ base specifier
 * represented by a cursor with kind CXCursor_CXXBaseSpecifier or
 * CXCursor_AccessSpecifier.
 */
func (c Cursor) AccessSpecifier() AccessSpecifier {
	o := C.clang_getCXXAccessSpecifier(c.c)
	return AccessSpecifier(o)
}

/**
 * \brief Determine the number of overloaded declarations referenced by a
 * \c CXCursor_OverloadedDeclRef cursor.
 *
 * \param cursor The cursor whose overloaded declarations are being queried.
 *
 * \returns The number of overloaded declarations referenced by \c cursor. If it
 * is not a \c CXCursor_OverloadedDeclRef cursor, returns 0.
 */
func (c Cursor) NumOverloadedDecls() int {
	o := C.clang_getNumOverloadedDecls(c.c)
	return int(o)
}

/**
 * \brief Retrieve a cursor for one of the overloaded declarations referenced
 * by a \c CXCursor_OverloadedDeclRef cursor.
 *
 * \param cursor The cursor whose overloaded declarations are being queried.
 *
 * \param index The zero-based index into the set of overloaded declarations in
 * the cursor.
 *
 * \returns A cursor representing the declaration referenced by the given
 * \c cursor at the specified \c index. If the cursor does not have an
 * associated set of overloaded declarations, or if the index is out of bounds,
 * returns \c clang_getNullCursor();
 */
func (c Cursor) OverloadedDecl(i int) Cursor {
	o := C.clang_getOverloadedDecl(c.c, C.uint(i))
	return Cursor{o}
}

/**
 * \brief For cursors representing an iboutletcollection attribute,
 *  this function returns the collection element type.
 *
 */
func (c Cursor) IBOutletCollectionType() Type {
	o := C.clang_getIBOutletCollectionType(c.c)
	return Type{o}
}

/**
 * \brief Describes how the traversal of the children of a particular
 * cursor should proceed after visiting a particular child cursor.
 *
 * A value of this enumeration type should be returned by each
 * \c CXCursorVisitor to indicate how clang_visitChildren() proceed.
 */
type ChildVisitResult uint32

const (
	/**
	 * \brief Terminates the cursor traversal.
	 */
	CVR_Break ChildVisitResult = C.CXChildVisit_Break

	/**
	 * \brief Continues the cursor traversal with the next sibling of
	 * the cursor just visited, without visiting its children.
	 */
	CVR_Continue = C.CXChildVisit_Continue

	/**
	 * \brief Recursively traverse the children of this cursor, using
	 * the same visitor and client data.
	 */
	CVR_Recurse = C.CXChildVisit_Recurse
)

/**
 * \brief Visitor invoked for each cursor found by a traversal.
 *
 * This visitor function will be invoked for each cursor found by
 * clang_visitCursorChildren(). Its first argument is the cursor being
 * visited, its second argument is the parent visitor for that cursor,
 * and its third argument is the client data provided to
 * clang_visitCursorChildren().
 *
 * The visitor should return one of the \c CXChildVisitResult values
 * to direct clang_visitCursorChildren().
 */
type CursorVisitor func(cursor, parent Cursor) (status ChildVisitResult)

/**
 * \brief Visit the children of a particular cursor.
 *
 * This function visits all the direct children of the given cursor,
 * invoking the given \p visitor function with the cursors of each
 * visited child. The traversal may be recursive, if the visitor returns
 * \c CXChildVisit_Recurse. The traversal may also be ended prematurely, if
 * the visitor returns \c CXChildVisit_Break.
 *
 * \param parent the cursor whose child may be visited. All kinds of
 * cursors can be visited, including invalid cursors (which, by
 * definition, have no children).
 *
 * \param visitor the visitor function that will be invoked for each
 * child of \p parent.
 *
 * \param client_data pointer data supplied by the client, which will
 * be passed to the visitor each time it is invoked.
 *
 * \returns a non-zero value if the traversal was terminated
 * prematurely by the visitor returning \c CXChildVisit_Break.
 */
func (c Cursor) Visit(visitor CursorVisitor) bool {
	forceEscapeVisitor = &visitor
	o := C._go_clang_visit_children(c.c, unsafe.Pointer(&visitor))
	if o != C.uint(0) {
		return false
	}
	return true
}

// forceEscapeVisitor is write-only: to force compiler to escape the address
// (else the address can become stale if the goroutine stack needs to grow
// and is forced to move)
// Explained by rsc in https://golang.org/issue/9125
var forceEscapeVisitor *CursorVisitor

//export GoClangCursorVisitor
func GoClangCursorVisitor(cursor, parent C.CXCursor, cfct unsafe.Pointer) (status ChildVisitResult) {
	fct := *(*CursorVisitor)(cfct)
	o := fct(Cursor{cursor}, Cursor{parent})
	return o
}

/**
 * \brief Retrieve a Unified Symbol Resolution (USR) for the entity referenced
 * by the given cursor.
 *
 * A Unified Symbol Resolution (USR) is a string that identifies a particular
 * entity (function, class, variable, etc.) within a program. USRs can be
 * compared across translation units to determine, e.g., when references in
 * one translation refer to an entity defined in another translation unit.
 */
func (c Cursor) USR() string {
	cstr := cxstring{C.clang_getCursorUSR(c.c)}
	defer cstr.Dispose()
	return cstr.String()
}

//FIXME
// /**
//  * \brief Construct a USR for a specified Objective-C class.
//  */
// CINDEX_LINKAGE CXString clang_constructUSR_ObjCClass(const char *class_name);

// /**
//  * \brief Construct a USR for a specified Objective-C category.
//  */
// CINDEX_LINKAGE CXString
//   clang_constructUSR_ObjCCategory(const char *class_name,
//                                  const char *category_name);

// /**
//  * \brief Construct a USR for a specified Objective-C protocol.
//  */
// CINDEX_LINKAGE CXString
//   clang_constructUSR_ObjCProtocol(const char *protocol_name);

// /**
//  * \brief Construct a USR for a specified Objective-C instance variable and
//  *   the USR for its containing class.
//  */
// CINDEX_LINKAGE CXString clang_constructUSR_ObjCIvar(const char *name,
//                                                     CXString classUSR);

// /**
//  * \brief Construct a USR for a specified Objective-C method and
//  *   the USR for its containing class.
//  */
// CINDEX_LINKAGE CXString clang_constructUSR_ObjCMethod(const char *name,
//                                                       unsigned isInstanceMethod,
//                                                       CXString classUSR);

// /**
//  * \brief Construct a USR for a specified Objective-C property and the USR
//  *  for its containing class.
//  */
// CINDEX_LINKAGE CXString clang_constructUSR_ObjCProperty(const char *property,
//                                                         CXString classUSR);

/** \brief For a cursor that is a reference, retrieve a cursor representing the
 * entity that it references.
 *
 * Reference cursors refer to other entities in the AST. For example, an
 * Objective-C superclass reference cursor refers to an Objective-C class.
 * This function produces the cursor for the Objective-C class from the
 * cursor for the superclass reference. If the input cursor is a declaration or
 * definition, it returns that declaration or definition unchanged.
 * Otherwise, returns the NULL cursor.
 */
func (c Cursor) Referenced() Cursor {
	o := C.clang_getCursorReferenced(c.c)
	return Cursor{o}
}

/**
 *  \brief For a cursor that is either a reference to or a declaration
 *  of some entity, retrieve a cursor that describes the definition of
 *  that entity.
 *
 *  Some entities can be declared multiple times within a translation
 *  unit, but only one of those declarations can also be a
 *  definition. For example, given:
 *
 *  \code
 *  int f(int, int);
 *  int g(int x, int y) { return f(x, y); }
 *  int f(int a, int b) { return a + b; }
 *  int f(int, int);
 *  \endcode
 *
 *  there are three declarations of the function "f", but only the
 *  second one is a definition. The clang_getCursorDefinition()
 *  function will take any cursor pointing to a declaration of "f"
 *  (the first or fourth lines of the example) or a cursor referenced
 *  that uses "f" (the call to "f' inside "g") and will return a
 *  declaration cursor pointing to the definition (the second "f"
 *  declaration).
 *
 *  If given a cursor for which there is no corresponding definition,
 *  e.g., because there is no definition of that entity within this
 *  translation unit, returns a NULL cursor.
 */
func (c Cursor) DefinitionCursor() Cursor {
	o := C.clang_getCursorDefinition(c.c)
	return Cursor{o}
}

/**
 * \brief Determine whether the declaration pointed to by this cursor
 * is also a definition of that entity.
 */
func (c Cursor) IsDefinition() bool {
	o := C.clang_isCursorDefinition(c.c)
	if o != C.uint(0) {
		return true
	}
	return false
}

/**
 * \brief Retrieve the canonical cursor corresponding to the given cursor.
 *
 * In the C family of languages, many kinds of entities can be declared several
 * times within a single translation unit. For example, a structure type can
 * be forward-declared (possibly multiple times) and later defined:
 *
 * \code
 * struct X;
 * struct X;
 * struct X {
 *   int member;
 * };
 * \endcode
 *
 * The declarations and the definition of \c X are represented by three
 * different cursors, all of which are declarations of the same underlying
 * entity. One of these cursor is considered the "canonical" cursor, which
 * is effectively the representative for the underlying entity. One can
 * determine if two cursors are declarations of the same underlying entity by
 * comparing their canonical cursors.
 *
 * \returns The canonical cursor for the entity referred to by the given cursor.
 */
func (c Cursor) CanonicalCursor() Cursor {
	o := C.clang_getCanonicalCursor(c.c)
	return Cursor{o}
}

/**
 * \brief Given a cursor pointing to a C++ method call or an ObjC message,
 * returns non-zero if the method/message is "dynamic", meaning:
 *
 * For a C++ method: the call is virtual.
 * For an ObjC message: the receiver is an object instance, not 'super' or a
 * specific class.
 *
 * If the method/message is "static" or the cursor does not point to a
 * method/message, it will return zero.
 */
func (c Cursor) IsDynamicCall() bool {
	o := C.clang_Cursor_isDynamicCall(c.c)
	if o != 0 {
		return true
	}
	return false
}

/**
 * \brief Given a cursor pointing to an ObjC message, returns the CXType of the
 * receiver.
 */
func (c Cursor) ReceiverType() Type {
	return Type{C.clang_Cursor_getReceiverType(c.c)}
}

// TODO: objective-c support
// /**
//  * \brief Given a cursor that represents a property declaration, return the
//  * associated property attributes. The bits are formed from
//  * \c CXObjCPropertyAttrKind.
//  *
//  * \param reserved Reserved for future use, pass 0.
//  */
// CINDEX_LINKAGE unsigned clang_Cursor_getObjCPropertyAttributes(CXCursor C,
//                                                              unsigned reserved);

// TODO: objective-c support
// /**
//  * \brief Given a cursor that represents an ObjC method or parameter
//  * declaration, return the associated ObjC qualifiers for the return type or the
//  * parameter respectively. The bits are formed from CXObjCDeclQualifierKind.
//  */
// CINDEX_LINKAGE unsigned clang_Cursor_getObjCDeclQualifiers(CXCursor C);

// TODO: objective-c support
// /**
//  * \brief Given a cursor that represents an ObjC method or property declaration,
//  * return non-zero if the declaration was affected by "@optional".
//  * Returns zero if the cursor is not such a declaration or it is "@required".
//  */
// CINDEX_LINKAGE unsigned clang_Cursor_isObjCOptional(CXCursor C);

/**
 * \brief Returns non-zero if the given cursor is a variadic function or method.
 */
func (c Cursor) IsVariadic() bool {
	o := C.clang_Cursor_isVariadic(c.c)
	if o != 0 {
		return true
	}
	return false
}

/**
 * \brief Given a cursor that represents a declaration, return the associated
 * comment's source range.  The range may include multiple consecutive comments
 * with whitespace in between.
 */
func (c Cursor) CommentRange() SourceRange {
	return SourceRange{C.clang_Cursor_getCommentRange(c.c)}
}

/**
 * \brief Given a cursor that represents a declaration, return the associated
 * comment text, including comment markers.
 */
func (c Cursor) RawCommentText() string {
	cstr := cxstring{C.clang_Cursor_getRawCommentText(c.c)}
	defer cstr.Dispose()
	return cstr.String()
}

/**
 * \brief Given a cursor that represents a documentable entity (e.g.,
 * declaration), return the associated \\brief paragraph; otherwise return the
 * first paragraph.
 */
func (c Cursor) BriefCommentText() string {
	cstr := cxstring{C.clang_Cursor_getBriefCommentText(c.c)}
	defer cstr.Dispose()
	return cstr.String()
}

/**
 * \brief Given a cursor that represents a documentable entity (e.g.,
 * declaration), return the associated parsed comment as a
 * \c CXComment_FullComment AST node.
 */
func (c Cursor) ParsedComment() Comment {
	return Comment{C.clang_Cursor_getParsedComment(c.c)}
}

/**
 * \brief Given a CXCursor_ModuleImportDecl cursor, return the associated module.
 */
func (c Cursor) Module() Module {
	return Module{C.clang_Cursor_getModule(c.c)}
}

/**
 * \defgroup CINDEX_CPP C++ AST introspection
 *
 * The routines in this group provide access information in the ASTs specific
 * to C++ language features.
 *
 * @{
 */

/**
 * \brief Determine if a C++ member function or member function template is
 * pure virtual.
 */
func (c Cursor) CXXMethod_IsPureVirtual() bool {
	o := C.clang_CXXMethod_isPureVirtual(c.c)
	if o != 0 {
		return true
	}
	return false
}

/**
 * \brief Determine if a C++ member function or member function template is
 * declared 'static'.
 */
func (c Cursor) CXXMethod_IsStatic() bool {
	o := C.clang_CXXMethod_isStatic(c.c)
	if o != 0 {
		return true
	}
	return false
}

/**
 * \brief Determine if a C++ member function or member function template is
 * explicitly declared 'virtual' or if it overrides a virtual method from
 * one of the base classes.
 */
func (c Cursor) CXXMethod_IsVirtual() bool {
	o := C.clang_CXXMethod_isVirtual(c.c)
	if o != 0 {
		return true
	}
	return false
}

/**
 * \brief Given a cursor that represents a template, determine
 * the cursor kind of the specializations would be generated by instantiating
 * the template.
 *
 * This routine can be used to determine what flavor of function template,
 * class template, or class template partial specialization is stored in the
 * cursor. For example, it can describe whether a class template cursor is
 * declared with "struct", "class" or "union".
 *
 * \param C The cursor to query. This cursor should represent a template
 * declaration.
 *
 * \returns The cursor kind of the specializations that would be generated
 * by instantiating the template \p C. If \p C is not a template, returns
 * \c CXCursor_NoDeclFound.
 */
func (c Cursor) TemplateCursorKind() CursorKind {
	o := C.clang_getTemplateCursorKind(c.c)
	return CursorKind(o)
}

/**
 * \brief Given a cursor that may represent a specialization or instantiation
 * of a template, retrieve the cursor that represents the template that it
 * specializes or from which it was instantiated.
 *
 * This routine determines the template involved both for explicit
 * specializations of templates and for implicit instantiations of the template,
 * both of which are referred to as "specializations". For a class template
 * specialization (e.g., \c std::vector<bool>), this routine will return
 * either the primary template (\c std::vector) or, if the specialization was
 * instantiated from a class template partial specialization, the class template
 * partial specialization. For a class template partial specialization and a
 * function template specialization (including instantiations), this
 * this routine will return the specialized template.
 *
 * For members of a class template (e.g., member functions, member classes, or
 * static data members), returns the specialized or instantiated member.
 * Although not strictly "templates" in the C++ language, members of class
 * templates have the same notions of specializations and instantiations that
 * templates do, so this routine treats them similarly.
 *
 * \param C A cursor that may be a specialization of a template or a member
 * of a template.
 *
 * \returns If the given cursor is a specialization or instantiation of a
 * template or a member thereof, the template or member that it specializes or
 * from which it was instantiated. Otherwise, returns a NULL cursor.
 */
func (c Cursor) SpecializedCursorTemplate() Cursor {
	o := C.clang_getSpecializedCursorTemplate(c.c)
	return Cursor{o}
}

/**
 * \brief Given a cursor that references something else, return the source range
 * covering that reference.
 *
 * \param C A cursor pointing to a member reference, a declaration reference, or
 * an operator call.
 * \param NameFlags A bitset with three independent flags:
 * CXNameRange_WantQualifier, CXNameRange_WantTemplateArgs, and
 * CXNameRange_WantSinglePiece.
 * \param PieceIndex For contiguous names or when passing the flag
 * CXNameRange_WantSinglePiece, only one piece with index 0 is
 * available. When the CXNameRange_WantSinglePiece flag is not passed for a
 * non-contiguous names, this index can be used to retreive the individual
 * pieces of the name. See also CXNameRange_WantSinglePiece.
 *
 * \returns The piece of the name pointed to by the given cursor. If there is no
 * name, or if the PieceIndex is out-of-range, a null-cursor will be returned.
 */
func (c Cursor) ReferenceNameRange(flags NameRefFlags, pieceIdx uint) SourceRange {
	o := C.clang_getCursorReferenceNameRange(c.c,
		C.uint(flags), C.uint(pieceIdx))
	return SourceRange{o}
}

/**
 * \brief Retrieve a completion string for an arbitrary declaration or macro
 * definition cursor.
 *
 * \param cursor The cursor to query.
 *
 * \returns A non-context-sensitive completion string for declaration and macro
 * definition cursors, or NULL for other kinds of cursors.
 */
func (c Cursor) CompletionString() CompletionString {
	return CompletionString{C.clang_getCursorCompletionString(c.c)}
}

type NameRefFlags uint32

const (
	/**
	 * \brief Include the nested-name-specifier, e.g. Foo:: in x.Foo::y, in the
	 * range.
	 */
	NR_WantQualifier = C.CXNameRange_WantQualifier

	/**
	 * \brief Include the explicit template arguments, e.g. <int> in x.f<int>, in
	 * the range.
	 */
	NR_WantTemplateArgs = C.CXNameRange_WantTemplateArgs

	/**
	 * \brief If the name is non-contiguous, return the full spanning range.
	 *
	 * Non-contiguous names occur in Objective-C when a selector with two or more
	 * parameters is used, or in C++ when using an operator:
	 * \code
	 * [object doSomething:here withValue:there]; // ObjC
	 * return some_vector[1]; // C++
	 * \endcode
	 */
	NR_WantSinglePiece = C.CXNameRange_WantSinglePiece
)

// TODO
//
// /**
//  * \brief Find references of a declaration in a specific file.
//  *
//  * \param cursor pointing to a declaration or a reference of one.
//  *
//  * \param file to search for references.
//  *
//  * \param visitor callback that will receive pairs of CXCursor/CXSourceRange for
//  * each reference found.
//  * The CXSourceRange will point inside the file; if the reference is inside
//  * a macro (and not a macro argument) the CXSourceRange will be invalid.
//  *
//  * \returns one of the CXResult enumerators.
//  */
// func (c Cursor) FindReferencesInFile(f File, visitor CursorAndRangeVisitor) Result {
// 	return Result(C.clang_findReferencesInFile(c.c, f.c, visitor.c))
// }

// TODO
//
// CINDEX_LINKAGE
// CXResult clang_findReferencesInFileWithBlock(CXCursor, CXFile,
//                                              CXCursorAndRangeVisitorBlock);

// EOF
