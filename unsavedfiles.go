package clang

// #include <stdlib.h>
// #include "clang-c/Index.h"
import "C"

import (
	"unsafe"
)

type (
	UnsavedFiles  map[string]string
	cUnsavedFiles []C.struct_CXUnsavedFile
)

func (us UnsavedFiles) to_c() (ret cUnsavedFiles) {
	ret = make(cUnsavedFiles, 0, len(us))
	for filename, contents := range us {
		ret = append(ret, C.struct_CXUnsavedFile{Filename: C.CString(filename), Contents: C.CString(contents), Length: C.ulong(len(contents))})
	}
	return ret
}

func (us cUnsavedFiles) Dispose() {
	for i := range us {
		C.free(unsafe.Pointer(us[i].Filename))
		C.free(unsafe.Pointer(us[i].Contents))
	}
}

func (us cUnsavedFiles) ptr() *C.struct_CXUnsavedFile {
	if len(us) > 0 {
		return &us[0]
	}
	return nil
}
