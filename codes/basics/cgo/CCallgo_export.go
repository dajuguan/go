package main

import "C"

//export AGoFunction
func AGoFunction(a *C.int) C.int {
	print("AGoFunction() called from C\n")
	*a = 10
	return C.int(11)
}

// Must be main function (could be nil) to make it a C shared library,
// , or the exported function called from C will actually execute nothing.
func main() {}
