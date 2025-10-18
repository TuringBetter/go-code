package main

import (
	"fmt"
	"unsafe"
)

type Person struct{}

type Person2 struct {
	person Person
}

type Person3 struct {
	person *Person
}

func main() {
	p1 := Person{}
	p2 := Person{}
	p3 := Person2{person: p1}
	p4 := Person3{person: &p1}
	fmt.Printf("&p1=%p,size=%d\n", &p1, unsafe.Sizeof(p1))
	fmt.Printf("&p2=%p,size=%d\n", &p2, unsafe.Sizeof(p2))
	fmt.Printf("&p3=%p,size=%d\n", &p3, unsafe.Sizeof(p3))
	fmt.Printf("&p4=%p,size=%d\n", &p4, unsafe.Sizeof(p4))
	fmt.Printf("p1==p2=%t\n", p1 == p2)
	fmt.Printf("p1==p4=%t\n", p1 == *p4.person)

}
