# Go 内存模型
## 1.1、介绍
### 1.1.1、忠告
### 1.1.2、非正式概览
## 1.2、内存模型
## 1.3、包含数据竞争程序的实现限制
## 1.4、同步
### 1.4.1、初始化
### 1.4.2、Goroutine 创建
### 1.4.3、Goroutine 销毁
### 1.4.4、Channel 通信
### 1.4.5、Lock 类型
sync包提供了两种lock数据类型，`sync.Mutex` 和 `sync.RWMutex`。
假设有 `n < m`，那么对于任意的 `sync.Mutex` 或 `sync.RWMutex` 变量 `l`，`l.Unlock()` 的第 n 次调用在 `l.Lock()` 的第 m 次调用返回之前同步。
```go
var l sync.Mutex
var a string

func f() {
	a = "hello, world"
	l.Unlock()
}

func main() {
	l.Lock()
	go f()
	l.Lock()
	print(a)
}
```
程序保证会打印 `hello, world`。对 `l.Unlock()`（f 中）的第一次调用在第二次调用 `l.Lock()`（main 中）返回之前同步，该调用在打印之前排序。

对于在 `sync.RWMutex` 变量 `l` 上对 `l.RLock` 的任何调用，有一个 `n`，使得对 `l.Unlock` 的第 `n` 次调用在 `l.RLock` 返回之前同步，并且对 `l.RUnlock` 的匹配调用在 `l.RLock` 返回之前同步。 从调用 `n+1` 返回到 `l.Lock`。

成功调用 `l.TryLock`（或 `l.TryRLock`）相当于调用 `l.Lock`（或 `l.RLock`）。不成功的调用根本没有同步效果。就内存模型而言，`l.TryLock`（或 `l.TryRLock`）可以被认为能够返回 false，即使互斥体 `l` 已解锁。

### 1.4.6、Once 类型
### 1.4.7、Atomic Values
### 1.4.8、Finalizers
### 1.4.9、附加机制
## 1.5、同步错误
## 1.6、编译错误
## 1.7、总结