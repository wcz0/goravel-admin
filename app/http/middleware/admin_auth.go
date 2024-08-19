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
		err := facades.Auth(ctx).Guard("admin").User(&user)
		if err != nil {
			ctx.Request().AbortWithStatusJson(http.StatusOK, response.Unauthorized)
			return
		}
		ctx.WithValue("user", user)
		ctx.Request().Next()
	}
}
