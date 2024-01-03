# x、Request结构
```go
type Request struct {
	Method string
	URL *url.URL
	Proto      string 
	ProtoMajor int
	ProtoMinor int
	Header Header
	Body io.ReadCloser
	GetBody func() (io.ReadCloser, error)
	ContentLength int64
	TransferEncoding []string
	Close bool
	Host string
	Form url.Values
	PostForm url.Values
	MultipartForm *multipart.Form
	Trailer Header
	RemoteAddr string
	RequestURI string
	TLS *tls.ConnectionState
	Cancel <-chan struct{}
	Response *Response
	ctx context.Context
}
```
## x.1、Request.Method
```go
// Method specifies the HTTP method (GET, POST, PUT, etc.).
// For client requests, an empty string means GET.
//
// Go's HTTP client does not support sending a request with
// the CONNECT method. See the documentation on Transport for
// details.
Method string
```
## x.2、Request.URL
```shell

```
## x.3、Request.Proto
```go

```
## x.4、Request.ProtoMajor
```go

```
## x.5、Request.ProtoMinor
```go

```
## x.6、Request.Header
```go
// Header contains the request header fields either received
// by the server or to be sent by the client.
//
// If a server received a request with header lines,
//
//	Host: example.com
//	accept-encoding: gzip, deflate
//	Accept-Language: en-us
//	fOO: Bar
//	foo: two
//
// then
//
//	Header = map[string][]string{
//		"Accept-Encoding": {"gzip, deflate"},
//		"Accept-Language": {"en-us"},
//		"Foo": {"Bar", "two"},
//	}
//
// For incoming requests, the Host header is promoted to the
// Request.Host field and removed from the Header map.
//
// HTTP defines that header names are case-insensitive. The
// request parser implements this by using CanonicalHeaderKey,
// making the first character and any characters following a
// hyphen uppercase and the rest lowercase.
//
// For client requests, certain headers such as Content-Length
// and Connection are automatically written when needed and
// values in Header may be ignored. See the documentation
// for the Request.Write method.
Header Header
```
## x.7、Request.Body
```go

```
## x.8、Request.GetBody
```go

```
## x.9、Request.ContentLength
```go

```
## x.10、Request.TransferEncoding
```go

```
## x.11、Request.Close
```go

```
## x.12、Request.Host
```go
// For server requests, Host specifies the host on which the
// URL is sought. For HTTP/1 (per RFC 7230, section 5.4), this
// is either the value of the "Host" header or the host name
// given in the URL itself. For HTTP/2, it is the value of the
// ":authority" pseudo-header field.
// It may be of the form "host:port". For international domain
// names, Host may be in Punycode or Unicode form. Use
// golang.org/x/net/idna to convert it to either format if
// needed.
// To prevent DNS rebinding attacks, server Handlers should
// validate that the Host header has a value for which the
// Handler considers itself authoritative. The included
// ServeMux supports patterns registered to particular host
// names and thus protects its registered Handlers.
//
// For client requests, Host optionally overrides the Host
// header to send. If empty, the Request.Write method uses
// the value of URL.Host. Host may contain an international
// domain name.
Host string
```
## x.13、Request.Form
```go

```
## x.14、Request.PostForm
```go

```
## x.15、Request.MultipartForm
```go
// MultipartForm is the parsed multipart form, including file uploads.
// This field is only available after ParseMultipartForm is called.
// The HTTP client ignores MultipartForm and uses Body instead.
MultipartForm *multipart.Form
```
如上描述是对 Request 结构的 MultipartForm 成员进行了说明。意为：`MultipartForm` 是解析后的 multipart 表单，包括文件上传。该字段仅在调用 `ParseMultipartForm` 后才可用。HTTP 客户端会忽略 `MultipartForm` 并使用 `Body` 代替。

## x.16、Request.Trailer
```go

```
## x.17、Request.RemoteAddr
```go

```
## x.18、Request.RequestURI
```go

```
## x.19、Request.TLS
```go

```
## x.20、Request.Cancel
```go

```
## x.21、Request.Response
```go

```
## x.22、Request.ctx
```go

```

# 参考
[^1]: https://www.jianshu.com/p/29e38bcc8a1d
[^2]: https://www.rfc-editor.org/rfc/rfc2388.txt