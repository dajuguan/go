# 常见并发模型

# 生产者消费者模型
生产者消费者模型，是常见的并发模型，生产者产生一系列数据，放在队列中；然后消费者去消费这些数据；如果消费者处理能力很强，那么就会消费者就会进去饥饿的等待状态，反之，生产者数据过剩，生产者就会被阻塞。
```
package main
import (
   ."fmt"
   "os"
   "os/signal"
   "syscall"
)

func producer(factor int, c chan int){
   for i := 0; i<10 ; i++{
      c <- factor*i
   }
}

func consumer(c chan int){
   for v := range c{
      Println(v)
   }
}

func main(){
   c := make(chan int, 10)
   go producer(3, c)
   go producer(5, c)
   go consumer(c)

    //Ctrl+c退出
   sig := make(chan os.Signal, 1)
   signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
   Printf("quit (%v)\n", <-sig)
}
```

# 管道串联实现素数筛

素数筛是一个经典的并发例子，通过串联不同素数的管道，来实现逐个输出素数，原理如下

![素数筛](/.gitbook/assets/字符串内存组织.png)

```
package main
import (
   ."fmt"
)

func generatePrime() chan int{
   c := make(chan int)
   go func(){
      for i :=2; ;i++{
         c <- i
         Println("in-----------------")
      }
   }()
   return c
}

func filter(c chan int, prime int) chan int{
   out := make(chan int)
   go func(){
      for{
         i := <-c 
         Println("out>>>>>>",prime,i)
         if i % prime != 0{
            out <- i  
            Println("输出>>>>>>",prime,i)        
         }
      }
   }()
   return out
}

func main(){
   ch := generatePrime()
   for i := 0; i< 10; i++{
      prime := <- ch
      Println("第",i,"个素数为：", prime)
      ch = filter(ch, prime)
   }
}
```
1. 生成自然数管道序列generatePrime
2. 过滤已有素数filter
3. 串联上管道ch = filter(ch, prime)

程序非常巧妙的实现了输出前10个素数。

# 基于select的多路复用

select操作类似switch语句，并可以用于channel的选择，当某个channel准备好时，就会执行它；当多个channle同时准备好时，会随机执行，直到执行完毕；没有准备好的channel时，会执行default语句。

比如，实现一个管道操作，超时退出的操作:

```
package main
import (
   ."fmt"
   "time"
   "math/rand"
)

func main(){
   ch := make(chan int)
   rand.Seed(time.Now().Unix())
   go func(){
      for{
         var i int = rand.Intn(4)
         Println("输入",i)
         time.Sleep(time.Second* time.Duration(i))
         ch <- i
      }
   }()
   select {
      case v := <-ch:
         Println(v)
      case <-time.After(time.Second*time.Duration(2)):{
         Println("Timeout")
         return
      }
   }
}
```

或者，实现一个输出整数的程序：
```
package main
import (
   ."fmt"
)

func main(){
   ch := make(chan int, 1)
   for i :=0;i<100;i++{
      select {
      case x:= <-ch:
         Println(x)
      case ch<-i:
      }
   }
}
```

注意，次数缓冲区大小需设置为1，才能保证select语句正确执行，不陷入死锁。

# 并发的安全退出

有时候我们需要通知Goroutine停止它在干的事情，特别是协程错误的时候，但是go语言并未提供一个直接接收goroutine停止的方法，因为这样会导致goroutine之间的共享变量处在未定义的状态。那么我们想停止一个或者福讴歌goroutine，改怎么办的？

## 单个协程退出程序

基于上面的select语句，可以容易的实现一个goroutine退出的程序：

```
package main
import (
   ."fmt"
   "time"
)

func worker(ch chan int){
   for{
      select {
      default:
         Println("Working")
      case <-ch:
         Println("Exiting")
         return
      }
   }

}

func main(){
   ch := make(chan int)
   go worker(ch)
   time.Sleep(time.Second)
   ch <- 1
}
```

但是这样做管道的发送和接收是一一对应的，如果创建与goroutine数量一致的channel，开销太大了，该怎么办呢？

## 多个协程退出

我们可以采用管道的close函数，来广播关闭管道的通知。所有从关闭管道接收的操作均会收到一个零值和一个可选的失败标志。

```
package main
import (
   ."fmt"
   "time"
)

func worker(ch chan int){
   for{
      select {
      default:
         Println("Working")
      case <-ch:
         Println("Exiting")
         return
      }
   }

}

func main(){
   ch := make(chan int)
   for i :=0; i< 5; i++{
      go worker(ch)
   }
   time.Sleep(time.Second)
   close(ch)
   for{}
}
```
可以看到，程序接收到了close的通知：
```
Working
Working
Working
Working
Exiting
Working
Working
Exiting
Working
Exiting
Working
Exiting
Exiting
```

## sycn.WaitGroup退出

不过这个程序依然不够健壮，每个goroutine程序收到close操作，希望做一些清理的工作，但是清理的工作不一定保证被完成，因为main线程没有等待goroutine线程退出的机制。这时，我们可以利用sync.WaitGroup来进行改进：

```
package main
import (
   ."fmt"
   "time"
   "sync"
)

func DoSth(){
   Println("Doing Cleaning")
}

func worker(wg *sync.WaitGroup, ch chan int){
   defer wg.Done()
   for{
      select {
      default:
         Println("Working")
      case <-ch:
         Println("Exiting")
         DoSth()
         return
      }
   }

}

func main(){
   var wg sync.WaitGroup
   ch := make(chan int)
   for i :=0; i< 5; i++{
      wg.Add(1)
      go worker(&wg, ch)
   }
   time.Sleep(time.Second)
   close(ch)
   wg.Wait()
}
```

现在每个工作goroutine的创建、运行、暂停、退出，都在main函数的**安全控制**之下了。

## context包退出

在Go1.7发布时，标准库增加了一个context包，用来简化对于处理单个请求的多个Goroutine之间与请求域的数据、超时和退出等操作，官方有博文对此做了专门介绍。我们使用context包，对上面的退出机制进行重新实现

```
package main
import (
   ."fmt"
   "time"
   "sync"
   "context"
)

func DoSth(){
   Println("Doing Cleaning")
}

func worker(ctx context.Context, wg *sync.WaitGroup) error{
   defer wg.Done()
   for{
      select {
      default:
         Println("Working")
      case <-ctx.Done():
         Println("Exiting")
         return ctx.Err()
      }
   }

}

func main(){
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
   var wg sync.WaitGroup

   for i :=0; i< 5; i++{
      wg.Add(1)
      go worker(ctx, &wg)
   }
   time.Sleep(time.Second)
   cancel()
   wg.Wait()
}
```

# Goroutine泄漏
Go语言是带内存自动回收特性的，因此内存一般不会泄漏。在前面素数筛的例子中，GenerateNatural和PrimeFilter函数内部都启动了新的Goroutine，当main函数不再使用管道时后台Goroutine有泄漏的风险。我们可以通过context包来避免这个问题，下面是改进的素数筛实现：

```
package main

import (
	"context"
	. "fmt"
	"time"
)

func generatePrime(ctx context.Context) chan int {
	c := make(chan int)
	go func() {
		for i := 2; ; i++ {
			select {
			case c <- i:
				//
			case <-ctx.Done():
				return
			}
			c <- i
		}
	}()
	return c
}

func filter(ctx context.Context, c chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for {
			select {
			case i := <-c:
				{
					if i%prime != 0 {
						out <- i
					}
				}
			case <-ctx.Done():
				return
			}

		}
	}()
	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	ch := generatePrime(ctx)
	for i := 0; i < 10; i++ {
		prime := <-ch
		Println("第", i, "个素数为：", prime)
		ch = filter(ctx, ch, prime)
	}
	cancel()
}
```

程序结束时，通过cancel来通知goroutine退出，这样就避免了goroutine的泄漏。
