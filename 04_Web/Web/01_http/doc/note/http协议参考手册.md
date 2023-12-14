## 3 Protocol Parameters（协议参数）
### 3.12 Range Units（范围单位）
HTTP/1.1 允许客户端请求响应中的部分（一定范围）响应实体。
HTTP/1.1 在 Range（第 14.35 节）和 Content-Range（第 14.16 节）标头字段中使用范围单位。一个实体可以根据不同的结构单元分解为多个子范围。

```shell
	range-unit       = bytes-unit | other-range-unit
	bytes-unit       = "bytes"
	other-range-unit = token
```
HTTP/1.1 定义的唯一范围单位是 `bytes`。HTTP/1.1 的实现可能会忽略其他单元指定的范围。HTTP/1.1 has been designed to allow implementations of applications that do not depend on knowledge of ranges.




## 14 Header Field Definitions（头字段定义）
### 14.5 Accept-Ranges

Accept-Ranges 响应头允许服务端指示请求资源的可接受范围：
```shell
	Accept-Ranges = "Accept-Ranges" ":" acceptable-ranges
	acceptable-ranges = 1#range-unit | "none"
```
接受字节范围请求的源服务器可能会发送：

```shell
	Accept-Ranges: bytes
```


但不要求这样做。客户端可以在没有收到所涉及资源的此标头的情况下生成字节范围请求。范围单位在第 3.12 节中定义。

不接受任何类型的资源范围请求的服务器可以发送：

```shell
	Accept-Ranges: none
```

以建议客户端不要尝试范围请求。



### 14.13 Content-Length【待补充】

Content-Length 实体头字段指示实体主体的大小，它是发送给接收者的十进制八位字节数，xxxx

```shell
	Content-Length    = "Content-Length" ":" 1*DIGIT
```

一个示例：

```shell
Content-Length: 3495
```

应用程序应该使用这个字段指示消息体的传输长度，除非第 4.4 节中的规则禁止这样做。任何大于或等于零的内容长度都是有效值。4.4 节描述了如果没有给出 Content-Length，如何确定消息体的长度。




### 14.16 Content-Range
Content-Range 实体标头与部分实体主体一起发送，以指定应在完整实体主体中的何处应用部分主体。

Range units 被定义在了 3.12 章节。
```shell
	Content-Range 			= "Content-Range" ":" content-range-spec        
	content-range-spec      = byte-content-range-spec
	byte-content-range-spec = bytes-unit SP                                  
							  byte-range-resp-spec "/"                                
							  ( instance-length | "*" )        
	byte-range-resp-spec 	= (first-byte-pos "-" last-byte-pos)
	instance-length         = 1*DIGIT
```
标头应该指示完整实体主体的总长度，除非该长度未知或难以确定。星号 `*` 字符表示生成响应时实例长度未知。

与 `byte-ranges-specifier` 值（参见第 14.35.1 节）不同，`byte-range-resp-spec` 必须仅指定一个范围，并且必须包含该范围的第一个和最后一个字节的绝对字节位置。
具有 `byte-range-resp-spec` 的 `byte-content-range-spec`，其 `last-byte-pos` 值小于其`first-byte-pos` 值，或其 `instance-length` 值小于或等于其 `last-byte-pos` 值，这都是无效的。无效 `byte-content-range-spec` 的接收者必须忽略它以及随之传输的任何内容。

发送状态码为 416（请求的范围不可满足）的服务器应该包含一个 `byte-range-resp-spec` 为 `*` 的 Content-Range 标头。实例长度指定为当前所选资源的长度。状态代码为 206（部分内容）的响应不得包含 `byte-range-resp-spec` 为 `*` 的 Content-Range 字段。

`byte-content-range-spec` 值的示例，假设实体总共包含 1234 个字节：

* 前500字节：
	bytes 0-500/1234
* 第二个500字节：
	bytes 500-999/1234
* 除前500字节外的所有内容：
	bytes 500-1233/1234
* 后500字节：
	bytes 734-1233/1234

当 HTTP 消息包含单个范围的内容时（例如：对单个范围请求的响应，或对一系列没有任何重叠间隙的范围请求的响应），该内容通过 Content-Range 标头和显示实际传输字节数的 Content-Length 标头进行传输。例如，
```shell
       HTTP/1.1 206 Partial content
       Date: Wed, 15 Nov 1995 06:25:24 GMT
       Last-Modified: Wed, 15 Nov 1995 04:58:08 GMT
       Content-Range: bytes 21010-47021/47022
       Content-Length: 26012
       Content-Type: image/gif
```
当 HTTP 消息包含多个范围的内容时（例如，对多个不重叠范围的请求的响应），它们作为多部分消息传输。 用于此目的的多部分媒体类型是附录 19.2 中定义的 `multipart/byteranges`。 有关兼容性问题，请参阅附录 19.6.3。

对单个范围请求的响应不得使用 `multipart/byteranges` 媒体类型发送。对多个范围请求的响应（其结果是单个范围）可以作为包含一个部分的 `multipart/byteranges` 媒体类型发送。无法解码`multipart/byteranges` 消息的客户端不得在单个请求中请求多个字节范围。

当客户端在一个请求中请求多个字节范围时，服务器应该按照它们在请求中出现的顺序返回它们。

如果服务器因为语法上无效而忽略 `byte-range-spec`，则服务器应该将请求视为不存在无效的 Range 标头字段。 （通常，这意味着返回包含完整实体的 200 响应）。

如果服务器接收到一个请求（除了包含 If-Range 请求头字段的请求之外），该请求带有不可满足的 Range 请求头字段（即所有其 `byte-range-spec` 值的 `first-byte-pos` 大于比所选资源的当前长度），它应该返回响应代码 416（请求的范围不可满足）（第 10.4.17 节）。

注意：对于不可满足的 Range 请求标头，客户端不能依赖服务器发送 416（请求的范围不可满足）响应而不是 200（OK）响应，因为并非所有服务器都实现此请求标头。

### 14.35 Range
#### 14.35.1、Byte Ranges（字节范围）
由于所有 HTTP 实体在 HTTP 消息中都表示为字节序列，因此字节范围的概念是对于任何 HTTP 实体都有意义（但是，并非所有客户端和服务器都需要支持字节范围操作）。HTTP 中的字节范围规范适用于实体主体中的字节序列（不一定与实体主体中的字节序列相同）。字节范围操作可以指定单个字节范围或单个实体内的一组范围。
```shell
	range-specifier 	  = byte-ranges-specifier
	byte-ranges-specifier = bytes-unit "=" byte-range-set
	byte-range-set 		  = 1#( byte-range-spec | suffix-byte-range-spec )
	byte-range-spec 	  = first-byte-pos "-" [last-byte-pos]
	first-byte-pos 	 	  = 1*DIGIT
	last-byte-pos  		  = 1*DIGIT
```
`byte-ranges-specifier` 中的 `first-byte-pos` 值给出了范围中第一个字节的字节偏移量。`last-byte-pos` 值给出范围内最后一个字节的字节偏移量；也就是说，指定的字节位置是包含在内的。字节偏移开始为零。
如果 `last-byte-pos` 值存在，它必须大于或等于该字节范围规范中的 `first-byte-pos`，否则 `byte-range-spec` 在语法上无效。包含一个或多个语法上无效的 `byte-range-set` 的 `byte-ranges-spec` 的接收者必须忽略包含 `byte-range-set` 的头字段。

如果 `last-byte-pos` 值不存在，或者该值大于或等于实体主体的当前长度，`last-byte-pos `被认为等于实体主体的当前长度（以字节为单位）减一。

通过选择 `last-byte-pos`，客户端可以在不知道实体大小的情况下限制检索的字节数。

```shell
	suffix-byte-range-spec = "-" suffix-length
    suffix-length = 1*DIGIT
```

`suffix-byte-range-spec` 用于指定实体主体的后缀，其长度由后缀长度值给定（也就是说，此形式指定实体主体的最后 N 个字节）。如果实体短于指定的后缀长度，使用整个实体主体。

语法上有效的 `byte-range-set` 至少包含一个字节范围规范，其第一个字节位置小于实体主体的当前长度或至少一个具有非零后缀长度的 `suffix-byte-range-spec`。否则，`byte-range-set` 是不可满足的。如果 `byte-range-set` 设置不满足，服务器应返回状态为 416（请求的范围无法满足）的响应。否则，服务器应该返回状态为 206（部分内容）的响应，其中包含实体主体的可满足范围。

`byte-ranges-specifier` 值的示例（假设实体长度为10000）：

* 前 500 字节（字节偏移 0-499，包含）：
	bytes=0-499
* 第二个 500 字节（字节偏移 500-999，包含）：
	bytes=500-999
* 后 500 字节（字节偏移 9500-9999，包含）：
	bytes=-500
	或
	bytes=9500-
* 第一个字节和最后一个字节（字节 0 和 9999）：
	bytes=0-0,-1
* 几个二合法但不和规范的大的第二个 500 字节（字节偏移 500-999，包含）：
	bytes=500-600,601-999
	bytes=500-700,601-999



#### 14.35.2、Range Retrieval Requests（范围检索请求）



# 附录

## 附录A：术语表

entity：实体；指 http 请求或响应。

entity-header：实体头；就是 http 请求或响应的 header。

entity-body：实体体；就是 http 请求或响应的 body。