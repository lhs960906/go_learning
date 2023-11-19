# 文件打开模式
文件有如下几种打开模式：
* os.O_CREATE：
* os.O_APPEND：
* os.O_RDONLY：
* os.O_WRONLY：
* os.O_RDRW：
* os.O_SYNC：
* os.O_TRUNC：
* os.O_EXCL：

多个模式可以使用 `|` 组合进行使用，比如：
* `os.O_RDONLY | os.O_CREATE`：代表以只读方式打开文件，如果文件不存在则进行创建