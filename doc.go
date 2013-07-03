// clang provides naive bindings to the CLang C-API.
//
// typical usage follows:
//
//  const excludeDeclarationsFromPCH = 1
//  const displayDiagnostics = 1
//  idx := clang.NewIndex(excludeDeclarationsFromPCH, displayDiagnostics)
//  defer idx.Dispose()
//
//  args := []string{}
//  tu := idx.Parse("somefile.cxx", args)
//  defer tu.Dispose()
//  fmt.Printf("translation unit: %s\n", tu.Spelling())
//
package clang

// EOF
