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
    if r.ActionOfQuickEdit(ctx) {
        return r.ControllerImpl.Service.QuickEdit(ctx)
    }
    if r.ActionOfQuickEditItem(ctx) {
        return r.ControllerImpl.Service.QuickEditItem(ctx)
    }
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
        r.CreateButton(ctx, r.form(ctx), true, "md", "", "drawer"),
    }
	HeaderToolbar = append(HeaderToolbar, r.BaseHeaderToolBar()...)
    notAdmin := "true"
	currentUser, ok := ctx.Value("admin_user").(*admin.AdminUser)
    if !ok || currentUser == nil {
        notAdmin = "true"
    } else if currentUser.IsAdministrator() {
        notAdmin = "false"
    }

	crud := r.BaseCRUD(ctx).HeaderToolbar(HeaderToolbar).Filter(
		r.BaseFilter().Body([]any{
			gamis.TextControl().Name("keyword").Label(tools.AdminLang(ctx, "keyword")).Size("md").Placeholder(tools.AdminLang(ctx, "admin_user.search_username")),
		}),
    ).Columns([]any{
        gamis.TableColumn().Name("id").Label("ID").Sortable(true),
        gamis.TableColumn().Name("avatar").Label(tools.AdminLang(ctx, "admin_user.avatar")).Set("type", "avatar").Set("src", "${avatar}"),
        gamis.TableColumn().Name("username").Label(tools.AdminLang(ctx, "username")),
        gamis.TableColumn().Name("name").Label(tools.AdminLang(ctx, "admin_user.name")),
        gamis.TableColumn().Name("roles").Label(tools.AdminLang(ctx, "admin_user.roles")).Set("type", "each").Set("items", map[string]any{
            "type":      "tag",
            "label":     "${name}",
            "className": "my-1",
        }),
        gamis.TableColumn().Name("enabled").Label(tools.AdminLang(ctx, "extensions.card.status")).QuickEdit(map[string]any{
            "type": "switch", "mode": "inline", "disabledOn": "${id == 1}", "saveImmediately": true,
            "trueValue": 1, "falseValue": 0,
        }),
        gamis.TableColumn().Name("createdAt").Label(tools.AdminLang(ctx, "created_at")).Set("type", "datetime").Set("format", "YYYY-MM-DD HH:mm:ss").Sortable(true),
        r.RowActions(ctx, r.form(ctx), []any{
            r.RowEditButton(ctx, r.form(ctx), true, "md", "", "drawer").HiddenOn("${administrator && " + notAdmin + "}"),
            r.RowDeleteButton(ctx, "").HiddenOn("${id == 1}"),
        }, "md"),
    })

	return r.BaseList(crud)
}

func (r *UserController) Show(ctx http.Context) http.Response {
	// 使用统一的验证错误处理方法
    if hasError, response := r.HandleValidationErrors(ctx, map[string]string{
        "id": "required|number",
    }, nil); hasError {
        return response
    }

	return nil
}

func (r *UserController) Store(ctx http.Context) http.Response {
    if r.ActionOfQuickEdit(ctx) {
        return r.ControllerImpl.Service.QuickEdit(ctx)
    }
    if r.ActionOfQuickEditItem(ctx) {
        return r.ControllerImpl.Service.QuickEditItem(ctx)
    }

	// 使用统一的验证错误处理方法
    if hasError, response := r.HandleValidationErrors(ctx, map[string]string{
        "username": "required|string|min_len:3|max_len:32",
        "password": "string|min_len:6|max_len:32",
        "confirmPassword": "string|min_len:6|max_len:32|eq_field:password",
        "name":     "required|string|max_len:50",
        "avatar":   "string|max_len:255",
        "enabled":  "int|in:0,1",
        "roles":    "array",
    }, map[string]string{
        "username.required": "用户名不能为空",
        "username.min_len": "用户名至少 3 位",
        "username.max_len": "用户名最多 32 位",
        "password.min_len": "密码至少 6 位",
        "password.max_len": "密码最多 32 位",
        "confirmPassword.eq_field": "两次输入的密码不一致",
        "name.required": "姓名不能为空",
        "name.max_len": "姓名最多 50 位",
        "avatar.max_len": "头像地址长度不能超过 255",
        "enabled.in": "状态值必须为 0 或 1",
    }); hasError {
        return response
    }

	return r.ControllerImpl.Service.Store(ctx)
}

func (r *UserController) Update(ctx http.Context) http.Response {
    if r.ActionOfQuickEdit(ctx) {
        return r.ControllerImpl.Service.QuickEdit(ctx)
    }
    if r.ActionOfQuickEditItem(ctx) {
        return r.ControllerImpl.Service.QuickEditItem(ctx)
    }
    // 使用统一的验证错误处理方法
    if hasError, response := r.HandleValidationErrors(ctx, map[string]string{
        "id":       "required|number",
        "username": "required|string|min_len:3|max_len:32",
        "name":     "required|string|max_len:50",
        "avatar":   "string|max_len:255",
        "enabled":  "int|in:0,1",
        "roles":    "array",
        "password": "string|min_len:6|max_len:32",
    }, map[string]string{
        "id.required": "ID 不能为空",
        "username.required": "用户名不能为空",
        "username.min_len": "用户名至少 3 位",
        "username.max_len": "用户名最多 32 位",
        "name.required": "姓名不能为空",
        "name.max_len": "姓名最多 50 位",
        "avatar.max_len": "头像地址长度不能超过 255",
        "enabled.in": "状态值必须为 0 或 1",
        "password.min_len": "密码至少 6 位",
        "password.max_len": "密码最多 32 位",
    }); hasError {
        return response
    }

    return r.ControllerImpl.Service.Update(ctx)
}

func (r *UserController) Destroy(ctx http.Context) http.Response {
    hasError, response := r.HandleValidationErrors(ctx, map[string]string{
        "id": "required|number",
    }, nil)
	if hasError {
		return response
	}
	return r.ControllerImpl.Service.Destroy(ctx)
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
			gamis.TextControl().Type("input-password").Label(tools.AdminLang(ctx, "admin_user.confirm_password")).Name("confirm_password"),
		})
	return a.SuccessData(ctx, gamis.Page().Body(form))
}

func (a *UserController) form(ctx http.Context) *renderers.Form {
	return a.BaseForm(ctx, false).Body([]any{
		gamis.ImageControl().Name("avatar").Label(tools.AdminLang(ctx, "admin_user.avatar")).Receiver(a.UploadImagePath(ctx)),
		gamis.TextControl().Name("username").Label(tools.AdminLang(ctx, "username")).Required(true),
		gamis.TextControl().Name("name").Label(tools.AdminLang(ctx, "admin_user.name")).Required(true),
		gamis.TextControl().Name("password").Label(tools.AdminLang(ctx, "password")).Type("input-password").Set("requiredOn", "${!id}"),
		gamis.TextControl().Name("confirmPassword").Label(tools.AdminLang(ctx, "confirm_password")).Type("input-password").Set("requiredOn", "${!id}"),
		gamis.SelectControl().Name("roles").Label(tools.AdminLang(ctx, "admin_user.roles")).
			Searchable(true).Multiple(true).LabelField("name").
			ValueField("id").
			JoinValues(false).
			Set("extractValue", true).
			DisabledOn("${id == 1}").
			Source(tools.GetAdmin("system/admin_roles/options")),
		gamis.SwitchControl().Name("enabled").Label(tools.AdminLang(ctx, "extensions.card.status")).
			OnText(tools.AdminLang(ctx, "extensions.status_map.enabled")).
			OffText(tools.AdminLang(ctx, "extensions.status_map.disabled")).
			DisabledOn("${id == 1}").
			Set("trueValue", 1).
			Set("falseValue", 0).
			Value(1),
	})
}

func (a *UserController) List(ctx http.Context) *renderers.Page {
	crud := a.BaseCRUD(ctx).HeaderToolbar([]any{}).Filter(a.BaseFilter().Body([]any{})).Columns([]any{})

	return a.BaseList(crud)
}

// func (a *UserController) detail(ctx http.Context) *renderers.Form {
// 	return a.BaseDetail(ctx).Body([]any{})
// }
