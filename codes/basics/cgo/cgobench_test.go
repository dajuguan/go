package main

import (
	"testing"
)

func BenchmarkGoCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NopGo(0, 1)
	}
}

// every cgo call c cost abot 30~40 ns
func BenchmarkCgoCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Nop(0, 1)
	}
}

func BenchmarkGoSlicePassToCAndSetmem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GoSlicePassToC()
	}
}

func BenchmarkCAllocAndSetmem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AllocCCombined()
	}
}

func BenchmarkGoSlicePassToCAndLLVMSetmemLoop(b *testing.B) {
	// compile llvm IR to obj first
	// llc -filetype=obj mymemset.ll -o mymemset.o
	for i := 0; i < b.N; i++ {
		LLVMGoSliceTouchC()
	}
}

func BenchmarkGoSlicePassToCAndLLVMFnptr(b *testing.B) {
	engine := NewEngine()
	engine.AddObjectFileByFilename("llvm_touch_c.o")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Execute()
	}
}
