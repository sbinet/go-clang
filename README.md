go-clang
========

[![Build Status](https://secure.travis-ci.org/sbinet/go-clang.png)](http://travis-ci.org/sbinet/go-clang)
[![Build Status](https://drone.io/github.com/sbinet/go-clang/status.png)](https://drone.io/github.com/sbinet/go-clang/latest)


Naive Go bindings to the C-API of ``CLang``.

Installation
------------

As there is no ``pkg-config`` entry for clang, you may have to tinker
a bit the various CFLAGS and LDFLAGS options, or pass them via the
shell:

```
$ CGO_CFLAGS="-I`llvm-config --includedir`" \
  CGO_LDFLAGS="-L`llvm-config --libdir`" \
  go get github.com/sbinet/go-clang
```

Example
-------

An example on how to use the AST visitor of ``CLang`` is provided
here:

 https://github.com/sbinet/go-clang/blob/master/go-clang-dump/main.go

``` go
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
```

which can be installed like so:

```
$ go get github.com/sbinet/go-clang/go-clang-dump
```

Limitations
-----------

- Only a subset of the C-API of ``CLang`` has been provided yet.
  More will come as patches flow in and time goes by.

- Go-doc documentation is lagging (but the doxygen docs from the C-API
  of ``CLang`` are in the ``.go`` files)


Documentation
-------------

Is available at ``godoc``:

 http://godoc.org/github.com/sbinet/go-clang

