package clang

// #include <stdlib.h>
// #include "clang-c/Index.h"
import "C"

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
	// A C or C++ union
	CK_UnionDecl = C.CXCursor_UnionDecl
	// A C++ class
	CK_ClassDecl = C.CXCursor_ClassDecl
	// An enumeration
	CK_EnumDecl = C.CXCursor_EnumDecl
	// A field (in C) or non-static data member (in C++) in a
	// struct, union, or C++ class.
	CK_FieldDecl = C.CXCursor_FieldDecl
	/** \brief An enumerator constant. */
	CK_EnumConstantDecl = C.CXCursor_EnumConstantDecl
	/** \brief A function. */
	CK_FunctionDecl = C.CXCursor_FunctionDecl
	/** \brief A variable. */
	CK_VarDecl = C.CXCursor_VarDecl
	/** \brief A function or method parameter. */
	CK_ParmDecl = C.CXCursor_ParmDecl
	/** \brief An Objective-C @interface. */
	CK_ObjCInterfaceDecl = C.CXCursor_ObjCInterfaceDecl
	/** \brief An Objective-C @interface for a category. */
	CK_ObjCCategoryDecl = C.CXCursor_ObjCCategoryDecl
	/** \brief An Objective-C @protocol declaration. */
	CK_ObjCProtocolDecl = C.CXCursor_ObjCProtocolDecl
	/** \brief An Objective-C @property declaration. */
	CK_ObjCPropertyDecl = C.CXCursor_ObjCPropertyDecl
	/** \brief An Objective-C instance variable. */
	CK_ObjCIvarDecl = C.CXCursor_ObjCIvarDecl
	/** \brief An Objective-C instance method. */
	CK_ObjCInstanceMethodDecl = C.CXCursor_ObjCInstanceMethodDecl
	/** \brief An Objective-C class method. */
	CK_ObjCClassMethodDecl = C.CXCursor_ObjCClassMethodDecl
	/** \brief An Objective-C @implementation. */
	CK_ObjCImplementationDecl = C.CXCursor_ObjCImplementationDecl
	/** \brief An Objective-C @implementation for a category. */
	CK_ObjCCategoryImplDecl = C.CXCursor_ObjCCategoryImplDecl
	/** \brief A typedef */
	CK_TypedefDecl = C.CXCursor_TypedefDecl
	/** \brief A C++ class method. */
	CK_CXXMethod = C.CXCursor_CXXMethod
	/** \brief A C++ namespace. */
	CK_Namespace = C.CXCursor_Namespace
	/** \brief A linkage specification, e.g. 'extern "C"'. */
	CK_LinkageSpec = C.CXCursor_LinkageSpec
	/** \brief A C++ constructor. */
	CK_Constructor = C.CXCursor_Constructor
	/** \brief A C++ destructor. */
	CK_Destructor = C.CXCursor_Destructor
	/** \brief A C++ conversion function. */
	CK_ConversionFunction = C.CXCursor_ConversionFunction
	/** \brief A C++ template type parameter. */
	CK_TemplateTypeParameter = C.CXCursor_TemplateTypeParameter
	/** \brief A C++ non-type template parameter. */
	CK_NonTypeTemplateParameter = C.CXCursor_NonTypeTemplateParameter
	/** \brief A C++ template template parameter. */
	CK_TemplateTemplateParameter = C.CXCursor_TemplateTemplateParameter
	/** \brief A C++ function template. */
	CK_FunctionTemplate = C.CXCursor_FunctionTemplate
	/** \brief A C++ class template. */
	CK_ClassTemplate = C.CXCursor_ClassTemplate
	/** \brief A C++ class template partial specialization. */
	CK_ClassTemplatePartialSpecialization = C.CXCursor_ClassTemplatePartialSpecialization
	/** \brief A C++ namespace alias declaration. */
	CK_NamespaceAlias = C.CXCursor_NamespaceAlias
	/** \brief A C++ using directive. */
	CK_UsingDirective = C.CXCursor_UsingDirective
	/** \brief A C++ using declaration. */
	CK_UsingDeclaration = C.CXCursor_UsingDeclaration
	/** \brief A C++ alias declaration */
	CK_TypeAliasDecl = C.CXCursor_TypeAliasDecl
	/** \brief An Objective-C @synthesize definition. */
	CK_ObjCSynthesizeDecl = C.CXCursor_ObjCSynthesizeDecl
	/** \brief An Objective-C @dynamic definition. */
	CK_ObjCDynamicDecl = C.CXCursor_ObjCDynamicDecl
	/** \brief An access specifier. */
	CK_CXXAccessSpecifier = C.CXCursor_CXXAccessSpecifier

	CK_FirstDecl = C.CXCursor_FirstDecl
	CK_LastDecl  = C.CXCursor_LastDecl

	/* References */
	CK_FirstRef          = C.CXCursor_FirstRef
	CK_ObjCSuperClassRef = C.CXCursor_ObjCSuperClassRef
	CK_ObjCProtocolRef   = C.CXCursor_ObjCProtocolRef
	CK_ObjCClassRef      = C.CXCursor_ObjCClassRef
	/**
	 * \brief A reference to a type declaration.
	 *
	 * A type reference occurs anywhere where a type is named but not
	 * declared. For example, given:
	 *
	 * \code
	 * typedef unsigned size_type;
	 * size_type size;
	 * \endcode
	 *
	 * The typedef is a declaration of size_type (CXCursor_TypedefDecl),
	 * while the type of the variable "size" is referenced. The cursor
	 * referenced by the type of size is the typedef for size_type.
	 */
	CK_TypeRef          = C.CXCursor_TypeRef
	CK_CXXBaseSpecifier = C.CXCursor_CXXBaseSpecifier
	/**
	 * \brief A reference to a class template, function template, template
	 * template parameter, or class template partial specialization.
	 */
	CK_TemplateRef = C.CXCursor_TemplateRef
	/**
	 * \brief A reference to a namespace or namespace alias.
	 */
	CK_NamespaceRef = C.CXCursor_NamespaceRef
	/**
	 * \brief A reference to a member of a struct, union, or class that occurs in
	 * some non-expression context, e.g., a designated initializer.
	 */
	CK_MemberRef = C.CXCursor_MemberRef
	/**
	 * \brief A reference to a labeled statement.
	 *
	 * This cursor kind is used to describe the jump to "start_over" in the
	 * goto statement in the following example:
	 *
	 * \code
	 *   start_over:
	 *     ++counter;
	 *
	 *     goto start_over;
	 * \endcode
	 *
	 * A label reference cursor refers to a label statement.
	 */
	CK_LabelRef = C.CXCursor_LabelRef

	/**
	 * \brief A reference to a set of overloaded functions or function templates
	 * that has not yet been resolved to a specific function or function template.
	 *
	 * An overloaded declaration reference cursor occurs in C++ templates where
	 * a dependent name refers to a function. For example:
	 *
	 * \code
	 * template<typename T> void swap(T&, T&);
	 *
	 * struct X { ... };
	 * void swap(X&, X&);
	 *
	 * template<typename T>
	 * void reverse(T* first, T* last) {
	 *   while (first < last - 1) {
	 *     swap(*first, *--last);
	 *     ++first;
	 *   }
	 * }
	 *
	 * struct Y { };
	 * void swap(Y&, Y&);
	 * \endcode
	 *
	 * Here, the identifier "swap" is associated with an overloaded declaration
	 * reference. In the template definition, "swap" refers to either of the two
	 * "swap" functions declared above, so both results will be available. At
	 * instantiation time, "swap" may also refer to other functions found via
	 * argument-dependent lookup (e.g., the "swap" function at the end of the
	 * example).
	 *
	 * The functions \c clang_getNumOverloadedDecls() and
	 * \c clang_getOverloadedDecl() can be used to retrieve the definitions
	 * referenced by this cursor.
	 */
	CK_OverloadedDeclRef = C.CXCursor_OverloadedDeclRef

	CK_LastRef = C.CXCursor_LastRef

	/* Error conditions */
	CK_FirstInvalid   = C.CXCursor_FirstInvalid
	CK_InvalidFile    = C.CXCursor_InvalidFile
	CK_NoDeclFound    = C.CXCursor_NoDeclFound
	CK_NotImplemented = C.CXCursor_NotImplemented
	CK_InvalidCode    = C.CXCursor_InvalidCode
	CK_LastInvalid    = C.CXCursor_LastInvalid

	/* Expressions */
	CK_FirstExpr = C.CXCursor_FirstExpr

	/**
	 * \brief An expression whose specific kind is not exposed via this
	 * interface.
	 *
	 * Unexposed expressions have the same operations as any other kind
	 * of expression; one can extract their location information,
	 * spelling, children, etc. However, the specific kind of the
	 * expression is not reported.
	 */
	CK_UnexposedExpr = C.CXCursor_UnexposedExpr

	/**
	 * \brief An expression that refers to some value declaration, such
	 * as a function, varible, or enumerator.
	 */
	CK_DeclRefExpr = C.CXCursor_DeclRefExpr

	/**
	 * \brief An expression that refers to a member of a struct, union,
	 * class, Objective-C class, etc.
	 */
	CK_MemberRefExpr = C.CXCursor_MemberRefExpr

	/** \brief An expression that calls a function. */
	CK_CallExpr = C.CXCursor_CallExpr

	/** \brief An expression that sends a message to an Objective-C
	  object or class. */
	CK_ObjCMessageExpr = C.CXCursor_ObjCMessageExpr

	/** \brief An expression that represents a block literal. */
	CK_BlockExpr = C.CXCursor_BlockExpr

	/** \brief An integer literal.
	 */
	CK_IntegerLiteral = C.CXCursor_IntegerLiteral

	/** \brief A floating point number literal.
	 */
	CK_FloatingLiteral = C.CXCursor_FloatingLiteral

	/** \brief An imaginary number literal.
	 */
	CK_ImaginaryLiteral = C.CXCursor_ImaginaryLiteral

	/** \brief A string literal.
	 */
	CK_StringLiteral = C.CXCursor_StringLiteral

	/** \brief A character literal.
	 */
	CK_CharacterLiteral = C.CXCursor_CharacterLiteral

	/** \brief A parenthesized expression, e.g. "(1)".
	 *
	 * This AST node is only formed if full location information is requested.
	 */
	CK_ParenExpr = C.CXCursor_ParenExpr

	/** \brief This represents the unary-expression's (except sizeof and
	 * alignof).
	 */
	CK_UnaryOperator = C.CXCursor_UnaryOperator

	/** \brief [C99 6.5.2.1] Array Subscripting.
	 */
	CK_ArraySubscriptExpr = C.CXCursor_ArraySubscriptExpr

	/** \brief A builtin binary operation expression such as "x + y" or
	 * "x <= y".
	 */
	CK_BinaryOperator = C.CXCursor_BinaryOperator

	/** \brief Compound assignment such as "+=".
	 */
	CK_CompoundAssignOperator = C.CXCursor_CompoundAssignOperator

	/** \brief The ?: ternary operator.
	 */
	CK_ConditionalOperator = C.CXCursor_ConditionalOperator

	/** \brief An explicit cast in C (C99 6.5.4) or a C-style cast in C++
	 * (C++ [expr.cast]), which uses the syntax (Type)expr.
	 *
	 * For example: (int)f.
	 */
	CK_CStyleCastExpr = C.CXCursor_CStyleCastExpr

	/** \brief [C99 6.5.2.5]
	 */
	CK_CompoundLiteralExpr = C.CXCursor_CompoundLiteralExpr

	/** \brief Describes an C or C++ initializer list.
	 */
	CK_InitListExpr = C.CXCursor_InitListExpr

	/** \brief The GNU address of label extension, representing &&label.
	 */
	CK_AddrLabelExpr = C.CXCursor_AddrLabelExpr

	/** \brief This is the GNU Statement Expression extension: ({int X=4; X;})
	 */
	CK_StmtExpr = C.CXCursor_StmtExpr

	/** \brief Represents a C1X generic selection.
	 */
	CK_GenericSelectionExpr = C.CXCursor_GenericSelectionExpr

	/** \brief Implements the GNU __null extension, which is a name for a null
	 * pointer constant that has integral type (e.g., int or long) and is the same
	 * size and alignment as a pointer.
	 *
	 * The __null extension is typically only used by system headers, which define
	 * NULL as __null in C++ rather than using 0 (which is an integer that may not
	 * match the size of a pointer).
	 */
	CK_GNUNullExpr = C.CXCursor_GNUNullExpr

	/** \brief C++'s static_cast<> expression.
	 */
	CK_CXXStaticCastExpr = C.CXCursor_CXXStaticCastExpr

	/** \brief C++'s dynamic_cast<> expression.
	 */
	CK_CXXDynamicCastExpr = C.CXCursor_CXXDynamicCastExpr

	/** \brief C++'s reinterpret_cast<> expression.
	 */
	CK_CXXReinterpretCastExpr = C.CXCursor_CXXReinterpretCastExpr

	/** \brief C++'s const_cast<> expression.
	 */
	CK_CXXConstCastExpr = C.CXCursor_CXXConstCastExpr

	/** \brief Represents an explicit C++ type conversion that uses "functional"
	 * notion (C++ [expr.type.conv]).
	 *
	 * Example:
	 * \code
	 *   x = int(0.5);
	 * \endcode
	 */
	CK_CXXFunctionalCastExpr = C.CXCursor_CXXFunctionalCastExpr

	/** \brief A C++ typeid expression (C++ [expr.typeid]).
	 */
	CK_CXXTypeidExpr = C.CXCursor_CXXTypeidExpr

	/** \brief [C++ 2.13.5] C++ Boolean Literal.
	 */
	CK_CXXBoolLiteralExpr = C.CXCursor_CXXBoolLiteralExpr

	/** \brief [C++0x 2.14.7] C++ Pointer Literal.
	 */
	CK_CXXNullPtrLiteralExpr = C.CXCursor_CXXNullPtrLiteralExpr

	/** \brief Represents the "this" expression in C++
	 */
	CK_CXXThisExpr = C.CXCursor_CXXThisExpr

	/** \brief [C++ 15] C++ Throw Expression.
	 *
	 * This handles 'throw' and 'throw' assignment-expression. When
	 * assignment-expression isn't present, Op will be null.
	 */
	CK_CXXThrowExpr = C.CXCursor_CXXThrowExpr

	/** \brief A new expression for memory allocation and constructor calls, e.g:
	 * "new CXXNewExpr(foo)".
	 */
	CK_CXXNewExpr = C.CXCursor_CXXNewExpr

	/** \brief A delete expression for memory deallocation and destructor calls,
	 * e.g. "delete[] pArray".
	 */
	CK_CXXDeleteExpr = C.CXCursor_CXXDeleteExpr

	/** \brief A unary expression.
	 */
	CK_UnaryExpr = C.CXCursor_UnaryExpr

	/** \brief ObjCStringLiteral, used for Objective-C string literals i.e. "foo".
	 */
	CK_ObjCStringLiteral = C.CXCursor_ObjCStringLiteral

	/** \brief ObjCEncodeExpr, used for in Objective-C.
	 */
	CK_ObjCEncodeExpr = C.CXCursor_ObjCEncodeExpr

	/** \brief ObjCSelectorExpr used for in Objective-C.
	 */
	CK_ObjCSelectorExpr = C.CXCursor_ObjCSelectorExpr

	/** \brief Objective-C's protocol expression.
	 */
	CK_ObjCProtocolExpr = C.CXCursor_ObjCProtocolExpr

	/** \brief An Objective-C "bridged" cast expression, which casts between
	 * Objective-C pointers and C pointers, transferring ownership in the process.
	 *
	 * \code
	 *   NSString *str = (__bridge_transfer NSString *)CFCreateString();
	 * \endcode
	 */
	CK_ObjCBridgedCastExpr = C.CXCursor_ObjCBridgedCastExpr

	/** \brief Represents a C++0x pack expansion that produces a sequence of
	 * expressions.
	 *
	 * A pack expansion expression contains a pattern (which itself is an
	 * expression) followed by an ellipsis. For example:
	 *
	 * \code
	 * template<typename F, typename ...Types>
	 * void forward(F f, Types &&...args) {
	 *  f(static_cast<Types&&>(args)...);
	 * }
	 * \endcode
	 */
	CK_PackExpansionExpr = C.CXCursor_PackExpansionExpr

	/** \brief Represents an expression that computes the length of a parameter
	 * pack.
	 *
	 * \code
	 * template<typename ...Types>
	 * struct count {
	 *   static const unsigned value = sizeof...(Types);
	 * };
	 * \endcode
	 */
	CK_SizeOfPackExpr = C.CXCursor_SizeOfPackExpr

	/** \brief Represents the "self" expression in a ObjC method.
	 */
	CK_ObjCSelfExpr = C.CXCursor_ObjCSelfExpr

	CK_LastExpr = C.CXCursor_LastExpr

	/* Statements */
	CK_FirstStmt = C.CXCursor_FirstStmt
	/**
	 * \brief A statement whose specific kind is not exposed via this
	 * interface.
	 *
	 * Unexposed statements have the same operations as any other kind of
	 * statement; one can extract their location information, spelling,
	 * children, etc. However, the specific kind of the statement is not
	 * reported.
	 */
	CK_UnexposedStmt = C.CXCursor_UnexposedStmt

	/** \brief A labelled statement in a function.
	 *
	 * This cursor kind is used to describe the "start_over:" label statement in
	 * the following example:
	 *
	 * \code
	 *   start_over:
	 *     ++counter;
	 * \endcode
	 *
	 */
	CK_LabelStmt = C.CXCursor_LabelStmt

	/** \brief A group of statements like { stmt stmt }.
	 *
	 * This cursor kind is used to describe compound statements, e.g. function
	 * bodies.
	 */
	CK_CompoundStmt = C.CXCursor_CompoundStmt

	/** \brief A case statment.
	 */
	CK_CaseStmt = C.CXCursor_CaseStmt

	/** \brief A default statement.
	 */
	CK_DefaultStmt = C.CXCursor_DefaultStmt

	/** \brief An if statement
	 */
	CK_IfStmt = C.CXCursor_IfStmt

	/** \brief A switch statement.
	 */
	CK_SwitchStmt = C.CXCursor_SwitchStmt

	/** \brief A while statement.
	 */
	CK_WhileStmt = C.CXCursor_WhileStmt

	/** \brief A do statement.
	 */
	CK_DoStmt = C.CXCursor_DoStmt

	/** \brief A for statement.
	 */
	CK_ForStmt = C.CXCursor_ForStmt

	/** \brief A goto statement.
	 */
	CK_GotoStmt = C.CXCursor_GotoStmt

	/** \brief An indirect goto statement.
	 */
	CK_IndirectGotoStmt = C.CXCursor_IndirectGotoStmt

	/** \brief A continue statement.
	 */
	CK_ContinueStmt = C.CXCursor_ContinueStmt

	/** \brief A break statement.
	 */
	CK_BreakStmt = C.CXCursor_BreakStmt

	/** \brief A return statement.
	 */
	CK_ReturnStmt = C.CXCursor_ReturnStmt

	/** \brief A GCC inline assembly statement extension.
	 */
	CK_GCCAsmStmt = C.CXCursor_GCCAsmStmt
	CK_AsmStmt    = CK_GCCAsmStmt

	/** \brief Objective-C's overall @try-@catc-@finall statement.
	 */
	CK_ObjCAtTryStmt = C.CXCursor_ObjCAtTryStmt

	/** \brief Objective-C's @catch statement.
	 */
	CK_ObjCAtCatchStmt = C.CXCursor_ObjCAtCatchStmt

	/** \brief Objective-C's @finally statement.
	 */
	CK_ObjCAtFinallyStmt = C.CXCursor_ObjCAtFinallyStmt

	/** \brief Objective-C's @throw statement.
	 */
	CK_ObjCAtThrowStmt = C.CXCursor_ObjCAtThrowStmt

	/** \brief Objective-C's @synchronized statement.
	 */
	CK_ObjCAtSynchronizedStmt = C.CXCursor_ObjCAtSynchronizedStmt

	/** \brief Objective-C's autorelease pool statement.
	 */
	CK_ObjCAutoreleasePoolStmt = C.CXCursor_ObjCAutoreleasePoolStmt

	/** \brief Objective-C's collection statement.
	 */
	CK_ObjCForCollectionStmt = C.CXCursor_ObjCForCollectionStmt

	/** \brief C++'s catch statement.
	 */
	CK_CXXCatchStmt = C.CXCursor_CXXCatchStmt

	/** \brief C++'s try statement.
	 */
	CK_CXXTryStmt = C.CXCursor_CXXTryStmt

	/** \brief C++'s for (* : *) statement.
	 */
	CK_CXXForRangeStmt = C.CXCursor_CXXForRangeStmt

	/** \brief Windows Structured Exception Handling's try statement.
	 */
	CK_SEHTryStmt = C.CXCursor_SEHTryStmt

	/** \brief Windows Structured Exception Handling's except statement.
	 */
	CK_SEHExceptStmt = C.CXCursor_SEHExceptStmt

	/** \brief Windows Structured Exception Handling's finally statement.
	 */
	CK_SEHFinallyStmt = C.CXCursor_SEHFinallyStmt

	/** \brief A MS inline assembly statement extension.
	 */
	CK_MSAsmStmt = C.CXCursor_MSAsmStmt

	/** \brief The null satement ";": C99 6.8.3p3.
	 *
	 * This cursor kind is used to describe the null statement.
	 */
	CK_NullStmt = C.CXCursor_NullStmt

	/** \brief Adaptor class for mixing declarations with statements and
	 * expressions.
	 */
	CK_DeclStmt = C.CXCursor_DeclStmt

	/** \brief OpenMP parallel directive.
	 */
	CK_OMPParallelDirective = C.CXCursor_OMPParallelDirective

	CK_LastStmt = C.CXCursor_LastStmt

	/**
	 * \brief Cursor that represents the translation unit itself.
	 *
	 * The translation unit cursor exists primarily to act as the root
	 * cursor for traversing the contents of a translation unit.
	 */
	CK_TranslationUnit = C.CXCursor_TranslationUnit

	/* Attributes */
	CK_FirstAttr = C.CXCursor_FirstAttr
	/**
	 * \brief An attribute whose specific kind is not exposed via this
	 * interface.
	 */
	CK_UnexposedAttr = C.CXCursor_UnexposedAttr

	CK_IBActionAttr           = C.CXCursor_IBActionAttr
	CK_IBOutletAttr           = C.CXCursor_IBOutletAttr
	CK_IBOutletCollectionAttr = C.CXCursor_IBOutletCollectionAttr
	CK_CXXFinalAttr           = C.CXCursor_CXXFinalAttr
	CK_CXXOverrideAttr        = C.CXCursor_CXXOverrideAttr
	CK_AnnotateAttr           = C.CXCursor_AnnotateAttr
	CK_LastAttr               = C.CXCursor_LastAttr

	/* Preprocessing */
	CK_PreprocessingDirective = C.CXCursor_PreprocessingDirective
	CK_MacroDefinition        = C.CXCursor_MacroDefinition
	CK_MacroExpansion         = C.CXCursor_MacroExpansion
	CK_MacroInstantiation     = C.CXCursor_MacroInstantiation
	CK_InclusionDirective     = C.CXCursor_InclusionDirective
	CK_FirstPreprocessing     = C.CXCursor_FirstPreprocessing
	CK_LastPreprocessing      = C.CXCursor_LastPreprocessing

	/* Extra Declarations */
	/**
	 * \brief A module import declaration.
	 */
	CK_ModuleImportDecl = C.CXCursor_ModuleImportDecl
	CK_FirstExtraDecl   = C.CXCursor_FirstExtraDecl
	CK_LastExtraDecl    = C.CXCursor_LastExtraDecl
)

func (c CursorKind) to_c() uint32 {
	return uint32(c)
}

func (c CursorKind) String() string {
	return c.Spelling()
}

/* for debug/testing */
func (c CursorKind) Spelling() string {
	cstr := cxstring{C.clang_getCursorKindSpelling(c.to_c())}
	defer cstr.Dispose()
	return cstr.String()
}

// EOF
