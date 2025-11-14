package admin

import (
	"goravel/app/models/admin"
	"goravel/app/services"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type PermissionController struct {
	*ControllerImpl[*services.AdminPermissionService]
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		ControllerImpl: NewAdminController[*services.AdminPermissionService](services.NewAdminPermissionService()),
	}
}

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
	}

	return p.ControllerImpl.Service.Store(ctx)
}

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
	}

	return p.ControllerImpl.Service.Update(ctx)
}

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
}
