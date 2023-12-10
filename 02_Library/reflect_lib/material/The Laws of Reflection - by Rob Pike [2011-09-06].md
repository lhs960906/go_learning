# 反思法则（Laws）
Rob Pike
6 September 2011

## 介绍
计算中的反射是程序检查其自身结构的能力，特别是通过类型； 这是元编程的一种形式。 这也是造成混乱的一个重要原因。
在本文中，我们试图通过解释反射在 Go 中的工作原理来澄清问题。 每种语言的反射模型都是不同的（许多语言根本不支持它），但本文是关于 Go 的，因此在本文的其余部分中，“反射”一词应理解为“Go 中的反射”。
2022 年 1 月添加的注释：这篇博文写于 2011 年，早于 Go 中的参数多态性（又称泛型）。 尽管本文中的任何重要内容都没有因为语言的发展而变得不正确，但它在一些地方进行了调整，以避免让熟悉现代 Go 的人感到困惑。

## 类型和接口
因为反射建立在类型系统之上，所以让我们首先回顾一下 Go 中的类型。
Go 是静态类型的。 每个变量都有一个静态类型，即在编译时已知并固定的一种类型：int、float32、*MyType、[]byte 等。 如果我们声明
```go
type MyInt int

var i int
var j MyInt
```
那么 i 的静态类型为 int，j 的静态类型为 MyInt。 变量 i 和 j 具有不同的静态类型，尽管它们具有相同的基础类型，但如果不进行转换，它们就无法相互分配。

类型的一个重要类别是接口类型，它表示固定的方法集。（在讨论反射时，我们可以忽略使用接口定义作为多态代码中的约束。）接口变量可以存储任何具体（非接口）值，只要该值实现接口的方法即可。 一对著名的例子是 io.Reader 和 io.Writer（来自 io 包的 Reader 和 Writer 类型）：
```go
// Reader is the interface that wraps the basic Read method.
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer is the interface that wraps the basic Write method.
type Writer interface {
    Write(p []byte) (n int, err error)
}
```
任何使用此签名实现 Read（或 Write）方法的类型都被称为实现 io.Reader（或 io.Writer）。出于本次讨论的目的，这意味着 io.Reader 类型的变量可以保存其类型具有 Read 方法的任何值：
```go
var r io.Reader
r = os.Stdin
r = bufio.NewReader(r)
r = new(bytes.Buffer)
// and so on
```
需要明确的是，无论 r 持有什么具体值，r 的类型始终是 io.Reader，因为 Go 是静态类型的，所以接口类型变量 r 的静态类型始终是 io.Reader。

接口类型的一个极其重要的例子是空接口：
```go
interface{}
```
或其等效别名：
```go
any
```
它表示空的方法集，并且可以满足任何值，因为每个值都有零个或多个方法。

有人说 Go 的接口是动态类型的，但这是误导。它们是静态类型的：接口类型的变量始终具有相同的静态类型，即使在运行时存储在接口变量中的值可能会更改类型，该值也将始终满足接口的要求。

我们需要精确地对待这一切，因为反射和接口密切相关。

## 接口的表示
Russ Cox 撰写了一篇关于 Go 中接口值表示的详细博客文章。 这里没有必要重复整个故事，但可以做一个简单的总结。

接口类型的变量存储一个（值，类型）对：
* 分配给该变量的具体值
* 该值的类型描述符
更准确地说，值是实现接口的底层具体数据项，类型描述了该项的完整类型。例如：
```go
var r io.Reader
tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
if err != nil {
    return nil, err
}
r = tty
```
r 示意性地包含（值，类型）对（tty，*os.File）。请注意，类型 *os.File 实现了 Read 以外的方法； 尽管接口值仅提供对 Read 方法的访问，但内部值携带有关该值的所有类型信息。这就是为什么我们可以这样做：
```go
var w io.Writer
w = r.(io.Writer)
```
此赋值中的表达式是类型断言；它断言 r 中的项也实现了 io.Writer，因此我们可以将其分配给 w。赋值后，w 将包含 (tty, *os.File) 对。这与 r 中保存的是同一对。接口的静态类型决定了可以使用接口变量调用哪些方法，即使内部的具体值可能具有更大的方法集。

接下来，我们可以这样做：
```go
var empty interface{}
empty = w
```
我们的空接口值 empty 将再次包含同一对（tty，*os.File）。 这很方便：空接口可以保存任何值，并包含我们可能需要的有关该值的所有信息。

这里我们不需要类型断言，因为静态地知道 w 满足空接口。在我们将值从 Reader 移动到 Writer 的示例中，我们需要显式地使用类型断言，因为 Writer 的方法不是 Reader 的子集。

## 反射法则之一
从根本上来说，反射只是一种检查存储在接口变量内的（值，类型）对的机制。 首先，我们需要了解反射包中的两种类型：Type 和 Value。 这两种类型可以访问接口变量的内容，两个简单的函数（称为reflect.TypeOf 和 reflect.ValueOf）可以从接口值中检索 reflect.Type 和 reflect.Value 部分。此外，从 reflect.Value 很容易获得相应的 reflect.Type，但现在让我们将 Value 和 Type 概念分开。

让我们从 reflect.TypeOf 开始：
```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    var x float64 = 3.4
    fmt.Println("type:", reflect.TypeOf(x))
}
```
这个程序将会打印：
```shell
type: float64
```
您可能想知道这里的接口在哪里，因为程序看起来像是将 float64 变量 x（而不是接口值）传递给reflect.TypeOf。 但它就在那里； 正如 godoc 报告的那样，reflect.TypeOf 的签名包含一个空接口：
```go
// TypeOf returns the reflection Type of the value in the interface{}.
func TypeOf(i interface{}) Type
```
当我们调用 reflect.TypeOf(x) 时，x 首先存储在一个空接口中，然后将其作为参数传递；Reflect.TypeOf 解压该空接口以恢复类型信息。

然，reflect.ValueOf 函数会恢复值（从这里开始，我们将省略样板文件，只关注可执行代码）：
```go
var x float64 = 3.4
fmt.Println("value:", reflect.ValueOf(x).String())
```
打印：
```go
value: <float64 Value>
```
这里我们显式调用 String 方法，因为默认情况下 fmt 包会挖掘 reflect.Value 来显示内部的具体值。而 String 方法则不会。

reflect.Type 和 reflect.Value 都有很多方法让我们检查和操作它们。一个重要的例子是 Value 有一个 Type 方法，该方法返回 reflect.Value 的 Type。 另一个是 Type 和 Value 都有一个 Kind 方法，该方法返回一个常量，指示存储的项目类型：reflect.Uint、reflect.Float64、reflect.Slice 等。此外，Value 类型具有的 Int 和 Float 等名称方法可以让我们获取存储在其中的值（如 int64 和 float64）：
```go
var x float64 = 3.4
v := reflect.ValueOf(x)
fmt.Println("type:", v.Type())
fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
fmt.Println("value:", v.Float())
```
打印：
```go
type: float64
kind is float64: true
value: 3.4
```
还有像 SetInt 和 SetFloat 这样的方法，但要使用它们，我们需要了解可设置性，这是反射第三定律的主题，如下所述。
反射库有几个值得特别指出的属性。首先，为了保持 API 简单，Value 的 `getter` 和 `setter` 方法对可以保存该值的最大类型进行操作：例如，对于所有有符号整数，int64。 也就是说，Value 的 Int 方法返回一个 int64，而 SetInt 的值则取一个 int64； 可能需要转换为涉及的实际类型：
```go
var x uint8 = 'x'
v := reflect.ValueOf(x)
fmt.Println("type:", v.Type())                            // uint8.
fmt.Println("kind is uint8: ", v.Kind() == reflect.Uint8) // true.
x = uint8(v.Uint())                                       // v.Uint returns a uint64.
```
第二个属性是反射对象的 Kind 描述的是基础类型，而不是静态类型。 如果反射对象包含用户定义的整数类型的值，如：
```go
type MyInt int
var x MyInt = 7
v := reflect.ValueOf(x)
fmt.Println("type:", v.Type())                        // src.MyInt
fmt.Println("kind is int: ", v.Kind() == reflect.Int) // true.
```
v 的种类仍然是 reflect.Int，即使 x 的静态类型是 MyInt 而不是 int。换句话说，Kind 无法区分 int 和 MyInt，但 Type 可以。

## 反射法则之二

## 反射法则之三



# 术语
静态类型：
接口类型变量/接口变量：
基础类型：