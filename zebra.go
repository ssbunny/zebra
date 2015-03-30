package zebra

import (
	"log"
	"net/http"
	"os"
	"strings"
)

type Zebra struct {
	name        string
	middlewares []Middleware
	server      *http.Server
	logger      *log.Logger // 框架自身使用的logger
}

func New() *Zebra {
	return NewWithServer(&http.Server{Addr: ":3000"})
}

func NewWithServer(server *http.Server) *Zebra {
	return &Zebra{
		name:        "zebra",
		middlewares: make([]Middleware, 0),
		server:      server,
		logger:      log.New(os.Stdout, "zebra ", log.LstdFlags|log.Lshortfile),
	}
}

func (z *Zebra) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(z.middlewares) > 0 {
		c := newContext(w, r)
		c.zebra = z
		var index int = 0
		for _, middleware := range z.middlewares {
			index++
			if !middleware.Excute(c) {
				break
			}
		}
		for index--; index >= 0; index-- {
			if v, ok := z.middlewares[index].(MiddlewareCallback); ok {
				v.Callback(c)
			}
		}
	}
}

func (z *Zebra) Use(m interface{}) {
	switch v := m.(type) {
	case http.Handler:
		z.middlewares = append(z.middlewares, Wrap(v))
	case Middleware:
		z.middlewares = append(z.middlewares, v)
	}
}

func (z *Zebra) Name() string {
	return z.name
}

func (z *Zebra) SetName(name string) {
	z.name = name
}

func (z *Zebra) SetLogger(logger *log.Logger) {
	z.logger = logger
}

func (z *Zebra) UseFullFeatures() {
}

func (z *Zebra) Run() {
	z.server.Handler = z
	z.server.ListenAndServe()
}

func (z *Zebra) RunOnAddr(addr string) {
	z.server.Handler = z
	z.server.Addr = addr
	z.server.ListenAndServe()
}

// Context 是 Middleware 的上下文对象
type Context struct {
	Method   int    // Http Method
	Urlpath  string // *http.Request 中的 URL Path
	res      http.ResponseWriter
	req      *http.Request
	zebra    *Zebra
	response *responseAdapter
	Transfer map[string]interface{}
}

// 创建 Context
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	c := Context{
		res: w,
		req: r,
		response: &responseAdapter{
			w:           w,
			wroteHeader: false,
		},
		Transfer: make(map[string]interface{}),
	}
	if r != nil {
		c.Method = HttpMethodCode(r.Method)
		c.Urlpath = r.URL.Path
	}
	return &c
}

// 获取原始的 http.ResponseWriter
func (c *Context) OriginalResponseWriter() http.ResponseWriter {
	return c.res
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.response
}

// 获取原始的 *http.Request
func (c *Context) OriginalRequest() *http.Request {
	return c.req
}

// TODO
func (c *Context) Request() *http.Request {
	return c.req
}

func (c *Context) Zebra() *Zebra {
	return c.zebra
}

func (c *Context) Logger() *log.Logger {
	return c.zebra.logger
}

/*
func (c *Context) Put(key string, val interface{}) {
	c.transfer[key] = val
}

func (c *Context) Get(key string) interface{} {
	return c.transfer[key]
}*/

// zebra中的一切组件均为一个 Middleware.
type Middleware interface {
	// Excute 用来执行当前 Middleware,
	// 返回 true 表示继续执行后续 Middleware,
	// 返回 false 阻止后续 Middleware 的执行。
	Excute(*Context) bool
}

type MiddlewareCallback interface {
	Middleware
	Callback(*Context)
}

// 一组表示Http Method的常量。其中Any表示任意一种Http Method.
const (
	MethodAny int = iota
	MethodGet
	MethodPost
	MethodPut
	MethodDelete
	MethodHead
	MethodOptions
)

var httpMethodCode = map[string]int{
	"GET":     MethodGet,
	"POST":    MethodPost,
	"PUT":     MethodPut,
	"OPTIONS": MethodOptions,
	"DELETE":  MethodDelete,
	"HEAD":    MethodHead,
}

// 通过请求的 HTTP Method 获取对应的 code.
func HttpMethodCode(method string) int {
	return httpMethodCode[strings.ToUpper(method)]
}

// 适配 server 端的 response, 中间件内部使用
type responseAdapter struct {
	status      int
	wroteHeader bool
	w           http.ResponseWriter
}

func (r *responseAdapter) Header() http.Header {
	return r.w.Header()
}

func (r *responseAdapter) Write(data []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.w.Write(data)
}

func (r *responseAdapter) WriteHeader(code int) {
	r.status = code
	if !r.wroteHeader {
		r.w.WriteHeader(code)
		r.wroteHeader = true
	}
}
