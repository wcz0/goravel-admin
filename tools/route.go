package tools

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
)

// RequestIs 判断请求是否是指定的路径
func RequestIs(ctx http.Context, path string) bool {
	return MatchPath(path, ctx.Request().Path())
}

// MatchPath 判断 requestPath 是否命中 pattern。
//
// pattern 支持尾部通配符（如 "/admin-api/system/*"）。两侧都会保留前导斜杠，
// 因为 ctx.Request().Path() 返回的是 "/admin-api/..." 这种完整路径；
// 之前调用方把 pattern 的前导斜杠 trim 掉了，导致通配符永不命中。
func MatchPath(pattern, requestPath string) bool {
	req := strings.TrimSuffix(requestPath, "/")
	if req == "" {
		req = "/"
	}
	pat := strings.TrimSuffix(pattern, "/")
	if pat == "" {
		pat = "/"
	}
	if req == pat {
		return true
	}
	// 支持模式通配符，如 "/system/*"
	if strings.HasSuffix(pat, "*") {
		base := strings.TrimSuffix(pat, "*")
		base = strings.TrimSuffix(base, "/")
		return base != "" && strings.HasPrefix(req, base)
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