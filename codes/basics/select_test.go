package basics

import (
	"testing"
	"time"
)

func TestSelectBlock(t *testing.T) {
	a := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		a <- 1
	}()

	println("blocking select")
	select {
	case v := <-a:
		println(v)
	}
	println("end")
}

func TestSelectNonBlocking(t *testing.T) {
	a := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		a <- 1
	}()

	println("non blocking select")
	select {
	case v := <-a:
		println(v)
	default:
		println("default")
	}
	println("end")
}

func TestSelectNonBlocking2(t *testing.T) {
	a := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		a <- 1
	}()

	b := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		b <- 1
	}()

	println("non blocking select")
	select {
	case v := <-a:
		println("a:", v)
	case v := <-b:
		println("b:", v)
	default:
		println("default")
	}
	println("end")
}
