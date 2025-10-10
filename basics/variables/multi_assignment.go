package main

import "fmt"

// 多重赋值示例
// 演示Go语言中多重赋值的特性和执行顺序
func main() {
	// 示例1：基本多重赋值
	i := 1
	arr := []int{0, 0, 0}

	// 注意：Go中多重赋值的执行顺序是从左到右
	// 这里会先计算右边的表达式，然后按顺序赋值
	i, arr[i], arr[i-1] = i+1, i+2, i+3

	fmt.Println("i =", i)     // 输出: i = 2
	fmt.Println("arr =", arr) // 输出: arr = [4 3 0]

	// 解释：
	// 1. 先计算右边：i+1=2, i+2=3, i-1=0, i+3=4
	// 2. 然后赋值：i=2, arr[1]=3, arr[0]=4
	// 3. 所以最终arr = [4, 3, 0]

	fmt.Println("\n--- 更多多重赋值示例 ---")

	// 示例2：交换变量值
	a, b := 10, 20
	fmt.Printf("交换前: a=%d, b=%d\n", a, b)
	a, b = b, a // 交换a和b的值
	fmt.Printf("交换后: a=%d, b=%d\n", a, b)

	// 示例3：函数返回值的多重赋值
	x, y := divide(10, 3)
	fmt.Printf("10除以3: 商=%d, 余数=%d\n", x, y)

	// 示例4：切片的多重赋值
	slice := []int{1, 2, 3, 4, 5}
	fmt.Println("原始切片:", slice)
	slice[0], slice[4] = slice[4], slice[0] // 交换首尾元素
	fmt.Println("交换后:", slice)
}

// 返回商和余数的函数
func divide(a, b int) (int, int) {
	return a / b, a % b
}
