package middleware

import (
	"goravel/app/models"
	"goravel/app/response"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Auth() http.Middleware {
	return func(ctx http.Context) {
		var user models.User
		err := facades.Auth(ctx).User(&user)
		if err != nil {
			response.NewUnauthorizedError().Response(ctx)
		}
		if user.Id == 0 {
			response.NewUnauthorizedError().Response(ctx)
		}
		ctx.WithValue("user", user)
		ctx.Request().Next()
	}
}
