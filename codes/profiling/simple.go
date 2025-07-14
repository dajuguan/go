package profiling

import (
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"testing"
	"time"
)

const (
	sleepTime   = 30 * time.Millisecond
	cpuTime     = 30 * time.Millisecond
	networkTime = 60 * time.Millisecond
)

func cpuIntensiveTask() {
	start := time.Now()
	for time.Since(start) <= cpuTime {
		// Spend some time in a hot loop to be a little more realistic than
		// spending all time in time.Since().
		for i := 0; i < 1000; i++ {
			_ = i
		}
	}
}

func weirdFunction() {
	time.Sleep(sleepTime)
}

func TestProfiling(t *testing.T) {
	// 创建 profile 文件
	f, _ := os.Create("cpu.prof")
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// 创建 trace 文件
	traceFile, _ := os.Create("trace.out")
	defer traceFile.Close()
	trace.Start(traceFile)
	defer trace.Stop()

	runTime := time.Second * 10
	now := time.Now()
	for {
		if time.Since(now) > runTime {
			break
		}
		// // Http request to a web service that might be slow.
		// slowNetworkRequest()
		// Some heavy CPU computation.
		cpuIntensiveTask()
		// Poorly named function that you don't understand yet.
		weirdFunction()
	}

}
