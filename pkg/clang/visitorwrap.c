/* helper functions to visit cursors
 */

#include "_cgo_export.h"
#include "clang-c/Index.h"
#include "go-clang.h"

unsigned
_go_clang_visit_children(CXCursor c)
{
  return clang_visitChildren(c, (CXCursorVisitor)&GoClangCursorVisitor, NULL);
}

/* EOF */
