package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"github.com/wcz0/goravel-admin/core"
	"github.com/wcz0/goravel-admin/response"
)

// Permission 返回后台权限中间件。
//
// Goravel 1.17 的 http.Middleware 是 func(Context) 类型（1.18 起改为带
// Signature/Handle 的接口），因此这里以闭包形式实现。
func Permission() http.Middleware {
	return func(ctx http.Context) {
		if core.NewPermission().PermissionIntercept(ctx) {
			_ = ctx.Response().Json(http.StatusOK, response.PermissionError).Abort()
			return
		}
		ctx.Request().Next()
	}
}
