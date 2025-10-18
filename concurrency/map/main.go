package main

import (
	"fmt"
	"time"
)

func main() {
	m := make(map[int]int)
	go func() {
		for i := 0; i < 1000; i++ {
			m[i] = i // 并发写
		}
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recover")
			}
		}()
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			_ = m[i] // 并发读
		}
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recover")
			}
		}()
	}()
	time.Sleep(10 * time.Second)
}
