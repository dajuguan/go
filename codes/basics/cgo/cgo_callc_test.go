package main

import (
	"testing"
)

func TestCallCFromGo(t *testing.T) {
	Random()
	Add(1, 2)
}

func BenchmarkCgoGoCallC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(1, 2)
	}
}
