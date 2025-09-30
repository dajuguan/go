package main

// #include <dlfcn.h>
// #include <stdlib.h>
// #include <string.h>
// #include <stdio.h>
// static void* load_symbol(const char* lib, const char* sym) {
//     void* handle = dlopen(lib, RTLD_LAZY);
//     return dlsym(handle, sym);
// }
// #include "mymem.h"
// extern void GoUseMemory(Memory* m);
import "C"
import (
	"fmt"
	"unsafe"

	llvm "tinygo.org/x/go-llvm"
)

func LLVMCallCMemThenPassMemToGo() {
	llvm.InitializeNativeTarget()
	llvm.InitializeNativeAsmPrinter()

	ctx := llvm.NewContext()
	mod := ctx.NewModule("mymodule")
	builder := ctx.NewBuilder()

	// Declare external hello_from_c : void(i32)
	int8Type := llvm.GlobalContext().Int8Type()
	int64Type := llvm.GlobalContext().Int64Type()
	// Build LLVM type for Memory struct
	sliceHeaderType := llvm.StructType([]llvm.Type{
		llvm.PointerType(int8Type, 0), // data
		int64Type,                     // len
		int64Type,                     // cap
	}, false)
	memoryStructType := llvm.StructType([]llvm.Type{
		sliceHeaderType, // store
		int64Type,       // lastGasCost
	}, false)

	//--------------------------------------------------------//
	// memory_new() type
	memNewType := llvm.FunctionType(
		llvm.PointerType(memoryStructType, 0), // returns Memory*
		nil,                                   // no args
		false,
	)
	memoryNew := llvm.AddFunction(mod, "memory_new", memNewType)

	//--------------------------------------------------------//
	memorySetType := llvm.FunctionType(
		llvm.GlobalContext().VoidType(), // returns void
		[]llvm.Type{
			llvm.PointerType(memoryStructType, 0), // Memory*
			int64Type,                             // offset
			int64Type,                             // size
			llvm.PointerType(int8Type, 0),         // const uint8_t*
		},
		false,
	)

	memorySet := llvm.AddFunction(mod, "memory_set", memorySetType)

	// GoUseMem--------------------------------------------------------//
	// Function type: void GoUseMemory(Memory*)
	goUseMemType := llvm.FunctionType(
		llvm.GlobalContext().VoidType(),                    // returns void
		[]llvm.Type{llvm.PointerType(memoryStructType, 0)}, // Memory* argument
		false,
	)
	goUseMemory := llvm.AddFunction(mod, "GoUseMemory", goUseMemType)
	goUseMemory.SetLinkage(llvm.ExternalLinkage)

	//  main function--------------------------------------------------------//
	mainFnType := llvm.FunctionType(llvm.GlobalContext().VoidType(), nil, false)
	mainFn := llvm.AddFunction(mod, "main", mainFnType)
	entry := ctx.AddBasicBlock(mainFn, "entry")
	builder.SetInsertPointAtEnd(entry)
	// call memory_new
	cmem := builder.CreateCall(memNewType, memoryNew, nil, "cmem") // Memory* from C
	// Call memory_set--------------------------------------------------------//
	int32Type := llvm.GlobalContext().Int32Type()
	value := llvm.ConstInt(int8Type, 40, false) // i8 xx
	valueArray := llvm.ConstArray(int8Type, []llvm.Value{value})
	globalVal := llvm.AddGlobal(mod, valueArray.Type(), "val")
	globalVal.SetInitializer(valueArray)
	ptrVal := builder.CreateGEP(valueArray.Type(), globalVal, []llvm.Value{
		llvm.ConstInt(int32Type, 0, false), // first index into global
		llvm.ConstInt(int32Type, 0, false), // second index into array element
	}, "val_ptr")
	builder.CreateCall(memorySetType, memorySet, []llvm.Value{
		cmem,                               // Memory*
		llvm.ConstInt(int64Type, 0, false), // offset 0
		llvm.ConstInt(int64Type, 1, false), // size 1
		ptrVal,                             // pointer to value 42
	}, "")

	// call goUseMemory
	builder.CreateCall(goUseMemType, goUseMemory, []llvm.Value{cmem}, "")
	// return void from main
	builder.CreateRetVoid()

	// Execution engine
	engine, _ := llvm.NewExecutionEngine(mod)
	engine.AddObjectFileByFilename("./mymem.o")
	engine.AddGlobalMapping(goUseMemory, unsafe.Pointer(C.GoUseMemory))

	// Run main
	engine.RunFunction(mainFn, nil)
	fmt.Println("Main returned:")
}
