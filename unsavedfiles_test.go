package clang_test

import (
	"testing"

	"github.com/sbinet/go-clang"
)

func TestUnsavedFiles(t *testing.T) {
	us := clang.UnsavedFiles{"hello.cpp": `
#include <stdio.h>
int main(int argc, char **argv) {
	printf("Hello world!\n");
	return 0;
}
`}

	idx := clang.NewIndex(0, 0)
	defer idx.Dispose()
	tu := idx.Parse("hello.cpp", nil, us, 0)
	if !tu.IsValid() {
		t.Fatal("TranslationUnit is not valid")
	}
	defer tu.Dispose()

	res := tu.CompleteAt("hello.cpp", 4, 1, us, 0)
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
}
