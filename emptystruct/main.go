package main

import (
	"fmt"
	"unsafe"
)

type People struct{}

func main() {
	var a1 uint8
	fmt.Println(unsafe.Sizeof(a1))

	// a and b are in the stack, the compiler optimizes a == b setting to false
	a := &struct{}{}
	b := &struct{}{}
	println(a, b, a == b)

	// fmt.Printf causes c and d to escase to the heap,
	// empty structs point to the same address called zerobase,
	// c == d comparison is true
	c := &struct{}{}
	d := &struct{}{}
	fmt.Printf("%p, %p, %v\n", c, d, c == d)
}
