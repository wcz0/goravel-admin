package middleware

import (
	"goravel/app/support/core"

	"github.com/goravel/framework/contracts/http"
)

func Authenticate() http.Middleware {
	return func(ctx http.Context) {
		permission := core.NewPermission()
		permission.AuthIntercept(ctx)
		permission.CheckUserStatus(ctx)
		ctx.Request().Next()
	}
}
