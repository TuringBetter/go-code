# CAS
CAS操作包含三个关键参数：

- 内存位置（V）：需要读写的内存值
- 预期的原值（A）：进行比较的值
- 新值（B）：拟写入的新值

CAS的操作逻辑是：如果内存位置V的值等于预期原值A，则将位置V的值修改为新值B，否则不做任何操作。整个操作是原子性的，在多线程环境下不会被中断。

```go
func CompareAndSwap(addr *T, old, new T) bool {
    if *addr == old {
        *addr = new
        return true
    }
    return false
}
```
# CAS 工作原理
CAS的核心思想是**乐观并发控制**：它假设操作不会发生冲突，每次操作时都不加锁，而是尝试更新。如果发现值已经被其他线程修改，则操作失败，通常会选择重试。
# 典型场景
计数器
```go
type Counter struct {
    value int32
}

func (c *Counter) Increment() {
    for {
        oldValue := atomic.LoadInt32(&c.value)
        newValue := oldValue + 1
        if atomic.CompareAndSwapInt32(&c.value, oldValue, newValue) {
            return
        }
    }
}

func (c *Counter) Value() int32 {
    return atomic.LoadInt32(&c.value)
}

```
自旋锁
```go
type SpinLock struct {
    flag int32
}

func (sl *SpinLock) Lock() {
    for !atomic.CompareAndSwapInt32(&sl.flag, 0, 1) {
        // 可加入runtime.Gosched()避免过度占用CPU
        runtime.Gosched()
    }
}

func (sl *SpinLock) Unlock() {
    atomic.StoreInt32(&sl.flag, 0)
}

```
无锁结构
标志状态更新
# 优缺点
## 优点
- 无锁操作：避免线程阻塞和上下文切换，提高并发性能
- 避免死锁：由于不使用锁，自然避免了死锁问题
- 轻量级：相比互斥锁消耗更少资源
- 高性能：在低竞争环境下，性能比锁更高
## 缺点
- ABA问题：值可能从A变为B又变回A，CAS会误认为没有变化。解决方法通常是使用版本号或标记指针。
- 自旋开销：在高竞争环境下，大量CPU时间可能浪费在重试操作上
- 只能保护单个变量：复杂操作需要多个CAS组合，增加实现复杂度
- 实现复杂：相比锁机制，无锁算法设计更加复杂