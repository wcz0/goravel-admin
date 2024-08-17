package middleware

import (
	"goravel/app/models"
	"goravel/app/response"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func AdminAuth() http.Middleware {
	return func(ctx http.Context) {
		var user models.AdminUser
		err := facades.Auth(ctx).User(&user)
		if err != nil {
			ctx.Response().Success().Json(response.Unauthorized)
		}
		if user.ID == 0 {
			ctx.Response().Success().Json(response.Unauthorized)
		}
		ctx.WithValue("user", user)
		ctx.Request().Next()
	}
}
