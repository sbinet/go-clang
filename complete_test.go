package clang

import (
	//	"fmt"
	"testing"
)

func TestCompleteAt(t *testing.T) {
	idx := NewIndex(0, 1)
	defer idx.Dispose()
	tu := idx.Parse("visitorwrap.c", nil, 0)
	if !tu.IsValid() {
		t.Fatal("TranslationUnit is not valid")
	}
	defer tu.Dispose()

	res := tu.CompleteAt("visitorwrap.c", 11, 16, 0)
	if !res.IsValid() {
		t.Fatal("CompleteResults are not valid")
	}
	defer res.Dispose()
	t.Logf("%+v", res)
	for _, r := range res.Results() {
		t.Logf("%+v", r)
		for _, c := range r.CompletionString.Chunks() {
			t.Logf("\t%+v", c)
		}
	}
}
