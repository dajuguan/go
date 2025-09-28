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
	defer s.Handle.Delete() // Clean up the handle after use
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.BenchmarkCallHandle()
	}
}

func BenchmarkCgoCCallGoOverheadHandle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExampleCallHandle()
	}
}

func BenchmarkCgoCCallGoOverheadRegistry(b *testing.B) {
	fnPtr := register(MyGoCallBack)
	defer unregister(fnPtr)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		RegistryCallGoFunc(fnPtr)
	}
}

func BenchmarkCgoCCallGoOverheadRegistryHandle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExampleCCallGo()
	}
}
