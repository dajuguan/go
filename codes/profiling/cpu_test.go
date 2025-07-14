package profiling

import (
	"fmt"
	"sync"
	"syscall"
	"testing"
	"time"
)

func GetCPU() int64 {
	usage := new(syscall.Rusage)
	syscall.Getrusage(syscall.RUSAGE_SELF, usage)
	return usage.Utime.Nano() + usage.Stime.Nano()
}

func TestCPU(t *testing.T) {
	cpuUsage := GetCPU()
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cpuIntensiveTask()
		}()
	}
	wg.Wait()

	time.Sleep(time.Millisecond * 500)

	// CPU time accumulates all go routine's consumed time
	fmt.Println("cpu time:", time.Duration(GetCPU()-cpuUsage))
}
