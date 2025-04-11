package basics

import (
	"sync"
	"testing"
	"time"
)

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
