package interfaces

import (
	"testing"
)

type testBackEnd struct {
	number int
}

func (b *testBackEnd) Number() int {
	return b.number
}

func TestFilterAPI(t *testing.T) {
	var backend Backend
	backend = &testBackEnd{number: 3}
	api := NewFilterAPI(backend)
	println("api.GetNumber()", api.GetNumber())
}
