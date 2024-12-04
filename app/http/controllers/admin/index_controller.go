package admin

import (
	"goravel/app/services"
	"goravel/app/support/core"
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type IndexController struct {
	*ControllerImpl[*services.AdminSettingService]
}

func NewIndexController() *IndexController {
	return &IndexController{
		ControllerImpl: NewAdminController[*services.AdminSettingService](services.NewAdminSettingService()),
	}
}

func (i *IndexController) SaveSetting(ctx http.Context) http.Response {
	if err := i.Service.SetMany(ctx.Request().All()); err != nil {
		return i.FailMsg(ctx, err.Error())
	}
	return i.Success(ctx)
}

func (i *IndexController) GetSetting(ctx http.Context) http.Response {
	config := facades.Config()
	localOptions := config.Get("admin.layout.local_options", map[string]string{
		"en":    "English",
		"zh_CN": "简体中文",
	})

	data := map[string]any{
		"nav":      getNav(),
		"assets":   getAssets(),
		"app_name": config.Get("app.name"),
		"locale":   i.Service.Get("admin_locale", "zh_CN", false),
		"layout":   config.Get("admin.layout"),
		"logo":     config.Env("APP_URL").(string) + config.Get("admin.logo").(string),

		"login_captcha":          config.GetBool("admin.auth.login_captcha"),
		"locale_options":         tools.Map2options(localOptions.(map[string]string)),
		"show_development_tools": config.Get("admin.show_development_tools"),
		"system_theme_setting":   i.Service.Get("system_theme_setting", nil, false),
		"enabled_extensions":     []string{},
	}
	return i.SuccessMsgData(ctx, "", data)
}

func (i *IndexController) NoContext(ctx http.Context) http.Response {
	return i.SuccessMsg(ctx, "")
}

func (i *IndexController) DownloadExport(c http.Context) http.Response {
	return i.Success(c)
}

func (i *IndexController) ImageUpload(c http.Context) http.Response {
	return i.UploadImage(c)
}

func (i *IndexController) FileUpload(c http.Context) http.Response {
	return i.UploadFile(c)
}

func (i *IndexController) RichFileUpload(c http.Context) http.Response {
	return i.UploadRich(c)

}

func (i *IndexController) GetUserSetting() {

}

func (i *IndexController) PutUserSetting() {

}

func (i *IndexController) GetCurrentUser() {

}

func (i *IndexController) SearchIcon() {

}

// GetMenus 获取菜单
func (i *IndexController) GetMenus(ctx http.Context) http.Response {
	instance, err := facades.App().Make("admin.menu")
	if err != nil {
		return i.FailMsg(ctx, "获取菜单失败")
	}
	menuInstance, ok := instance.(*core.Menu)
	if !ok {
		return i.FailMsg(ctx, "获取菜单失败2")
	}
	return i.SuccessData(ctx, menuInstance.All(ctx))
}

func (i *IndexController) GetDashBoard() {

}

func getAssets() any {
	return map[string]any{
		"css":     []string{},
		"js":      []string{},
		"scripts": []string{},
		"styles":  []string{},
	}
}

func getNav() any {
	return map[string]any{
		"appendNav":  nil,
		// "prependNav": nil,
		"prependNav": map[string]any{
			"type": "flex",
			"items": []any{
				map[string]any{
					"type": "alert",
					"className": "m-0 mr-5",
					"body": "Version: <b>v0.0.2</b>",
				},
				map[string]any{
					"type": "button",
					"actionType": "url",
					"icon": "fa-brands fa-github",
					"blank": true,
					"url": "https://github.com/wcz0/goravel-admin",
					"tooltip": "Goravel-Admin",
					"className": "mr-2 rounded-full",
				},
			},
		},
	}
}
