package main

import (
	"fmt"
	"reflect"
)

// 1. 当结构体中的字段都是可比较的，那么这个结构体就可以用==和!=进行比较
// 1.1 结构体名称需要相同，字段名称和类型需要相同，字段顺序需要相同
// 1.2 匿名结构体只要字段相同，顺序相同，就可以比较
// 2. 不可比较的类型有：slice, map, function
// 3. 不可比较的结构体可使用reflect.DeepEqual进行比较，但是性能较差，且无法处理循环引用

// 可比较的结构体示例
type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Name string
	Age  int
}

// 不可比较的结构体示例（包含slice字段）
type Company struct {
	Name      string
	Employees []string // slice字段，不可比较
}

// 不可比较的结构体示例（包含map字段）
type Config struct {
	Host string
	Port int
	Data map[string]interface{} // map字段，不可比较
}

// 不可比较的结构体示例（包含function字段）
type Handler struct {
	Name string
	Func func() // function字段，不可比较
}

func StructCompareExamples() {
	fmt.Println("=== 结构体比较示例 ===")

	// 1. 可比较的结构体比较
	fmt.Println("\n1. 可比较的结构体比较:")
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 30}
	p3 := Person{Name: "Bob", Age: 25}

	fmt.Printf("p1 == p2: %t\n", p1 == p2) // true
	fmt.Printf("p1 == p3: %t\n", p1 == p3) // false
	fmt.Printf("p1 != p3: %t\n", p1 != p3) // true

	// 1.1 不同结构体类型不能比较（编译错误）
	// emp := Employee{Name: "Alice", Age: 30}
	// fmt.Printf("p1 == emp: %t\n", p1 == emp) // 编译错误

	// 1.2 匿名结构体比较
	fmt.Println("\n1.2 匿名结构体比较:")
	anon1 := struct {
		Name string
		Age  int
	}{Name: "Charlie", Age: 35}

	anon2 := struct {
		Name string
		Age  int
	}{Name: "Charlie", Age: 35}

	// anon3 := struct {
	// 	Age  int
	// 	Name string
	// }{Name: "Charlie", Age: 35} // 字段顺序不同

	fmt.Printf("anon1 == anon2: %t\n", anon1 == anon2) // true
	// fmt.Printf("anon1 == anon3: %t\n", anon1 == anon3) // 编译错误：字段顺序不同

	// 2. 不可比较的结构体示例
	fmt.Println("\n2. 不可比较的结构体示例:")
	company1 := Company{Name: "TechCorp", Employees: []string{"Alice", "Bob"}}
	company2 := Company{Name: "TechCorp", Employees: []string{"Alice", "Bob"}}

	// 以下代码会编译错误，因为Company包含slice字段
	// fmt.Printf("company1 == company2: %t\n", company1 == company2)
	fmt.Println("company1 == company2: 编译错误，因为包含slice字段")

	// 3. 使用reflect.DeepEqual进行比较
	fmt.Println("\n3. 使用reflect.DeepEqual进行比较:")
	fmt.Printf("DeepEqual(company1, company2): %t\n", reflect.DeepEqual(company1, company2)) // true

	// 修改slice内容
	company2.Employees = append(company2.Employees, "Charlie")
	fmt.Printf("修改后 DeepEqual(company1, company2): %t\n", reflect.DeepEqual(company1, company2)) // false

	// map字段比较
	config1 := Config{
		Host: "localhost",
		Port: 8080,
		Data: map[string]interface{}{"key1": "value1"},
	}
	config2 := Config{
		Host: "localhost",
		Port: 8080,
		Data: map[string]interface{}{"key1": "value1"},
	}
	fmt.Printf("DeepEqual(config1, config2): %t\n", reflect.DeepEqual(config1, config2)) // true

	// function字段比较
	handler1 := Handler{
		Name: "handler1",
		Func: func() { fmt.Println("Hello") },
	}
	handler2 := Handler{
		Name: "handler1",
		Func: func() { fmt.Println("Hello") },
	}
	fmt.Printf("DeepEqual(handler1, handler2): %t\n", reflect.DeepEqual(handler1, handler2)) // false，因为function不能比较

	// 4. 循环引用问题
	fmt.Println("\n4. 循环引用问题:")
	type Node struct {
		Value int
		Next  *Node
	}

	node1 := &Node{Value: 1, Next: nil}
	node1.Next = node1 // 创建循环引用

	node2 := &Node{Value: 1, Next: nil}
	node2.Next = node2 // 创建循环引用

	// DeepEqual无法处理循环引用，会导致栈溢出
	fmt.Println("DeepEqual(node1, node2): 会导致栈溢出，无法处理循环引用")
}

// 比较结构体字段的辅助函数
func CompareStructFields() {
	fmt.Println("\n=== 结构体字段比较辅助函数 ===")

	// 比较两个Person结构体的字段
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 25}

	fmt.Printf("p1: %+v\n", p1)
	fmt.Printf("p2: %+v\n", p2)

	// 逐个字段比较
	if p1.Name == p2.Name {
		fmt.Println("Name字段相同")
	} else {
		fmt.Println("Name字段不同")
	}

	if p1.Age == p2.Age {
		fmt.Println("Age字段相同")
	} else {
		fmt.Println("Age字段不同")
	}
}

// 自定义比较函数示例
func (p Person) Equals(other Person) bool {
	return p.Name == other.Name && p.Age == other.Age
}

func CustomCompareExample() {
	fmt.Println("\n=== 自定义比较函数示例 ===")

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Alice", Age: 30}
	p3 := Person{Name: "Bob", Age: 30}

	fmt.Printf("p1.Equals(p2): %t\n", p1.Equals(p2)) // true
	fmt.Printf("p1.Equals(p3): %t\n", p1.Equals(p3)) // false
}

// 主函数，运行所有示例
func main() {
	StructCompareExamples()
	CompareStructFields()
	CustomCompareExample()
}
