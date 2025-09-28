package main

/*
#include <stdlib.h>
#include <string.h>

static inline int add(int a, int b) {
	return a + b;
}
static inline void nop() { }


// 单独的分配/释放/写入
static inline void* alloc_c(size_t n) {
    return malloc(n);
}
static inline void free_c(void* p) {
    free(p);
}
static inline void touch_c(void* p, size_t n) {
	for (int i = 0; i < n; i++) {
        memset(p, 42, i);
    }
}

// 一次性完成 alloc+touch+free
static inline void alloc_touch_free_c(size_t n, int repeat) {
    void* p = malloc(n);
	touch_c(p, repeat);

    free(p);
}

static inline void llvm_touch_c(void* p, int num, int n);
// Tell cgo to link this object file
#cgo LDFLAGS: llvm_touch_c.o
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func Random() {
	fmt.Println("Calling C.random from Go:", C.random())
}

func Add(a, b int) {
	// fmt.Println("Calling C.add from Go:", C.add(C.int(a), C.int(b)))
	C.add(C.int(a), C.int(b))
}

func Nop(a, b int) {
	C.nop()
}

func NopGo(a, b int) {
	_ = a + b
}

const (
	N      = 128
	REPEAT = 10
)

func AllocCCombined() {
	C.alloc_touch_free_c(N, REPEAT)
}

func GoSlicePassToC() {
	buf := make([]byte, N)
	C.touch_c(unsafe.Pointer(&buf[0]), REPEAT)
}

func LLVMGoSliceTouchC() {
	buf := make([]byte, N)
	C.llvm_touch_c(unsafe.Pointer(&buf[0]), 42, REPEAT)
}
