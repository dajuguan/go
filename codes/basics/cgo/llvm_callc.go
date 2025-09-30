package main

// #include <dlfcn.h>
// static void* load_symbol(const char* lib, const char* sym) {
//     void* handle = dlopen(lib, RTLD_LAZY);
//     return dlsym(handle, sym);
// }
import "C"
import (
	"fmt"

	llvm "tinygo.org/x/go-llvm"
)

func LLVMCallC() {
	llvm.InitializeNativeTarget()
	llvm.InitializeNativeAsmPrinter()

	ctx := llvm.NewContext()
	mod := ctx.NewModule("mymodule")
	builder := ctx.NewBuilder()

	// Declare external squre
	int32Type := llvm.GlobalContext().Int32Type()
	fnType := llvm.FunctionType(int32Type, []llvm.Type{int32Type}, false)
	// shoule use the same name, if we don't link it in runtime with AddGlobalMapping
	fn := llvm.AddFunction(mod, "square", fnType)

	// Declare cmalloc
	fnMallocType := llvm.FunctionType(int32Type, []llvm.Type{}, false)
	fnMalloc := llvm.AddFunction(mod, "cmalloc", fnMallocType)

	// main function
	mainFnType := llvm.FunctionType(int32Type, nil, false)
	mainFn := llvm.AddFunction(mod, "main", mainFnType)
	entry := ctx.AddBasicBlock(mainFn, "entry")
	builder.SetInsertPointAtEnd(entry)

	x := llvm.ConstInt(int32Type, 1, false)
	ret := builder.CreateCall(fnType, fn, []llvm.Value{x}, "")
	fmt.Println("ret:", ret)
	ret = builder.CreateCall(fnMallocType, fnMalloc, []llvm.Value{}, "")
	builder.CreateRet(ret)

	// Execution engine
	engine, _ := llvm.NewExecutionEngine(mod)

	// Load symbol from shared lib (no cgo "import C" trickery, just dlopen/dlsym)
	// sym := C.load_symbol(C.CString("./libmylib.so"), C.CString("square"))
	// engine.AddGlobalMapping(fn, unsafe.Pointer(sym))

	// Only object file or dynamic linked object is needed in llvm.
	// Don't use archive file, cause it can automatically link libc, we don't need to worry about it.
	engine.AddObjectFileByFilename("./mylib.o")

	// Run main
	res := engine.RunFunction(mainFn, nil)
	fmt.Println("Main returned:", res.Int(false))
}
