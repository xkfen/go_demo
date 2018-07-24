package reflect_demo

import (
	"reflect"
	"fmt"
	"go_demo/util"
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

func GetInfo1(info interface{})(data interface{}, err error){
	if user, ok := info.(UserInfo); ok {
		//s1.f()
		//s1.g()
		fmt.Println(util.StringifyJson(user))
		data = UserInfo{}
	}

	if work, ok := info.(WorkInfo); ok {
		fmt.Println(util.StringifyJson(work))
		data = WorkInfo{}
	}
	return
}