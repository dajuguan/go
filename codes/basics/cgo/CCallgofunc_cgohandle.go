package main

import (
	"fmt"
	"runtime/cgo"
)

/*
#include <stdint.h>
extern void go_callback_handle(uintptr_t handle, int arg);
static inline void CallGoFunctionWithHandle(uintptr_t handle) {
    go_callback_handle(handle, 5);
}
*/
import "C"

// Mustn't have space before export.
//
//export go_callback_handle
func go_callback_handle(handle C.uintptr_t, arg C.int) {
	fn := cgo.Handle(handle).Value().(func(C.int))
	fn(arg)
}

func MyCgoCallback(x C.int) {
	fmt.Println("callback with", x)
}

func ExampleCallHandle() {
	h := cgo.NewHandle(MyCgoCallback)
	C.CallGoFunctionWithHandle(C.uintptr_t(h))
	h.Delete() // Clean up the handle after use
}
