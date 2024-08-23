package middleware

import (
	"goravel/app/core"
	"goravel/app/response"

	"github.com/goravel/framework/contracts/http"
)

func Authenticate() http.Middleware {
	return func(ctx http.Context) {
		permission := core.NewPermission()
		if permission.AuthIntercept(ctx) {
			ctx.Request().AbortWithStatusJson(http.StatusOK, response.Unauthorized)
			return
		}
		permission.CheckUserStatus(ctx)
		ctx.Request().Next()
	}
}
