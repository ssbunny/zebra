package zebra

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// zebra中的一切组件均为一个 Middleware.
type Middleware interface {
	// Excute 用来执行当前 Middleware,
	// 返回 true 表示继续执行后续 Middleware,
	// 返回 false 阻止后续 Middleware 的执行。
	Excute(*Context) bool
}

type Callback interface {
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

// Context 是 Middleware 的上下文对象
type Context struct {
	index   int    // Middleware 的注册顺序
	Method  int    // Http Method
	Urlpath string // *http.Request 中的 URL Path
	res     http.ResponseWriter
	req     *http.Request
	zebra   *Zebra
}

// 创建 Context
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	c := Context{
		index: 0,
		res:   w,
		req:   r,
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

// 获取原始的 *http.Request
func (c *Context) OriginalRequest() *http.Request {
	return c.req
}

func (c *Context) Zebra() *Zebra {
	return c.zebra
}

func (c *Context) Logger() *log.Logger {
	return c.zebra.logger
}

type Zebra struct {
	name        string
	middlewares []Middleware
	callbacks   []Callback
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
		logger:      log.New(os.Stdout, "zebra ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (z *Zebra) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(z.middlewares) > 0 || len(z.callbacks) > 0 {
		c := newContext(w, r)
		c.zebra = z
		for i, middleware := range z.middlewares {
			c.index = i
			if !middleware.Excute(c) {
				break
			}
		}
		for _, callback := range z.callbacks {
			callback.Callback(c)
		}
	}
}

func (z *Zebra) Use(m interface{}) {
	if v, ok := m.(http.Handler); ok {
		z.middlewares = append(z.middlewares, Wrap(v))
		return
	}
	if v, ok := m.(Middleware); ok {
		z.middlewares = append(z.middlewares, v)
	}
	if v, ok := m.(Callback); ok {
		z.callbacks = append(z.callbacks, v)
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

func (z *Zebra) Run() {
	z.server.Handler = z
	z.server.ListenAndServe()
}

func (z *Zebra) RunOnAddr(addr string) {
	z.server.Handler = z
	z.server.Addr = addr
	z.server.ListenAndServe()
}
