package main

import "fmt"

func maintest() {
	s := []int{1, 2, 3, 4, 5}

	fmt.Println(len(s))

	func(s *[]int) {
		*s = append(*s, 6)
		fmt.Println(len(*s))
	}(&s)

	fmt.Println(len(s))
}
