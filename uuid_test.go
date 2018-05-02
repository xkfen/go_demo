package go_demo

import (
	"testing"
	"github.com/satori/go.uuid"
	"fmt"
)

func TestGenerateUuid(t *testing.T){
	// 创建:NewV4 returns random generated UUID
	u1, err := uuid.NewV4()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("UUIDv4: %s\n", u1)

	// 解析
	u2, err := uuid.FromString("f5394eef-e576-4709-9e4b-a7c231bd34a4")
	if err != nil {
		fmt.Printf("Something gone wrong: %s", err)
		return
	}
	fmt.Printf("Successfully parsed: %s", u2)
}