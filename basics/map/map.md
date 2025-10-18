| 类型                | 是否可作为 key | 说明                          |
|---------------------|:-------------:|-------------------------------|
| int / string / bool | ✅             | 内建可比较类型                |
| float / complex     | ✅             | 虽可比较，但注意精度陷阱      |
| pointer / chan      | ✅             | 指向的值可比较即可           |
| interface           | ✅             | 包含的值必须全部可比较        |
| array               | ✅             | 所有元素类型可比较时可用      |
| struct              | ✅ / ❌        | 取决于字段是否都可比较        |
| slice / map / func  | ❌             | 天然不可比较                  |

通道值是可比较的。两个通道值相等当且仅当它们是由同一个make调用创建的或者两个值都为nil。

```go
ci1 := make(chan int, 1)
ci2 := ci1
ci3 := make(chan int, 1)

fmt.Println(ci1 == ci2) // 输出: true
fmt.Println(ci1 == ci3) // 输出: false
```