package basics

import (
	"fmt"
	"testing"
	"unsafe"
)

// Test memory pointer might be changed if the memory is unpinned.
func TestUnpinMemory(t *testing.T) {
	// 1. 分配一个大切片，让它在heap上
	s := make([]int, 10_000_000)
	s[0] = 42

	// 2. 打印切片 header 信息
	sliceHeader := (*struct {
		ptr uintptr
		len int
		cap int
	})(unsafe.Pointer(&s))
	fmt.Printf("Before GC: ptr=%x\n", sliceHeader.ptr)

	news := make([]int, 10_000_000)
	s = append(s, news...)

	// 3. 再次打印切片 header 信息
	fmt.Printf("After GC: ptr=%x\n", sliceHeader.ptr)

	// 6. 访问切片数据
	fmt.Println("s[0] =", s[0])
}
