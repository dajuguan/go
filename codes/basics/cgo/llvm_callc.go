package main

// #include <dlfcn.h>
// static void* load_symbol(const char* lib, const char* sym) {
//     void* handle = dlopen(lib, RTLD_LAZY);
//     return dlsym(handle, sym);
// }
import "C"
import (
	"fmt"
	"unsafe"

	llvm "tinygo.org/x/go-llvm"
)

func LLVMCallC() {
	llvm.InitializeNativeTarget()
	llvm.InitializeNativeAsmPrinter()

	ctx := llvm.NewContext()
	mod := ctx.NewModule("mymodule")
	builder := ctx.NewBuilder()

	// Declare external hello_from_c : void(i32)
	int32Type := llvm.GlobalContext().Int32Type()
	fnType := llvm.FunctionType(int32Type, []llvm.Type{int32Type}, false)
	fn := llvm.AddFunction(mod, "square_from_c", fnType)

	// main function
	mainFnType := llvm.FunctionType(int32Type, nil, false)
	mainFn := llvm.AddFunction(mod, "main", mainFnType)
	entry := ctx.AddBasicBlock(mainFn, "entry")
	builder.SetInsertPointAtEnd(entry)

	x := llvm.ConstInt(int32Type, 2, false)
	ret := builder.CreateCall(fnType, fn, []llvm.Value{x}, "")
	builder.CreateRet(ret)

	// Execution engine
	engine, _ := llvm.NewExecutionEngine(mod)

	// Load symbol from shared lib (no cgo "import C" trickery, just dlopen/dlsym)
	sym := C.load_symbol(C.CString("./libmylib.so"), C.CString("square"))
	engine.AddGlobalMapping(fn, unsafe.Pointer(sym))

	// Run main
	res := engine.RunFunction(mainFn, nil)
	fmt.Println("Main returned:", res.Int(false))
}
