package basics

import (
	"fmt"
	"testing"
)

type NestedStruct struct {
	Num    int
	Groups map[int]int
}

type OuterStrcuct struct {
	nest NestedStruct
}

func TestStructInitialization(t *testing.T) {
	o := OuterStrcuct{}
	fmt.Println(o.nest.Num)
	fmt.Println(o.nest.Groups == nil)
	// read not exist key is safe
	fmt.Println(o.nest.Groups[1])

	// write to not exist key will panic
	// o.nest.Groups[0] = 1
}
