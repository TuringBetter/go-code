package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

// 使用unsafe直接转换（Go 1.20+推荐）
// string到[]byte的零拷贝转换
func StringToBytes1(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// []byte到string的零拷贝转换
func BytesToString1(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// 使用reflect.Header（兼容旧版本）
func StringToBytes2(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}

	return *(*[]byte)(unsafe.Pointer(&bh))
}

func BytesToString2(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}

	return *(*string)(unsafe.Pointer(&sh))
}
func main() {
	s := "hello, gopher"
	ps := unsafe.StringData(s)

	b := []byte(s) // 标准转换,会有损耗
	pb := unsafe.SliceData(b)

	fmt.Printf("ps=%p pb=%p equal=%v\n", ps, pb, ps == pb)
	// 输出：equal=false，证明指针地址不同

	// 对于字符串读取，使用strings.Reader
	s1 := "large string data"
	r := strings.NewReader(s1) // 零拷贝
	io.Copy(os.Stdout, r)

	// 对于字符串构建，使用strings.Builder
	var builder strings.Builder
	builder.Grow(1024) // 预分配空间
	builder.WriteString("prefix")
	result := builder.String() // 仅一次分配
	fmt.Println(result)

}

/*
零拷贝转换虽然高效，但风险也很大：

破坏string不可变性：Go语言规定string是不可变的，但通过零拷贝转换得到的[]byte是可变的。如果修改这些字节，可能破坏语言规范。

内存安全问题：如果原string或[]byte已被回收，访问转换后的数据可能导致程序崩溃。

兼容性问题：这种方法依赖于Go内部实现，未来版本如有变化可能失效。
*/
