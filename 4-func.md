# 函数

函数对应操作序列，是程序的基本组成元素。在Go语言中函数是第一类对象，可以将函数保存在变量中。

## 函数声明
包括：
- 函数名
- 形参列表
- 返回值列表（可省略）
- 函数体
  
```
func f(i int, j string){

}
```
实参通过传值方式传递给形参，因此对形参的修改一般情况下不会影响实参。但是，如果实参包括引用类型，比如指针,slice,map,function,channel等类型，实参可能会由于参数的间接引用被修改。

## 递归

Go语言支持递归，也就是函数自己调用自己。在Go语言中，使用的是可变栈，不同与大多是使用固定大小的函数调用栈64k或者2M不等，Go语言栈的大小按需增加，这使得我们再使用递归时不必太考虑内存溢出和安全问题。但是也不能超过**1G**(64位系统）。

```
package main
import "fmt"

func Fibonach(n int) int{
   if n == 1 || n == 2 {
      return 1
   }
   return Fibonach(n-1) + Fibonach(n-2)
}

func main() {
   for i := 1; i < 10; i++ {
      fmt.Println(Fibonach(i))
   }
}
```

在Go1.4以前，Go的动态栈采用的是分段式的动态栈，通俗地说就是采用一个链表来实现动态栈，每个链表的节点内存位置不会发生变化。但是链表实现的动态栈对某些导致跨越链表不同节点的热点调用的性能影响较大，因为相邻的链表节点它们在内存位置一般不是相邻的，这会增加CPU高速缓存命中失败的几率。为了解决热点调用的CPU缓存命中率问题，Go1.4之后改用连续的动态栈实现，也就是采用一个类似动态数组的结构来表示栈。不过连续动态栈也带来了新的问题：当连续栈动态增长时，需要将之前的数据移动到新的内存空间，这会导致之前栈中全部变量的地址发生变化。

> 虽然Go语言运行时会自动更新引用了地址变化的栈变量的指针，但最重要的一点是要明白**Go语言中指针不再是固定不变的了**（因此不能随意将指针保持到数值变量中，Go语言的地址也不能随意保存到不在GC控制的环境中，因此使用CGO时不能在C语言中长期持有Go语言对象的地址）。

因为，Go语言函数的栈会自动调整大小，所以普通Go程序员已经很少需要关心栈的运行机制的。在Go语言规范中甚至故意没有讲到栈和堆的概念。我们无法知道函数参数或局部变量到底是保存在栈中还是堆中，我们只需要知道它们能够正常工作就可以了。看看下面这个例子：

```
func f(x int) *int {
    return &x
}

func g() int {
    x = new(int)
    return *x
}
```

第一个函数直接返回了函数参数变量的地址——这似乎是不可以的，因为如果参数变量在栈上的话，函数返回之后栈变量就失效了，返回的地址自然也应该失效了。但是Go语言的编译器和运行时比我们聪明的多，它会保证指针指向的变量在合适的地方。第二个函数，内部虽然调用new函数创建了*int类型的指针对象，但是依然不知道它具体保存在哪里。

> 对于有C/C++编程经验的程序员需要强调的是：不用关心Go语言中函数栈和堆的问题，编译器和运行时会帮我们搞定；同样不要假设变量在内存中的位置是固定不变的，指针随时可能会变化，特别是在你不期望它变化的时候。

## 多返回值

在Go中，一个函数可以由多可返回值。很多标准库中的函数返回两个值，一个是期望得到的值，一个是函数出错时的信息。

多用多返回值的函数式，返回给调用者的每个值，都必须显示调用。如果某个值不需要使用，可以将它分配给空标志_。

```
package main
import "fmt"

func Divide(a, b int) (float32, error) {
   if b == 0 {
      return 0.0, fmt.Errorf("分母不能为0")
   }
   return float32(a)/float32(b),nil
}

func main() {
   res, err := Divide(1, 0)
   fmt.Println(res, err)
   res, err = Divide(1, 2)
   res, _ = Divide(1, 0)   //不使用error参数
}
```
## 函数作为值

在Go中，函数作为first-class values:和其他值一样，有类型，可以复制给其他变量，传递给函数，从函数返回。对函数值(function value)的调用与函数调用类似.栗子如下，函数作为参数值，实现过滤；函数作为返回值，实现柯里化curryAdd(1)(2):

```
package main
import "fmt"

//函数作为参数值，实现过滤
type test_isEven func(int) bool
func isEven(i int) bool{
   if i%2 == 0 {
      return true
   }
   return false
}

func Filter(slice []int, f test_isEven) []int{
   res := slice[:0]
   for _, value := range(slice){
      if f(value){
         res = append(res, value)
      }
   }
   return res
}

### 柯里化
//函数作为返回值，实现柯里化
func curryAdd(int) func(int) int{
   var x int
   return func(y int) int{
      x += y
      return x
   }
}
func main() {
   a := []int{1,2,3,4,5,6,7}
   fmt.Println(Filter(a, isEven))
   fmt.Println(curryAdd(1)(2))
}
```

## 可变参数
函数参数可变的函数成为可变参数函数。声明可变函数的参数，要在最后一个参数类型之前加上省略号..."，以表示该函数可接受人类数量的该类型参数。

```
package main
import "fmt"

// 可变数量的参数
// more 对应 []int 切片类型
func Sum(a int, more ...int) int {
   for _, v := range more {
       a += v
   }
   return a
}

func main() {
   a := []int{1,2,3}
   fmt.Println(Sum(0,a...))
}
```
当可变参数是接口类型时，由于编译时不会检查参数类型，参数传入是是否解构，结果可能完全不同。

```
func main() {
    var a = []interface{}{123, "abc"}

    Print(a...) // 123 abc
    Print(a)    // [123 abc]
}

func Print(a ...interface{}) {
    fmt.Println(a...)
}
```

## 匿名函数
Go语言的函数分为：
- 具名函数:一般对应包级函数，是匿名函数的特例
- 匿名函数：匿名函数引用了外部作用域的变量就成了闭包函数，闭包函数式函数式变成的核心。

```
package main
import "fmt"
//具名函数
func Add(a, b int) int{
   return a + b
}

//匿名函数
var add = func(a, b int) int{
   return a + b
}

func main() {
   fmt.Println(Add(1,2), add(1,2))
}
```

## 给返回值命名

不仅函数的参数可以有名字，也可以给函数的返回值命名：
如果返回值命名了，可以通过名字来修改返回值，也可以通过defer语句在return语句之后修改返回值,结果为43：
```
package main
import "fmt"

func Inc() (v int) {
   defer func(){ v++ } ()
   return 42
}

func main() {
   v := Inc()
   fmt.Println(v)
}
```

## 闭包

在Go中，匿名函数引用了外部函数的局部变量，这种函数值我们成为闭包，比如上面的Inc函数中的defer引用了外部的v。

闭包带来了隐含的问题:

```
package main
import "fmt"

var Println = fmt.Println
func main() {
   for i := 0 ; i < 3; i ++{
      defer func(){Println(i)}()
   }
}
```
上面函数输出结果都是3，这是因为延迟语句引用的是同一个迭代变量i，循环结束后的值变为3，所以最终输出为3.

解决办法是为每轮迭代的defer生产独有的变量。解决办法有两种：

```
package main
import "fmt"

var Println = fmt.Println
func main() {
   for i := 0 ; i < 3; i ++{
      i := i                  //定义局部变量
      defer func(){Println(i)}()
   }
   //输出 2 1 0
   for i := 0 ; i < 3; i ++{                 
      defer func(i int){Println(i)}(i)
   }
   //输出 2 1 0

}
```

第一种方法是在循环体内部再定义一个局部变量，这样每次迭代defer语句的闭包函数捕获的都是不同的变量，这些变量的值对应迭代时的值。第二种方式是将迭代变量通过闭包函数的参数传入，defer语句会马上对调用参数求值。两种方式都是可以工作的。不过一般来说,在for循环内部执行defer语句并不是一个好的习惯，此处仅为示例，不建议使用。

## 切片作为函数参数

切片作为函数参数的时候，如果以切片为参数调用函数时，有时候会给人一种参数采用了传引用的方式的假象：因为在被调用函数内部可以修改传入的切片的元素。其实，任何可以通过函数参数修改调用参数的情形，都是因为函数参数中显式或隐式传入了指针参数。函数参数传值的规范更准确说是只针对数据结构中固定的部分传值，例如字符串或切片对应结构体中的指针和字符串长度结构体传值，但是并不包含指针间接指向的内容。将切片类型的参数替换为类似reflect.SliceHeader结构体就很好理解切片传值的含义了：

```
package main
import "fmt"

var Println = fmt.Println
func twice(x []int) {
   for i := range x {
       x[i] *= 2
   }
}

type IntSliceHeader struct {
   Data []int
   Len  int
   Cap  int
}

func twice2(x IntSliceHeader) {
   for i := 0; i < x.Len; i++ {
       x.Data[i] *= 2
   }
}
func main() {
   a := []int{1,2,3}
   twice(a)
   Println(a)
   b :=  IntSliceHeader{[]int{1,2,3},3,8}
   twice2(b)
   Println(b)
}
```

