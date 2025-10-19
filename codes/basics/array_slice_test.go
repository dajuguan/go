package basics

import "testing"

func TestArray(t *testing.T) {
	var a [5]int
	println(a[0])
	b := a
	b[0] = 1
	println(a[0], b[0])
}
