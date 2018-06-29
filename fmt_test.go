package go_demo

import (
	"testing"
	"fmt"
	"os"
)

// 测试golang的格式化输出

type Point struct {
	x, y int
}

func TestFmt(t *testing.T){
	p := Point{1, 2}
	// %v : 打印结构体实例 输出：{1 2}
	fmt.Printf("%v\n", p)
	// % + v 如果值是结构体，%+v格式化输出内容将打印包括结构体的字段名
	fmt.Printf("%+v\n", p) // 输出：{x:1 y:2}
	// % # v :输出这个值的go语法表示，例如，值的运行源代码片段
	fmt.Printf("%#v\n", p) // 输出：go_demo.Point{x:1, y:2}
	// %T 打印值的类型
	fmt.Printf("%T\n", p) // 输出：go_demo.Point
	// %t 格式化bool值
	fmt.Printf("%t\n", true)
	// %d 标准十进制
	fmt.Printf("%d\n", 123)
	// %b 二进制
	fmt.Printf("%b\n", 14)
	// %c 输出给定整数对应的字符
	fmt.Printf("%c\n", 33) // 输出： !
	// %x 十六进制
	fmt.Printf("%x\n", 78) // 输出：4e
	// %o 八进制
	//fmt.Printf("%o", 12)
	// %f 浮点型 十进制
	fmt.Printf("%f\n", 78.9)
	// %e与%E 科学计数表示
	fmt.Printf("%e\n", 12000000.0)
	fmt.Printf("%E\n", 1200000.0)
	// %s 基本字符串输出
	fmt.Printf("%s\n", "\"string\"")
	// %q 输出双引号
	fmt.Printf("%q\n", "\"string\"")
	// %x 输出使用base-64 编码的字符串，每个字节使用两个字符表示
	fmt.Printf("%x\n", "hex this")
	// %p 输出指针的值
	fmt.Printf("%p\n", &p) // 输出：0xc420020430
	// 当输出数字的时候，你将经常想要控制输出结果的宽度和精度，可以使用在 % 后面使用数字来控制输出宽度。默认结果使用右对齐并且通过空格来填充空白部分。
	fmt.Printf("%6d%6d\n", 12, 345)
	// (右对齐)可以指定浮点型的输出宽度，同时也可以通过 宽度.精度 的语法来指定输出的精度。
	fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)
	//要最左对齐，使用 - 标志。
	fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45)
	//你也许也想控制字符串输出时的宽度，特别是要确保他们在类表格输出时的对齐。这是基本的右对齐宽度表示。
	fmt.Printf("|%6s|%6s|\n", "foo", "b")
	//要左对齐，和数字一样，使用 - 标志。
	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")
	//到目前为止，我们已经看过 Printf了，它通过 os.Stdout输出格式化的字符串。Sprintf 则格式化并返回一个字符串而不带任何输出。
	s := fmt.Sprintf("a %s", "string")
	fmt.Println(s)
	//你可以使用 Fprintf 来格式化并输出到 io.Writers而不是 os.Stdout。
	fmt.Fprintf(os.Stderr, "an %s\n", "error")
}