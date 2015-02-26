#ifndef _GO_CLANG
#define _GO_CLANG 1

#include "clang-c-index.h"

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
