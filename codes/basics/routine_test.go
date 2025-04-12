package basics

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGoRoutineVar(t *testing.T) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 10)
		go func() {
			fmt.Println(i)
		}()
	}

	for i := 0; i < 3; i++ {
		i := i
		go func() {
			println("true:", i)
		}()
	}
	time.Sleep(time.Second)
	// the go rountine will run after the above codes completes
}
func TestGoRoutine(t *testing.T) {
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		println(i)
	}
	// the go rountine will run after the above codes completes
	go func() {
		println("immediate run")
	}()
}

func TestNestedRoutines(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			time.Sleep(time.Second)
			println(i)

			wg.Add(1)
			// the go rountine will run after the above codes completes
			go func() {
				defer wg.Done()

				println("immediate run", i)
			}()
		}(i)
	}

	wg.Wait()

}
