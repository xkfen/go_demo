package go_demo

import (
	"testing"
	"runtime"
	"sync"
	"fmt"
)

func TestSync(t *testing.T){
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{}
	wg.Add(12)
	for i := 0; i < 6; i++ {
		go func() {
			fmt.Println("T1:", i)
			wg.Done()
		}()
	}
	for i := 0; i < 6; i++ {
		go func(i int) {
			fmt.Println("T2:",i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}