package main

// #cgo LDFLAGS: -L. -lmymem
// #include "mymem.h"
import "C"
import (
	"fmt"
	"unsafe"
)

//export AGoFunction
func AGoFunction(a *C.int) C.int {
	print("AGoFunction() called from C\n")
	*a = 10
	return C.int(11)
}

//export GoUseMemory
func GoUseMemory(cmem *C.Memory) {
	// Recover Go slice from Memory pointer
	slice := unsafe.Slice((*byte)(cmem.store.data), cmem.store.len)
	fmt.Println("Go recovered slice:", slice)
}

// Must be main function (could be nil) to make it a C shared library,
// , or the exported function called from C will actually execute nothing.
func main() {}
