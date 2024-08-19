package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Casbin() http.Middleware {
	return func(ctx http.Context) {
		// route := ctx.Request().Path()
		// user := ctx.Value("user").(models.User)
		// if user.Id == 0 {
		// 	ctx.Request().AbortWithStatusJson(http.StatusOK, response.Unauthorized)
		// 	return
		// }
		// userId := strconv.FormatUint(user.Id, 10)
		// method := ctx.Request().Method()
		// facades.Enforcer().Enforce(userId, route, method)
		config := facades.Config()
		if config.GetBool("admin.auth.enable") {
			ctx.Request().Next()
			return
		}
		ctx.Request().Next()
	}
}
