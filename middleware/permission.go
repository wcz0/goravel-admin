package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"github.com/wcz0/goravel-admin/core"
	"github.com/wcz0/goravel-admin/response"
)

// permission 校验当前用户是否拥有访问该路由的权限。
type permission struct{}

// Permission 返回后台权限中间件。
func Permission() http.Middleware {
	return &permission{}
}

func (m *permission) Signature() string {
	return "admin.permission"
}

func (m *permission) Handle(ctx http.Context) {
	if core.NewPermission().PermissionIntercept(ctx) {
		_ = ctx.Response().Json(http.StatusOK, response.PermissionError).Abort()
		return
	}
	ctx.Request().Next()
}
