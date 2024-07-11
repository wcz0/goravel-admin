package admin

import (
	"goravel/app/http/controllers"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type IndexController struct {
	*controllers.Controller
}

func NewIndexController() *IndexController {
	return &IndexController{
		Controller: controllers.NewController(),
	}
}

func (i *IndexController) SaveSetting(ctx http.Context) http.Response {

	return i.Success(ctx)

}

func (i *IndexController) GetSetting(ctx http.Context) http.Response {
	config := facades.Config()
	data := map[string]any{
		"nav":  getNav(),
		"assets": getAssets(),
		"app_name": config.Get("app.name"),
		"locale":   config.Get("app.locale", "zh-CN"),
		"layout":   config.Get("admin.layout"),
		"logo":     config.Get("admin.logo"),

		"login_captcha":          config.Get("admin.auth.login_captcha"),
		"show_development_tools": config.Get("admin.show_development_tools"),
		"system_theme_setting":   config.Get("admin.system_theme_setting"),
		"enabled_extensions":     config.Get("admin.enabled_extensions"),
	}
	return i.DataSuccess(ctx, data)
}

func (i *IndexController) NoContext(c http.Context) http.Response {
	return i.Success(c)
}

func (i *IndexController) DownloadExport(c http.Context) http.Response {
	return i.Success(c)
}

func (i *IndexController) ImageUpload(c http.Context) http.Response {
	return i.Success(c)

}

func (i *IndexController) FileUpload(c http.Context) http.Response {
	return i.Success(c)

}

func (i *IndexController) RichFileUpload(c http.Context) http.Response {
	return i.Success(c)

}

func (i *IndexController) GetUserSetting() {

}

func (i *IndexController) PutUserSetting() {

}

func (i *IndexController) GetCurrentUser() {

}

func (i *IndexController) SearchIcon() {

}

func (i *IndexController) GetMenus(c http.Context) http.Response {
	return i.Success(c)

}

func (i *IndexController) GetDashBoard() {

}


func getAssets() any {
	return map[string]any{
		"css": "",
		"js": "",
		"scripts": "",
		"styles": "",
	}
}

func getNav() any {
	return map[string]any{
		"appendNav": "",
		"prependNav": "",
	}
}
