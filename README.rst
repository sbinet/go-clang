go-clang
========

Naive Go bindings to the C-API of ``CLang``.

Installation
------------

As there is no ``pkg-config`` entry for clang, you may have to tinker
a bit the various CFLAGS and LDFLAGS options, or pass them via the
shell:

::

  $ CGO_CFLAGS="-I/somewhere" \
  CGO_LDFLAGS="-L/somewhere/else" \
  go get bitbucket.org/binet/go-clang/pkg/clang


Example
-------

An example on how to use the AST visitor of ``CLang`` is provided
here:

 

Limitations
-----------

- Only a subset of the C-API of ``CLang`` has been provided yet.
  More will come as patches flow in and time goes by.

- Go-doc documentation is lagging (but the doxygen docs from the C-API
  of ``CLang`` are in the ``.go`` files)
