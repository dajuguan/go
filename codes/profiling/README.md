# go profiling
no need to install flamegraph tool, because pprof has enshrined flamegraph web interface.
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

# Ref
- https://eddycjy.gitbook.io/golang/di-9-ke-gong-ju/go-tool-pprof


## go routine scheduler
https://github.com/ethstorage/zk-decoder/tree/main/golang/cmd/hash_bench
tasktest -c 1 ./hash_bench -n 100000000 -r 100000 -t 1
- 1和17的核，共用ALU，超线程
- 1-16是真实核，其他的是共享相应的核心
先用Go跑出计算性能，然后优化后是2-3倍
- hyper threaded
- 影响