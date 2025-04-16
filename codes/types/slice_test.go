package types

import (
	"fmt"
	"testing"
)

func sliceAsArg(s []int) {
	s[0] = 1
}

func sliceAppend(s []int) {
	_ = append(s, 1)
}

func TestSlice(t *testing.T) {
	a := []int{0, 2, 3}
	sliceAsArg(a)
	fmt.Println(a)
}

func TestSliceAppend(t *testing.T) {
	a := make([]int, 0, 5)
	a = append(a, 1, 2, 3)
	sliceAppend(a)
	// slice a has been changed because capacity is enough
	fmt.Println(a[:]) // this will print the same, because len is not changed, but its underlying memory has been changed
	fmt.Println(a[:4])

	b := make([]int, 0)
	b = append(b, 1, 2, 3)
	sliceAppend(b)
	// slice b is the same, so it'll throw error
	// fmt.Println(b[:4])
}

func TestSliceAsStruct(t *testing.T) {
	type S struct {
		slice []int
	}
	s := []int{0, 0, 0}
	a := S{slice: s}
	fmt.Println("slice before:", s)
	fmt.Println("struct after:", a)
	s[0] = 5
	fmt.Println("slice after:", s)
	fmt.Println("struct after:", a)
}
