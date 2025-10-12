package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		fmt.Println("child start:make panic")
		panic("panic")
	}()
	time.Sleep(1 * time.Second)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover")
		}
	}()
}
