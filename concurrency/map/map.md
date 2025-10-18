# 为什么 recover 无法捕获 map 并发读写错误？
## 两种不同的"异常"机制
1. 普通Panic：可以通过panic()函数触发，也会由一些运行时错误（如空指针解引用、越界访问等）产生。这种 panic 可以被 recover 捕获。
2. Runtime Throw：这是由 Go 运行时系统触发的致命错误（fatal error），不走 panic/recover 机制，会直接终止进程。
## 特殊情况：nil map的并发访问
| Map 状态        | 读写类型   | 运行时行为                | 是否可 recover    |
|-----------------|-----------|--------------------------|------------------|
| nil map         | 读        | 安全，返回零值           | 无需 recover      |
| nil map         | 写        | 普通 panic               | ✅ 可 recover     |
| 已初始化 map    | 并发读写  | Runtime Throw (fatal)    | ❌ 不可 recover   |
## 并发访问map
1. 使用互斥锁（sync.Mutex）保护
```go
var m = make(map[int]int)
var mu sync.Mutex

// 写操作
mu.Lock()
m[key] = value
mu.Unlock()

// 读操作
mu.Lock()
value := m[key]
mu.Unlock()
```
2. 使用sync.Map（Go 1.9+）适合读多写少
```go
var sm sync.Map

// 写操作
sm.Store(key, value)

// 读操作
if value, ok := sm.Load(key); ok {
    // 使用value
}
```
## 其他recover无法捕获的错误
1. 栈溢出（stack overflow）：通常由于无限递归或过大的局部变量导致
2. 内存不足：无法分配请求的内存
3. 其他运行时致命错误：如数据竞争、线程调度问题等