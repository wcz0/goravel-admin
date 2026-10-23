package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/wcz0/gamis/v6"
	"github.com/wcz0/gamis/v6/renderers"

	"github.com/wcz0/goravel-admin/services"
	"github.com/wcz0/goravel-admin/tools"
)

// ExtensionsController 负责「开发者工具 - 扩展」。
//
// Go 版没有运行时加载扩展代码的能力，因此这里只维护扩展名单与启停状态：
// 扩展本体的控制器、路由、菜单需要在宿主编译期注册，前端能看到的
// enabled_extensions 由这张表决定。
type ExtensionsController struct {
	*ControllerImpl[*services.AdminExtensionService]
}

func NewExtensionsController() *ExtensionsController {
	return &ExtensionsController{ControllerImpl: NewAdminController[*services.AdminExtensionService](services.NewAdminExtensionService())}
}

func (e *ExtensionsController) Index(ctx http.Context) http.Response {
	if e.ActionOfGetData(ctx) {
		return e.ControllerImpl.Service.List(ctx)
	}

	return e.SuccessData(ctx, e.list(ctx))
}

func (e *ExtensionsController) Store(ctx http.Context) http.Response {
	return e.ControllerImpl.Service.Store(ctx)
}

// Show 返回单条扩展记录（Resource 路由要求实现）。
func (e *ExtensionsController) Show(ctx http.Context) http.Response {
	return e.SuccessData(ctx, map[string]any{})
}

func (e *ExtensionsController) Update(ctx http.Context) http.Response {
	return e.ControllerImpl.Service.Enable(ctx)
}

func (e *ExtensionsController) Destroy(ctx http.Context) http.Response {
	return e.ControllerImpl.Service.Destroy(ctx)
}

// Enable 启用扩展。
func (e *ExtensionsController) Enable(ctx http.Context) http.Response {
	return e.ControllerImpl.Service.Enable(ctx)
}

// Disable 停用扩展。
func (e *ExtensionsController) Disable(ctx http.Context) http.Response {
	return e.ControllerImpl.Service.Disable(ctx)
}

func (e *ExtensionsController) list(ctx http.Context) *renderers.Page {
	header := []any{
		e.CreateButton(ctx, e.form(ctx), true, "md", "", ""),
	}
	header = append(header, e.BaseHeaderToolBar()...)

	crud := e.BaseCRUD(ctx).
		HeaderToolbar(header).
		FooterToolbar([]string{"statistics"}).
		BulkActions([]any{e.BulkDeleteButton(ctx)}).
		Columns([]any{
			gamis.TableColumn().Name("id").Label("ID").Sortable(true),
			gamis.TableColumn().Name("name").Label(tools.AdminLang(ctx, "extensions.name")),
			gamis.TableColumn().Name("is_enabled").Label(tools.AdminLang(ctx, "extensions.enable")).
				QuickEdit(map[string]any{"type": "switch", "saveImmediately": true}),
			gamis.TableColumn().Name("created_at").Label(tools.AdminLang(ctx, "created_at")).Type("datetime").Sortable(true),
			e.RowActions(ctx, e.form(ctx), []any{
				e.RowDeleteButton(ctx, ""),
			}, "md"),
		})

	return e.BaseList(crud)
}

func (e *ExtensionsController) form(ctx http.Context) *renderers.Form {
	return e.BaseForm(ctx, false).Body([]any{
		gamis.TextControl().Name("name").
			Label(tools.AdminLang(ctx, "extensions.name")).
			Required(true).
			Description(tools.AdminLang(ctx, "extensions.name_desc")),
	})
}
