package main

import "fmt"

func main() {
	// utf-8 编码
	// ASCII字符占一个字节
	// 中文占三个字节
	str := "hello世界"

	// 使用普通for循环遍历 基于len()索引遍历
	fmt.Println("使用普通for循环遍历:")
	for i := 0; i < len(str); i++ {
		fmt.Printf("%c", str[i])
	}

	fmt.Println("\n使用range遍历:")
	// 使用range遍历 自动解码utf-8 每次迭代返回一个rune
	for _, char := range str {
		fmt.Printf("%c", char)
	}
	fmt.Println()
	// byte uint8别名
	// rune int32别名
}
