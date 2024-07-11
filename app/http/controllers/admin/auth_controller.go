package admin

import (
	"goravel/app/http/controllers"
	"goravel/app/services"
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/gamis"
)

type AuthController struct {
	*controllers.Controller
	*services.AdminUserService
}

func NewAuthController() *AuthController {
	return &AuthController{
		Controller:      controllers.NewController(),
		AdminUserService: services.NewAdminUserService(),
	}
}

func (a *AuthController) Login(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"username": "required|max_len:32",
		"password": "required|min_len:5|max_len:32",
	})
	if err != nil {
		return a.MsgError(ctx, err.Error())
	}
	if validator.Fails() {
		return a.MsgError(ctx, validator.Errors().All())
	}
	return a.AdminUserService.Login(ctx)
}

func (a *AuthController) Logout(ctx http.Context) http.Response {
	facades.Auth().Guard("admin").Logout(ctx)
	return a.Success(ctx)
}

func (a *AuthController) Register(ctx http.Context) http.Response {
	return a.Success(ctx)
}

func (a *AuthController) LoginPage(ctx http.Context) http.Response {

	form := gamis.Form().PanelClassName("border-none").Id("login-form").Title("").Api(tools.GetAdmin("login")).InitApi("/no-content").Body([]any{
		gamis.TextControl().Name("username").Placeholder("请输入用户名").Required(true),
		gamis.TextControl().Type("input-password").Name("password").Placeholder("请输入密码").Required(true),
		// captcha
		gamis.VanillaAction().ActionType("submit").Label("登录").ClassName("w-full"),
	}).Actions([]any{})

	failAction := map[string]any{}
	event := map[string]any{
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
	}

	for k, v := range failAction {
		event[k] = v
	}
	form.OnEvent(event)

	card := gamis.Card().ClassName("w-96 m:w-full").Body([]any{
		gamis.Flex().Justify("space-between").ClassName("px-2.5 pb-2.5").Items([]any{
			gamis.Image().Src(facades.Config().GetString("admin.logo")).Width(40).Height(40),
			gamis.Tpl().ClassName("font-medium").Tpl("<div style=\"font-size: 24px\">" + facades.Config().GetString("admin.name") + "</div>"),
		}),
		form,
	})

	return a.DataSuccess(ctx, gamis.Page().ClassName("login-bg").Css(map[string]any{
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

func (a *AuthController) GetUserSetting(c http.Context) http.Response {
	return a.Success(c)
}
