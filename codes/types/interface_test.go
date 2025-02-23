package types

import "testing"

type Animal interface {
	walk()
}

type Lion struct {
	name string
}

func (l *Lion) walk() {
	l.name = "changed lion"
	println("lioner pointer walk")
}

type Duck struct {
	name string
}

func (duck Duck) walk() {
	duck.name = "changed duck"
	println("duck walk")
}

func walk(animal Animal) {
	animal.walk()
}

func TestInterf(t *testing.T) {
	a := &Lion{name: "lion"}
	walk(a)
	println("a.name:", a.name)
	b := Duck{name: "duck"}
	walk(b)
	println("b.name:", b.name)
}
