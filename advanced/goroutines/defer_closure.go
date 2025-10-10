package main

import (
	"fmt"
	"time"
)

// defer和闭包示例
// 演示defer语句与闭包的交互，以及常见的陷阱
func main() {
	fmt.Println("=== defer与闭包的基本示例 ===")

	// 示例1：defer中的闭包陷阱
	fmt.Println("示例1 - defer闭包陷阱:")
	for i := 0; i < 3; i++ {
		defer func() {
			fmt.Printf("defer中的i: %d\n", i)
		}()
	}
	// 输出结果：所有defer都会打印3，因为闭包捕获的是变量的引用

	fmt.Println("\n示例2 - 正确的defer闭包使用:")
	for i := 0; i < 3; i++ {
		defer func(val int) { // 通过参数传递值
			fmt.Printf("defer中的val: %d\n", val)
		}(i)
	}

	fmt.Println("\n示例3 - 另一种正确方式:")
	for i := 0; i < 3; i++ {
		i := i // 创建局部变量
		defer func() {
			fmt.Printf("defer中的局部i: %d\n", i)
		}()
	}

	// 等待defer执行
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n=== defer的执行顺序 ===")
	defer fmt.Println("defer 1")
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3")
	fmt.Println("正常执行")
	// 输出顺序：正常执行 -> defer 3 -> defer 2 -> defer 1

	fmt.Println("\n=== defer与返回值 ===")
	fmt.Println("返回值测试:", testReturn())
}

// 测试defer对返回值的影响
func testReturn() (result int) {
	defer func() {
		result++ // defer可以修改命名返回值
	}()
	return 1 // 实际返回2
}
