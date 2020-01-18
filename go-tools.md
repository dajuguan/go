# Go开发工具介绍
## go vet检查错误
vet命令可以帮组开发人员检测代码的常见错误，它可以捕捉下面4种错误

- Printf类函数调用时，类型匹配错误的参数
- 定义常用方法时，方法签名的错误
- 错误的结构标签
- 没有知名字段名的结构变量

## go fmt格式化代码
自动把代码格式化微Go源码类似的风格，不用再为括号的换行问题，tab缩进多少个空格而中伦不休。
```
package main

import "fmt"
func main(){
      fmt.Println("I am the last..")}

func init(){
   fmt.Println("I am the first!")
}
```
运行```go gmt 文件名```命令后：

```
package main

import "fmt"
func main(){
      fmt.Println("I am the last..")}

func init(){
   fmt.Println("I am the first!")
}
```
直接大变身！很漂亮有木有

## go doc查看文档
# 直接查看
比如想查看archive/tar包相关的文档,可以运行命令
```go doc tar```
可以看到下面的接口函数文档描述：
```
package tar // import "archive/tar"

Package tar implements access to tar archives. It aims to cover most of the
variations, including those produced by GNU and BSD tars.

References:

    http://www.freebsd.org/cgi/man.cgi?query=tar&sektion=5
    http://www.gnu.org/software/tar/manual/html_node/Standard.html
    http://pubs.opengroup.org/onlinepubs/9699919799/utilities/pax.html

const TypeReg = '0' ...
var ErrWriteTooLong = errors.New("archive/tar: write too long") ...
var ErrHeader = errors.New("archive/tar: invalid tar header")
type Header struct{ ... }
    func FileInfoHeader(fi os.FileInfo, link string) (*Header, error)
type Reader struct{ ... }
    func NewReader(r io.Reader) *Reader
type Writer struct{ ... }
    func NewWriter(w io.Writer) *Writer
```

# 本地查看
可以在本地启动go的文档服务器，只需运行如下中断命令
```godoc -http=:6060```,就可以在```http://localhost:6060/```，看到如下界面：
![go本地文档](.gitbook/assets/go本地文档.png)

go文档也支持，使用简单规则来自动包含代码文档
比如我们修改包那一节的test文件
```
//test.go
package test
import "fmt"
//输出语句
func Say(){
    fmt.Println("Hi,I am exported")
}
```
运行```godoc test```,会输出：
```
PACKAGE DOCUMENTATION

package test
    import "test"

    test.go

FUNCTIONS

func Say()
    输出语句
```

## go get导入第三方库
# 自动安装
不像java的maven仓库或者nodejs有npm可以去搜索第三方安装包，一般而言我们需要找go第三方的库时，直接上github搜索就可以了。比如安装```https://github.com/sirupsen/logrus```这个包,可以直接运行:
```go get github.com/sirupsen/logrus```
如果没有反应则一切正常，否则如果报错
```
package golang.org/x/sys/unix: unrecognized import path "golang.org/x/sys/unix" (https fetch: Get https://golang.org/x/sys/unix?go-get=1: dial tcp 216.239.37.1:443: i/o timeout)
```
那么说明访问不了golang，需要先执行
```
cd $GOPATH/src
mkdir /golang.org
mkdir /golang.org/x
cd $GOPATH/src/golang.org/x
git clone https://github.com/golang/crypto.git
go get -u golang.org/x/crypto/ssh/terminal
git clone https://github.com/golang/sys.git
go get -u golang.org/x/sys/unix
```
然后再运行
```go get github.com/sirupsen/logrus```
即可

## 依赖管理
Go语言的依赖管理是非常失望的，如上一节使用go get会发现有很多的问题
- 如果想定位到特定的第三方库版本号，是比较困难的；
- 依赖的完整性无法校验，基于域名的package，域名变化或者子路径变化，都会无法使用
- 第三方包没有安全审计，很容易引入bug

好在社区有一定的解决方案，参见[https://github.com/golang/go/wiki/PackageManagementTools]

## 小结