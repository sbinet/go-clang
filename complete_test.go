package clang_test

import (
	"strings"
	"testing"

	"github.com/sbinet/go-clang"
)

func TestCompleteAt(t *testing.T) {
	idx := clang.NewIndex(0, 0)
	defer idx.Dispose()
	tu := idx.Parse("visitorwrap.c", nil, nil, 0)
	if !tu.IsValid() {
		t.Fatal("TranslationUnit is not valid")
	}
	defer tu.Dispose()

	const lineno = 10 // ie: call to clang_visitChildren
	res := tu.CompleteAt("visitorwrap.c", lineno, 16, nil, 0)
	if !res.IsValid() {
		t.Fatal("CompleteResults are not valid")
	}
	defer res.Dispose()
	if n := len(res.Results()); n < 10 {
		t.Errorf("Expected more results than %d", n)
	}
	t.Logf("%+v", res)
	for _, r := range res.Results() {
		t.Logf("%+v", r)
		for _, c := range r.CompletionString.Chunks() {
			t.Logf("\t%+v", c)
		}
	}

	diags := res.Diagnostics()
	defer diags.Dispose()
	ok := false
	for _, d := range diags {
		if strings.Contains(d.Spelling(), "_cgo_export.h") {
			ok = true
		}
		t.Log(d.Severity(), d.Spelling())
	}
	if !ok {
		t.Errorf("Expected to find a diagnostic regarding _cgo_export.h")
	}
}
