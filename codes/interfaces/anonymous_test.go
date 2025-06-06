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
	a int
}

// It'll override the anomynous function
func (p pet) walk() {
	println("pet walk")
}

func TestAnonymousInterface(t *testing.T) {
	d := dog{}
	p := pet{d, 1}
	// pet auto inherit Animal's method
	p.walk()
}
