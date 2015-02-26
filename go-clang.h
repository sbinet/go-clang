#ifndef _GO_CLANG
#define _GO_CLANG 1

/*
 * include our own clang-c/Index.h.
 * It should be exactly the same than the upstream one, except that:
 *  - CXComment
 *  - CXCursor
 *  - CXIdxLoc
 *  - CXSourceLocation
 *  - CXSourceRange
 *  - CXString
 *  - CXTUResourceUsage
 *  - CXToken
 *  - CXType
 * have been modified to hide the 'void *field[x]' fields from the Go GC.
 * Not hiding these fields confuses the Go GC during garbage collection and pointer scanning,
 * making it think the heap/stack has been somehow corrupted.
 */
#include "clang-c/Index.h"

inline static
CXCursor _go_clang_ocursor_at(CXCursor *c, int idx) {
  return c[idx];
}

inline static
CXPlatformAvailability
_goclang_get_platform_availability_at(CXPlatformAvailability* array, int idx) {
  return array[idx];
}

unsigned _go_clang_visit_children(CXCursor c, void *fct);

CXPlatformAvailability
_goclang_get_platform_availability_at(CXPlatformAvailability* array, int idx);

#endif /* !_GO_CLANG */
