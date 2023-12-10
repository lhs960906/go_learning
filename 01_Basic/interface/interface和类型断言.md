# 一、类型元数据

## 1.1、for all 类型

在go语言中，int8/int16/int32/int64/int/byte/string/slice/func/map... 等属于内置类型，而我们通过如下方法定义的则属于自定义类型：

```go
type T int

type T struct {
	name string
}

type I interface {
	Name() string
}
```

不管是自定义类型（包括结构体和接口）还是内置类型，都有对应的类型描述信息，它们称为类型的元数据，这些类型元数据共同构成了 Go 语言的类型系统。

类型编号、类型的大小，对齐边界，是否自定义类型等，这是每种类型的元数据都必须要记录的信息。它们被定义在了 runtime 包的 type.go 文件中，由 _type 类型进行描述。
```go
type _type struct {
	size       uintptr	// 数据类型共占用的空间大小
	ptrdata    uintptr  // 含有所有指针类型前缀大小
	hash       uint32	// 类型hash值；避免在哈希表中计算
	tflag      tflag	// 额外类型信息标志
	align      uint8	// 该类型变量对齐方式
	fieldAlign uint8	// 该类型结构字段对齐方式
	kind       uint8	// 类型编号
	equal func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata    *byte		// gc数据
	str       nameOff	// 类型名字的偏移
	ptrToThis typeOff
}
```

一言以蔽之，就是 go 的类型系统都是以 `runtime._time` 基础，并扩展除了诸如 slicetype、structtype 等其他类型。

## 1.2、内置类型slicetype

在某些具体类型的元数据信息中，还可能包含各种类型额外需要描述的信息。以 slice 这种内置类型为例，它的类型元数据结构被定义在了 runtime 包的 type.go 文件的 slicetype 结构中。

```go
type slicetype struct {
	typ  _type	// 通用元数据
	elem *_type	// 额外的描述信息
```

slicetype 中的 typ 就是每种类型都要有的通用元数据， 而 elem 则是 slice 这种类型需要额外描述的元数据信息。因为 slice 中需要存储元素，所以 slice 必须记录其内部元素的类型元数据信息，而 elem 指向的就是 slice 中存储元素的类型元数据信息（比如对于 []string 这种类型的 slice 来说，elem 则会指向 string 类型元数据信息）。

## 1.3、自定义类型

如果你自定义了一个类型，那么还需要一个 uncommontype 类型（runtime 包的 type.go 文件）去描述这个自定义类型的元数据信息。

```go
type uncommontype struct {
	pkgpath nameOff	// 记录类型所在的包路径
	mcount  uint16  // 记录类型关联的方法个数
	xcount  uint16  // 记录导出方法的个数
	moff    uint32  // 记录uncomontype结构体和方法元数据method结构偏离的字节数
	_       uint32  // unused
}
```

例如，我们基于 `[]string` 定义了一个新的自定义类型 myslice，并且绑定了方法 `Len()` 和 `Cap()`。
```go
type myslice []string

func (ms myslice) Len() {
	fmt.Println(len(ms))
}

func (ms myslice) Cap() {
	fmt.Println(cap(ms))
}
```
那么 myslice 的类型元数据将会是如下结构：

![image-20230617111709481](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617111709481.png)



# 二、接口

## 2.1、接口相关概念

### 2.1.1、接口值

概念上讲一个接口的值，接口值，由两个部分组成，一个具体的类型和那个类型的值。它们被称为接口的动态类型和动态值。

### 2.1.2、接口实现

一个类型如果拥有一个接口需要的所有方法，那么这个类型就实现了这个接口。例如，`*os.File`类型实现了io.Reader，Writer，Closer，和ReadWriter接口。`*bytes.Buffer`实现了Reader，Writer，和ReadWriter这些接口。但是它没有实现Closer接口因为它不具有Close方法。Go的程序员经常会简要的把一个具体的类型描述成一个特定的接口类型。举个例子，`*bytes.Buffer `是 io.Writer，`*os.Files`是 io.ReadWriter。

## 2.2、空接口

### 2.2.1、空接口结构

空接口的结构信息被定义在 runtime 包中 runtime2.go 的 eface 结构中：

```go
type eface struct {
    _type *_type
    data  unsafe.Pointer
}
```

空接口类型无任何方法，并且可以接受任意类型的数据，所以它只需要记录这个接口的动态类型是啥，动态值在哪就行。这个 eface 结构中的 _type 指向的就是接口的**动态类型元数据**，data 就指向接口的**动态值**。

所以，一个空接口的结构就大概如下图所示：

![image-20230616214538197](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230616214538197.png)

### 2.2.2、空接口具体示例

我们逐行分析如下代码的行为：

```go
var ei interface{}
f, _ := os.Open("eggo.txt")
ei = f
```

首先我们声明了一个空接口 ei：

```go
var ei interface{}
```

在被赋值之前， _type 和 data 都是 nil，如下图：

![image-20230615170628996](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230615170628996.png)

现在，我们创建一个 \*os.File 类型的变量 f：

```go
f, _ := os.Open("eggo.txt")
```

在了解类型元数据之后，我们可以知道，\*os.File 的类型元数据将会是如下结构：

![image-20230617113522352](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617113522352.png)

如果我们将 f 赋值给这个空接口 ei：

```go
ei = f
```

那么 ei 的动态类型 _type 将会指向 *os..File 的类型元数据，ei 的动态值将会变为 f，指向一个 os.File 对象。那么整个结构就会如下图所示：

![image-20230617120707892](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617120707892.png)

## 2.3、非空接口

### 2.3.1、非空接口结构

非空接口的结构信息被定义在 runtime 包中 runtime2.go 的 iface 结构中：

```go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}
```

和空接口类型一样，iface 结构中的 data unsafe.Pointer 指向接口的动态值，而接口要求的方法列表以及接口的动态类型信息则存储在 itab 结构中：

```go
type itab struct {
    inter *interfacetype // 描述接口类型元数据的指针-interfacetype指针，接口需要实现的方法列表存储在此结构中
    _type *_type		 // 指向接口的动态类型元数据
    hash  uint32		 // 是从动态类型元数_type中拷贝出来的类型哈希值，用于快速判断类型是否相等
    _     [4]byte	
    fun   [1]uintptr     // 记录动态类型实现的接口方法地址。其大小可变。如果fun[0]==0则表示动态类型_type *_type没有实现inter *interfacetype接口
}
```

itab 结构的第一个字段 inter interfacetype 存储接口类型的元数据：

```go
type interfacetype struct {
	typ     _type		
	pkgpath name		// 包路径
	mhdr    []imethod	// 接口需要实现的方法列表就记录在此
}
```

所以，一个非空接口的结构就大致如下图：

![image-20230616214409907](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230616214409907.png)

### 2.3.2、非空接口具体示例

假设我们有如下代码：

```go
var nei io.ReadWriter
f, _ := os.Open("eggo.txt")
nei = f
```

首先我们声明了一个非空接口 nei：

```go
var nei io.ReadWriter
```

在被赋值之前， nei 的 `tab  *itab` 和 `data unsafe.Pointer` 指向的都是 nil。紧接着我们创建了一个 *os.File 的对象，并将其赋值给了 nei，那么 nei 将会变为如下结构：

![image-20230617121214865](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617121214865.png)

需要注意的，go 会将动态类型元数据中接口要求的方法（这里是 Read、Writer）地址拷贝到 itab 的 fun 中，这样可以快速定义到实现的方法地址，而无需每次都去动态类型元数据中进行查找。

一旦接口类型 inter 和接口的动态类型 _type 确定了，那么接口的 itab 结构也就不会再改变了。所以实际上 Go 会以 `接口类型的hash ^ 接口的动态类型hash` 的结果为 key，以 itab 结构体指针为 value，将这些 itab 组合起来构造一个 hashtable 作为 itab 缓存，以便存储和查询 itab。如果 key 有对应的 val，那么可以直接拿来使用，反之在 hashtable 中加入新的 itab 进行缓存。

![image-20230617172639667](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617172639667.png)

明白了空接口和接口的数据结构，理解了接口的动态类型和动态值在赋值前的变化，之后我们就可以聊聊类型断言了。



# 三、类型断言

前面说过，接口可以分为空接口和非空接口两类。相对于接口这种抽象类型，int、string、slice、map、struct 等都是具体类型。

类型断言 `x.(T)` 作用于接口值之上，断言类型 x 可以是空接口或非空接口，而断言的目标类型 T 可以是具体类型或非空接口类型。这样就会出现四种类型断言：

* 1）空接口.(具体类型)
* 2）非空接口.(具体类型)
* 3）空接口.(非空接口)
* 4）非空接口.(非空接口)

## 3.1、四种类型断言的具体操作

### 3.1.1、空接口ei.(具体类型T)

我们先给出该种类型断言的结论：

```
res, ok := 空接口ei.(具体类型T)
```

> 判断空接口 ei 的动态类型元数据是否为 T 类型
>
> * 判断成功：res 为 T 类型变量，值为 ei 的动态值，ok 为 true
> * 判断失败：res 为 T 类型的零值，ok 为 false

### 3.1.2、非空接口nei.(具体类型T)

我们先给出该种类型断言的结论：

```go
res, ok := 非空接口nei.(具体类型T)
```

> 判断非空接口 nei 的动态类型是否为是否为具体类型 T
>
> * 断言成功：res 为 T 类型变量，值为 nei 的动态值，ok 为 true
> * 断言失败：res 为 T 类型的零值，ok 为 false

### 3.1.3、空接口ei.(非空接口T)

我们先给出该种类型断言的结论：

```go
res, ok := 空接口ei.(非空接口T)
```

> 判断空接口 ei 的动态类型是否实现了非空接口 T
>
> * 断言成功：res 为 T 类型变量，其 tab 指向 T 类型和 ei 的动态类型对应的 itab，其动态值和 ei 的相同，ok 为 true
> * 断言失败：res 为 T 类型的零值，ok 为 false

### 3.1.4、非空接口nei.(非空接口T)

我们先给出该种类型断言的结论：

```go
res, ok := 非空接口nei.(非空接口T)
```

> 判断非空接口nei 的动态类型是否实现了非空接口 T
>
> * 断言成功：res 为 T 类型变量，其 tab 指向 T 类型和 rw 的动态类型对应的 itab，其动态值和 rw 的相同，ok 为 true
>
> * 断言失败：res 为 T 类型的零值，ok 为 false

## 3.2、四种断言的底层原理

### 3.2.1、空接口ei.(具体类型T)

这种类型的断言是要判断：空接口 ei 的动态类型元数据是否为 T 类型

#### 3.2.1.1、断言成功

现有如下操作：

```go
var ei interface{}
f, _ := os.Open("eggo.txt")
ei = f
res, ok := ei.(*os.File)
```

在执行完前 3 行代码之后，运用我们前面所学的只是，我们知道空接口 ei 的结构将会变为如下：

![image-20230617125132938](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617125132938.png)

最后一行代码做了一个类型断言操作，它会去判断 ei 的动态类型是否是 *os.File 类型，是则成功，反之失败。在这个例子中，ei 的动态类型就是 *os.File，那么类型断言成功，res 将是 *os.File 类型，值为 ei 的动态值 f。

#### 3.2.1.2、断言失败

假设代码变为如下：

```go
var ei interface{}
f := "eggo.txt"
ei = f
res, ok := ei.(*os.File)
```

在执行完前 3 行代码之后，ei 的结构将如下图所示：

![image-20230617162331098](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617162331098.png)

因为断言的目标类型是 *os.File，而 ei 的动态类型为 string，所以此时断言失败，res 为 *os.File 的类型零值，ok 为 false。

### 3.2.2、非空接口nei.(具体类型T)

这种类型的断言是要判断：非空接口 nei 的动态类型是否为是否为具体类型 T

#### 3.2.2.1、断言成功

现有如下操作：

```go
var nei io.ReadWriter
f, _ := os.Open("eggo.txt")
nei = f
res, ok := nei.(*os.File)
```

我们首先还是看前 3 行代码执行完之后 io.ReadWriter 接口的结构：

![image-20230617170623635](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617170623635.png)

紧接着第 4 行代码会进行类型断言，前面我们说过，程序中的 itab 都会被缓存进 hashtable 中，所以上面的断言操作会去 hashtable 中用 <io.ReadWriter, *os.File> 为 key 去寻找对应的 itab。如下图：

![image-20230617172558459](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617172558459.png)

在这个例子中，找到的 itab 其动态类型为 *os.File，和断言的目标类型 *os.File 一致，此时类型断言成功。res 为 *os.File 类型，其值为 f（f 是一个指针，指向对应的 os.File）。

#### 3.2.2.2、断言失败

假如非空接口的动态类型变为不是 *os.File，如下代码所示：

```go
var nei io.ReadWriter
f := eggo{name: "eggo"}
nei = f
r, ok := nei.(*os.File)

type eggo struct {
    name string
}

func (e *eggo) Write(b [byte]) (n int, err, error) {
    return len(e.name), nil
}
```

那么当前 3 行代码执行完毕之后，nei 将为如下结构：

![image-20230617173030374](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617173030374.png)

当执行第 4 行代码的类型断言操作时，会去 hashtable 中用 <io.ReadWriter, eggo> 为 key 去寻找对应的 itab。如下图：

![image-20230617174117174](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617174117174.png)

在这个例子中，找到的 itab 其动态类型为 eggo 类型，和断言的目标类型 *os.File 不一致，此时类型断言失败。res 为 *os.File 类型，其值为 *os.File 的零值 nil。

### 3.2.3、空接口ei.(非空接口T)

这种类型的断言是要判断：空接口 ei 的动态类型是否实现了非空接口 T

#### 3.2.3.1、断言成功

现有如下操作：

```go
var ei interface{}
f, _ := os.Open("eggo.txt")
ei = f
res, ok := ei.(io.ReadWriter)
```

结合[2.1、空接口结构](2.1、空接口结构)相关的知识，我们可以知道，前 3 行代码执行完之后，我们将得到如下结构的 ei：

![image-20230617183629284](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617183629284.png)

当执行第 4 行代码的类型断言操作时，通过检查 ei 动态类型元数据 *os.File 中的方法列表，和 io.ReadWriter 接口类型元数据中的 mhdr 的方法列表，我们可以知道 ei 的动态类型是否实现了 io.ReadWriter 接口。但其实也不必每次都这么比较，结合我们之前在[2.2.2、非空接口具体示例]()中提到的 itab 缓存相关知识，我们可以 <io.ReadWriter, *os.File> 为 key 去查找 itab 缓存，如果能找到，只需要检查 itab 中 fun[0] 是否为 0 即可判断 *os.File 是否实现了 io.ReadWriter 接口（fun[0] 说明动态类型并未实现接口），如果找不到，就需要去对比 ei 动态类型元数据 *os.File 中的方法列表和 io.ReadWriter 接口类型元数据中的 mhdr 的方法列表。

在本例中，*os.File 类型是实现了 io.ReadWriter 接口的，所以类型断言成功，res 为 io.ReadWriter 类型，其 itab 指向 itab 缓存中的 key 为 <io.ReadWriter, *os.File> 的 itab，其动态值为 f（f 是一个指针，指向对应的 os.File）。

#### 3.2.3.2、断言失败

这里如果我们自定义一个 eggo 类型，并且 eggo 类型只实现了 Write 方法：

```go
var ei interface{}
f, _ := eggo{ name: "eggo" }
ei = &f
res, ok := ei.(io.ReadWriter)

type eggo struct {
    name string
}

func (e *eggo) Write(b [byte]) (n int, err, error) {
    return len(e.name), nil
}
```

那么在前 3 行代码执行完毕之后有如下 ei 结构：

![image-20230617183507352](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617183507352.png)

当执行第 4 行代码的类型断言操作时，我们知道 eggo 只有一个 Write 方法，并没有实现 io.ReadWriter 接口，所以类型断言失败。res 为 io.ReadWriter 类型的变量，其动态值为 io.ReadWriter 类型的零值 nil。

### 3.2.4、非空接口nei.(非空接口T)

这种类型的断言是要判断：非空接口 nei 的动态类型是否实现了非空接口 T

#### 3.2.4.1、断言成功

现有如下操作：

```go
var nei io.Writer
f, _ := os.Open("eggo.txt")
nei = f
res, ok := nei.(io.ReadWriter)
```

前 3 行代码执行完毕之后，nei 的结构：

![image-20230617185909072](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617185909072.png)

同样地，第 4 行代码的断言操作通过去 itab 缓存中查找 key 为 <io.ReaderWriter, *os.File> 的 itab，若存在，且 fun[0] != 0，那么断言成功。如果不存在，会再去检查 nei 动态类型元数据 *os.File 中的方法列表和 io.ReadWriter 接口中的方法列表，并缓存 itab 信息。

在本例中，nei 的动态类型 *os.File 有 Write 和 Read 方法，因此它实现了 io.ReadWriter 接口，所以类型断言成功。res 为 io.ReadWriter 类型，其 itab 指向缓存中的 itab，其值为 nei 的动态值 f。

#### 3.2.4.2、断言失败

假如代码变为如下：

```go
var nei io.Writer
f, _ := eggo{name: "eggo"}
nei = &f
res, ok := nei.(io.ReadWriter)

type eggo struct {
    name string
}

func (e *eggo) Write(b [byte]) (n int, err, error) {
    return len(e.name), nil
}
```

同样还是看前 3 行代码执行完毕之后，nei 的结构为：

![image-20230617190955043](https://raw.githubusercontent.com/lhs960906/image_repo/main/PicGoimage-20230617190955043.png)

其中，因为 *eggo 并未实现 Write 方法，所以这个新的 itab 其 fun[0] 被置为 0，并将被放入 itab 缓存中，key 为 <io.ReadWriter, *eggo>。

第 4 行代码的断言操作通过去 itab 缓存中查找 key 为 <io.ReaderWriter, *eggo> 的 itab，此时是可以找到的，但因为其 fun[0] = 0，所以断言失败。res 为 io.ReaderWriter 类型，其值为 io.ReaderWriter 类型的零值（即 itab 为 nil，动态值也为 nil）。
