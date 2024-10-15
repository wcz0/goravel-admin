package admin

import (
	"goravel/app/services"
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"
	"github.com/wcz0/gamis"
)

type UserController struct {
	*Controller[*services.AdminUserService]
}

func NewUserController() *UserController {
	return &UserController{
		Controller: NewAdminController[*services.AdminUserService](services.NewAdminUserService()),
	}
}

func (u *UserController) PutUserSetting(c http.Context) http.Response {
	return u.Success(c)
}

func (r *UserController) Index(ctx http.Context) http.Response {
	return nil
}

func (r *UserController) Show(ctx http.Context) http.Response {
	return nil
}

func (r *UserController) Store(ctx http.Context) http.Response {
	return nil
}

func (r *UserController) Update(ctx http.Context) http.Response {
	return nil
}

func (r *UserController) Destroy(ctx http.Context) http.Response {
	return nil
}

func (a *UserController) GetUserSetting(ctx http.Context) http.Response {
	form := gamis.Form().
		Title("").
		PanelClassName("px-48 m:px-0").
		Mode("horizontal").
		InitApi("/current-user").
		Api("put:" + tools.GetAdmin("/current-user")).
		Body([]any{
			gamis.ImageControl().Label(tools.AdminLang(ctx, "admin_user.avatar")).Name("avatar").Receiver(tools.GetAdmin("upload_image")),
			gamis.TextControl().Label(tools.AdminLang(ctx, "admin_user.name")).Name("name").Required(true),
			gamis.TextControl().Type("input-password").Label(tools.AdminLang(ctx, "admin_user.old_password")).Name("old_password"),
			gamis.TextControl().Type("input-password").Label(tools.AdminLang(ctx, "admin_user.password")).Name("password"),
			gamis.TextControl().Type("input-password").Label(tools.AdminLang(ctx, "admin_user.confirm_password")).Name("confirm_password").Required(true),
		})
	return a.SuccessData(ctx, gamis.Page().Body(form))
}
