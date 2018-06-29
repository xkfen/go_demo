package go_demo

import (
	"testing"
	"fmt"
	"strconv"
	"github.com/stretchr/testify/assert"
	"strings"
)

// 字符串测试
func TestString(t *testing.T){
	// int类型所占的位数  64位
	fmt.Printf("%d\n",strconv.IntSize)
	// 字符串转int
	intStr := "21"
	intValue, err := strconv.Atoi(intStr)
	fmt.Printf("%d\n",intValue)
	floatStr := "21.4"
	f, err := strconv.ParseFloat(floatStr, 64)
	fmt.Printf("%f\n", f)
	// int 转 string
	num := 100
	fmt.Printf("%s\n",strconv.Itoa(num))
	// 字符串比较
	str1 := "test1   "
	str2 := "test2"
	com := strings.Compare(str1, str2)
	fmt.Printf("%d\n",com)
	// 查找，包含
	flag := strings.Contains(str1, str2)
	fmt.Println(flag)
	// 查找位置，找不到返回-1
	index := strings.Index(str1, "t")
	lastIndex := strings.LastIndex(str1, "2")
	fmt.Println(index)
	fmt.Println(lastIndex)
	// 统计给定字符串出现的次数
	fmt.Println(strings.Count(str1, "e"))
	// 替换：在s字符串中, 把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换
	fmt.Println(strings.Replace(str1, "1", "2324", -1))
	// 去除字符串头尾字符串
	fmt.Println(strings.Trim(str1, "1"))
	// 去除字符串左边字符串
	fmt.Println(strings.TrimLeft(str1, "t"))
	// 删除开头末尾的空格
	fmt.Println(strings.TrimSpace(str1))
	// 字符串  首字母大写
	fmt.Println(strings.Title(str1))
	// 小写
	fmt.Println(strings.ToLower(str1))
	// 大写
	fmt.Println(strings.ToUpper(str1))
	//前缀 后缀
	fmt.Println(strings.HasPrefix("Gopher", "Go")) // true
	fmt.Println(strings.HasSuffix("Amigo", "go"))  // true

	fieldsStr := "  hello   it's  a  nice day today    "
	//根据空白符分割,不限定中间间隔几个空白符
	fieldsSlece := strings.Fields(fieldsStr)
	fmt.Println(fieldsSlece) //[hello it's a nice day today]

	//根据特定字符分割
	slice01 := strings.Split("q,w,e,r,t,y,", ",")
	fmt.Println(slice01)      //[q w e r t y ]
	fmt.Println(cap(slice01)) //7  最后多个空""

	//拼接
	//Join 用于将元素类型为 string 的 slice, 使用分割符号来拼接组成一个字符串：
	var str08  = strings.Join(fieldsSlece, ",")
	fmt.Println("Join拼接结果=" + str08) //hello,it's,a,nice,day,today


	assert.NoError(t, err)
}