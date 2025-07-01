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
