package interfaces

import "testing"

type Human interface {
	speak()
}

type Man struct{}

func (m *Man) speak() {
	println("man speak")
}

type Women struct{}

func (m *Women) speak() {
	println("woman speak")
}

func TestDuckType(t *testing.T) {
	m := Man{}
	w := Women{}
	animals := []Human{&m, &w}
	for _, people := range animals {
		people.speak()
	}

}
