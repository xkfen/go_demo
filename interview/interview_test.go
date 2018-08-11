package interview

import (
	"testing"
	"fmt"
	"runtime"
	"sync"
)

// 测试interface的nil值:interface顶层由type和data组成，只有当type和data都是nil的时候，整个interface才是nil
// 答案：no
func TestInterfaceNil(t *testing.T) {
	var s []int = nil
	var f interface{} = s
	if f == nil {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}

// 测试defer和panic的发生顺序
/**
答案：
打印后
打印中
打印前
触发异常
 */
/**
解析：defer是先进后出，出现panic的时候，会先按照defer的现金后出顺序执行defer，最后才会执行panic
 */

func TestDeferPanic(t *testing.T) {
	defer func() {
		fmt.Println("打印前")
	}()
	defer func() {
		fmt.Println("打印中")
	}()
	defer func() {
		fmt.Println("打印后")
	}()
	panic("触发异常")
}

// 测试foreach的坑
type student struct {
	Name string
	Age  int
}

func TestForEach(t *testing.T) {
	m := make(map[string]*student)
	stus := []student{
		{
			Name: "zhou",
			Age:  21,
		},
		{
			Name: "li",
			Age:  22,
		},
		{
			Name: "zhao",
			Age:  24,
		},
	}
	/**
	这是错误的写法.
	这样的写法初学者经常会遇到的，很危险！
	与Java的foreach一样，都是使用副本的方式。
	所以m[stu.Name]=&stu实际上一致指向同一个指针，
	最终该指针的值为遍历的最后一个struct的值拷贝。
	 */
	//for _, stu := range stus {
	//	m[stu.Name] = &stu
	//}

	/**
		错误写法最终打印的结果为：
		li => zhao
		zhao => zhao
		zhou => zhao
	 */
	//for k,v := range m {
	//	fmt.Println(k,"=>",v.Name)
	//}

	// 下面看正确的写法
	// 正确写法1
	for _, stu := range stus {
		tmp := stu
		m[stu.Name] = &tmp
	}
	// 正确写法2
	for i := 0; i < len(stus); i++ {
		m[stus[i].Name] = &stus[i]
	}
	for k,v := range m {
		fmt.Println(k,"=>",v.Name)
	}
}


// 测试goroutine执行顺序
/**
答案：协程goroutine随机机型的输出，并不知道最终的输出结果是什么
但是可以确定的是，A始终输出的是循环的最后一个元素，也就是10
而B最终的输出结果是0-9
解析：
第一个go func中i是外部for的一个变量，地址不变化。遍历完成后，最终i=10。 故go func执行时，i的值始终是10。
第二个go func中i是函数参数，与外部for中的i完全是两个变量。 尾部(i)将发生值拷贝，go func内部指向值拷贝地址。
 */
func TestGoroutine(t *testing.T){
	// 最大cpu运行
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10 ; i++ {
		go func() {
			fmt.Println("A:", i)
			wg.Done()
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("B:", i)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

// 测试go的组合继承
type People struct {

}
func (p *People) ShowA(){
	fmt.Println("showA")
	p.showB()
}
func (p *People) showB(){
	fmt.Println("showB")
}
type Teacher struct {
	People
}

func (t *Teacher)ShowB(){
	fmt.Println("teacher showB")

}
/**
答案：
showA
showB
 */
 /**
 解析：
 这是Golang的组合模式，可以实现OOP的继承。 被组合的类型People所包含的方法虽然升级成了外部类型Teacher这个组合类型的方法（一定要是匿名字段），但它们的方法(ShowA())调用时接受者并没有发生变化。 此时People类型并不知道自己会被什么类型组合，当然也就无法调用方法时去使用未知的组合者Teacher类型的功能
  */
func TestGoInherit(t *testing.T){
	tmp := Teacher{}
	tmp.ShowA()
}

// select的随机性
/**
解析：
select会随机选择一个可用通用做收发操作。
所以代码是有肯触发异常，也有可能不会。
单个chan如果无缓冲时，将会阻塞。但结合 select可以在多个chan间等待执行。有三点原则：

select 中只要有一个case能return，则立刻执行。
当如果同一时间有多个case均能return则伪随机方式抽取任意一个执行。
如果没有一个case能return则可以执行”default”块
 */
func TestSelect(t *testing.T){
	runtime.GOMAXPROCS(1)
	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)
	int_chan <- 1
	string_chan <- "hello"
	select{
	case value := <- int_chan:
		fmt.Println(value)
	case value := <- string_chan:
		panic(value)
	}
}