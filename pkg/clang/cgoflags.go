package clang

// #cgo darwin LDFLAGS: -L/opt/local/libexec/llvm-3.1/lib
// #cgo linux LDFLAGS: -L/usr/lib/llvm
// #cgo LDFLAGS: -lclang
// #cgo CFLAGS: -I/opt/local/libexec/llvm-3.1/include
import "C"

//EOF
