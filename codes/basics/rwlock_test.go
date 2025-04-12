package basics

import (
	"sync"
	"testing"
	"time"
)

var (
	mu   sync.RWMutex
	data = map[int]string{
		1: "foo",
		2: "t",
	}
)

func read(key int) {
	mu.RLock()
	defer mu.RUnlock()

	println("read:", key, data[key])
	time.Sleep(time.Millisecond * 400)

}

func write(key int, val string) {
	mu.Lock()
	defer mu.Unlock()

	println("write:", key, val)
	time.Sleep(time.Millisecond * 400)
	data[key] = val
}

func TestRLock(t *testing.T) {
	// 并发读锁可以同时获取
	go read(1)
	go read(1)
	go read(1)
	go write(1, "test")
	go read(1)

	time.Sleep(time.Second)
}
