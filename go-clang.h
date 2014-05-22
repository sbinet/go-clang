#ifndef _GO_CLANG
#define _GO_CLANG 1

unsigned _go_clang_visit_children(CXCursor c, void *fct);

CXPlatformAvailability
_goclang_get_platform_availability_at(CXPlatformAvailability* array, int idx);

#endif /* !_GO_CLANG */
