## 基本语法

```go
//go:linkname localname [importpath.name]
```

## Pull模式（拉取模式）
在Pull模式下，当前包使用go:linkname指令主动链接到其他包中的私有函数。
```go
// 在foo/foo.go中
package foo
import (
    _ "unsafe"
    _ "github.com/example/bar"
)

//go:linkname Add github.com/example/bar.add
func Add(a, b int) int

// 在bar/bar.go中
package bar
func add(a, b int) int { 
    return a + b
}
```
## Push模式（推送模式）
```go
// 在bar/bar.go中
package bar
import _ "unsafe"

//go:linkname div github.com/example/foo.Div
func div(a, b int) int {
    return a / b
}

// 在foo/foo.go中
package foo
import _ "github.com/example/bar"

func Div(a, b int) int
```
## Handshake模式（握手模式）
```go
// 在bar/bar.go中
package bar
import _ "unsafe"

//go:linkname hello
func hello(name string) string {
    return "Hello " + name + "!"
}

// 在foo/foo.go中
package foo
import (
    _ "unsafe"
    _ "github.com/example/bar"
)

//go:linkname Hello github.com/example/bar.hello
func Hello(name string) string
```