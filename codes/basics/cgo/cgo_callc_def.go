package main

/*
#include <stdlib.h>

static inline int add(int a, int b) {
	return a + b;
}
*/
import "C"
import "fmt"

func Random() {
	fmt.Println("Calling C.random from Go:", C.random())
}

func Add(a, b int) {
	fmt.Println("Calling C.add from Go:", C.add(C.int(a), C.int(b)))
}
