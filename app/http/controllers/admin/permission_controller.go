package admin

import (
	"goravel/app/services"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/gamis"
	"github.com/wcz0/gamis/renderers"
)

type PermissionController struct {
	*Controller[*services.AdminPermissionService] // 继承
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		Controller: NewAdminController[*services.AdminPermissionService](services.NewAdminPermissionService()),
	}
}

func (r *PermissionController) Show(ctx http.Context) http.Response {
	return nil
}

func (r *PermissionController) Store(ctx http.Context) http.Response {
	validation, err := ctx.Request().Validate(map[string]string{
		"parent_id": "number",
		"name":      "required|string",
		"value":     "required|string",
		"method":    "required|string",
	})
	if err != nil {
		return r.FailMsg(ctx, err.Error())
	}
	if validation.Fails() {
		return r.FailMsg(ctx, validation.Errors().All())
	}
	return r.Controller.Service.Store(ctx)
}

func (r *PermissionController) Update(ctx http.Context) http.Response {
	validation, err := ctx.Request().Validate(map[string]string{
		"id":        "required|number",
		"parent_id": "number",
		"name":      "required|string",
		"value":     "required|string",
		"method":    "required|string",
	})
	if err != nil {
		return r.FailMsg(ctx, err.Error())
	}
	if validation.Fails() {
		return r.FailMsg(ctx, validation.Errors().All())
	}
	return r.Controller.Service.Update(ctx)
}

func (r *PermissionController) Destroy(ctx http.Context) http.Response {
	validation, err := ctx.Request().Validate(map[string]string{
		"id": "required|number",
	})
	if err != nil {
		return r.FailMsg(ctx, err.Error())
	}
	if validation.Fails() {
		return r.FailMsg(ctx, validation.Errors().All())
	}
	return r.Controller.Service.Destroy(ctx)
}

func (r *PermissionController) Index(ctx http.Context) http.Response {
	if r.ActionOfGetData(ctx) {
		return r.Controller.Service.List(ctx)
	}
	// return r.DataSuccess(ctx, list(ctx))
	return r.Success(ctx)
}

// 返回列表页面
func (p *PermissionController) list(ctx http.Context) *renderers.Page {
	// var autoBtn any;
	// 自动生成权限 按钮开关
	if facades.Config().GetBool("admin.show_auto_generate_permission_button") {
		// autoBtn = gamis.AjaxAction().Label("自动生成权限").Level("success").
		// 	ConfirmText("确定要自动生成权限吗?").Api(tools.GetAdmin("system/_admin_permissions_auto_generate"))
	}

	// 自动生成权限
	// crud := gamis.CRUDTable().PerPage(20).AffixHeader(false).FilterTogglable(true).FilterDefaultVisible(false).
	// 	Api(tools.GetAdmin("system/_admin_permissions_auto_generate"))

	// crud := p.Controller.BaseCRUD(ctx).Api(p.GetListGetDataPath(ctx)).QuickSaveApi()

	return gamis.Page().Title("权限列表")
}
