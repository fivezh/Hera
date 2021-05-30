package main

import (
	"fmt"
	"unsafe"
)

type Foo struct {
	a int8
	b int64
	c int32
}

func main() {
	a := Foo{}
	fmt.Printf("size of Foo is %d", unsafe.Sizeof(a))
}
