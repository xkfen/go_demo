package reflect_demo

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"encoding/json"
)

func TestGetUserInfo(t *testing.T) {
	info, err := GetInfo(UserInfo{})
	assert.NoError(t, err)
	byteData, _ := json.Marshal(info)
	fmt.Println(string(byteData))
}

func TestGetWorkInfo(t *testing.T) {
	info, err := GetInfo(WorkInfo{})
	assert.NoError(t, err)
	byteData, _ := json.Marshal(info)
	fmt.Println(string(byteData))
}

func TestGetInfo1(t *testing.T) {
	info, err := GetInfo1(UserInfo{})
	assert.NoError(t, err)
	byteData, _ := json.Marshal(info)
	fmt.Println(string(byteData))
}