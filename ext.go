package zebra

import (
	"net/http"
	"strings"
)

type mount struct {
	prefix     string
	middleware Middleware
}

// Mount 函数用于将一个 Middleware 挂载到一个路径前缀上。
// 该 Middleware 只会在指定 prefix 的URL子段下执行。
func Mount(prefix string, m Middleware) Middleware {
	if "/" == prefix {
		return m
	}
	return &mount{prefix, m}
}

func (m *mount) Excute(cxt *Context) bool {
	// TODO 这个判断不对
	if has := strings.HasPrefix(cxt.Urlpath, m.prefix); has {
		return m.middleware.Excute(cxt)
	}
	return true
}

type middlewareWrapper struct {
	handler http.Handler
}

func Wrap(handler http.Handler) Middleware {
	return &middlewareWrapper{handler}
}

func (wrapper *middlewareWrapper) Excute(cxt *Context) bool {
	wrapper.handler.ServeHTTP(cxt.OriginalResponseWriter(), cxt.OriginalRequest())
	return true
}
