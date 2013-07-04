package clang_test

import (
	"testing"

	"github.com/sbinet/go-clang"
)

func TestReparse(t *testing.T) {
	us := clang.UnsavedFiles{"hello.cpp": "int world();"}

	idx := clang.NewIndex(0, 0)
	defer idx.Dispose()
	tu := idx.Parse("hello.cpp", nil, us, 0)
	if !tu.IsValid() {
		t.Fatal("TranslationUnit is not valid")
	}
	defer tu.Dispose()
	ok := false
	tu.ToCursor().Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.Spelling() == "world" {
			ok = true
			return clang.CVR_Break
		}
		return clang.CVR_Continue
	})
	if !ok {
		t.Error("Expected to find 'world', but didn't")
	}
	us["hello.cpp"] = "int world2();"
	tu.Reparse(us, 0)

	ok = false
	tu.ToCursor().Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if s := cursor.Spelling(); s == "world2" {
			ok = true
			return clang.CVR_Break
		} else if s == "world" {
			t.Errorf("'world' should no longer be part of the translationunit, but it still is")
		}
		return clang.CVR_Continue
	})
	if !ok {
		t.Error("Expected to find 'world2', but didn't")
	}
}

// EOF
