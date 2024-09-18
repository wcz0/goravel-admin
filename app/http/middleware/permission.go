package middleware

import (
	"goravel/app/support/core"
	"goravel/app/response"

	"github.com/goravel/framework/contracts/http"
)

func Permission() http.Middleware {
	return func(ctx http.Context) {
		p := core.NewPermission()
		if p.PermissionIntercept(ctx) {
			ctx.Request().AbortWithStatusJson(http.StatusOK, response.PermissionError)
			return
		}
		ctx.Request().Next()
	}
}
