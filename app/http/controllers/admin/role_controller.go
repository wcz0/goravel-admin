package admin

import (
	"goravel/app/services"
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"
	"github.com/wcz0/gamis"
	"github.com/wcz0/gamis/renderers"
)

type RoleController struct {
	*ControllerImpl[*services.AdminRoleService] // 继承
}

func NewRoleController() *RoleController {
	return &RoleController{
		ControllerImpl: NewAdminController[*services.AdminRoleService](services.NewAdminRoleService()),
	}
}

func (r *RoleController) Index(ctx http.Context) http.Response {
	if r.ActionOfGetData(ctx) {
		return r.ControllerImpl.Service.List(ctx)
	}
	return r.SuccessData(ctx, r.list(ctx))
}

func (r *RoleController) list(ctx http.Context) *renderers.Page {
	HeaderToolbar := []any{
		r.CreateButton(ctx, r.form(ctx), true, "md", "", ""),
	}
	HeaderToolbar = append(HeaderToolbar, r.BaseHeaderToolBar()...)

	crud := r.BaseCRUD(ctx).HeaderToolbar(HeaderToolbar).Filter(
		r.BaseFilter().Body([]any{
			gamis.TextControl().Name("keyword").Label(tools.AdminLang(ctx, "keyword")).Size("md").Placeholder("搜索角色名称或标识"),
		}),
	).Columns([]any{
		gamis.TableColumn().Name("id").Label("ID").Sortable(true),
		gamis.TableColumn().Name("name").Label(tools.AdminLang(ctx, "admin_role.name")).Sortable(true),
		gamis.TableColumn().Name("slug").Label(tools.AdminLang(ctx, "admin_role.slug")).Sortable(true),
		gamis.TableColumn().Name("description").Label(tools.AdminLang(ctx, "admin_role.description")),
		gamis.TableColumn().Name("created_at").Label(tools.AdminLang(ctx, "created_at")).Set("type", "datetime").Sortable(true),
		r.RowActions(ctx, r.form(ctx), []any{
			r.RowEditButton(ctx, r.form(ctx), true, "md", "", ""),
			r.RowDeleteButton(ctx, "").HiddenOn("${slug == 'administrator'}"),
		}, "md"),
	})

	return r.BaseList(crud)
}

func (r *RoleController) Show(ctx http.Context) http.Response {
	// 使用统一的验证错误处理方法
	if hasError, response := r.HandleValidationErrors(ctx, map[string]string{
		"id": "required|number",
	}); hasError {
		return response
	}
	
	return r.ControllerImpl.Service.Show(ctx)
}

func (r *RoleController) Store(ctx http.Context) http.Response {
	// 使用统一的验证错误处理方法
	if hasError, response := r.HandleValidationErrors(ctx, map[string]string{
		"name":        "required|string|max:255",
		"slug":        "required|string|max:255|unique:admin_roles,slug",
		"description": "string|max:500",
		"permissions": "array",
	}); hasError {
		return response
	}

	return r.ControllerImpl.Service.Store(ctx)
}

func (r *RoleController) Update(ctx http.Context) http.Response {
	// 使用统一的验证错误处理方法
	if hasError, response := r.HandleValidationErrors(ctx, map[string]string{
		"name":        "required|string|max:255",
		"slug":        "required|string|max:255",
		"description": "string|max:500",
		"permissions": "array",
	}); hasError {
		return response
	}

	return r.ControllerImpl.Service.Update(ctx)
}

func (r *RoleController) Destroy(ctx http.Context) http.Response {
	// 使用统一的验证错误处理方法
	if hasError, response := r.HandleValidationErrors(ctx, map[string]string{
		"id": "required|integer|min:1",
	}); hasError {
		return response
	}

	return r.ControllerImpl.Service.Destroy(ctx)
}

func (r *RoleController) form(ctx http.Context) *renderers.Form {
	return r.BaseForm(ctx, false).Body([]any{
		gamis.TextControl().Name("name").Label(tools.AdminLang(ctx, "admin_role.name")).Required(true).
			Placeholder(tools.AdminLang(ctx, "admin_role.name")),
		gamis.TextControl().Name("slug").Label(tools.AdminLang(ctx, "admin_role.slug")).Required(true).
			Placeholder(tools.AdminLang(ctx, "admin_role.slug")).
			Description("角色标识，用于权限控制，只能包含字母、数字、下划线和连字符"),
		gamis.TextareaControl().Name("description").Label(tools.AdminLang(ctx, "admin_role.description")).
			Placeholder(tools.AdminLang(ctx, "admin_role.description")).
			Description("可选，描述该角色的用途和职责"),
		gamis.TreeSelectControl().Name("permissions").Label(tools.AdminLang(ctx, "admin_role.permissions")).
			Searchable(true).Multiple(true).LabelField("name").
			ValueField("id").
			JoinValues(false).
			ExtractValue("").
			Options(r.Service.PermissionOptions(ctx)).
			Description("选择该角色拥有的权限"),
	})
}