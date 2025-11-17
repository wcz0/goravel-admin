package admin

import (
	"goravel/app/http/controllers"
	"goravel/app/models/admin"
	"goravel/app/services"
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/gamis"
)

type AuthController struct {
	*controllers.Controller
	*services.AdminUserService
	*services.AdminSettingService
	cacheKey string
}

func NewAuthController() *AuthController {
	return &AuthController{
		Controller:          controllers.NewController(),
		AdminUserService:    services.NewAdminUserService(),
		AdminSettingService: services.NewAdminSettingService(),
		cacheKey:            "app_setting_",
	}
}

func (a *AuthController) Login(ctx http.Context) http.Response {
	// 使用统一的验证错误处理方法
	if hasError, response := a.HandleValidationErrors(ctx, map[string]string{
		"username": "required|min_len:3|max_len:32",
		"password": "required|min_len:5|max_len:32",
	})
	if err != nil {
		return a.Controller.FailMsg(ctx, err.Error())
	}
	if validator.Fails() {
		return a.Controller.FailMsg(ctx, validator.Errors().All())
	}
	return a.AdminUserService.Login(ctx)
}

func (a *AuthController) Logout(ctx http.Context) http.Response {
	facades.Auth(ctx).Guard("admin").Logout()
	return a.Controller.Success(ctx)
}

func (a *AuthController) Register(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"username": "required|string|min_len:max:32",
		"password": "required|string|min_len:5|max_len:32",
	})
	if err != nil {
		return a.Controller.FailMsg(ctx, err.Error())
	}
	if validator.Fails() {
		return a.Controller.FailMsg(ctx, validator.Errors().All())
	}
	password, err := facades.Hash().Make(ctx.Request().Input("password"))
	if err != nil {
		return a.Controller.Fail(ctx)
	}
	return a.Controller.SuccessData(ctx, password)
}

// LoginPage 登录页面 需要写入设置 开启amis页面登录
func (a *AuthController) LoginPage(ctx http.Context) http.Response {
	form := gamis.Form().
		PanelClassName("border-none").
		Id("login-form").
		Title("").
		Api(tools.GetAdmin("/login")).
		InitApi("/no-content").
		Body([]any{
			gamis.TextControl().Name("username").Placeholder(tools.AdminLang(ctx, "username")).Required(true),
			gamis.TextControl().Type("input-password").Name("password").Placeholder(tools.AdminLang(ctx, "password")).Required(true),
			// captcha
			gamis.CheckboxControl().Name("remember_me").Option(tools.AdminLang(ctx, "remember_me")).Value(true),
			gamis.VanillaAction().ActionType("submit").Label(tools.AdminLang(ctx, "login")).Level("primary").ClassName("w-full"),
		}).
		Actions([]any{}).
		OnEvent(map[string]any{
			"inited": map[string]any{
				"actions": []any{
					map[string]any{
						"actionType": "custom",
						"script": `
let loginParams = localStorage.getItem(window.$owl.getCacheKey('loginParams'))
if(loginParams){
	loginParams = JSON.parse(decodeURIComponent(window.atob(loginParams)))
	doAction({
		actionType: 'setValue',
		componentId: 'login-form',
		args: { value: loginParams }
	})
}
`,
					},
				},
			},
			"submitSucc": map[string]any{
				"actions": []any{
					map[string]any{
						"actionType": "custom",
						"script": `
let _data = {}
if(event.data.remember_me){
	_data = { username: event.data.username, password: event.data.password }
}
window.$owl.afterLoginSuccess(_data, event.data.result.data.token)
`,
					},
				},
			},
			"submitFail": map[string]any{
				"actions": []any{
					map[string]any{
						"actionType":  "reload",
						"componentId": "captcha-service",
					},
				},
			},
		})

	card := gamis.Card().ClassName("w-96 m:w-full").Body([]any{
		gamis.Flex().Justify("space-between").ClassName("px-2.5 pb-2.5").Items([]any{
			gamis.Image().Src(facades.Config().GetString("admin.logo")).Width(40).Height(40),
			gamis.Tpl().ClassName("font-medium").Tpl("<div style=\"font-size: 24px\">" + facades.Config().GetString("admin.name") + "</div>"),
		}),
		form,
	})

	return a.Controller.Json(ctx, gamis.Page().ClassName("login-bg").Css(map[string]any{
		".captcha-box .cxd-Image--thumb": map[string]any{
			"padding":                    "0",
			"cursor":                     "pointer",
			"border":                     "var(--Form-input-borderWidth) solid var(--Form-input-borderColor)",
			"border-top-right-radius":    "4px",
			"border-bottom-right-radius": "4px",
		},
		".cxd-Image-thumb": map[string]any{
			"width": "auto",
		},
		".login-bg": map[string]any{
			"background": "var(--owl-body-bg)",
		},
	}).Body(
		gamis.Wrapper().ClassName("h-screen w-full flex items-center justify-center").Body(card),
	))
}

func (a *AuthController) Get(key string, default_ any, fresh bool) any {
	var adminSetting admin.AdminSetting
	if fresh {
		value := facades.Orm().Query().Where("key", key).Select("values").First(&adminSetting)
		if value == nil {
			return default_
		}
	}
	value, err := facades.Cache().RememberForever(a.cacheKey, func() (any, error) {
		err := facades.Orm().Query().Where("key", key).Select("values").First(&adminSetting)
		if err != nil {
			return nil, err
		}
		return adminSetting.Values, nil
	})
	if err != nil {
		return default_
	}
	if value != nil {
		return value
	} else {
		return default_
	}
}

func (a *AuthController) CurrentUser(ctx http.Context) http.Response {
	if !facades.Config().GetBool("admin.auth.enable") {
		return a.Controller.Success(ctx)
	}
	userInfo := ctx.Value("admin_user").(*admin.AdminUser)
	menus := gamis.DropdownButton().HideCaret("true").Trigger("hover").
		Label(userInfo.Name).
		ClassName("h-full w-full").
		BtnClassName("navbar-user w-full").
		MenuClassName("min-w-0").
		Set("icon", userInfo.Avatar).
		Buttons([]any{
			gamis.VanillaAction().
				IconClassName("pr-2").
				Icon("fa fa-user-gear").
				Label(tools.AdminLang(ctx, "user_setting")).
				OnClick("window.location.href = '#/user_setting'"),
			gamis.VanillaAction().
				IconClassName("pr-2").
				Label(tools.AdminLang(ctx, "logout")).
				Icon("fa-solid fa-right-from-bracket").
				OnClick("window.$owl.logout()"),
		})
	return a.Controller.SuccessData(ctx, struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
		Menus  any    `json:"menus"`
	}{
		Name:   userInfo.Name,
		Avatar: userInfo.Avatar,
		Menus:  menus,
	})
}
