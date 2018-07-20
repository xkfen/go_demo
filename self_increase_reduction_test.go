package go_demo

import (
	"testing"
	"fmt"
)

func TestSelfIncreaseReduction(t *testing.T){
	v := 1
	//v++
	v--
	fmt.Printf("%d\n", v)

}