package middleware

import (
	"goravel/app/models"
	"goravel/app/response"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/wcz0/goravel-authz/facades"
)

func Casbin() http.Middleware {
	return func(ctx http.Context) {
		route := ctx.Request().Path()
		user := ctx.Value("user").(models.User)
		if user.Id == 0 {
			response.NewUnauthorizedError().Response(ctx)
		}
		userId := strconv.FormatUint(user.Id, 10)
		method := ctx.Request().Method()
		facades.Enforcer().Enforce(userId, route, method)
		ctx.Request().Next()
	}
}
