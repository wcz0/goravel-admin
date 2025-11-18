package admin

import (
	"goravel/app/models/admin"
	"goravel/app/services"
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"

	"github.com/wcz0/gamis"
	"github.com/wcz0/gamis/renderers"
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
func (p *PermissionController) page(_ http.Context) interface{} {
    return map[string]interface{}{
        "title": "权限列表",
    }
}

func (p *PermissionController) Show(ctx http.Context) http.Response {
	return p.ControllerImpl.Service.Show(ctx)
}

func (p *PermissionController) Store(ctx http.Context) http.Response {
    if p.ActionOfQuickEdit(ctx) {
        return p.ControllerImpl.Service.QuickEdit(ctx)
    }
    if p.ActionOfQuickEditItem(ctx) {
        return p.ControllerImpl.Service.QuickEditItem(ctx)
    }
    // 验证输入
    validation, err := ctx.Request().Validate(map[string]string{
        "parent_id":  "number",
        "name":       "required|string|max:255",
        "value":      "required|string|max:255",
        "http_method": "string",
        "http_path":   "string",
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
    hasError, response := p.HandleValidationErrors(ctx, map[string]string{
        "id":           "required|number",
        "parent_id":    "number",
        "name":         "required|string|max:255",
        "value":        "required|string|max:255",
        "http_method":  "string",
        "http_path":    "string",
        "custom_order": "number",
    }, nil)
    if hasError {
        return response
    }

    return p.ControllerImpl.Service.Update(ctx)
}

func (p *PermissionController) Destroy(ctx http.Context) http.Response {
    hasError, response := p.HandleValidationErrors(ctx, map[string]string{
        "id": "required|number",
    }, nil)
    if hasError {
        return response
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
    return p.SuccessData(ctx, p.page(ctx))
}

func (p *PermissionController) Options(ctx http.Context) http.Response {
    return p.ControllerImpl.Service.PermissionOptions(ctx)
}

func (p *PermissionController) TreeOptions(ctx http.Context) http.Response {
    return p.ControllerImpl.Service.PermissionTreeOptions(ctx)
}

// 返回列表页面
func (p *PermissionController) list(ctx http.Context) *renderers.Page {
    HeaderToolbar := []any{
        p.CreateButton(ctx, p.form(ctx), true, "md", "", ""),
    }
    HeaderToolbar = append(HeaderToolbar, p.BaseHeaderToolBar()...)

    crud := p.BaseCRUD(ctx).HeaderToolbar(HeaderToolbar).Filter(
        p.BaseFilter().Body([]any{
            gamis.TextControl().Name("keyword").Label(tools.AdminLang(ctx, "keyword")).Size("md").Placeholder(tools.AdminLang(ctx, "admin_permission.name")),
        }),
    ).Columns([]any{
        gamis.TableColumn().Name("id").Label("ID").Sortable(true),
        gamis.TableColumn().Name("name").Label(tools.AdminLang(ctx, "admin_permission.name")),
        gamis.TableColumn().Name("slug").Label(tools.AdminLang(ctx, "admin_permission.slug")),
        gamis.TableColumn().Name("http_method").Label(tools.AdminLang(ctx, "admin_permission.http_method")),
        gamis.TableColumn().Name("http_path").Label(tools.AdminLang(ctx, "admin_permission.http_path")),
        gamis.TableColumn().Name("custom_order").Label(tools.AdminLang(ctx, "order")).Sortable(true),
        gamis.TableColumn().Name("created_at").Label(tools.AdminLang(ctx, "created_at")).Set("type", "datetime").Sortable(true),
        p.RowActions(ctx, p.form(ctx), []any{
            p.RowEditButton(ctx, p.form(ctx), true, "md", "", ""),
            p.RowDeleteButton(ctx, ""),
        }, "md"),
    })

    return p.BaseList(crud)
}

func (p *PermissionController) form(ctx http.Context) *renderers.Form {
    // 使用模型的HTTP方法选项
    permission := &admin.AdminPermission{}
    httpMethodOptions := permission.GetHttpMethodOptions()

    return p.BaseForm(ctx, false).Body([]any{
        gamis.TreeSelectControl().Name("parent_id").Label(tools.AdminLang(ctx, "parent")).
            LabelField("label").ValueField("value").Multiple(false).
            Source(tools.GetAdmin("system/permissions/tree_options")).
            Placeholder(tools.AdminLang(ctx, "please_select")),
        gamis.TextControl().Name("name").Label(tools.AdminLang(ctx, "admin_permission.name")).Required(true).
            Placeholder(tools.AdminLang(ctx, "admin_permission.name")),
        gamis.TextControl().Name("value").Label(tools.AdminLang(ctx, "admin_permission.slug")).Required(true).
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
}
