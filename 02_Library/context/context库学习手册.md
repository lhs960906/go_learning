# 概述
context 于 Go 1.7 引入 context

Contxt 接口
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <- chan struct{}
    Err() error
    Value(key any) any
}
```

# Context 三大功能
Cancel
Deadline
上下文值

# Context 键值


# Contxt 经典错误


参考
# [Context 接口](https://www.bilibili.com/video/BV19C4y1o79b/?spm_id_from=333.337.search-card.all.click&vd_source=a7c46d733a95cbe7a825e7d12877b8b9)