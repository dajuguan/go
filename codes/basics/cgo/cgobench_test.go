package main

import (
	"testing"
)

func BenchmarkGoCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = i + 1 // 空的 Go 函数调用
	}
}

// every cgo call c cost abot 30~40 ns
func BenchmarkCgoCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Nop(0, 1)
	}
}
