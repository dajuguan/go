package main

/*
#include <stdlib.h>
#include <string.h>
#cgo LDFLAGS: -L. -lmymem
#include "mymem.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Memory struct {
	store       []byte
	lastGasCost uint64
}

func SetGoMemInC() {
	// Step 1: allocate Memory in C
	cmem := C.memory_new()
	defer C.memory_free(cmem) // ensure free at end

	goslice := unsafe.Slice((*byte)(cmem.store.data), cmem.store.len)
	fmt.Println("Before size :", len(goslice))
	fmt.Println("Before value:", goslice)

	// Step 2: prepare Go data
	val := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Step 3: copy Go data into C memory
	cval := C.CBytes(val)
	defer C.free(cval)

	// Step 4: write into memory (may trigger realloc)
	C.memory_set(cmem, 0, C.size_t(len(val)), (*C.uint8_t)(cval))

	// Step 5: recover Go slice from C pointer
	goslice = unsafe.Slice((*byte)(cmem.store.data), cmem.store.len)

	fmt.Println("After size :", len(goslice))
	fmt.Println("After value:", goslice)
}
