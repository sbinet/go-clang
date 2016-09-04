/* helper functions to visit cursors
 */

#include "_cgo_export.h"
#include "go-clang.h"

unsigned
_go_clang_visit_children(CXCursor c, uintptr_t callback_id)
{
  return clang_visitChildren(c, (CXCursorVisitor)&GoClangCursorVisitor, (CXClientData)callback_id);
}

/* EOF */
