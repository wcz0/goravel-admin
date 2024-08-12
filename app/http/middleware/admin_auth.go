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
			response.NewUnauthorizedError().Response(ctx)
		}
		if user.ID == 0 {
			response.NewUnauthorizedError().Response(ctx)
		}
		ctx.WithValue("adminUser", user)
		ctx.Request().Next()
	}
}
