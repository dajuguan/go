# 协程(goroutine)与通道(channel)

# 并发与并行

并发指的是同一时间处理多件事情，这些事情一般是独立的，一般是抢占式的，主要通过提高单个CPU的利用率；
并行指的是同一时间同时处理一件或多件事情，一般而言这些事情是一个大任务的并列子事情，一般是非抢占式的，需要多核CPU的支持。

Go语言在语言层面支持了并发，通过goroutine提供一种用户态线程，我们一般把它成为**协程**。协程在某种程度上可以看做是轻量级线程，它由应用程序而非操作系统来创建和管理，使用开销较低（一般为4K）, 程序员可以轻松的创建很多goroutine，并且go语言会负责保证尽可能公平地调度使用CPU。

调度器主要分为4个部分，前三个定义在runtime.h中，Sched定义在proc.c中:
- M(work thread):OS Thread，由os管理
- G(goroutine):goroutine实体，包含调用栈，重要的调度信息，比如channel等
- P(processor):衔接M和G的调度上下文，负责等待执行的G和M的对接，通过GOMAXPROCS()设置P的数量，来决定有多少个goroutine可以并发。
- Sched:负责调度

在Go中，OS Thread和User Thread之间是多堆垛的关系，使得多个goroutine可以再多个OS Threads是执行。
- 即保证了上下文切换的效率
- 又能利用多核的优势

每个goroutine都会被特定的P维护，每个P会把自己所维护的所有的goroutine放到一个G队列中，然后M每次挑选一个有效P，然后执行其中的goroutine。

默认情况下，P的数量和M一样。所以当我们创建多个goroutine时，他们会被分配到不同的P中，而M不唯一，当M随机挑选P是，也就等于随机挑选了goroutine来执行。

![协程](/.gitbook/assets/goroutine.png)

所以多个goroutine之间的执行顺序是不确定的，因为gorouine进入的P是不确定的，而P被执行的顺序也是不确定的。

# goroutine

在Go语言中，只需要在代码块前面加上关键字go，即可创建一个goroutine。

```
package main
import (
   ."fmt"
   "time"
)

func main(){
   for i := 0; i < 10; i++ {
      go func(i int){
         Println(i)
      }(i)
   }
   time.Sleep(time.Second*1)
}
```
可以运行看到输出的结果为
```
1
3
8
2
6
0
9
7
5
4
```
并且每次的顺序是随机的。
在这里读者们注意到了为什么要使用time.Sleep，这是因为在Go中，goroutine和main函数使用的不是一个协程，他们之间是并发的关系，谁都有可能发生，并且main函数不会等待其他goroutine运行完才推出，main函数主体运行完就会退出。如果不加这一句，可能得不到输出。但是加上这一句我们就一定能得到正确的输出吗？答案是否的的，因为每个goroutine实际等待的时间我们不能预知，如果有外界的因素干扰，有可能等1秒，但是其他的goroutine被阻塞了，还是得不到正确的输出结果。

所以这个时候就需要锁，来进行显式的同步。