package interfaces

import "testing"

type Animal interface {
	walk()
}

type dog struct {
}

func (d dog) walk() {
	println("wangwang")
}

// nest anonymous structure
type pet struct {
	Animal
}

func TestAnonymousInterface(t *testing.T) {
	d := dog{}
	p := pet{d}
	// pet auto inherit Animal's method
	p.walk()
}
