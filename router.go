package zebra

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

// 一个路由路径中间件
type route struct {
	path     string             // 路由配置的原始路径
	handlers []func(c *Captain) // 处理器
	method   int                // http method
	pattern  *regexp.Regexp     // 匹配实际path用的正则
	router   *Router            // 所属Router
}

// 路由中间件 Router .
type Router struct {
	routers []*route
	captain *Captain
}

// 创建 Router 中间件
func NewRouter() *Router {
	router := Router{
		routers: make([]*route, 0),
	}
	router.captain = newCaptain(&router)
	return &router
}

func (r *route) excute(url string, method int) {
	captain := r.router.captain
	if (r.method == MethodAny || r.method == method) &&
		(r.path == url || makePathParams(url, r, captain)) {
		for _, h := range r.handlers {
			h(captain)
		}
	}
}

func (r *Router) Excute(c *Context) bool {
	if len(r.routers) > 0 {
		for _, route := range r.routers {
			route.excute(c.Urlpath, c.Method)
		}
		if r.captain != nil && len(r.captain.output) > 0 {
			r.captain.write(c)
		}
	}
	return true
}

func (r *Router) Any(path string, h ...func(c *Captain)) {
	r.addRoute(MethodAny, path, h)
}

func (r *Router) Get(path string, h ...func(c *Captain)) {
	r.addRoute(MethodGet, path, h)
}

func (r *Router) Post(path string, h ...func(c *Captain)) {
	r.addRoute(MethodPost, path, h)
}

func (r *Router) Put(path string, h ...func(c *Captain)) {
	r.addRoute(MethodPut, path, h)
}

func (r *Router) Delete(path string, h ...func(c *Captain)) {
	r.addRoute(MethodDelete, path, h)
}

func (r *Router) Options(path string, h ...func(c *Captain)) {
	r.addRoute(MethodOptions, path, h)
}

func (r *Router) register(route *route) {
	route.router = r
	r.routers = append(r.routers, route)
}

func (r *Router) addRoute(method int, path string, handlers []func(c *Captain)) *route {
	route := newRoute(method, path, handlers)
	r.register(route)
	return route
}

var (
	regAdvanced = regexp.MustCompile(`:[A-Za-z_]+[0-9]*\{.+\}`)
	regCommon   = regexp.MustCompile(`:[A-Za-z_]+[0-9]*`)
)

// 辅助方法，创建route
func newRoute(method int, path string, handlers []func(c *Captain)) *route {
	r := route{
		path:     path,
		method:   method,
		handlers: handlers,
	}
	path = regAdvanced.ReplaceAllStringFunc(path, func(m string) string {
		sp := strings.Index(m, "{")
		regStr := m[sp+1 : len(m)-1]
		_, err := regexp.Compile(regStr)
		if err != nil {
			panic("路由配置中使用了错误的正则: " + regStr)
		}
		return fmt.Sprintf(`(?P<%s>%s)`, m[1:sp], regStr)
	})
	path = regCommon.ReplaceAllStringFunc(path, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})
	path += `\/?`
	r.pattern = regexp.MustCompile(path)
	return &r
}

func makePathParams(reqUrl string, route *route, c *Captain) (ret bool) {
	matches := route.pattern.FindStringSubmatch(reqUrl)
	if len(matches) > 0 && matches[0] == reqUrl {
		for i, name := range route.pattern.SubexpNames() {
			if len(name) > 0 {
				c.paths[name] = matches[i]
				ret = true
			}
		}
	}
	return
}

// Router 上下文对象
type Captain struct {
	router *Router
	paths  map[string]string
	exporter
}

func newCaptain(router *Router) *Captain {
	return &Captain{
		paths:  make(map[string]string),
		router: router,
	}
}

func (c *Captain) Path(pathname string) string {
	return c.paths[pathname]
}

func (c *Captain) write(cxt *Context) {
	res := cxt.OriginalResponseWriter()
	res.Header().Set("Content-Type", c.mimeType)
	res.Header().Set("X-Powered-By", "Zebra")
	if c.status == 0 {
		c.status = http.StatusOK
	}
	res.WriteHeader(c.status)
	res.Write(c.output)
}

// 输出方法
type exporter struct {
	output   []byte
	mimeType string
	status   int
}

func (e *exporter) Status(code int) {
	e.status = code
}

func (e *exporter) WriteJSON(data interface{}) {
	result, err := json.Marshal(data)
	if err != nil {
		panic("JSON解析时遇到无法解析的数据类型: " + reflect.TypeOf(data).Name())
	}
	e.output = result
	e.mimeType = "application/json"
}

func FullFeatures(zebra *Zebra) {
	zebra.Use(ServeFavicon("./favicon.ico"))
}
