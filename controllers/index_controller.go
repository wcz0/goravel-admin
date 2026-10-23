package controllers

import (
	gohttp "net/http"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/goravel-admin/core"
	"github.com/wcz0/goravel-admin/core/icons"
	"github.com/wcz0/goravel-admin/services"
	"github.com/wcz0/goravel-admin/tools"
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
	localeOptions := config.Get("admin.layout.locale_options", map[string]string{
		"en":    "English",
		"zh_CN": "简体中文",
	})

	appUrl := config.GetString("app.url")
	logo := config.GetString("admin.logo")

	data := map[string]any{
		"nav":      getNav(),
		"assets":   getAssets(),
		"app_name": config.Get("app.name"),
		"locale":   i.Service.Get("admin_locale", "zh_CN", false),
		"layout":   config.Get("admin.layout"),
		"logo":     strings.TrimSuffix(appUrl, "/") + logo,

		"login_captcha":          config.GetBool("admin.auth.login_captcha"),
		"locale_options":         tools.Map2options(localeOptions.(map[string]string)),
		"show_development_tools": config.GetBool("admin.show_development_tools"),
		"system_theme_setting":   i.Service.Get("system_theme_setting", nil, false),
		"enabled_extensions":     services.NewAdminExtensionService().EnabledNames(),
	}
	return i.SuccessMsgData(ctx, "", data)
}

func (i *IndexController) NoContext(ctx http.Context) http.Response {
	return i.SuccessMsg(ctx, "")
}

// PageSchema 返回「页面管理」里登记的自定义页面结构。
//
// 前端在菜单命中 url_type=4 时会带上 component 里的 sign 请求这里，
// 拿到 schema 后直接交给 amis 渲染。
func (i *IndexController) PageSchema(ctx http.Context) http.Response {
	schema := services.NewAdminPageService().Get(ctx.Request().Input("sign"))
	if schema == nil {
		return i.FailMsg(ctx, tools.AdminLang(ctx, "pages.schema_not_found"))
	}

	return i.SuccessData(ctx, schema)
}

// IconifySearch 按关键字搜索图标名，供后台图标选择器使用。
func (i *IndexController) IconifySearch(ctx http.Context) http.Response {
	names := icons.Search(ctx.Request().Input("query"))
	items := make([]map[string]string, 0, len(names))
	for _, name := range names {
		items = append(items, map[string]string{"icon": name})
	}

	return i.SuccessData(ctx, map[string]any{
		"items": items,
		"total": len(items),
	})
}

// DownloadExport 校验签名后返回导出文件。
//
// 下载地址是拼接出来的公开链接，不依赖登录态，所以必须校验 HMAC 签名与
// 有效期，避免被拿来读取任意文件。
func (i *IndexController) DownloadExport(c http.Context) http.Response {
	path := c.Request().Input("path")
	expires := c.Request().Input("expires")
	signature := c.Request().Input("signature")

	if err := tools.VerifyExport(path, expires, signature); err != nil {
		return c.Response().Data(gohttp.StatusForbidden, "text/plain; charset=utf-8", []byte(err.Error()))
	}

	disk := facades.Config().GetString("admin.upload.disk")
	if disk == "" {
		disk = "public"
	}
	storage := facades.Storage().Disk(disk)
	if !storage.Exists(path) {
		return c.Response().Data(gohttp.StatusNotFound, "text/plain; charset=utf-8", []byte("文件不存在"))
	}

	// 本地磁盘优先走 Download，能带上正确的文件名；其他驱动回退为内存下发。
	if localPath := storage.Path(path); localPath != "" {
		return c.Response().Download(localPath, filepath.Base(path))
	}

	content, err := storage.Get(path)
	if err != nil {
		return c.Response().Data(gohttp.StatusNotFound, "text/plain; charset=utf-8", []byte("文件不存在"))
	}

	return c.Response().Data(gohttp.StatusOK,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		[]byte(content))
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
