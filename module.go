package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

// Module describes a C++ Module
type Module struct {
	c C.CXModule
}

/**
 * \param Module a module object.
 *
 * \returns the module file where the provided module object came from.
 */
func (m Module) ASTFile() File {
	return File{C.clang_Module_getASTFile(m.c)}
}

/**
 * \param Module a module object.
 *
 * \returns the parent of a sub-module or NULL if the given module is top-level,
 * e.g. for 'std.vector' it will return the 'std' module.
 */
func (m Module) Parent() Module {
	return Module{C.clang_Module_getParent(m.c)}
}

/**
 * \param Module a module object.
 *
 * \returns the name of the module, e.g. for the 'std.vector' sub-module it
 * will return "vector".
 */
func (m Module) Name() string {
	o := cxstring{C.clang_Module_getName(m.c)}
	defer o.Dispose()
	return o.String()
}

/**
 * \param Module a module object.
 *
 * \returns the full name of the module, e.g. "std.vector".
 */
func (m Module) FullName() string {
	o := cxstring{C.clang_Module_getFullName(m.c)}
	defer o.Dispose()
	return o.String()
}
