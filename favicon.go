package zebra

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var maxAge int = 86400

func ServeFavicon(path string) Middleware {
	return &favicon{path}
}

type favicon struct {
	path string
}

func (f *favicon) Excute(cxt *Context) bool {

	if "/favicon.ico" != cxt.Urlpath {
		return true
	}

	if len(f.path) <= 0 {
		f.path = "favicon.ico"
	}

	w := cxt.OriginalResponseWriter()
	if cxt.Method != MethodGet && cxt.Method != MethodHead {
		w.Header().Set("Allow", "GET, HEAD, OPTIONS")
		if cxt.Method == MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return false
	}

	if data, error := ioutil.ReadFile(f.path); error == nil {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(data)
	} else {
		logger := cxt.Logger()
		// TODO 好像没必要这样
		logger.Panicf("找不到指定的favicon.ico:%s\n", f.path)
	}
	return false

}
