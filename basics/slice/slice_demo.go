package main

import "fmt"

func main() {
	// 演示slice传值的情况
	s := []int{1, 2, 3, 4, 5}
	fmt.Println("原始slice:", s, "长度:", len(s))

	// 1. 传值方式 - 修改元素内容
	modifyElements(s)
	fmt.Println("修改元素后:", s)

	// 2. 传值方式 - append操作
	appendElements(s)
	fmt.Println("append后:", s, "长度:", len(s))

	// 3. 传指针方式 - append操作
	appendElementsByPointer(&s)
	fmt.Println("传指针append后:", s, "长度:", len(s))
}

// 传值方式修改slice元素
func modifyElements(s []int) {
	// 这会修改外部slice，因为共享底层数组
	s[0] = 999
	fmt.Println("函数内修改元素:", s)
}

// 传值方式append
func appendElements(s []int) {
	// 这不会影响外部slice，因为s是副本
	s = append(s, 6, 7, 8)
	fmt.Println("函数内append:", s, "长度:", len(s))
}

// 传指针方式append
func appendElementsByPointer(s *[]int) {
	// 这会修改外部slice，因为操作的是原始slice
	*s = append(*s, 6, 7, 8)
	fmt.Println("函数内传指针append:", *s, "长度:", len(*s))
}
