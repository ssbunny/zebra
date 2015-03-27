package zebra

import (
	"log"
	"net/http"
	"net/textproto"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	combined = `:remote-addr - :remote-user [:date[web]] ":method :url HTTP/:http-version" :status :res[content-length] ":referrer" ":user-agent"`
	common   = `:remote-addr - :remote-user [:date[web]] ":method :url HTTP/:http-version" :status :res[content-length]`
)

var tokenReg = regexp.MustCompile(`:([-\w]{2,})(?:\[([^\]]+)\])?`)

// 服务访问日志
type AccessLogger struct {
	format    string
	output    string
	logger    *log.Logger
	immediate bool
}

// 创建服务访问日志中间件
func NewLogger() *AccessLogger {
	return &AccessLogger{
		format:    combined,
		immediate: false,
		// TODO 默认写到文件里
		logger: log.New(os.Stdout, "", 0),
	}
}

// TODO 不应该阻塞主程序，日志压入队列中异步执行
func (logger *AccessLogger) Excute(cxt *Context) bool {
	cxt.Transfer["reqtime"] = time.Now()
	if logger.immediate {
		logger.log(cxt)
	}
	return true
}

func (logger *AccessLogger) Callback(cxt *Context) {
	if !logger.immediate {
		logger.log(cxt)
	}
}

func (logger *AccessLogger) parse(cxt *Context) string {
	return tokenReg.ReplaceAllStringFunc(logger.format, func(m string) string {
		if strings.HasPrefix(m, ":req") {
			return cxt.Request().Header.Get(plainToken(m))
		}
		if strings.HasPrefix(m, ":res") {
			// TODO 一定程度上讲，在Handler里几乎不可能得到太多有用的 response header
			// 大多数头信息是在 ServeHTTP 之后计算得到，日志组件的生命周期无法渗透其中
			// 如果zebra定义于一个完整的Server而非Handler则显得过于庞杂
			// 也可以考虑在扩展的ResponseWriter中计算各中间件需要的头信息，在Handler
			// 阶段就写入，但是这样对所设计的中间件的执行生命周期有一定的破坏性。
			return cxt.ResponseWriter().Header().Get(plainToken(m))
		}
		return tokenTransmit(m, cxt)
	})
}

func (logger *AccessLogger) log(cxt *Context) {
	out := logger.parse(cxt)
	logger.logger.Println(out)
}

// 设置日志格式，默认为 Combined 格式
func (logger *AccessLogger) SetFormat(format string) {
	if len(format) > 0 {
		if format == "combined" {
			logger.format = combined
		} else if format == "common" {
			logger.format = common
		} else {
			logger.format = format
		}
	}
}

// 开启即时模式：请求时输出日志而不是响应时输出
func (logger *AccessLogger) EnableImmediate() {
	logger.immediate = true
}

// 使用给定的 Logger 记录服务访问日志
func (logger *AccessLogger) SetLogger(l *log.Logger) {
	logger.logger = l
}

func tokenTransmit(token string, cxt *Context) string {
	w := cxt.ResponseWriter()
	req := cxt.Request()

	switch token {
	case ":url":
		return req.URL.Path // TODO 是否应该带上QueryString?
	case ":method":
		return req.Method
	case ":remote-addr":
		return lookupIP(req)
	case ":remote-user":
		return lookupUser(req)
	case ":http-version":
		return strconv.Itoa(req.ProtoMajor) + "." + strconv.Itoa(req.ProtoMinor)
	case ":status":
		return lookupStatus(w)
	case ":user-agent":
		return req.UserAgent()
	case ":response-time":
		return lookupResponseTime(cxt)
	case ":date[unix]":
		return time.Now().Format(time.UnixDate)
	case ":date[web]":
		return time.Now().Format(time.RFC1123)
	case ":date":
		return time.Now().Format("2006-01-02 15:04:05")
	case ":referrer":
		return lookupReferrer(req)
	}
	return ""
}

func lookupReferrer(req *http.Request) string {
	r := req.Referer()
	if len(r) <= 0 {
		r = req.Header.Get("Referrer")
	}
	if len(r) <= 0 {
		r = "-"
	}
	return r
}

func lookupResponseTime(cxt *Context) string {
	if v, ok := cxt.Transfer["reqtime"]; ok {
		d := time.Since(v.(time.Time))
		return strconv.FormatInt((int64)(d/time.Millisecond), 10)
	}
	return "-"
}

func lookupUser(req *http.Request) string {
	if req.URL != nil && req.URL.User != nil {
		return req.URL.User.Username()
	}
	return "-"
}

func lookupStatus(w http.ResponseWriter) string {
	var s int
	if s = w.(*responseAdapter).status; s == 0 {
		return "-"
	}
	return strconv.Itoa(s)
}

func lookupIP(req *http.Request) string {
	addr := req.Header.Get("X-Real-IP")
	if len(addr) <= 0 {
		addr = req.Header.Get("X-Forwarded-For")
		if len(addr) <= 0 {
			addr = req.RemoteAddr
		} else {
			addr = "-"
		}
	}
	return addr
}

func plainToken(m string) string {
	t := m[strings.Index(m, "[")+1 : strings.LastIndex(m, "]")]
	return textproto.CanonicalMIMEHeaderKey(t)
}
