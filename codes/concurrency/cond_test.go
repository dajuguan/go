package concurrency

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var done = false

func read(name string, c *sync.Cond) {
	c.L.Lock()
	for !done {
		c.Wait() // 等价于 c.L.unLock(), wait for signal, c.L.Lock(), so must call unlock after
	}
	c.L.Unlock()
	fmt.Println("read:", name)
}

func write(c *sync.Cond) {
	time.Sleep(time.Second)
	c.L.Lock()
	done = true
	c.L.Unlock()
	c.Broadcast()
}

func TestCond(t *testing.T) {
	var mu sync.Mutex
	c := sync.NewCond(&mu)
	go read("1", c)
	go read("2", c)
	go read("3", c)
	write(c)
}
