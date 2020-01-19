# 数组、字符串和切片

Go语言中数组、字符串和切片是三种密切相关且使用的最为频繁的数据结构，只有不能满足需求时才会考虑更加复杂的自定义数据结构。
- 数组是一种值类型，虽然数组的元素可以修改，但数组本身的赋值和函数传参是通过**整体赋值的方式**也就是**传值**，而不是引用传递。
- 字符串对应的也是字节数组，只是它是只读的，不能像数组可以修改；字符串的赋值只复制数据地址和长度，不会导致底层数据的复制，也是**传值方式**
- 切片的结构类似于字符串，但是更为灵活。它的底层数据也是对应数据类型的数组，但是每个切片有独立的长度和容量信息，并且切片的赋值和传参都是采用**传值方式**。

## 数组

### 定义

数组的元数据包括：
- 长度
- 0个或多个元素
- 类型 (可以使字符串、结构体、函数、接口、管道等)
  
不同长度或者类型的数组无法直接赋值。

可以由四种定义方式：
直接定义：
```
var a[3]init
var aa[3]string
//0,0,0
```
初始化赋值，按顺序指定初始化的值，长度自动计算
```
var b = [...]int{1,2,3}
var bb = [2]string("hello", "你好") //支持中文
//1,2,3
```
索引方式，索引与值对应，长度以最大索引为准，其他的为0
```
var c = [...]int{2:3, 1:2}
//0,2,3
```
混合方式,前两个元素按顺序初始化，第三四个初始化为0，第五个元素索引初始化为5，最后一个元素接着按顺序初始化：
```
var d = [...]int{1,2,4:5,6}
//1,2,0,0,5,6
```
从下图可以看出，数组数据在内存中连续存储(每个int占8位）：
![数组存储](/.gitbook/assets/数组的存储.png)

数组的名字代表的就是整个数组，而不是像c代表第一个袁术的指针。因此当数组被赋值或者被值传递的时候，实际上会复制整个数据。如果数组较大会产生比较大的开销，为了避免复制数组带来的开销，可以采用指针传递

```
package main
import "fmt"
func main() {
   var a = [...]int{1,2}
   var b = a
   var c = &a
   fmt.Println(a,b,c)  
   c[0] = 3
   fmt.Println(a,b,c)  //b的值没有发生变化，可见赋值是整体复制
}
```

### 遍历

可以使用fmt.Printf函数来打印数组的类型和详细信息，同时可用下述三种for循环来遍历：
```
package main
import "fmt"
func main() {
   var a = [...]int{1,2}
   for i,v := range(a){
      fmt.Println(i,a[i],v)
   } //i代表索引，v代表值
   for _, v := range(a){
      fmt.Println(v)
   } //_不使用索引，v代表值
   for i := 0; i < len(a); i++ {
      fmt.Println(a[i])
   }
   fmt.Printf("a: %T\n", a)
   fmt.Printf("a: %v\n", a)
}
```

# 字符串

## 内存结构

字符串是一个不可改变的字节序列，且单个元素不可修改，可以认为是一个只读的字节数组。Go语言采用UTF8编码字符串，也可以采用GBK等非UTF8编码的表示字符串，但是不支持用for range遍历非UTF8编码的字符串。
字符串数据结构包含两部分：
- Data uintptr:字符串指向的底层字节数组
- Len int:字符串的字节的长度

可以看下字符串“Hello, world"本身对应的内存结构

![字符串内存结构](.gitbook/assets/字符串内存组织.png)

可以发现，"Hello, world"的底层数据和以下数组是一样的:

```
var data = [...]byte{
   'h','e','l','o',',',' ','w','o','r','l','d'
}
```
## 访问和遍历

字符串虽然不是切片，但是支持切片操作，不同位置的切片底层对应的是同一块内存数据.但是赋值的时候，采用的是复制，而不是引用复制，举例：
```
package main
import "fmt"
func main() {
   s := "hello, world"
   hello := s[:5]
   world := s[7:]
   fmt.Println(s,hello, world)
   hello = "heelo"
   fmt.Println(s,hello, world)
   s = "hemlo, world"
   fmt.Println(s,hello, world)
}
```
可以看出hello改变后，原数组s并未改变，且s改变后，也米有影响hello。同时也可以看出可以采用Println来直接打印

可以通过len函数返回字符串的长度

```
s := "h"
fmt.Println(len(s))
s = "陈"
fmt.Println(len(s))
fmt.Printf("%#v\n",[]byte("陈"))
```
可以看出英文字符占据1个字节，中文字符占据3个字节，与UTF8编码格式一致。并且可以看到底层的中文数据“陈”对应的16进制数据为```[]byte{0xe9, 0x99, 0x88}```

如果想遍历数组的单个袁术可以采用range或者传统的下标点方式,不过range遍历输出的是int32型的编码，需要通过string转化为Unicode字符。
```
package main
import "fmt"
func main() {
   s := "hello, world, 陈"
   for _, c:= range (s){
      fmt.Println(c, string(c))
   }
   fmt.Println("---------------------------")
   for _, c:= range []byte(s){
      fmt.Printf("%v\n",c)
   }
}
```

# 切片slice

## 内存结构

简单的书，切片是一种简化版的动态数组。数组的类型和操作不够灵活，因此切片使用的更为广泛，理解其原理和用法就比较重要。

切片的结构定义, reflect.SliceHeader:
```
type SliceHeader struct{
   Data unitptr
   Len  int
   Cap  int
}
```
可以看出切片的开头部分和Go字符串一样，但是多了Cap表示切片指向的内存空间的最大元素个数容量（不是字节数）。

下图是```x := []int{2,3,5,7,11} 和 y := x[1:3]```两个切片对应的内存结构。

![切片的内存结构](/.gitbook/assets/切片的内存结构.png)

## 定义

```
package main
import "fmt"
func main() {
   var (
      a []int               // nil切片, 和 nil 相等, 一般用来表示一个不存在的切片
      b = []int{}           // 空切片, 和 nil 不相等, 一般用来表示一个空的集合
      c = []int{1, 2, 3}    // 有3个元素的切片, len和cap都为3
      d = c[:2]             // 有2个元素的切片, len为2, cap为3
      e = c[0:2:cap(c)]     // 有2个元素的切片, len为2, cap为3
      f = c[:0]             // 有0个元素的切片, len为0, cap为3
      g = make([]int, 3)    // 有3个元素的切片, len和cap都为3
      h = make([]int, 2, 3) // 有2个元素的切片, len为2, cap为3
      i = make([]int, 0, 3) // 有0个元素的切片, len为0, cap为3
   )
   fmt.Println(a,b,c,d,e,f,g,h,i)
   d = []int{4,4,4}
   fmt.Println(c,d)
}
```
可以看出，改变d之后,c并没有变化,切片的也是采用的赋值，而没有改变底层的数据。

## 遍历

与数组遍历类似

```
package main
import "fmt"
func main() {
   var a = []int{1,2,3,4}
   for i := range(a){
      fmt.Println(i, a[i])
   }
   for i, v := range(a){
      fmt.Println(i, v)
   }
   for i := 0; i < len(a); i++ {
      fmt.Println(i, a[i])
   }
}
```

## 添加、插入和删除元素

### 追加

可以在头部或者尾部添加，如下例
```
package main
import "fmt"
func main() {
   //头部添加
   var a []int
   fmt.Println(len(a),cap(a))
   a = append(a, 1)
   a = append(a, 2, 3)             //追加多个元素
   a = append(a, []int{4,5,6}...)  //追加一个切片，切片需要用...解包
   fmt.Println(a)
   fmt.Println(len(a),cap(a))
   //尾部添加
   a = append([]int{-1,0}, a...)
   fmt.Println(a)
   fmt.Println(len(a),cap(a))
   a = append(a,7)
   fmt.Println(len(a),cap(a))
}
```
但是再尾部添加需要重新分配内存，且导致已有的数据重新复制一下。
并且我们可以看出，当容量不够时，会动态增加容量，对于int型在这个例子中一次增加的容量为8（据资料实际上只2的指数次增长）.

### 插入：

采用copy和append组合来避免创建中间的临时切片，完成插入:

```
package main
import "fmt"
func main() {
   //头部添加
   var a = []int{1,2,3}
   insert_slice(&a,-1,4)
   fmt.Println(a)
}

//i为插入的位置，v为插入的值
func insert_slice(x *[]int , i int, v int){
   if i > len(*x) || i < 0 {
      fmt.Println("插入位置越界")
      return 
   }
   *x = append(*x, 0)
   copy((*x)[i+1:], (*x)[i:])
   (*x)[i] = v
}
```
读者可以思考如何插入一个切片？🤔

### 删除切片

#### 删除尾部

```
package main
import "fmt"
func main() {
   var a = []int{1,2,3,4,5,6,7}
   a = a[:len(a)-1]  //删除尾部
   fmt.Println(a)
   N := 2
   a = a[:len(a)-N]  //删除尾部N个
   fmt.Println(a)
}
```

#### 删除头部

```
package main
import "fmt"
func main() {
   var a = []int{1,2,3,4,5,6,7}
   a = a[1:]  //删除头部
   fmt.Println(a)
   N := 2
   a = a[N:]  //删除头部N个
   fmt.Println(a)
}
```
上面的方法实际上对数据重新复制，指针会发生变化；当然也可以用append操作，不移动数据指针，只讲数据往前移动

```
package main
import "fmt"
func main() {
   var a = []int{1,2,3,4,5,6,7}
   a = append(a[:0], a[1:]...)  //删除第一个
   fmt.Println(a)
   N := 2
   a = append(a[:0], a[N:]...)  //删除头部N个
   fmt.Println(a)
}
```

####中间删除

也可以用append完成

```
package main
import "fmt"
func main() {
   var a = []int{1,2,3,4,5,6,7}
   delete_slice(&a, 8, 1)
   fmt.Println(a)
}

func delete_slice( a *[]int , i int, N int){
   if i < 0 || i > len(*a){
      fmt.Println("删除位置越界")
      return
   }
   if i+N > len(*a){
      fmt.Println("删除长度越界")
      return
   }
   *a = append((*a)[:i],(*a)[i+N:]...)
}
```

## 一些切片技巧Trick

对于切片而言，len为0但是cap不为0的切片是很有用的特性，可以简洁高效的实现一些功能：

### 利用0长切片实现元素的删除

```
package main
import "fmt"
func main() {
   s := []byte("Hell o, wor ld")
   fmt.Println(string(s))
   a := TrimSpace(s)
   fmt.Println(string(a))
}

func TrimSpace(s []byte) []byte{
   b := s[:0]
   for _, v := range(s){
      if v != ' '{
         b = append(b, v)
      }
   }
   return b
}
```

### 根据过滤条件原地删除切片元素

```
package main
import "fmt"
func main() {
   s := []byte("Hell o, wor ld")
   fmt.Println(string(s))
   a := Filter(s, func(x byte) bool{
      if x == 'o'{
         return true
      }
      return false
   })
   fmt.Println(string(a))
}

func Filter(s []byte, fn func(x byte) bool) []byte{
   b := s[:0]
   for _, v := range(s){
      if !fn(v){
         b = append(b, v)
      }
   }
   return b
}
```

## 避免切片内存泄漏

如前面所说，切片操作并不会复制底层的数据。底层的数组会被保存在内存中，直到它不再被引用。但是有时候可能会因为一个小的内存引用而导致底层整个数组处于被使用的状态，这会延迟自动内存回收器对底层数组的回收。

例如，FindPhoneNumber函数加载整个文件到内存，然后搜索第一个出现的电话号码，最后结果以切片方式返回。

```
func FindPhoneNumber(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    return regexp.MustCompile("[0-9]+").Find(b)
}
```

这段代码返回的[]byte指向保存整个文件的数组。因为切片引用了整个原始数组，导致自动垃圾回收器不能及时释放底层数组的空间。一个小的需求可能导致需要长时间保存整个文件数据。这虽然这并不是传统意义上的内存泄漏，但是可能会拖慢系统的整体性能。

要修复这个问题，可以将感兴趣的数据复制到一个新的切片中（数据的传值是Go语言编程的一个哲学，虽然传值有一定的代价，但是换取的好处是切断了对原始数据的依赖）：

```
func FindPhoneNumber(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    b = regexp.MustCompile("[0-9]+").Find(b)
    return append([]byte{}, b...)
}
```

类似的问题，在删除切片元素时可能会遇到。假设切片里存放的是指针对象，那么下面删除末尾的元素后，被删除的元素依然被切片底层数组引用，从而导致不能及时被自动垃圾回收器回收（这要依赖回收器的实现方式）：

```
var a []*int{ ... }
a = a[:len(a)-1]    // 被删除的最后一个元素依然被引用, 可能导致GC操作被阻碍
保险的方式是先将需要自动内存回收的元素设置为nil，保证自动回收器可以发现需要回收的对象，然后再进行切片的删除操作：
```

```
var a []*int{ ... }
a[len(a)-1] = nil // GC回收最后一个元素内存
a = a[:len(a)-1]  // 从切片删除最后一个元素
```

当然，如果切片存在的周期很短的话，可以不用刻意处理这个问题。因为如果切片本身已经可以被GC回收的话，切片对应的每个元素自然也就是可以被回收的了。

## 切片类型强制转换

为了安全，当两个切片类型[]T和[]Y的底层原始切片类型不同时，Go语言是无法直接转换类型的。不过安全都是有一定代价的，有时候这种转换是有它的价值的——可以简化编码或者是提升代码的性能。比如在64位系统上，需要对一个[]float64切片进行高速排序，我们可以将它强制转为[]int整数切片，然后以整数的方式进行排序（因为float64遵循IEEE754浮点数标准特性，当浮点数有序时对应的整数也必然是有序的）。

下面的代码通过两种方法将[]float64类型的切片转换为[]int类型的切片：

```
// +build amd64 arm64

import "sort"

var a = []float64{4, 2, 5, 7, 2, 1, 88, 1}

func SortFloat64FastV1(a []float64) {
    // 强制类型转换
    var b []int = ((*[1 << 20]int)(unsafe.Pointer(&a[0])))[:len(a):cap(a)]

    // 以int方式给float64排序
    sort.Ints(b)
}

func SortFloat64FastV2(a []float64) {
    // 通过 reflect.SliceHeader 更新切片头部信息实现转换
    var c []int
    aHdr := (*reflect.SliceHeader)(unsafe.Pointer(&a))
    cHdr := (*reflect.SliceHeader)(unsafe.Pointer(&c))
    *cHdr = *aHdr

    // 以int方式给float64排序
    sort.Ints(c)
}
```

第一种强制转换是先将切片数据的开始地址转换为一个较大的数组的指针，然后对数组指针对应的数组重新做切片操作。中间需要unsafe.Pointer来连接两个不同类型的指针传递。需要注意的是，Go语言实现中非0大小数组的长度不得超过2GB，因此需要针对数组元素的类型大小计算数组的最大长度范围（[]uint8最大2GB，[]uint16最大1GB，以此类推，但是[]struct{}数组的长度可以超过2GB）。

第二种转换操作是分别取到两个不同类型的切片头信息指针，任何类型的切片头部信息底层都是对应reflect.SliceHeader结构，然后通过更新结构体方式来更新切片信息，从而实现a对应的[]float64切片到c对应的[]int类型切片的转换。

通过基准测试，我们可以发现用sort.Ints对转换后的[]int排序的性能要比用sort.Float64s排序的性能好一点。不过需要注意的是，这个方法可行的前提是要保证[]float64中没有NaN和Inf等非规范的浮点数（因为浮点数中NaN不可排序，正0和负0相等，但是整数中没有这类情形）。





