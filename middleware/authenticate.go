package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"github.com/wcz0/goravel-admin/core"
)

// Authenticate 返回后台登录中间件。
//
// Goravel 1.17 的 http.Middleware 是 func(Context) 类型（1.18 起改为带
// Signature/Handle 的接口），因此这里以闭包形式实现。
func Authenticate() http.Middleware {
	return func(ctx http.Context) {
		permission := core.NewPermission()
		permission.AuthIntercept(ctx)
		permission.CheckUserStatus(ctx)
		ctx.Request().Next()
	}
}
