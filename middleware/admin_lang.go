package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/goravel-admin/services"
)

// AdminLang 返回后台多语言中间件。
//
// Goravel 1.17 的 http.Middleware 是 func(Context) 类型（1.18 起改为带
// Signature/Handle 的接口），因此这里以闭包形式实现。
func AdminLang() http.Middleware {
	return func(ctx http.Context) {
		config := facades.Config()
		setting := services.NewAdminSettingService()
		locale, ok := setting.Get("admin_locale", config.GetString("app.locale"), false).(string)
		if !ok || locale == "" {
			locale = config.GetString("app.locale")
		}
		facades.App().SetLocale(ctx, locale)
		ctx.Request().Next()
	}
}
