package main

// #include <stdint.h>
// typedef void (*func)(void* p, int num, int n);
// static  void execute(uint64_t f, void* p, int n) { ((func)f)(p, 42, n); }
import "C"

import (
	"unsafe"

	"tinygo.org/x/go-llvm"
)

type NativeEngine struct {
	engine llvm.ExecutionEngine
	ptr    uint64
}

func NewEngine() *NativeEngine {
	ctx := llvm.NewContext()
	module := ctx.NewModule("evm_module")

	engine, _ := llvm.NewMCJITCompiler(module, llvm.MCJITCompilerOptions{})
	return &NativeEngine{
		engine: engine,
	}
}

// must initialize it first!
func init() {
	llvm.InitializeNativeTarget()
	llvm.InitializeNativeAsmPrinter()
}

func (n *NativeEngine) AddObjectFileByFilename(fileName string) {
	llvm.InitializeAllAsmParsers()
	n.engine.AddObjectFileByFilename(fileName)
	ptr := n.engine.GetFunctionAddress("llvm_touch_c")
	n.ptr = ptr
}

func (n *NativeEngine) Execute() {
	buf := make([]byte, N)
	C.execute(C.uint64_t(n.ptr), unsafe.Pointer(&buf[0]), REPEAT)
}
