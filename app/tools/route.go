package tools

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
)

// RequestIs 判断请求是否是指定的路径
func RequestIs(ctx http.Context, path string) bool {
	requestPath := strings.TrimSuffix(ctx.Request().Path(), "/")
	path = strings.TrimSuffix(path, "/")
	if requestPath == path {
		return true
	}
	if strings.HasSuffix(requestPath, "*") {
		requestPath = strings.TrimSuffix(requestPath, "*")
		requestPath = strings.TrimSuffix(requestPath, "/")
		if strings.HasPrefix(path, requestPath) {
			return true
		}
	}
	return false
}

/**
* RequestMethodIs 判断请求是否是指定的方法
*/
func RequestMethodIs(ctx http.Context, method string) bool {
	requestMethod := ctx.Request().Method()
	method = strings.ToUpper(method)
	if method == requestMethod {
		return true
	}
	if strings.HasSuffix(requestMethod, "*") {
		return true
	}
	return false
}