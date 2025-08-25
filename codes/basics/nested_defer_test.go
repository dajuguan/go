package basics

import (
	"fmt"
	"testing"
)

func TestNestedDefer(t *testing.T) {
	defer func() {
		println("1")
	}()

	if true {
		fmt.Println("a")
		defer func() {
			fmt.Println(2)
		}()
	}
	fmt.Println("main ends")
}

type deferStruct struct {
	ret int
}

// ret will be set before return
func (a *deferStruct) deferReturn() (ret int) {
	defer func() { a.ret = ret }()
	return 1
}

func TestDeferReturn(t *testing.T) {
	s := deferStruct{}
	s.deferReturn()
	println(s.ret)
}
