# Go语言中类型(Type)转换、断言三两事

> 不少PHP出身的同学，对Golang中类型转换心存迷茫，这篇文章，整理Golang种类型相关的方方面面。

内容概要：

- 类型转换Type Conversions
- 类型断言Type Assertions
- 接口类型转换Interface Conversions

> **注意**：文中大量引用golang.org官方文档和链接，如链接无法访问，需自备梯子。

## 一、类型转换Type Conversions

> 这里所谓的类型转换，和其他语言的类型转换相似，但需注意Golang**多数情况为显式类型转换**，也存在某些隐式类型转换会在后文注意阐述。

### 1.1 显式类型转换(explicit conversion)

显式类型转换的实例：`T(x)`，将表达式`x`转换为类型T。
显式类型转换的语法规范：

```go
Conversion = Type "(" Expression [ "," ] ")" .
```

[Type](https://golang.org/ref/spec#Type)、[Expression](https://golang.org/ref/spec#Expression)可参考规范，

#### 1.1.1 类型T加括号避免歧义

注意：当类型以`*`、`->`、`func`开始时，注意使用括号，防止歧义。

```Golang
*Point(p)        // same as *(Point(p))
(*Point)(p)      // p is converted to *Point
<-chan int(c)    // same as <-(chan int(c))
(<-chan int)(c)  // c is converted to <-chan int
func()(x)        // function signature func() x
(func())(x)      // x is converted to func()
(func() int)(x)  // x is converted to func() int
func() int(x)    // x is converted to func() int (unambiguous)
```

- `*Point(p)`等价于`*(Point(p))`，p转换为Point类型，然后取指针值
- `(*Point)(p)`，将p转换为*Point类型
- `<-chan int(c)`等价于`<-(chan int(c))`，c转换为chan int型通道，然后读通道
- `(<-chan int)(c)`，将c转换为`<-chan int`只读通道
- `func()(x)`，函数签名为`func() x`，无返回值的函数x声明
- `(func() int)(x)`等价于`func() int(x)`，将x转换为`func() int`类型，且后者无括号也无歧义
- 综上几个实例，类型转换时，正确使用括号避免歧义

#### 1.1.2 常量类型转换

将常量转换为对应的类型常量（typed constant，也就是类型明确的常量）
> Converting a constant yields a typed constant as result.

```Golang
uint(iota)               // iota value of type uint
float32(2.718281828)     // 2.718281828 of type float32
complex128(1)            // 1.0 + 0.0i of type complex128
float32(0.49999999)      // 0.5 of type float32
float64(-1e-1000)        // 0.0 of type float64
string('x')              // "x" of type string
string(0x266c)           // "♬" of type string
MyString("foo" + "bar")  // "foobar" of type MyString
string([]byte{'a'})      // not a constant: []byte{'a'} is not a constant
(*int)(nil)              // not a constant: nil is not a constant, *int is not a boolean, numeric, or string type
int(1.2)                 // illegal: 1.2 cannot be represented as an int
string(65.0)             // illegal: 65.0 is not an integer constant
```

> 注意上述实例中[iota](https://golang.org/ref/spec#Iota)是常量声明中的关键词

注意：

- 显示类型转换要求：x必须是可被类型T表征的，[表征的官方含义](https://golang.org/ref/spec#Representability)
- int型的x到string类型，是允许的，因为涉及unicode/utf8编码值转换为字符串
  - `"\u65e5" == "日" == "\xe6\x97\xa5"`，中文`日`对应的Unicode编码(\u开头)为`"\u65e5`，对应的utf8编码(\x开头)为`\xe6\x97\xa5`
- 其他类型不兼容，是不能做显式类型转换的
- 注意，类型转换时涉及到`representable`、`assignable`都可在官方spec中找到完整介绍

#### 1.1.3 非常量值的类型转换

下列描述可将非常量(non-constant)值x转换为类型T：

- x可赋值成类型T，[可赋值assignable](https://golang.org/ref/spec#Assignability)
- 忽略struct的tags，x的类型和T类型具有相同的底层类型（underlying types，通过type类型别名后，底层类型相同）
- 忽略struct的tags，x的类型和T都是指针类型（不是被定义的类型，什么是定义的类型，详见类型别名和类型定义的区别），指针指向的类型具有相同的底层类型
- x的类型、T都是整型、浮点型的指针类型
- x的类型、T都是`complex`类型
- x是整型、字节或rune的切片，T是字符串类型（后文有各类型转换成字符串类型的详细说明）
- x是字符串类型，T是字节或rune的切片类型

**总结**：

- 上述列举的几种情况，均是允许的类型转换，非允许则禁止
- struct类型的tags标签在类型比较时将会忽略，不影响类型等价的判断，下面实例中`*Person`和`data`类型忽略标签，底层结构是等价的

```Golang
type Person struct {
    Name    string
    Address *struct {
        Street string
        City   string
    }
}

var data *struct {
    Name    string `json:"name"`
    Address *struct {
        Street string `json:"street"`
        City   string `json:"city"`
    } `json:"address"`
}

var person = (*Person)(data)  // ignoring tags, the underlying types are identical
```

- 数值型、字符串型间相互转换时，会有特定的规则，此类的类型转换可能改变x的值、并带来性能损耗；其他类型的转换，仅改变类型，而不会影响x代表的值
- Go在语言级没有自带指针间的转换方法，`unsafe`包在严格约束下有具体实现。

### 1.2 数值类型间的转换

非常量(non-constant)的数值类型进行类型转换时，遵循以下规则：

- 整数间类型转换时，如果x是有符号整型，转换类型后则扩充符号位；无符号型，则补零；长度变小，则裁剪字节。

```Golang
v := uint16(0x10F0) // v是无符号型 0x 10F0
int8(v) // 有符号型 0x F0，裁剪，只保留低8位，值为-16，也就是 - 0x 10
uint32(int8(v)) // 有符号型转换为无符号型，扩展符号位的1；也就是从 0x F0，变为 0x FF FF FF F0
```

- 将浮点型转换为整型时，小数位抹除

```Golang
// 这种是将常量转换，不是本小节关注的非常量的数值类型转换
fmt.Println(int(11.6)) // 语法报错：constant 11.6 truncated to integer

foo := 11.6
fmt.Println(int(foo)) // 结果为：11
```

- 将浮点型或整型转换为浮点型，或complex转换为另一种complex时，结果值四舍五入到目标类型指定的精度。**一句话：涉及浮点类型转换时，应仔细考虑精度变化的影响**

**总结**：在所有非常量类型转换时，涉及浮点型、complex类型值时，如果结果类型无法表示当前值，转换仍然成功，但最终结果值是多少会受到具体实现的影响。
换言之，如果不能完整表示，那结果将出乎你的预料。请不要给出任何预料，除非你非常清楚发生了什么。

### 1.3 字符串类型转换(转换到string、从string转换到其他类型)

- 将 「有符号或无符号」 整型转换为字符串，将产生一个包含此整型值代表的UTF-8字符串。超出有效 Unicode 编码的值将转换为`\uFFFD`

```Golang
string('a')       // "a"
string(-1)        // "\ufffd" == "\xef\xbf\xbd"
string(0xf8)      // "\u00f8" == "ø" == "\xc3\xb8"
type MyString string
MyString(0x65e5)  // "\u65e5" == "日" == "\xe6\x97\xa5"
```

> 说明：int或uint值 所代表的编码值，对应字符串，而其存储字节会是另外的形式
> 上面实例中， int值(0x65e5，Unicode编码，等价表示\u65e5) == 字符串(日) == 存储字节(\xe6\x97\xa5)

- 将字节切片(`a slice of bytes`)转换为字符串，将产生一个包含所有切片元素值的字符串。注意 `byte` 时，实际 `byte` 中存储的是 `utf-8` 编码后数据，而不是 `Unicode` 编码值。

> Converting a slice of bytes to a string type yields a string whose successive bytes are the elements of the slice.

```Golang
string([]byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'})   // "hellø"
string([]byte{})                                     // ""
string([]byte(nil))                                  // ""
type MyBytes []byte
string(MyBytes{'h', 'e', 'l', 'l', '\xc3', '\xb8'})  // "hellø"
```

- 将rune切片(a slice of runes)转换为字符串，将产生由所有切片元素值组合成的字符串。注意，rune时，切片中是 `Unicode` 值。

> Converting a slice of runes to a string type yields a string that is the concatenation of the individual rune values converted to strings.

```Golang
string([]rune{0x767d, 0x9d6c, 0x7fd4})   //"\u767d\u9d6c\u7fd4" == "白鵬翔"
string([]rune{})                         // ""
string([]rune(nil))                      // ""
type MyRunes []rune
string(MyRunes{0x767d, 0x9d6c, 0x7fd4})  //"\u767d\u9d6c\u7fd4" == "白鵬翔"
```

- 将字符串转换为切片，将字符串逐字节转换为切片元素

```Golang
[]byte("hellø")   // []byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'}
[]byte("")        // []byte{}
MyBytes("hellø")  // []byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'}
```

- 将字符串转换为rune切片时，rune切片为包含字符串Unicode编码值的切片

```Golang
[]rune(MyString("白鵬翔"))  // []rune{0x767d, 0x9d6c, 0x7fd4}
[]rune("")                 // []rune{}
MyRunes("白鵬翔")           // []rune{0x767d, 0x9d6c, 0x7fd4}
```

### 小节总结

- 注意集中表达形式
  - `0x767d`：这种是Unicode码值，也就是一个Int值
  - `\u767d`：Unicode编码格式表示
  - `\xc3 \xb8`：UTF-8编码表示
  - 如 `日` 的三种表示法分别是，`string(0x65e5)  // "\u65e5" == "日" == "\xe6\x97\xa5"`
- 字节向字符串转换
  - int值代表Unicode编码，转换为对应UTF-8编码后输出字符串
  - 切片时，逐一按字节转换、拼装后输出字符串
  - byte切片时，存Unicode编码值，转换为字符串
  - rune编码时，存UTF-8编码后结果，转换为字符串
  - 逐一多字节表示的UTF-8
- 字符串向字节转换
  - 字符串每个字节转换为字节值

### 小节练习

```Golang
// 字符表示，对应UTF-8编码后的表示
hi := string([]byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'})
// 字节值，对应每个的Unicode编码值
hi = string([]rune{0x68, 0x65, 0x6c, 0x6c, 0xf8})
fmt.Println(hi)
```

## 二、类型断言Type Assertions

类型断言，不同于类型转换。
断言，只能用在接口(`interface`)类型上，是 `Go` 语言中动态类型的一个体现。

语法：将接口类型x进行断言（类型为T则转换为类型T，否则会失败出 `panic` 或赋值给第二返回值），返回值为类型T的新变量值
注意：类型T需为接口实际类型或T可转换的另一个接口类型

```Golang
x.(T)
v := x.(T)
v, ok := x.(T)
```

断言推荐用法：

```Golang
str, ok := value.(string)
if ok {
    fmt.Printf("string value is: %q\n", str)
} else {
    fmt.Printf("value is not a string\n")
}

if str, ok := value.(string); ok {
    return str
} else if str, ok := value.(Stringer); ok {
    return str.String()
}
```

另外，`Go` 语言中动态类型，可用 [Type switches](https://golang.org/doc/effective_go.html#type_switch)来动态判断，也是接口、断言组合使用的常用方式。

实例：

```Golang
type Stringer interface {
    String() string
}

var value interface{} // Value provided by caller.
switch str := value.(type) {
case string:
    return str
case Stringer:
    return str.String()
}
```

## 三、接口类型转换Interface Conversions



## 参考文档

- [类型转换Conversions](https://golang.org/ref/spec#Conversions)
- [接口类型转换Interface Conversions](https://golang.org/doc/effective_go.html#conversions)
- [接口类型转换和断言Interface conversions and type assertions](https://golang.org/doc/effective_go.html#interface_conversions)
- [可赋值assignable](ttps://golang.org/ref/spec#Assignability)
- [可表征representable](https://golang.org/ref/spec#Representability)