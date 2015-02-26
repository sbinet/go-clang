package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

import (
	"time"
)

// A particular source file that is part of a translation unit.
type File struct {
	c C.CXFile
}

// Name retrieves the complete file and path name of the given file.
func (c File) Name() string {
	cstr := cxstring{C.clang_getFileName(c.c)}
	defer cstr.Dispose()
	return cstr.String()
}

// ModTime retrieves the last modification time of the given file.
func (c File) ModTime() time.Time {
	// time_t is in seconds since epoch
	sec := C.clang_getFileTime(c.c)
	const nsec = 0
	return time.Unix(int64(sec), nsec)
}
