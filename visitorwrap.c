/* helper functions to visit cursors
 */

#include "_cgo_export.h"
#include "go-clang.h"

unsigned
_go_clang_visit_children(CXCursor c, void *fct)
{
  return clang_visitChildren(c, (CXCursorVisitor)&GoClangCursorVisitor, fct);
}

/* EOF */
