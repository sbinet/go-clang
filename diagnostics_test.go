package clang_test

import (
	"strings"
	"testing"

	"github.com/sbinet/go-clang"
)

func TestDiagnostics(t *testing.T) {
	idx := clang.NewIndex(0, 0)
	defer idx.Dispose()
	tu := idx.Parse("visitorwrap.c", nil, nil, 0)
	if !tu.IsValid() {
		t.Fatal("TranslationUnit is not valid")
	}
	defer tu.Dispose()

	diags := tu.Diagnostics()
	defer diags.Dispose()
	ok := false
	for _, d := range diags {
		if strings.Contains(d.Spelling(), "_cgo_export.h") {
			ok = true
		}
		t.Log(d)
		t.Log(d.Severity(), d.Spelling())
		t.Log(d.Format(clang.Diagnostic_DisplayCategoryName | clang.Diagnostic_DisplaySourceLocation))
	}
	if !ok {
		t.Errorf("Expected to find a diagnostic regarding _cgo_export.h")
	}
}
