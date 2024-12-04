package admin

import (
	"goravel/app/models/admin"
	"goravel/app/services"
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"
	"github.com/wcz0/gamis"
	"github.com/wcz0/gamis/renderers"
)

type UserController struct {
	*ControllerImpl[*services.AdminUserService]
}

func NewUserController() *UserController {
	return &UserController{
		ControllerImpl: NewAdminController[*services.AdminUserService](services.NewAdminUserService()),
	}
}

func (u *UserController) PutUserSetting(c http.Context) http.Response {
	return u.Success(c)
}

func (r *UserController) Index(ctx http.Context) http.Response {
	if r.ActionOfGetData(ctx) {
		return r.ControllerImpl.Service.List(ctx)
	}
	if r.ActionOfExport(ctx) {
		return r.ControllerImpl.Service.Export(ctx)
	}
	return r.SuccessData(ctx, r.list(ctx))
}

func (r *UserController) list(ctx http.Context) *renderers.Page {
	HeaderToolbar := []any{
		r.CreateButton(ctx, r.form(ctx), true, "md", "", ""),
	}
	HeaderToolbar = append(HeaderToolbar, r.BaseHeaderToolBar()...)
	isAdmin := "false"
	currentUser, ok := ctx.Value("admin_user").(*admin.AdminUser)
	if !ok || currentUser == nil {
		isAdmin = "false"
	}

	if currentUser.IsAdministrator() {
		isAdmin = "true"
	}

	crud := r.BaseCRUD(ctx).HeaderToolbar(HeaderToolbar).Filter(
		r.BaseFilter().Body([]any{
			gamis.TextControl().Name("keyword").Label(tools.AdminLang(ctx, "keyword")).Size("md").Placeholder(tools.AdminLang(ctx, "admin_user.search_username")),
		}),
	).Columns([]any{
		gamis.TableColumn().Name("id").Label("ID").Sortable(true),
		gamis.TableColumn().Name("Avatar").Label(tools.AdminLang(ctx, "admin_user.avatar")).Set("type", "avatar").Set("src", "${avatar}"),
		gamis.TableColumn().Name("Username").Label(tools.AdminLang(ctx, "username")),
		gamis.TableColumn().Name("Name").Label(tools.AdminLang(ctx, "admin_user.name")),
		gamis.TableColumn().Name("AdminRoles").Label(tools.AdminLang(ctx, "admin_user.roles")).Set("type", "each").Set("items", map[string]any{
			"type":      "tag",
			"label":     "${name}",
			"className": "my-1",
		}),
		gamis.TableColumn().Name("Enabled").Label(tools.AdminLang(ctx, "extensions.card.status")).QuickEdit(map[string]any{
			"type": "switch", "mode": "inline", "disabledOn": "${id == 1}", "saveImmediately": true}),
		gamis.TableColumn().Name("created_at").Label(tools.AdminLang(ctx, "created_at")).Set("type", "datetime").Sortable(true),
		r.RowActions(ctx, r.form(ctx), []any{
			r.RowEditButton(ctx, r.form(ctx), true, "md", "", "").HiddenOn("${administrator && " + isAdmin + "}"),
			r.RowDeleteButton(ctx, "").HiddenOn("${id == 1}"),
		}, "md"),
	})

	return r.BaseList(crud)
}

func (r *UserController) Show(ctx http.Context) http.Response {
	return nil
}

func (r *UserController) Store(ctx http.Context) http.Response {
	if r.ActionOfQuickEdit(ctx) {
		return r.ControllerImpl.Service.QuickEdit(ctx)
	}
	if r.ActionOfQuickEditItem(ctx){
		return r.ControllerImpl.Service.QuickEditItem(ctx)
	}
	return r.ControllerImpl.Service.Store(ctx)
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

func (a *UserController) form(ctx http.Context) *renderers.Form {
	return a.BaseForm(ctx, false).Body([]any{
		gamis.ImageControl().Name("Avatar").Label(tools.AdminLang(ctx, "admin_user.avatar")).Receiver(a.UploadImagePath(ctx)),
		gamis.TextControl().Name("Username").Label(tools.AdminLang(ctx, "username")).Required(true),
		gamis.TextControl().Name("Name").Label(tools.AdminLang(ctx, "admin_user.name")).Required(true),
		gamis.TextControl().Name("Password").Label(tools.AdminLang(ctx, "password")).Type("input-password"),
		gamis.TextControl().Name("confirmPassword").Label(tools.AdminLang(ctx, "confirm_password")).Type("input-password"),
		gamis.SelectControl().Name("AdminRoles").Label(tools.AdminLang(ctx, "admin_user.roles")).
			Searchable(true).Multiple(true).LabelField("name").
			ValueField("id").
			JoinValues(false).
			ExtractValue("").
			DisabledOn("${id == 1}").
			Options(a.Service.RoleOptions(ctx)),
		gamis.SwitchControl().Name("enabled").Label(tools.AdminLang(ctx, "extensions.card.status")).
			OnText(tools.AdminLang(ctx, "extensions.status_map.enabled")).
			OffText(tools.AdminLang(ctx, "extensions.status_map.disabled")).
			DisabledOn("${id == 1}").
			Value(true),
	})
}

func (a *UserController) List(ctx http.Context) *renderers.Page {
	crud := a.BaseCRUD(ctx).HeaderToolbar([]any{}).Filter(a.BaseFilter().Body([]any{})).Columns([]any{})

	return a.BaseList(crud)
}

// func (a *UserController) detail(ctx http.Context) *renderers.Form {
// 	return a.BaseDetail(ctx).Body([]any{})
// }
