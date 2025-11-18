package tools

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
)

// RequestIs 判断请求是否是指定的路径
func RequestIs(ctx http.Context, path string) bool {
    req := strings.TrimSuffix(ctx.Request().Path(), "/")
    pat := strings.TrimSuffix(path, "/")
    if pat == "" { pat = "/" }
    if req == pat { return true }
    // 支持模式通配符，如 "/system/*"
    if strings.HasSuffix(pat, "*") {
        base := strings.TrimSuffix(pat, "*")
        base = strings.TrimSuffix(base, "/")
        if strings.HasPrefix(req, base) { return true }
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