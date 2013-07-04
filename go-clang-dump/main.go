// go-clang-dump shows how to dump the AST of a C/C++ file via the Cursor
// visitor API.
//
// ex:
// $ go-clang-dump -fname=foo.cxx
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sbinet/go-clang"
)

var fname *string = flag.String("fname", "", "the file to analyze")

func main() {
	fmt.Printf(":: go-clang-dump...\n")
	flag.Parse()
	fmt.Printf(":: fname: %s\n", *fname)
	fmt.Printf(":: args: %v\n", flag.Args())
	if *fname == "" {
		flag.Usage()
		fmt.Printf("please provide a file name to analyze\n")
		os.Exit(1)
	}
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	nidx := 0
	args := []string{}
	if len(flag.Args()) > 0 && flag.Args()[0] == "-" {
		nidx = 1
		args = make([]string, len(flag.Args()[nidx:]))
		copy(args, flag.Args()[nidx:])
	}

	tu := idx.Parse(*fname, args, nil, 0)

	defer tu.Dispose()

	fmt.Printf("tu: %s\n", tu.Spelling())
	cursor := tu.ToCursor()
	fmt.Printf("cursor-isnull: %v\n", cursor.IsNull())
	fmt.Printf("cursor: %s\n", cursor.Spelling())
	fmt.Printf("cursor-kind: %s\n", cursor.Kind().Spelling())

	tu_fname := tu.File(*fname).Name()
	fmt.Printf("tu-fname: %s\n", tu_fname)

	fct := func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			fmt.Printf("cursor: <none>\n")
			return clang.CVR_Continue
		}
		fmt.Printf("%s: %s (%s)\n",
			cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR())
		switch cursor.Kind() {
		case clang.CK_ClassDecl, clang.CK_EnumDecl,
			clang.CK_StructDecl, clang.CK_Namespace:
			return clang.CVR_Recurse
		}
		return clang.CVR_Continue
	}

	cursor.Visit(fct)

	fmt.Printf(":: bye.\n")
}

// EOF
