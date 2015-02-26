package clang

// #include <stdlib.h>
// #include "go-clang.h"
// #include "clang-c/CXCompilationDatabase.h"
//
import "C"

import (
	"fmt"
	"unsafe"
)

/**
 * A compilation database holds all information used to compile files in a
 * project. For each file in the database, it can be queried for the working
 * directory or the command line used for the compiler invocation.
 *
 * Must be freed by \c clang_CompilationDatabase_dispose
 */
type CompilationDatabase struct {
	c C.CXCompilationDatabase
}

/**
 * \brief Contains the results of a search in the compilation database
 *
 * When searching for the compile command for a file, the compilation db can
 * return several commands, as the file may have been compiled with
 * different options in different places of the project. This choice of compile
 * commands is wrapped in this opaque data structure. It must be freed by
 * \c clang_CompileCommands_dispose.
 */
type CompileCommands struct {
	c C.CXCompileCommands
}

/**
 * \brief Represents the command line invocation to compile a specific file.
 */
type CompileCommand struct {
	c C.CXCompileCommand
}

/**
 * \brief Error codes for Compilation Database
 */
type CompilationDatabaseError int

func (err CompilationDatabaseError) Error() string {
	switch err {
	case CompilationDatabase_CanNotLoadDatabase:
		return "go-clang: can not load database"
	default:
		return fmt.Sprintf("go-clang: unknown compilationdatabase error (%d)", int(err))
	}
}

const (
	/*
	 * \brief No error occurred
	 */
	CompilationDatabase_NoError CompilationDatabaseError = C.CXCompilationDatabase_NoError

	/*
	 * \brief Database can not be loaded
	 */
	CompilationDatabase_CanNotLoadDatabase = C.CXCompilationDatabase_CanNotLoadDatabase
)

/**
 * \brief Creates a compilation database from the database found in directory
 * buildDir. For example, CMake can output a compile_commands.json which can
 * be used to build the database.
 *
 * It must be freed by \c clang_CompilationDatabase_dispose.
 */
func NewCompilationDatabase(builddir string) (CompilationDatabase, error) {
	var db CompilationDatabase

	c_dir := C.CString(builddir)
	defer C.free(unsafe.Pointer(c_dir))
	var c_err C.CXCompilationDatabase_Error
	c_db := C.clang_CompilationDatabase_fromDirectory(c_dir, &c_err)
	if c_err == C.CXCompilationDatabase_NoError {
		return CompilationDatabase{c_db}, nil
	}
	return db, CompilationDatabaseError(c_err)
}

/**
 * \brief Free the given compilation database
 */
func (db *CompilationDatabase) Dispose() {
	C.clang_CompilationDatabase_dispose(db.c)
}

/**
 * \brief Find the compile commands used for a file. The compile commands
 * must be freed by \c clang_CompileCommands_dispose.
 */
func (db *CompilationDatabase) GetCompileCommands(fname string) CompileCommands {
	c_fname := C.CString(fname)
	defer C.free(unsafe.Pointer(c_fname))
	c_cmds := C.clang_CompilationDatabase_getCompileCommands(db.c, c_fname)
	return CompileCommands{c_cmds}
}

/**
 * \brief Get all the compile commands in the given compilation database.
 */
func (db *CompilationDatabase) GetAllCompileCommands() CompileCommands {
	c_cmds := C.clang_CompilationDatabase_getAllCompileCommands(db.c)
	return CompileCommands{c_cmds}
}

/**
 * \brief Get the number of CompileCommand we have for a file
 */
func (cmds CompileCommands) GetSize() int {
	return int(C.clang_CompileCommands_getSize(cmds.c))
}

/**
 * \brief Get the I'th CompileCommand for a file
 *
 * Note : 0 <= i < clang_CompileCommands_getSize(CXCompileCommands)
 */
func (cmds CompileCommands) GetCommand(idx int) CompileCommand {
	c_cmd := C.clang_CompileCommands_getCommand(cmds.c, C.unsigned(idx))
	return CompileCommand{c_cmd}
}

/**
 * \brief Get the working directory where the CompileCommand was executed from
 */
func (cmd CompileCommand) GetDirectory() string {
	c_str := cxstring{C.clang_CompileCommand_getDirectory(cmd.c)}
	defer c_str.Dispose()
	return c_str.String()
}

/**
 * \brief Get the number of arguments in the compiler invocation.
 *
 */
func (cmd CompileCommand) GetNumArgs() int {
	return int(C.clang_CompileCommand_getNumArgs(cmd.c))
}

/**
 * \brief Get the I'th argument value in the compiler invocations
 *
 * Invariant :
 *  - argument 0 is the compiler executable
 */
func (cmd CompileCommand) GetArg(idx int) string {
	c_str := cxstring{C.clang_CompileCommand_getArg(cmd.c, C.unsigned(idx))}
	defer c_str.Dispose()
	return c_str.String()
}

// /**
//  * \brief Get the number of source mappings for the compiler invocation.
//  */
// func (cmd CompileCommand) GetNumMappedSources() int {
// 	return int(C.clang_CompileCommand_getNumMappedSources(cmd.c))
// }

// /**
//  * \brief Get the I'th mapped source path for the compiler invocation.
//  */
// func (cmd CompileCommand) GetMappedSourcePath(idx int) string {
// 	c_str := cxstring{C.clang_CompileCommand_getMappedSourcePath(cmd.c, C.unsigned(idx))}
// 	defer c_str.Dispose()
// 	return c_str.String()
// }
