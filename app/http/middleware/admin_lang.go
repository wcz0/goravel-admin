package middleware

import (
	"goravel/app/services"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func AdminLang() http.Middleware {
	return func(ctx http.Context) {
		config := facades.Config()
		services := services.NewAdminSettingService()
		facades.App().SetLocale(ctx, services.Get("admin_locale", config.GetString("app.locale"), false).(string))
		ctx.Request().Next()
	}
}
