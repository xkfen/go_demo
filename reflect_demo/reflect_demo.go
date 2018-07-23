package reflect_demo

import (
	"errors"
	"reflect"
	"fmt"
)

type UserInfo struct {
	Id uint
	Name string
}

type WorkInfo struct {
	Id uint
	Name string
}

func GetInfo(info interface{})(data interface{}, err error){
	v := reflect.ValueOf(info)
	switch v.Kind() {
	case reflect.Struct:
		if "UserInfo" == v.Type().Name() {
			fmt.Println("userinfo")
			data = UserInfo{}
			return
		}else if "WorkInfo" == v.Type().Name() {
			fmt.Println("workinfo")
			data = WorkInfo{}
			return
		}
	}
	return
}