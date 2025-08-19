package main

/*
extern void go_callback(int fnId, int arg);
// fn is the registed func in Go
static inline void CallGoFunction(int fnId) {
    go_callback(fnId, 5);
}
*/
import "C"

import (
	"fmt"
	"sync"
)

var mu sync.Mutex
var index int
var fns = make(map[int]func(C.int))

//export go_callback
func go_callback(fnId C.int, arg C.int) {
	fn := lookup(int(fnId))
	fn(arg)
	fmt.Printf("Call Go's fnId:%d done!\n", fnId)
}

func lookup(i int) func(C.int) {
	mu.Lock()
	defer mu.Unlock()
	return fns[i]
}

func register(fn func(C.int)) int {
	mu.Lock()
	defer mu.Unlock()
	index++
	for fns[index] != nil {
		index++
	}
	fns[index] = fn
	return index
}

func MyGoCallBack(arg C.int) {
	fmt.Println("MyGoCallback called with arg:", arg)
}

func ExampleCCallGo() {
	fnPtr := register(MyGoCallBack)
	// Can't directly pass MyGoCallBack to C because the pointer passing rules:
	// Go code cannot pass a pointer from the Go heap (especially one that points to Go memory) to C and let C hold it for a long time.
	// A Go function value (closure or function variable) is actually an object allocated by the Go runtime, which encapsulates the instruction address and its execution context. Therefore, Go cannot directly treat a function value as a pointer and pass it to C; otherwise, it would violate the pointer passing rule.

	C.CallGoFunction(C.int(fnPtr))
}
