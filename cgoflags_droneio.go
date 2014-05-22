//+build droneio

package clang

// #cgo linux CFLAGS: -I/usr/lib/llvm-3.4/include
// #cgo linux LDFLAGS: -L/usr/lib/llvm-3.4/lib
import "C"
