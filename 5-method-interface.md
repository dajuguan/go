# 面向对象-方法与接口

方法是面向对象的一个特性，在c++中方法对应一个类对象的成员函数；在Go语言中方法确是关联到**类型**的，这样在编译阶段就可以完成方法的静态绑定。

一个面向对象的程序会用方法来表达其属性对应的操作，而不用直接操作其对象。

# 方法实现

下面我们实现一组C语言风格的File函数：

```
// 文件对象
type File struct {
    fd int
}

// 打开文件
func OpenFile(name string) (f *File, err error) {
    // ...
}

// 关闭文件
func CloseFile(f *File) error {
    // ...
}

// 读文件数据
func ReadFile(f *File, int64 offset, data []byte) int {
    // ...
}
```
其中OpenFile类似构造函数用于打开文件对象，CloseFile类似析构函数用于关闭文件对象，ReadFile则类似普通的成员函数，这三个函数都是普通的函数。CloseFile和ReadFile作为普通函数，需要占用包级空间中的名字资源。不过CloseFile和ReadFile函数只是针对File类型对象的操作，这时候我们更希望这类函数和操作对象的类型紧密绑定在一起。

Go语言中的做法是，将CloseFile和ReadFile函数的第一个参数移动到函数名的开头：

```
package main
type File struct {
   fd int
}

// 关闭文件
func (f *File) CloseFile()  {
   // ...
}

// 读文件数据
func (f *File) ReadFile(offset int64 , data []byte)  {
   // ...
}

func ReadFile(int){  //可以看出File不需要占用名字空间
}
func CloseFile(int){
}

func main() {
}

```

这样的话，CloseFile和ReadFile函数就成了File类型独有的方法了（而不是File对象方法）。它们也不再占用包级空间中的名字资源，同时File类型已经明确了它们操作对象，因此方法名字一般简化为Close和Read：

```
// 关闭文件
func (f *File) CloseFile() error {
    // ...
}

// 读文件数据
func (f *File) ReadFile(int64 offset, data []byte) int {
    // ...
}
```

从代码上只是改变了参数的位置，但是从编程范式来看，Go已经进入了面向对象的行列。可以为任何类型添加一个或多个方法。每种类型对应的方法必须和类型在同一个包中，因此无法给int这种内置类型添加方法（方法的定义和类型的定义不在一个包中）。对于给定的类型，方法的名字必须是唯一的，方法和函数也*不支持重载*。

# 方法表达式

方法是由函数演变而来，只是将函数的第一个对象参数移动到了函数名前面了而已。因此我们依然可以按照原始的过程式思维来使用方法。通过方法表达式的特性可以将方法还原为普通类型的函数：

```
// 不依赖具体的文件对象
// func CloseFile(f *File) error
var CloseFile = (*File).Close

// 不依赖具体的文件对象
// func ReadFile(f *File, int64 offset, data []byte) int
var ReadFile = (*File).Read

// 文件处理
f, _ := OpenFile("foo.dat")
ReadFile(f, 0, data)
CloseFile(f)
```

# 方法值

在有些场景更关心一组相似的操作：比如Read读取一些数组，然后调用Close关闭。此时的环境中，用户并不关心操作对象的类型，只要能满足通用的Read和Close行为就可以了。不过在方法表达式中，因为得到的ReadFile和CloseFile函数参数中含有File这个特有的类型参数，这使得File相关的方法无法和其它不是File类型但是有着相同Read和Close方法的对象无缝适配。这种小困难难不倒我们Go语言码农，我们可以通过结合闭包特性来消除方法表达式中第一个参数类型的差异：

```
// 先打开文件对象
f, _ := OpenFile("foo.dat")

// 绑定到了 f 对象
// func Close() error
var Close = func Close() error {
    return (*File).Close(f)
}

// 绑定到了 f 对象
// func Read(int64 offset, data []byte) int
var Read = func Read(int64 offset, data []byte) int {
    return (*File).Read(f, offset, data)
}

// 文件处理
Read(0, data)
Close()
```

这刚好是方法值也要解决的问题。我们用方法值特性可以简化实现：

```
// 先打开文件对象
f, _ := OpenFile("foo.dat")

// 方法值: 绑定到了 f 对象
// func Close() error
var Close = f.Close

// 方法值: 绑定到了 f 对象
// func Read(int64 offset, data []byte) int
var Read = f.Read

// 文件处理
Read(0, data)
Close()
```

# 继承

Go语言不支持传统面向对象的继承特性，而是通过在结构体中内置匿名的成员来实现继承：

```
package main
import (
   "image/color"
   "fmt"
)

type Point struct{ X, Y float64 }

type ColoredPoint struct {
    Point
    Color color.RGBA
}

func main() {
   var cp ColoredPoint
   cp.X = 1
   fmt.Println(cp.Point.X) // "1"
   cp.Point.Y = 2
   fmt.Println(cp.Y)       // "2"
}
```

虽然我们可以将ColoredPoint定义为一个有三个字段的扁平结构的结构体，但是我们这里将Point嵌入到ColoredPoint来提供X和Y这两个字段。

通过嵌入匿名的成员，我们不仅可以继承匿名成员的内部成员，而且可以继承匿名成员类型所对应的方法。我们一般会将Point看作基类，把ColoredPoint看作是它的继承类或子类。不过这种方式继承的方法并不能实现C++中虚函数的多态特性。所有继承来的方法的接收者参数依然是那个匿名成员本身，而不是当前的变量。

```
type Cache struct {
    m map[string]string
    sync.Mutex
}

func (p *Cache) Lookup(key string) string {
    p.Lock()
    defer p.Unlock()

    return p.m[key]
}
```

Cache结构体类型通过嵌入一个匿名的sync.Mutex来继承它的Lock和Unlock方法. 但是在调用p.Lock()和p.Unlock()时, p并不是Lock和Unlock方法的真正接收者, 而是会将它们展开为p.Mutex.Lock()和p.Mutex.Unlock()调用. 这种展开是编译期完成的, 并没有运行时代价.

# 接口与多态

一般静态编程语言都有过于严格的类型系统，以使得编译器可以深入的检查程序是否有明显的错误；但是过于严格的类型系统，却使得变成过于复杂，程序猿成了编译系统的接锅侠。

Go语言试图在安全和灵活的编程之间取得一个平衡，它在提供类型检查的同时，还通过接口类型实现了对鸭子类型的支持，使得**安全动态**的编程变得相对容易。

## 接口约定

接口把具有共性的方法定义在一起，任何其他类型只要实现了这些方法就是实现了这个接口。
接口是抽象的类型。
- 它不暴露它所代表的的对象的内部值或对象支持的操作等实现细节，对象更具灵活性和适应性；
- 在使用第三方包时，为了不破坏这些类型的原有定义，可以创建新的接口类型满足已经存在的类型
- 接口中的方法必须全部实现，才能实现接口
- 任何提供了改接口实现代码的类型都隐式的实现了该接口，而不用显示声明

举个例子，在游戏中每个人有枪和弓箭这两种武器，可以切换武器来进行射击，这个时候可以把武器作为一种接口类型，来进行设计：

```
//武器接口类，可以执行射击动作
//test.go
package test

type Weapon interface{
    Shot()
}
```

```
package main
import (
   "fmt"
   "test"
)

//武器结构体，不同的武器有不同的名字，但是都可以射击
type weapon struct{
   name string
   test.Weapon
}

func (w weapon) Shot(){
   fmt.Printf(" use a %s shot\n", w.name)
}

type men struct{
   name string
   wp    weapon
}

func (s *men)ChangeWeapon(wp weapon){
   s.wp = wp
}

func (s men)shot(){
   fmt.Printf(s.name)
   s.wp.Shot()
}


func main(){
   bob := men{name:"Bob"}
   gun := weapon{name:"gun"}
   arrow := weapon{name:"arrow"}
   bob.ChangeWeapon(gun)
   bob.shot()
   bob.ChangeWeapon(arrow)
   bob.shot()  

   tom := men{name:"Tom"}
   tom.ChangeWeapon(gun)
   tom.shot()
}
```


接口在Go中无处不在，以fmt的Printf为例，它除了内置的类型外，还可以向任何自定义的输出流打印，可以是文件、标准输甚至是网络。它的实现基于fmt包中Fprintf这个函数

```
package fmt

func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)
func Printf(format string, args ...interface{}) (int, error) {
    return Fprintf(os.Stdout, format, args...)
}
```

其中io.Writer用于输出的接口，error是内置的错误接口，它们的定义如下：

```
type io.Writer interface {
    Write(p []byte) (n int, err error)
}

type error interface {
    Error() string
}
```

我们可以自定义一个对象，来实现字符串输出时全部变为小写：
```
package main

import (
   "fmt"
   "io"
   "bytes"
   "os"
)

type LowerWriter struct{
   io.Writer
}
func (p *LowerWriter) Write(data []byte)(n int, err error){
   return p.Writer.Write(bytes.ToLower(data))
}
func main() {
   s := LowerWriter{os.Stdout}
   fmt.Fprintln(&s, "Hello,World")
}
```

## 接口组合

接口之间还可以通过已有的接口，来定义新的接口，比如:

```
package main
import "fmt"

type Car interface{
   Drive()
}

type Plane interface{
   Fly()
}

type FlyingCar interface{
   Car
   Plane
}
type MyFlyingCar struct{
   FlyingCar
}
func (this MyFlyingCar) Drive(){
   fmt.Println("I can drive")
}

func (this MyFlyingCar) Fly(){
   fmt.Println("I can fly")
}


func main(){
   var mycar FlyingCar
   mycar = new(MyFlyingCar)
   mycar.Drive()
   mycar.Fly()
}
```
> 需要注意的是，编译的时候并不会检查时候完全实现了所有接口，如果调用没有实现的接口，程序就会报错！

## 接口转换

在Go语言，不同基础类型之间是不支持隐式转换的，比如不能把int类型赋值给int64类型的变量。但是接口之间的转换非常的灵活，可以隐式的进行接口与对象，接口与接口之间的转换。

```
var (
    a io.ReadCloser = (*os.File)(f) // 隐式转换, *os.File 满足 io.ReadCloser 接口
    b io.Reader     = a             // 隐式转换, io.ReadCloser 满足 io.Reader 接口
    c io.Closer     = a             // 隐式转换, io.ReadCloser 满足 io.Closer 接口
    d io.Reader     = c.(io.Reader) // 显式转换, io.Closer 不满足 io.Reader 接口
)
```

有时候对象和接口之间太灵活了，导致我们需要人为地限制这种无意之间的适配。常见的做法是定义一个含特殊方法来区分接口。比如runtime包中的Error接口就定义了一个特有的RuntimeError方法，用于避免其它类型无意中适配了该接口：

```
type runtime.Error interface {
    error

    // RuntimeError is a no-op function but
    // serves to distinguish types that are run time
    // errors from ordinary errors: a type is a
    // run time error if it has a RuntimeError method.
    RuntimeError()
}
```

不过这种做法只是君子协定，如果有人刻意伪造一个runtime.Error接口也是很容易的。再严格一点的做法是给接口定义一个私有方法。只有满足了这个私有方法的对象才可能满足这个接口，而私有方法的名字是包含包的绝对路径名的，因此只能在包内部实现这个私有方法才能满足这个接口。测试包中的testing.TB接口就是采用类似的技术：

```
type testing.TB interface {
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    ...

    // A private method to prevent users implementing the
    // interface and so future additions to it will not
    // violate Go 1 compatibility.
    private()
}
```

不过这种通过私有方法禁止外部对象实现接口的做法也是有代价的：首先是这个接口只能包内部使用，外部包正常情况下是无法直接创建满足该接口对象的；其次，这种防护措施也不是绝对的，恶意的用户依然可以绕过这种保护机制。

在前面的 [#继承](#继承) 中我们讲到，通过在结构体中嵌入匿名类型成员，可以继承匿名类型的方法。其实这个被嵌入的匿名成员不一定是普通类型，也可以是接口类型。我们可以通过嵌入匿名的testing.TB接口来伪造私有的private方法，因为接口方法是**延迟绑定**，编译时private方法是否真的存在并不重要。

```
package main

import (
    "fmt"
    "testing"
)

type TB struct {
    testing.TB
}

func (p *TB) Fatal(args ...interface{}) {
    fmt.Println("TB.Fatal disabled!")
}

func main() {
    var tb testing.TB = new(TB)
    tb.Fatal("Hello, playground")
}
```

我们在自己的TB结构体类型中重新实现了Fatal方法，然后通过将对象隐式转换为testing.TB接口类型（因为内嵌了匿名的testing.TB对象，因此是满足testing.TB接口的），然后通过testing.TB接口来调用我们自己的Fatal方法。

这种通过嵌入匿名接口或嵌入匿名指针对象来实现继承的做法其实是一种纯虚继承，我们继承的只是接口指定的规范，真正的实现在运行的时候才被注入。比如，我们可以模拟实现一个gRPC的插件：

```
type grpcPlugin struct {
    *generator.Generator
}

func (p *grpcPlugin) Name() string { return "grpc" }

func (p *grpcPlugin) Init(g *generator.Generator) {
    p.Generator = g
}

func (p *grpcPlugin) GenerateImports(file *generator.FileDescriptor) {
    if len(file.Service) == 0 {
        return
    }

    p.P(`import "google.golang.org/grpc"`)
    // ...
}
```

构造的grpcPlugin类型对象必须满足generate.Plugin接口（在"github.com/golang/protobuf/protoc-gen-go/generator"包中）：

```
type Plugin interface {
    // Name identifies the plugin.
    Name() string
    // Init is called once after data structures are built but before
    // code generation begins.
    Init(g *Generator)
    // Generate produces the code generated by the plugin for this file,
    // except for the imports, by calling the generator's methods
    // P, In, and Out.
    Generate(file *FileDescriptor)
    // GenerateImports produces the import declarations for this file.
    // It is called after Generate.
    GenerateImports(file *FileDescriptor)
}
```

generate.Plugin接口对应的grpcPlugin类型的GenerateImports方法中使用的p.P(...)函数却是通过Init函数注入的generator.Generator对象实现。这里的generator.Generator对应一个具体类型，但是如果generator.Generator是接口类型的话我们甚至可以传入直接的实现。

# 嵌入与聚合

在结构体中嵌入匿名的字段叫做，嵌入或内嵌；在结构体中的字段还包含类型名那么叫做，聚合。

在c语言中只考虑结构体和接口中嵌入的组合方式，有以下三种

## 在接口中嵌入接口

```
type Writer interface{
    Write()
}
type Reader interface{
    Read()
}
type Teacher interface{
    Reader
    Writer
}
```

## 在结构体中嵌入结构体

```
type Human struct{
    name string
}
type Writer interface{
    Write()
}
type Reader interface{
    Read()
}
type Teacher interface{
    Reader
    Writer
    Human
}
```

## 在结构体中嵌入接口

```
package main
import "fmt"

type Walker interface{
   Walk()
}

type Men struct{
   name string
   Walker
}

type Women struct{
   age int
}

func (m Men) Walk(){
   fmt.Println(m.name, " is walking")
}

func (m Women) Walk(){
   fmt.Println("A women aged", m.age,"is walking")
}

type Student struct{
   Men
   int
}

func main(){
   //直接对接口赋0值
   Bob := Men{name:"Bob"}
   Bob.Walk()

   //用实现了改接口的结构体对其赋值
   Alice := Men{"Alice",Women{20} }
   Alice.Walk()
   Alice.Walker.Walk()
}
```

可以看出直接调用Walk调用的是Men结构体自身的Walk方法，而不是其接口Walker的Walk方法


Go语言通过几种简单特性的组合，就轻易就实现了鸭子面向对象和虚拟继承等高级特性，真的是不可思议。






