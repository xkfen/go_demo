package go_demo

import (
	"testing"
	"fmt"
)

type S struct{}

func (s S) F() {}

type IF interface {
	F()
}

func InitType() S {
	var s S
	return s
}

func InitPointer() *S {
	var s *S
	return s
}
func InitEfaceType() interface{} {
	var s S
	return s
}

func InitEfacePointer() interface{} {
	var s *S
	return s
}

func InitIfaceType() IF {
	var s S
	return s
}

func InitIfacePointer() IF {
	var s *S
	return s
}

func TestNil(t *testing.T) {
	//fmt.Println(InitType() == nil)
	fmt.Println(InitPointer() == nil)
	fmt.Println(InitEfaceType() == nil)
	fmt.Println(InitEfacePointer() == nil)
	fmt.Println(InitIfaceType() == nil)
	fmt.Println(InitIfacePointer() == nil)
}

const N  = 3
func TestTemporaryPointer(t *testing.T){
	m := make(map[int]*int)

	for i := 0; i < N; i++ {
		m[i] = &i //A
		fmt.Println(m[i])
	}

	for _, v := range m {
		fmt.Println(*v)
	}
}

func f1() {
	defer fmt.Println("f1-begin")
	f2()
	defer fmt.Println("f1-end")
}

func f2() {
	defer fmt.Println("f2-begin")
	f3()
	defer fmt.Println("f2-end")
}

func f3() {
	defer fmt.Println("f3-begin")
	//panic(0)
	//defer fmt.Println("f3-end")
}

func TestDeffer(t *testing.T) {
	f1()
}