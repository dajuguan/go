package basics

import (
	"testing"
	"time"
)

// 同时向无缓冲区ch加入多个数据
func TestChDeadlock(t *testing.T) {
	// 创建一个无缓冲的通道
	ch := make(chan func())

	go func() {
		println("executing")
		for fn := range ch {
			println("executing fn...")
			fn()
		}
	}()

	for i := 0; i < 5; i++ {
		i := i
		println("add fi:", i, "len ch:", len(ch))
		ch <- func() {
			println("i:", i)
			time.Sleep(time.Millisecond * 10)
			for j := 0; j < 5; j++ {
				j := j
				println("add fj:", j, "len ch:", len(ch))
				ch <- func() {
					println("j:", j)
				}
			}
		}
	}

	close(ch)
}
