package clang

// #cgo darwin LDFLAGS: -L/opt/local/libexec/llvm-3.0/lib
// #cgo linux LDFLAGS: -L/usr/lib/llvm
// #cgo LDFLAGS: -lclang
// #cgo CFLAGS: -I/opt/local/libexec/llvm-3.0/include
import "C"

//EOF
