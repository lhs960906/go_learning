# 一、Cobra 是什么
Cobra 是 Go 的 CLI 框架。它包含一个用于创建功能强大的现代 CLI 应用程序的库，以及一个用于快速生成基于 Cobra 的应用程序和命令文件的工具。
Cobra 由 Go 项目成员和 hugo 作者 spf13 创建，已经被许多流行的 Go 项目采用，比如 GitHub CLI 和 Docker CLI。

# 二、安装和使用
安装：
```shell
$ go get -u github.com/spf13/cobra/cobra
```
使用：
```go
import (
    "github.com/spf13/cobra/cobra"
)
```

# 三、API 使用示例
## 3.1、Flag 相关 API
### 3.1.1、Persistent Flags 和 Local Flags
如果你想给你的命令绑定一个 Persistent Flags，可以使用如下 API：
```go
rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
```
如果你想给你的命令绑定一个 Local Flags：可以使用如下 API：
```go
localCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
```
### 3.1.2、必传 Flags
默认情况下，Flags 是可选的。相反，如果您希望命令在未设置 Flag 时报告错误，请根据需要将其 Flag 为 required：
```go
# 对于本地标记
rootCmd.Flags().StringVarP(&Region, "region", "r", "", "AWS region (required)")
rootCmd.MarkFlagRequired("region")

# 对于持久化标记
rootCmd.PersistentFlags().StringVarP(&Region, "region", "r", "", "AWS region (required)")
rootCmd.MarkPersistentFlagRequired("region")
```
### 3.1.3、Flags 组
#### 3.1.3.1、共存
如果你有不同的 Flags 需要被一起提供（例如：如果你提供了 --username 你就必须同时提供 --password），Cobra 可以强制执行该要求：
```go
rootCmd.Flags().StringVarP(&u, "username", "u", "", "Username (required if password is set)")
rootCmd.Flags().StringVarP(&pw, "password", "p", "", "Password (required if username is set)")
rootCmd.MarkFlagsRequiredTogether("username", "password")
```
#### 3.1.3.2、互斥
如果不同的 Flags 表示互斥的选项，例如将输出格式指定为 --json 或 --yaml 但不能同时指定两者，您还可以防止一起提供不同的 Flags：
```go
rootCmd.Flags().BoolVar(&ofJson, "json", false, "Output in JSON")
rootCmd.Flags().BoolVar(&ofYaml, "yaml", false, "Output in YAML")
rootCmd.MarkFlagsMutuallyExclusive("json", "yaml")
```
#### 3.1.3.3、存一
如果您希望一组中至少存在一个 Flag，则可以使用 MarkFlagsOneRequired。这可以与 MarkFlagsMutuallyExclusive 结合使用，以强制执行给定组中的一个 Flag：
```go
rootCmd.Flags().BoolVar(&ofJson, "json", false, "Output in JSON")
rootCmd.Flags().BoolVar(&ofYaml, "yaml", false, "Output in YAML")
rootCmd.MarkFlagsOneRequired("json", "yaml")
rootCmd.MarkFlagsMutuallyExclusive("json", "yaml")
```
在这些情况下：
* 本地和持久 Flags 都可以使用
    * 注意：该组仅在定义了每个 Flag 的命令上强制执行
* 一个 Flag 可能出现在多个组中
* 一个组可以包含任意数量的 Flag

## 3.4、help 命令

## 3.5、USAGE 信息

## 3.6、不合法命令出现时的建议信息
当发生 "未知命令" 错误时，Cobra 将打印自动建议。当发生拼写错误时，这使得 Cobra 的行为与 git 命令类似。例如：
```shell
$ hugo srever
Error: unknown command "srever" for "hugo"

Did you mean this?
        server

Run 'hugo --help' for usage.
```
建议是根据现有子命令自动生成的，并使用 Levenshtein distance 的实现。每个匹配最小距离 2（忽略大小写）的已注册命令都将显示为建议。

如果您需要禁用建议或调整命令中的字符串距离，请使用：
```go
command.DisableSuggestions = true
```
或者：
```go
command.SuggestionsMinimumDistance = 1
```
还可以使用 SuggestFor 属性显式设置建议给定命令的名称。这允许对字符串距离不接近、但在您的命令集中有意义但您不需要别名的字符串提供建议。例子：
```shell
$ kubectl remove
Error: unknown command "remove" for "kubectl"

Did you mean this?
        delete

Run 'kubectl help' for usage.
```


# 附录
## 附录A：术语表
### flags
标志提供修饰符来控制操作命令的运行方式。
在 Cobra 中，有如下几种类型的 Flag：
* Persistent Flags
* Local Flags

其中 `Persistent Flags` 是意味着该 Flag 将可用于分配给它的命令以及该命令下的每个子命令。对于全局 Flag，将 Flag 分配为根上的 Persistent Flags，也即是对所有命令都生效的 Flag，例如 Cobra 中的 `--help` 就是这样一个全局 Flag。

## 附录B：参考
* [^1]: [Cobra in github](https://github.com/spf13/cobra)
* [^2]: [Cobra 用户手册](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md)