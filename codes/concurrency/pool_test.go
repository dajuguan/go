package concurrency

import (
	"fmt"
	"sync"
	"testing"
)

func TestPool(t *testing.T) {
	var pool = sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new object")
			return make([]byte, 1024)
		},
	}

	// 放入多个对象
	var wg sync.WaitGroup

	// 多次 Get
	for i := 0; i < 5; i++ {
		wg.Add(1)
		i := i
		go func(i int) {
			defer wg.Done()
			buf := pool.Get().([]byte)
			buf = make([]byte, i+1)
			buf[0] = byte(i)
			pool.Put(buf)
		}(i)
	}

	wg.Wait()
	fmt.Println("Creating object pool done===")

	for i := 0; i < 5; i++ {
		buf := pool.Get().([]byte)
		fmt.Println("Got object of length", buf)
		pool.Put(buf)
	}
	// Creating new object”被打印少于 5 次，说明之前的对象被复用了。
}
