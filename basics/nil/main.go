package main

import "fmt"

func main() {
	var a *int
	fmt.Println(a == nil)
}

/*
只有引用类型可以赋值为nil slice map function channel
值类型不能赋值为nil int bool string array struct
*/
