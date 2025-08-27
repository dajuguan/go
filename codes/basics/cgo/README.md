
## 核心原理
在Go中调用C函数时，runtime.cgocall中调用entersyscall脱离调度器管理。runtime.asmcgocall切换到m的g0栈，于是得到C的运行环境。

在C中调用Go函数时，crosscall2解决gcc编译到6c编译之间的调用协议问题。cgocallback切换回goroutine栈。runtime.cgocallbackg中调用exitsyscall恢复Go的运行环境。
## References
- https://go.dev/wiki/cgo
- https://tiancaiamao.gitbooks.io/go-internals/content/zh/09.2.html