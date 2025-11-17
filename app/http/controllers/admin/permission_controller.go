package admin

import (
	"goravel/app/models/admin"
	"goravel/app/services"
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"
<<<<<<< HEAD
	"github.com/goravel/framework/facades"
=======

	"github.com/wcz0/gamis"
	"github.com/wcz0/gamis/renderers"
>>>>>>> 08c77dc3ed68fd34ac5aa196c797580ff3c72dcb
)

type PermissionController struct {
	*ControllerImpl[*services.AdminPermissionService]
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		ControllerImpl: NewAdminController[*services.AdminPermissionService](services.NewAdminPermissionService()),
	}
}

<<<<<<< HEAD
func (p *PermissionController) Index(ctx http.Context) http.Response {
	if p.ActionOfGetData(ctx) {
		return p.ControllerImpl.Service.List(ctx)
	}
	return p.SuccessData(ctx, p.list(ctx))
}

// 权限列表页面
func (p *PermissionController) list(_ http.Context) interface{} {
	return map[string]interface{}{
		"title": "权限列表",
	}
}

func (p *PermissionController) Show(ctx http.Context) http.Response {
	return p.ControllerImpl.Service.Show(ctx)
}

func (p *PermissionController) Store(ctx http.Context) http.Response {
	// 验证输入
	validation, err := ctx.Request().Validate(map[string]string{
		"parent_id":  "number",
		"name":       "required|string|max:255",
		"value":      "required|string|max:255",
		"http_method": "required|array",
		"http_path":   "required|array",
	})
	if err != nil {
		return p.FailMsg(ctx, err.Error())
	}
	if validation.Fails() {
		return p.FailMsg(ctx, validation.Errors().All())
=======
func (r *PermissionController) Show(ctx http.Context) http.Response {
	// 验证请求参数
	hasError, response := r.HandleValidationErrors(ctx, map[string]string{
		"id": "required|number",
	})
	if hasError {
		return response
	}
	return r.ControllerImpl.Service.Show(ctx)
}

func (r *PermissionController) Store(ctx http.Context) http.Response {
	hasError, response := r.HandleValidationErrors(ctx, map[string]string{
		"parent_id":    "number",
		"name":         "required|string|max:255",
		"slug":         "required|string|max:255",
		"http_method":  "string",
		"http_path":    "string",
		"custom_order": "number",
	})
	if hasError {
		return response
>>>>>>> 08c77dc3ed68fd34ac5aa196c797580ff3c72dcb
	}

	return p.ControllerImpl.Service.Store(ctx)
}

<<<<<<< HEAD
func (p *PermissionController) Update(ctx http.Context) http.Response {
	// 验证输入
	validation, err := ctx.Request().Validate(map[string]string{
		"id":         "required|number",
		"parent_id":  "number",
		"name":       "required|string|max:255",
		"value":      "required|string|max:255",
		"http_method": "required|array",
		"http_path":   "required|array",
	})
	if err != nil {
		return p.FailMsg(ctx, err.Error())
	}
	if validation.Fails() {
		return p.FailMsg(ctx, validation.Errors().All())
=======
func (r *PermissionController) Update(ctx http.Context) http.Response {
	hasError, response := r.HandleValidationErrors(ctx, map[string]string{
		"id":           "required|number",
		"parent_id":    "number",
		"name":         "required|string|max:255",
		"slug":         "required|string|max:255",
		"http_method":  "string",
		"http_path":    "string",
		"custom_order": "number",
	})
	if hasError {
		return response
>>>>>>> 08c77dc3ed68fd34ac5aa196c797580ff3c72dcb
	}

	return p.ControllerImpl.Service.Update(ctx)
}

<<<<<<< HEAD
func (p *PermissionController) Destroy(ctx http.Context) http.Response {
	// 验证输入
	validation, err := ctx.Request().Validate(map[string]string{
		"id": "required|number",
	})
	if err != nil {
		return p.FailMsg(ctx, err.Error())
	}
	if validation.Fails() {
		return p.FailMsg(ctx, validation.Errors().All())
=======
func (r *PermissionController) Destroy(ctx http.Context) http.Response {
	hasError, response := r.HandleValidationErrors(ctx, map[string]string{
		"id": "required|number",
	})
	if hasError {
		return response
>>>>>>> 08c77dc3ed68fd34ac5aa196c797580ff3c72dcb
	}

	return p.ControllerImpl.Service.Destroy(ctx)
}

// GetHttpMethodOptions 获取HTTP方法选项
func (p *PermissionController) GetHttpMethodOptions(ctx http.Context) http.Response {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	return p.SuccessData(ctx, methods)
}

// GetHttpPathOptions 获取HTTP路径选项
func (p *PermissionController) GetHttpPathOptions(ctx http.Context) http.Response {
	paths := []string{"/", "/admin-api/*"}
	return p.SuccessData(ctx, paths)
}

// GetRolesByPermission 根据权限获取角色
func (p *PermissionController) GetRolesByPermission(ctx http.Context) http.Response {
	permissionId := ctx.Request().Input("permission_id", "")
	if permissionId == "" {
		return p.FailMsg(ctx, "权限ID不能为空")
	}
<<<<<<< HEAD

	var permission admin.AdminPermission
	if err := facades.Orm().Query().Find(&permission, permissionId); err != nil {
		return p.FailMsg(ctx, "权限不存在")
	}

	var roles []admin.AdminRole
	if err := facades.Orm().Query().Model(&permission).Association("AdminRoles").Find(&roles); err != nil {
		return p.FailMsg(ctx, "获取角色失败: " + err.Error())
	}

	return p.SuccessData(ctx, roles)
}

// GetMenusByPermission 根据权限获取菜单
func (p *PermissionController) GetMenusByPermission(ctx http.Context) http.Response {
	permissionId := ctx.Request().Input("permission_id", "")
	if permissionId == "" {
		return p.FailMsg(ctx, "权限ID不能为空")
	}

	var permission admin.AdminPermission
	if err := facades.Orm().Query().Find(&permission, permissionId); err != nil {
		return p.FailMsg(ctx, "权限不存在")
	}

	var menus []admin.AdminMenu
	if err := facades.Orm().Query().Model(&permission).Association("AdminMenus").Find(&menus); err != nil {
		return p.FailMsg(ctx, "获取菜单失败: " + err.Error())
	}

	return p.SuccessData(ctx, menus)
=======
	return r.SuccessData(ctx, r.list(ctx))
}

// 返回列表页面
func (r *PermissionController) list(ctx http.Context) *renderers.Page {
	HeaderToolbar := []any{
		r.CreateButton(ctx, r.form(ctx), true, "md", "", ""),
	}
	HeaderToolbar = append(HeaderToolbar, r.BaseHeaderToolBar()...)

	crud := r.BaseCRUD(ctx).HeaderToolbar(HeaderToolbar).Filter(
		r.BaseFilter().Body([]any{
			gamis.TextControl().Name("keyword").Label(tools.AdminLang(ctx, "keyword")).Size("md").Placeholder(tools.AdminLang(ctx, "admin_permission.name")),
		}),
	).Columns([]any{
		gamis.TableColumn().Name("id").Label("ID").Sortable(true),
		gamis.TableColumn().Name("name").Label(tools.AdminLang(ctx, "admin_permission.name")),
		gamis.TableColumn().Name("slug").Label(tools.AdminLang(ctx, "admin_permission.slug")),
		gamis.TableColumn().Name("http_method").Label(tools.AdminLang(ctx, "admin_permission.http_method")).,
		gamis.TableColumn().Name("http_path").Label(tools.AdminLang(ctx, "admin_permission.http_path")),
		gamis.TableColumn().Name("custom_order").Label(tools.AdminLang(ctx, "order")).Sortable(true),
		gamis.TableColumn().Name("created_at").Label(tools.AdminLang(ctx, "created_at")).Set("type", "datetime").Sortable(true),
		r.RowActions(ctx, r.form(ctx), []any{
			r.RowEditButton(ctx, r.form(ctx), true, "md", "", ""),
			r.RowDeleteButton(ctx, ""),
		}, "md"),
	})

	return r.BaseList(crud)
}

func (r *PermissionController) form(ctx http.Context) *renderers.Form {
	// 使用模型的HTTP方法选项
	permission := &admin.AdminPermission{}
	httpMethodOptions := permission.GetHttpMethodOptions()

	return r.BaseForm(ctx, false).Body([]any{
		gamis.SelectControl().Name("parent_id").Label(tools.AdminLang(ctx, "parent")).
			Searchable(true).LabelField("name").ValueField("id").
			Options(r.Service.PermissionOptions(ctx)).
			Placeholder(tools.AdminLang(ctx, "please_select")),
		gamis.TextControl().Name("name").Label(tools.AdminLang(ctx, "admin_permission.name")).Required(true).
			Placeholder(tools.AdminLang(ctx, "admin_permission.name")),
		gamis.TextControl().Name("slug").Label(tools.AdminLang(ctx, "admin_permission.slug")).Required(true).
			Placeholder(tools.AdminLang(ctx, "admin_permission.slug")),
		gamis.CheckboxesControl().Name("http_method").Label(tools.AdminLang(ctx, "admin_permission.http_method")).
			Options(httpMethodOptions).
			Description(tools.AdminLang(ctx, "admin_permission.http_method_description")),
		gamis.TextareaControl().Name("http_path").Label(tools.AdminLang(ctx, "admin_permission.http_path")).
			Placeholder(tools.AdminLang(ctx, "admin_permission.http_path_placeholder")).
			Description(tools.AdminLang(ctx, "admin_permission.http_path_description")).
			Set("minRows", 3).
			Set("maxRows", 10),
		gamis.NumberControl().Name("custom_order").Label(tools.AdminLang(ctx, "order")).
			Value(0).Min(0).Placeholder(tools.AdminLang(ctx, "order_placeholder")),
	})
>>>>>>> 08c77dc3ed68fd34ac5aa196c797580ff3c72dcb
}
