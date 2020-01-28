# 同步与锁
Go语言更希望通过高级的同步，比如channel来实现，但是仍然提供了低级别(操作系统级别)的同步原语，比如互斥锁、同步携程组、Once等;当然，还有最低级(硬件级别)的atomic操作。

# sync.WaitGroup

sync.WaitGroup用于协程组的同步。

```
type WaitGroup struct {
        // Has unexported fields.
}
func (wg *WaitGroup) Add(delta int)
func (wg *WaitGroup) Done()
func (wg *WaitGroup) Wait()
```
- 主协程调用Add来设置协程需要等待的数量；
- 然后每一个子协程运行，并在结束时调用Done函数
- 同时，可以使用Wait来阻塞线程，保证所有的协程执行完

> 注意:WaitGroup在第一次使用完之后，千万不要拷贝，否则会引发```panic: sync: negative WaitGroup counter```。

# 原子操作
一个操作或者多个操作在执行过程中不被中断，叫做原子操作(atomatic).原子操作被封装为不可分割的整体，要么同时执行成功，要么同时执行失败，**中间状态对外不可见**。从线程的角度来看，多线程并发访问一个共享资源时，那么同一时刻只有一个线程被能对改资源操作。

有些朋友可能不知道，在 Go（甚至是大部分语言）中，一条普通的赋值语句其实不是一个原子操作。例如，在32位机器上写int64类型的变量就会有中间状态，因为它会被拆成两次写操作（MOV）——写低 32 位和写高 32 位，如下图所示：

![64位变量的赋值操作](/.gitbook/assets/原子操作.svg)

如果一个线程刚写完低32位，还没来得及写高32位时，另一个线程读取了这个变量，那它得到的就是一个毫无逻辑的中间变量，这很有可能使我们的程序出现诡异的 Bug。

# 互斥锁

原子操作都是通过互斥的方式实现的，如果想使用粗粒度的原子操作（操作系统调度），可以使用Mutex来实现。

比如，我们实现两个线程，对1-50的加法：

```
package main
import (
   ."fmt"
   "sync"
)

var total struct{
   sync.Mutex
   v int
}

func worker(wg *sync.WaitGroup){
   defer wg.Done()
   for i := 1; i < 10; i++{
      total.Lock()
      total.v += i
      total.Unlock()
   }
}

func main(){
   var wg sync.WaitGroup 
   wg.Add(2)
   go worker(&wg)
   go worker(&wg)
   wg.Wait()
   Println(total.v)
}
```


# sync/atomic

Mutex由操作系统实现，而atomic包中的原子操作则由底层硬件直接提供支持。在 CPU 实现的指令集里，有一些指令被封装进了atomic包，这些指令在执行的过程中是不允许中断（interrupt）的，因此原子操作可以在lock-free的情况下保证并发安全，并且它的性能也能做到随 CPU 个数的增多而线性扩展。(参见[喵叔没画画：atomic.Value 的前世今生](https://blog.betacat.io/post/golang-atomic-value-exploration/))

> Mutexes do not scale. Atomic load do.

用互斥锁来保护一个数值型的共享资源，麻烦且效率低下，我们可以使用atomic库，对上面的例子进行改写：

```
package main
import (
   ."fmt"
   "sync"
   "sync/atomic"
)

var total uint64

func worker(wg *sync.WaitGroup){
   defer wg.Done()
   var i uint64
   for i = 1; i < 10; i++{
      atomic.AddUint64(&total, i)
   }
}

func main(){
   var wg sync.WaitGroup 
   wg.Add(2)
   go worker(&wg)
   go worker(&wg)
   wg.Wait()
   Println(total)
}
```
atomic.AddUint64函数调用保证了total的读取、更新和保存是一个原子操作，因此在多线程中访问也是安全的。

# 单例模式

使用互斥锁配合原子操作可实现比较高效的单例模式，在性能敏感的地方采用原子检测标志位，以降低互斥锁的使用频率来提高性能。

比如Go语言标准库的sync.Once:

```
package main
import (
   ."fmt"
   "sync"
   "sync/atomic"
)

type Once struct {
   m sync.Mutex
   done uint32
}

func (o* Once) Do(f func()){
   if atomic.LoadUint32(&o.done) == 0 { //使用原子标志位，加快检测

      o.m.Lock()   //确保返回时，f已经执行完毕，所以需要加锁
      defer o.m.Unlock()
      if o.done == 0 {
         defer atomic.StoreUint32(&o.done, 1)  
         f()
      }
   }
}

func main(){
   var o Once

   var wg sync.WaitGroup
   for i := 0; i < 5; i++ {
       wg.Add(1)
       go func(i int){
         defer wg.Done()
         o.Do(func() {
            Println("第", i, "次执行...")
        })
       }(i)
   }
   wg.Wait()
}
```
可以看出输出只执行了一次。

# atomic.Value
sync/atomic包对基本的数值类型及复杂对象的读写都提供了原子操作的支持。atomic.Value原子对象提供了Load和Store两个原子方法，分别用于加载和保存数据，返回值和参数都是interface{}类型，因此可以用于任意的自定义复杂类型。

atomic.Value类型对外暴露的方法就两个：
- v.Store(c) - 写操作，将原始的变量c存放到一个atomic.Value类型的v里。
- c = v.Load() - 读操作，从线程安全的v中读取上一步存放的内容。

比如实现一个简单的自动读取更新配置信息的生产者消费者模型：

```
package main
import (
   ."fmt"
   "sync"
   "sync/atomic"
   "time"
   "math/rand"
)

type config struct{
   path string
   id int
}
func loadConfig()(c config){
   return config{path:"/111", id:rand.Int()}
}

func main(){
   var conf atomic.Value
   var wg sync.WaitGroup
   wg.Add(2)

    //重新生成配置
   go func(){
      for{      
         conf.Store(loadConfig())
         time.Sleep(time.Second*3)

      }
      // defer wg.Done()
   }()
    //每个消费者读取配置
   go func(){
      for {
         time.Sleep(time.Second*1)
         for i:=1; i < 10; i++{
            go func(){
               c := conf.Load()
               Println(c)
            }()
         }
      }
      // defer wg.Done()
   }()

   wg.Wait()
}
```

# 线程同步与顺序一致性模型

在Go语言中，之前我们使用Mutex或者atomic同步原语保证了同步，但是这个保障的前提是**顺序一致性模型**。也就是说，在同一goroutine内，不论编译器或者处理器如何对指令进行重排序，但是执行的结果是按照代码的书写顺序执行的；但是在不同的goroutine内是得不到保障的，比如:

```
package main
import (
   ."fmt"
)

func main(){
   var done bool
   var i int
   go func(){
      done = true
      Println("Hello, Go!")
   }()
   for {
      if done == true {
         Println("Hello, World:",i)
         return 
      }
      i++
   }
}
```
```
➜  gopath go run 1hello.go
Hello, Go!
Hello, World: 1079582
➜  gopath go run 1hello.go
Hello, World: 29235
```

需要等到i很大的时候，这是因为main线程和go func协程的顺序一致性得不到保证，这个时候就需要明确的同步原语来保证:

```
func main(){
   done := make(chan bool)
   var i int
   go func(){
      done <- true
      Println("Hello, Go!")
   }()
   for {
      
      if <-done == true {
         Println("Hello, World:",i)
         return 
      }
      i++
   }
}
```

或者使用sync.Mutex来保证同步:

```
package main
import (
   ."fmt"
   "sync"
)

func main(){
   var mu sync.Mutex
   mu.Lock()
   go func(){
      Println("Hello, Go!")
      mu.Unlock()
   }()
   mu.Lock()
   go func(){
      Println("Hello, World!")
      mu.Unlock()
   }()
   mu.Lock()
}
```
```
Hello, Go!
Hello, World!
```



# 基于channel的同步

## 同步控制
channel是goroutine之间同步的主要方法，在无缓冲的channel上的发送操作总是在对应的接受操作完成前发生：

```
package main
import (
   ."fmt"
)

var msg string
func main(){
   c := make(chan int)
   go func(){
      msg = "Hello, Go"
      c <- 1
   }()
   <-c
   Println(msg)
}
```

可保证打印出“Hello, Go”。该程序首先对msg进行写入，然后在c管道上发送同步信号，随后从c接收对应的同步信号，最后执行Println函数。

若在关闭Channel后继续从中接收数据，接收者就会收到该Channel返回的零值。因此在这个例子中，用close(c)关闭管道代替c <- 1依然能保证该程序产生相同的行为。

```
package main
import (
   ."fmt"
)

var msg string
func main(){
   c := make(chan int)
   go func(){
      msg = "Hello, Go"
      close(c)
   }()
   <-c
   Println(msg)
}
```

## 使用channel控制最大并发数量

```
package main
import (
   ."fmt"
   "time"
)

var limit = make(chan int, 3)

func work()(score[10] func()){
   for i:=0; i< 10; i++{
      score[i] = func(){
         time.Sleep(time.Second)
         Println("Hi")
      }
   }
   return score
}

func main(){
   for _,w := range work(){
      go func(){
         limit <- 0  
         w()
         <- limit
      }()
   }
   for{}
}

```

通过设置channel的缓冲区大小为3，我们保证goroutine每次只能执行3个