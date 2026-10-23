package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"github.com/wcz0/goravel-admin/core"
)

// authenticate 校验后台登录态与用户状态。
type authenticate struct{}

// Authenticate 返回后台登录中间件。
func Authenticate() http.Middleware {
	return &authenticate{}
}

func (m *authenticate) Signature() string {
	return "admin.auth"
}

func (m *authenticate) Handle(ctx http.Context) {
	permission := core.NewPermission()
	permission.AuthIntercept(ctx)
	permission.CheckUserStatus(ctx)
	ctx.Request().Next()
}
