## 3 Protocol Parameters(协议参数)
### 3.12 Range Units(Range单元)
```shell
HTTP/1.1允许客户端请求响应中的部分(一定范围)响应实体。
HTTP/1.1在Range(第14.35节)和Content-Range(第14.16节)标头字段中使用范围单位。一个实体可以根据不同的结构单元分解为多个子范围。
	range-unit       = bytes-unit | other-range-unit
	bytes-unit       = "bytes"
	other-range-unit = token
HTTP/1.1定义的唯一范围单位是"bytes"。HTTP/1.1的实现可能会忽略其他单元指定的范围。HTTP/1.1 has been designed to allow implementations of applications that do not depend on knowledge of ranges.
```


## 14 Header Field Definitions(头字段定义)
### 14.5 Accept-Ranges
```shell
Accept-Ranges 响应头允许服务端指示请求资源的可接受范围：
	Accept-Ranges = "Accept-Ranges" ":" acceptable-ranges
	acceptable-ranges = 1#range-unit | "none"
接受字节范围请求的源服务器可能会发送：
	Accept-Ranges: bytes
但不要求这样做。客户端可以在没有收到所涉及资源的此标头的情况下生成字节范围请求。范围单位在第 3.12 节中定义。

不接受任何类型的资源范围请求的服务器可以发送：
	Accept-Ranges: none
以建议客户端不要尝试范围请求。
```

### 14.16 Content-Range
Content-Range实体标头与部分实体主体一起发送，以指定应在完整实体主体中的何处应用部分主体。

Range units被定义在了3.12章节。
```shell
	Content-Range 			= "Content-Range" ":" content-range-spec        
	content-range-spec      = byte-content-range-spec
	byte-content-range-spec = bytes-unit SP                                  
							  byte-range-resp-spec "/"                                
							  ( instance-length | "*" )        
	byte-range-resp-spec 	= (first-byte-pos "-" last-byte-pos)
	instance-length         = 1*DIGIT
```
标头应该指示完整实体主体的总长度，除非该长度未知或难以确定。星号"*"字符表示生成响应时实例长度未知。

与byte-ranges-specifier值(参见第14.35.1节)不同，byte-range-resp-spec必须仅指定一个范围，并且必须包含该范围的第一个和最后一个字节的绝对字节位置。
具有 `byte-range-resp-spec` 的 `byte-content-range-spec`，其 `last-byte-pos` 值小于其`first-byte-pos` 值，或其 `instance-length` 值小于或等于其 `last-byte-pos` 值，这都是无效的。无效 ``byte-content-range-spec` 的接收者必须忽略它以及随之传输的任何内容。

发送状态码为 416（请求的范围不可满足）的服务器应该包含一个 `byte-range-resp-spec` 为 `*` 的 Content-Range 标头。实例长度指定为当前所选资源的长度。


### 14.35 Range
#### 14.35.1、Byte Ranges（字节范围）

rfc2616 中 14.35 部分定义了 Range 头部。我们来看相关说明：

```shell
由于所有 HTTP 实体在 HTTP 消息中都表示为字节序列，因此字节范围的概念是对于任何 HTTP 实体都有意义。（但是，并非所有客户端和服务器都需要支持字节范围操作。）HTTP 中的字节范围规范适用于实体主体中的字节序列（不一定与实体主体中的字节序列相同）。字节范围操作可以指定单个字节范围或单个实体内的一组范围。
	range-specifier 	  = byte-ranges-specifier
	byte-ranges-specifier = bytes-unit "=" byte-range-set
	byte-range-set 		  = 1#( byte-range-spec | suffix-byte-range-spec )
	byte-range-spec 	  = first-byte-pos "-" [last-byte-pos]
	first-byte-pos 	 	  = 1*DIGIT
	last-byte-pos  		  = 1*DIGIT
byte-ranges-specifier中的first-byte-pos值给出了范围中第一个字节的字节偏移量。last-byte-pos值给出范围内最后一个字节的字节偏移量；也就是说，指定的字节位置是包含在内的。字节偏移开始为零。
如果last-byte-pos值存在，它必须大于或等于该字节范围规范中的first-byte-pos，否则byte-range-spec在语法上无效。包含一个或多个语法上无效的byte-range-set的byte-ranges-spec的接收者必须忽略包含byte-range-set的头字段。

如果last-byte-pos值不存在，或者该值大于或等于实体主体的当前长度，last-byte-pos被认为等于实体主体的当前长度（以字节为单位）减一。

通过选择last-byte-pos，客户端可以在不知道实体大小的情况下限制检索的字节数。
	suffix-byte-range-spec = "-" suffix-length
    suffix-length = 1*DIGIT
suffix-byte-range-spec用于指定实体主体的后缀，其长度由后缀长度值给定。（也就是说，此形式指定实体主体的最后 N 个字节。）如果实体短于指定的后缀长度，
使用整个实体主体。

语法上有效的byte-range-set至少包含一个字节范围规范，其第一个字节位置小于实体主体的当前长度或至少一个具有非零后缀长度的suffix-byte-range-spec。否则，byte-range-set是不可满足的。如果byte-range-set设置不满足，服务器应返回状态为 416（请求的范围无法满足）的响应。否则，服务器应该返回状态为 206（部分内容）的响应，其中包含实体主体的可满足范围。

byte-ranges-specifier值的示例（假设实体长度为10000）：
* 前500字节（字节偏移0-499，包含）：
	bytes=0-499
* 第二个500字节（字节偏移500-999，包含）：
	bytes=500-999
* 后500字节（字节偏移9500-9999，包含）：
	bytes=-500
	或
	bytes=9500-
* 第一个字节和最后一个字节（字节0和9999）：
	bytes=0-0,-1
* 几个二合法但不和规范的大的第二个500字节（bytes offset 500-999，包含）：
	bytes=500-600,601-999
	bytes=500-700,601-999
```

> 注释：
>
> * 1#( byte-range-spec | suffix-byte-range-spec )：表示

### 14.35.2、Range Retrieval Requests（范围检索请求）