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

## Test profiling
Add code:
```
import (
    _ "net/http/pprof"
    "net/http"
    "log"
)

func init() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```
访问: http://localhost:6060/debug/pprof/heap?debug=1
可视化(8080为可视化端口，会从6060拉数据): go tool pprof -http=:8080  http://localhost:6060/debug/pprof/heap 

## 解读go heap profiling

```
4599	allocs (for history objects)
0	block
0	cmdline
70	goroutine
4599	heap (for live objects and history objects)
```
最外面的表示采了多少次样，和具体的memory/goroutine数目无关

### Heap 
The numbers in the beginning of each entry ("1: 262144 [4: 376832]") represent number of currently live objects, amount of memory occupied by live objects, total number of allocations and amount of memory occupied by all allocations, respectively.

```
heap profile: 13: 1065160 [2815614: 125248125936] @ heap/1048576
第一行: 总共有13个对象分配(4+4+1*5),共分配了1065160 bytes(425984+425984+106496+106496+96+80+24)内存。[2815614(总分配次数):: 125248125936(总分配内存bytes)]
4: 425984 [2677: 285089792] @ 0x459332 0x458d25 0x458eca 0x79d393 0x472e41
#	0x79d392	6.824/kvraft.(*KVServer).applyCmds+0x392	/home/po/now/MIT6.824/6.824/src/kvraft/server.go:71

4: 425984 [7522: 801062912] @ 0x50474e 0x50417d 0x504359 0x5f5cb1 0x5f6c08 0x5f6609 0x5ff9bf 0x79da96 0x79d4bb 0x472e41
#	0x50474d	bytes.growSlice+0x8d				/home/po/.gvm/gos/go1.21.13/src/bytes/buffer.go:249
#	0x50417c	bytes.(*Buffer).grow+0x13c			/home/po/.gvm/gos/go1.21.13/src/bytes/buffer.go:151
#	0x504358	bytes.(*Buffer).Write+0x58			/home/po/.gvm/gos/go1.21.13/src/bytes/buffer.go:179
#	0x5f5cb0	encoding/gob.(*Encoder).writeMessage+0x3b0	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/encoder.go:82
#	0x5f6c07	encoding/gob.(*Encoder).EncodeValue+0x447	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/encoder.go:253
#	0x5f6608	encoding/gob.(*Encoder).Encode+0x68		/home/po/.gvm/gos/go1.21.13/src/encoding/gob/encoder.go:176
#	0x5ff9be	6.824/labgob.(*LabEncoder).Encode+0x3e		/home/po/now/MIT6.824/6.824/src/labgob/labgob.go:34
#	0x79da95	6.824/kvraft.(*KVServer).snapshot+0x55		/home/po/now/MIT6.824/6.824/src/kvraft/server.go:192
#	0x79d4ba	6.824/kvraft.(*KVServer).applyCmds+0x4ba	/home/po/now/MIT6.824/6.824/src/kvraft/server.go:85

1: 106496 [532: 56655872] @ 0x5dda1b 0x5eb965 0x5eb8de 0x5ebbc7 0x5ec00e 0x5ebe45 0x5ffaee 0x6008db 0x79a625 0x79fd6a 0x79fd44 0x79fd43 0x79fd42 0x79e21a 0x472e41
#	0x5dda1a	internal/saferio.ReadData+0x5a			/home/po/.gvm/gos/go1.21.13/src/internal/saferio/io.go:36
#	0x5eb964	encoding/gob.(*Decoder).readMessage+0x44	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:103
#	0x5eb8dd	encoding/gob.(*Decoder).recvMessage+0xbd	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:91
#	0x5ebbc6	encoding/gob.(*Decoder).decodeTypeSequence+0x46	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:148
#	0x5ec00d	encoding/gob.(*Decoder).DecodeValue+0x16d	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:227
#	0x5ebe44	encoding/gob.(*Decoder).Decode+0x124		/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:204
#	0x5ffaed	6.824/labgob.(*LabDecoder).Decode+0x4d		/home/po/now/MIT6.824/6.824/src/labgob/labgob.go:55
#	0x6008da	6.824/labrpc.(*ClientEnd).Call+0x2fa		/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:120
#	0x79a624	6.824/kvraft.(*Clerk).Request+0x1c4		/home/po/now/MIT6.824/6.824/src/kvraft/client.go:63
#	0x79fd69	6.824/kvraft.(*Clerk).PutAppend+0x249		/home/po/now/MIT6.824/6.824/src/kvraft/client.go:85
#	0x79fd43	6.824/kvraft.(*Clerk).Append+0x223		/home/po/now/MIT6.824/6.824/src/kvraft/client.go:92
#	0x79fd42	6.824/kvraft.Append+0x222			/home/po/now/MIT6.824/6.824/src/kvraft/test_test.go:41
#	0x79fd41	6.824/kvraft.GenericTest.func1+0x221		/home/po/now/MIT6.824/6.824/src/kvraft/test_test.go:216
#	0x79e219	6.824/kvraft.run_client+0xd9			/home/po/now/MIT6.824/6.824/src/kvraft/test_test.go:57

1: 106496 [1080: 115015680] @ 0x50474e 0x50417d 0x504359 0x5f5cb1 0x5f6c08 0x5ffa65 0x6032ab 0x6028e5 0x601c38 0x472e41
#	0x50474d	bytes.growSlice+0x8d				/home/po/.gvm/gos/go1.21.13/src/bytes/buffer.go:249
#	0x50417c	bytes.(*Buffer).grow+0x13c			/home/po/.gvm/gos/go1.21.13/src/bytes/buffer.go:151
#	0x504358	bytes.(*Buffer).Write+0x58			/home/po/.gvm/gos/go1.21.13/src/bytes/buffer.go:179
#	0x5f5cb0	encoding/gob.(*Encoder).writeMessage+0x3b0	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/encoder.go:82
#	0x5f6c07	encoding/gob.(*Encoder).EncodeValue+0x447	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/encoder.go:253
#	0x5ffa64	6.824/labgob.(*LabEncoder).EncodeValue+0x64	/home/po/now/MIT6.824/6.824/src/labgob/labgob.go:39
#	0x6032aa	6.824/labrpc.(*Service).dispatch+0x3aa		/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:511
#	0x6028e4	6.824/labrpc.(*Server).dispatch+0x1e4		/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:430
#	0x601c37	6.824/labrpc.(*Network).processReq.func1+0x57	/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:247

1: 96 [524: 50304] @ 0x40a6ec 0x492e88 0x79d7de 0x79d802 0x4dafa7 0x4da079 0x60326e 0x6028e5 0x601c38 0x472e41
#	0x492e87	time.NewTimer+0x27				/home/po/.gvm/gos/go1.21.13/src/time/sleep.go:87
#	0x79d7dd	time.After+0x1fd				/home/po/.gvm/gos/go1.21.13/src/time/sleep.go:157
#	0x79d801	6.824/kvraft.(*KVServer).Request+0x221		/home/po/now/MIT6.824/6.824/src/kvraft/server.go:155
#	0x4dafa6	reflect.Value.call+0xce6			/home/po/.gvm/gos/go1.21.13/src/reflect/value.go:596
#	0x4da078	reflect.Value.Call+0xb8				/home/po/.gvm/gos/go1.21.13/src/reflect/value.go:380
#	0x60326d	6.824/labrpc.(*Service).dispatch+0x36d		/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:506
#	0x6028e4	6.824/labrpc.(*Server).dispatch+0x1e4		/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:430
#	0x601c37	6.824/labrpc.(*Network).processReq.func1+0x57	/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:247

1: 80 [416: 33280] @ 0x492ea8 0x79d7de 0x79d802 0x4dafa7 0x4da079 0x60326e 0x6028e5 0x601c38 0x472e41
#	0x492ea7	time.NewTimer+0x47				/home/po/.gvm/gos/go1.21.13/src/time/sleep.go:88
#	0x79d7dd	time.After+0x1fd				/home/po/.gvm/gos/go1.21.13/src/time/sleep.go:157
#	0x79d801	6.824/kvraft.(*KVServer).Request+0x221		/home/po/now/MIT6.824/6.824/src/kvraft/server.go:155
#	0x4dafa6	reflect.Value.call+0xce6			/home/po/.gvm/gos/go1.21.13/src/reflect/value.go:596
#	0x4da078	reflect.Value.Call+0xb8				/home/po/.gvm/gos/go1.21.13/src/reflect/value.go:380
#	0x60326d	6.824/labrpc.(*Service).dispatch+0x36d		/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:506
#	0x6028e4	6.824/labrpc.(*Server).dispatch+0x1e4		/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:430
#	0x601c37	6.824/labrpc.(*Network).processReq.func1+0x57	/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:247

1: 24 [114: 2736] @ 0x492e88 0x79d7de 0x79d802 0x4dafa7 0x4da079 0x60326e 0x6028e5 0x601c38 0x472e41
#	0x492e87	time.NewTimer+0x27				/home/po/.gvm/gos/go1.21.13/src/time/sleep.go:87
#	0x79d7dd	time.After+0x1fd				/home/po/.gvm/gos/go1.21.13/src/time/sleep.go:157
#	0x79d801	6.824/kvraft.(*KVServer).Request+0x221		/home/po/now/MIT6.824/6.824/src/kvraft/server.go:155
#	0x4dafa6	reflect.Value.call+0xce6			/home/po/.gvm/gos/go1.21.13/src/reflect/value.go:596
#	0x4da078	reflect.Value.Call+0xb8				/home/po/.gvm/gos/go1.21.13/src/reflect/value.go:380
#	0x60326d	6.824/labrpc.(*Service).dispatch+0x36d		/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:506
#	0x6028e4	6.824/labrpc.(*Server).dispatch+0x1e4		/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:430
#	0x601c37	6.824/labrpc.(*Network).processReq.func1+0x57	/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:247

0: 0 [85: 348160] @ 0x5dda1b 0x5eb965 0x5eb8de 0x5ebbc7 0x5ec00e 0x5ebe45 0x5ffaee 0x79dc38 0x79d1c5 0x472e41
#	0x5dda1a	internal/saferio.ReadData+0x5a			/home/po/.gvm/gos/go1.21.13/src/internal/saferio/io.go:36
#	0x5eb964	encoding/gob.(*Decoder).readMessage+0x44	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:103
#	0x5eb8dd	encoding/gob.(*Decoder).recvMessage+0xbd	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:91
#	0x5ebbc6	encoding/gob.(*Decoder).decodeTypeSequence+0x46	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:148
#	0x5ec00d	encoding/gob.(*Decoder).DecodeValue+0x16d	/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:227
#	0x5ebe44	encoding/gob.(*Decoder).Decode+0x124		/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:204
#	0x5ffaed	6.824/labgob.(*LabDecoder).Decode+0x4d		/home/po/now/MIT6.824/6.824/src/labgob/labgob.go:55
#	0x79dc37	6.824/kvraft.(*KVServer).restoreState+0x117	/home/po/now/MIT6.824/6.824/src/kvraft/server.go:205
#	0x79d1c4	6.824/kvraft.(*KVServer).applyCmds+0x1c4	/home/po/now/MIT6.824/6.824/src/kvraft/server.go:56

0: 0 [32: 5632] @ 0x5ea035 0x5eab6a 0x5eaebd 0x5ec036 0x5ebe45 0x5ffaee 0x6008db 0x79a625 0x7a11aa 0x7a117d 0x7a117c 0x7a117b 0x79e21a 0x472e41
#	0x5ea034	encoding/gob.(*Decoder).compileDec+0x194		/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decode.go:1150
#	0x5eab69	encoding/gob.(*Decoder).getDecEnginePtr+0x129		/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decode.go:1189
#	0x5eaebc	encoding/gob.(*Decoder).decodeValue+0xfc		/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decode.go:1234
#	0x5ec035	encoding/gob.(*Decoder).DecodeValue+0x195		/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:229
#	0x5ebe44	encoding/gob.(*Decoder).Decode+0x124			/home/po/.gvm/gos/go1.21.13/src/encoding/gob/decoder.go:204
#	0x5ffaed	6.824/labgob.(*LabDecoder).Decode+0x4d			/home/po/now/MIT6.824/6.824/src/labgob/labgob.go:55
#	0x6008da	6.824/labrpc.(*ClientEnd).Call+0x2fa			/home/po/now/MIT6.824/6.824/src/labrpc/labrpc.go:120
#	0x79a624	6.824/kvraft.(*Clerk).Request+0x1c4			/home/po/now/MIT6.824/6.824/src/kvraft/client.go:63
#	0x7a11a9	6.824/kvraft.(*Clerk).PutAppend+0x2e9			/home/po/now/MIT6.824/6.824/src/kvraft/client.go:85
#	0x7a117c	6.824/kvraft.(*Clerk).Append+0x2bc			/home/po/now/MIT6.824/6.824/src/kvraft/client.go:92
#	0x7a117b	6.824/kvraft.Append+0x2bb				/home/po/now/MIT6.824/6.824/src/kvraft/test_test.go:41
#	0x7a117a	6.824/kvraft.GenericTestLinearizability.func1+0x2ba	/home/po/now/MIT6.824/6.824/src/kvraft/test_test.go:359
#	0x79e219	6.824/kvraft.run_client+0xd9				/home/po/now/MIT6.824/6.824/src/kvraft/test_test.go:57

...
# HeapIdle = 60088320
# HeapInuse = 10559488 (当前使用内存bytes)
``



## go routine scheduler
https://github.com/ethstorage/zk-decoder/tree/main/golang/cmd/hash_bench
taskset -c 1 ./hash_bench -n 100000000 -r 100000 -t 1
- 1和17的核，共用ALU，超线程
- 1-16是真实核，其他的是共享相应的核心
先用Go跑出计算性能，然后优化后是2-3倍
- hyper threaded
- 影响


# Ref
- https://eddycjy.gitbook.io/golang/di-9-ke-gong-ju/go-tool-pprof
- [Go程序内存泄露问题快速定位](https://www.hitzhangjie.pro/blog/2021-04-14-go%E7%A8%8B%E5%BA%8F%E5%86%85%E5%AD%98%E6%B3%84%E9%9C%B2%E9%97%AE%E9%A2%98%E5%BF%AB%E9%80%9F%E5%AE%9A%E4%BD%8D/)
- https://lrita.github.io/2017/05/26/golang-memory-pprof/
- [intel:Debugging performance issues in Go programs](https://web.archive.org/web/20140703183759/https://software.intel.com/en-us/blogs/2014/05/10/debugging-performance-issues-in-go-programs)