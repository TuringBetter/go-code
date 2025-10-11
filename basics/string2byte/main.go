package main

import "fmt"

func main() {
	// 在Go语言中，字符串本质上就是只读的字节序列，而[]byte是可变的字节序列，它们在内存中的布局是兼容的。当使用...语法时，Go会将字符串中的字节展开，相当于将字符串转换为[]byte后再追加。
	data1 := append([]byte("Hello"), " world"...)
	fmt.Println(string(data1))
	// 这与显示转换为[]byte后再追加是等价的
	data2 := append([]byte("Hello"), []byte(" world")...)
	fmt.Println(string(data2))
}
