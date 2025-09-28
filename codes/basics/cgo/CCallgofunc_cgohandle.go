package main

import (
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
	// fmt.Println("callback with", x)
	x = x + 1
}

type S struct {
	val    int
	Handle cgo.Handle
}

func (s *S) MyCgoCallback(x C.int) {
	// fmt.Println("callback with", x, "in struct with val", s.val)
	x += 1
}

func ExampleCallHandle() {
	h := cgo.NewHandle(MyCgoCallback)
	C.CallGoFunctionWithHandle(C.uintptr_t(h))
	h.Delete() // Clean up the handle after use
}

func ExampleCallStructHandle() {
	s := &S{val: 42}
	h := cgo.NewHandle(s.MyCgoCallback)
	C.CallGoFunctionWithHandle(C.uintptr_t(h))
	h.Delete() // Clean up the handle after use
}

func NewSWithHandle() *S {
	s := &S{val: 42}
	h := cgo.NewHandle(s.MyCgoCallback)
	s.Handle = h
	return s
}

func (s *S) BenchmarkCallHandle() {
	C.CallGoFunctionWithHandle(C.uintptr_t(s.Handle))
}
