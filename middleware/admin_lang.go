package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/goravel-admin/services"
)

// adminLang 根据后台设置切换当前请求的语言环境。
type adminLang struct{}

// AdminLang 返回后台多语言中间件。
func AdminLang() http.Middleware {
	return &adminLang{}
}

func (m *adminLang) Signature() string {
	return "admin.lang"
}

func (m *adminLang) Handle(ctx http.Context) {
	config := facades.Config()
	setting := services.NewAdminSettingService()
	locale, ok := setting.Get("admin_locale", config.GetString("app.locale"), false).(string)
	if !ok || locale == "" {
		locale = config.GetString("app.locale")
	}
	facades.App().SetLocale(ctx, locale)
	ctx.Request().Next()
}
