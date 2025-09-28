package main

import (
	"testing"
)

func TestCCallGoWithCgoHandle(t *testing.T) {
	ExampleCallHandle()
	ExampleCallStructHandle()
}

func BenchmarkCgoCCallGoOverhead(b *testing.B) {
	s := NewSWithHandle()
	b.ResetTimer()
	defer s.Handle.Delete() // Clean up the handle after use

	for i := 0; i < b.N; i++ {
		s.BenchmarkCallHandle()
	}
}

func BenchmarkCgoCCallGoOverheadHandle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExampleCallHandle()
	}
}
