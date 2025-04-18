# go profiling
## add code
```
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

```

## visualize the profiling results
```
go run simple.go
# view flamegraph without IO/sleep time
go tool pprof -http=0.0.0.0:8080 cpu.prof
# view trace
go tool trace -http=0.0.0.0:8080 trace.out
```