
## 核心原理
在Go中调用C函数时，runtime.cgocall中调用entersyscall脱离调度器管理。runtime.asmcgocall切换到m的g0栈，于是得到C的运行环境。

在C中调用Go函数时，crosscall2解决gcc编译到6c编译之间的调用协议问题。cgocallback切换回goroutine栈。runtime.cgocallbackg中调用exitsyscall恢复Go的运行环境。

## benchmake CGO
```
taskset -c 1 go test -run=^$ -bench BenchmarkCgoC github.com/dajuguan/go/codes/basics/cgo -v 

BenchmarkCgoCCallGoOverhead                 14120484                86.15 ns/op  (Go=>c=>go, exclude handle (cgo type) init)
BenchmarkCgoCCallGoOverheadHandle           2772770                 410.7 ns/op (Go=>c=>go, include handle (cgo type) init)
BenchmarkCgoCCallGoOverheadRegistry         13506428                88.21 ns/op (Go=>c=>go, exclude handle(registry type) init)
BenchmarkCgoCCallGoOverheadRegistryHandle   9940838                 119.7 ns/op (Go=>c=>go, include handle(registry type) init)
BenchmarkCgoCall                           34297670                 35.16 ns/op
```

- 创建Cgo handle的开销是320ns (410-86 = 314), 用registry创建handle的时间只需要30ns (119-88=31)
- 通过Cgo handle从C调用Go的overhead是50ns （86-35 = 51）
- 通过Cgo, 从Go call C的overhead是~40ns

## benchmark
```
taskset -c 1 go test -run=^$ -bench=. -v -tags=llvm20

BenchmarkGoSlicePassToCAndSetmem-28             15674406                74.87 ns/op
BenchmarkCAllocAndSetmem-28                     23777948                52.18 ns/op
BenchmarkGoSlicePassToCAndLLVMSetmemLoop-28     18198148                64.05 ns/op
BenchmarkGoSlicePassToCAndLLVMFnptr-28         18215420                66.29 ns/op
```

可以看出：
- 采用Funcptr和直接内联调用的性能是一样的
- 在C中分配内存要比使用Go传递内存快一些
- llvm是setmem相对于直接在go中set met还是要快一些，但是还是没有直接调用c函数set met来的更快

## References
- https://go.dev/wiki/cgo
- https://tiancaiamao.gitbooks.io/go-internals/content/zh/09.2.html
- [cockroachlabs cgo over head](https://www.cockroachlabs.com/blog/the-cost-and-complexity-of-cgo/)