# 一、反射是什么

# 二、反射的原理


# 三、反射相关API
```go
func TypeOf(i interface{}) Type
func ValueOf(i interface{}) Value
```

# 四、反射相关示例

# 五、结论
* Go 中的接口是静态类型，接口类型的变量始终具有相同的静态类型，即使在运行时存储在接口变量中的值可能会更改具体类型，该值也将始终满足接口静态类型的要求。

# 参考
[The Laws of Reflection - by Rob Pike](https://go.dev/blog/laws-of-reflection)
[Go Data Structures: Interfaces - by Russ Cox](https://research.swtch.com/interfaces)
# 术语
静态类型：
接口变量：
接口值：