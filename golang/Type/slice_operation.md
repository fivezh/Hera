# 切片相关操作

## 切片合并

PHP中数组合并：`$a = array_merge($a, $b)`
Golang中合并两个数据：`s1 = append(s1, s2...)`

```golang
package main

import "fmt"

func main() {
    s1 := []int{0, 1, 2, 3}
    s2 := []int{4, 5, 6, 7}

    s1 = append(s1, s2...)
    fmt.Println(s1)
}
// [0 1 2 3 4 5 6 7]

append([]int{1,2}, []int{3,4}...)
等价于：
func foo(is ...int) {
    for i := 0; i < len(is); i++ {
        fmt.Println(is[i])
    }
}

func main() {
    foo([]int{9,8,7,6,5}...)
}
```

## 切片中的`...`符号

> If the final argument is assignable to a slice type []T, it is passed unchanged as the value for a ...T parameter if the argument is followed by .... In this case no new slice is created.
> 这句比较难理解，传递的是值？还是新的切片

文档：https://golang.org/ref/spec#Passing_arguments_to_..._parameters

## 参考

- https://segmentfault.com/q/1010000011354818
- https://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go
