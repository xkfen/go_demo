package sync

import (
	"testing"
	"sync"
	"fmt"
)
/**
在这个例子中，goroutine正在运行一个已经关闭迭代变量salutation的闭包，它有一个字符串类型。 当我们的循环迭代时，salutation被分配给切片中的下一个字符串值。 由于运行时调度器安排的goroutine可能会在将来的任何时间点运行，因此不确定在goroutine内将打印哪些值。
在这个例子中，循环在任何goroutines开始运行之前退出，所以salutation转移到堆中，并保存对字符串切片“good day”中最后一个值的引用。所以会看到“good day”打印三次
Go运行时足够敏锐地知道对str变量的引用仍然保留，因此会将内存传输到堆中，以便goroutine可以继续访问它。
 */
func TestSync1(t *testing.T){
	var wg sync.WaitGroup
	strs := []string{"katy", "darling", "good day"}
	for _, str := range strs {
		// 循环在任何goroutines开始运行之前退出
		fmt.Println("循环")
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(str)
		}()
	}
	wg.Wait()
}

/**
编写该循环的正确方法是将salutation的副本传递给闭包，以便在运行goroutine时，它将对来自其循环迭代的数据进行操作：
 */
func TestSync2(t *testing.T){
	var wg sync.WaitGroup
	strs := []string{"katy", "darling", "good day"}
	for _, str := range strs {
		wg.Add(1)
		tmp := str
		go func() {
			defer wg.Done()
			fmt.Println(tmp)
		}()
	}
	wg.Wait()
}

func TestSync3(t *testing.T){
	var wg sync.WaitGroup
	strs := []string{"katy", "darling", "good day"}
	for i, _ := range strs {
		wg.Add(1)
		tmp := strs[i]
		go func() {
			defer wg.Done()
			fmt.Println(tmp)
		}()
	}
	wg.Wait()
}

func TestSync4(t *testing.T){
	var wg sync.WaitGroup
	strs := []string{"katy", "darling", "good day"}
	for _, str := range strs {
		wg.Add(1)
		go func(str string) {
			// 闭包函数的参数str是任意的参数
			defer wg.Done()
			fmt.Println(str)
		}(str) // 这里这个str是循环变量的副本
	}
	wg.Wait()
}