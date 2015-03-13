package zebra

import (
	"log"
	"net/http"
	"os"
)

type accessLogger log.Logger

var (
	req *http.Request
)

func (l *accessLogger) Excute(cxt *Context) bool {

	req = cxt.OriginalRequest()

	addr := req.Header.Get("X-Real-IP")
	if len(addr) <= 0 {
		addr = req.Header.Get("X-Forwarded-For")
		if len(addr) <= 0 {
			addr = req.RemoteAddr
		}
	}

	logger := (*log.Logger)(l)
	logger.Printf("请求 [%s] %s %s [IP %s]", cxt.Zebra().Name(), req.Method, cxt.Urlpath, addr)
	return true
}

func (l *accessLogger) Callback(cxt *Context) {
	w := cxt.OriginalResponseWriter()
	typ := w.Header().Get("Content-Type")

	logger := (*log.Logger)(l)
	logger.Printf("响应 [%s] %s \n", cxt.Zebra().Name(), typ)
}

func NewLogger() *accessLogger {
	// TODO 默认写到文件里
	return (*accessLogger)(log.New(os.Stdout, "zebra-", log.LstdFlags|log.Lmicroseconds))
}

func NewLoggerWith(l *log.Logger) *accessLogger {
	return (*accessLogger)(l)
}
