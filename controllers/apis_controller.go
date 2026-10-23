package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/wcz0/gamis/v6"
	"github.com/wcz0/gamis/v6/renderers"

	"github.com/wcz0/goravel-admin/services"
	"github.com/wcz0/goravel-admin/tools"
)

// ApisController 负责「开发者工具 - 动态 API」。
//
// 与 owl-admin 的差异：模板是编译期注册的 handler，不能在后台粘贴代码即时生效，
// 因此这里只做接口记录（路径 + 模板 + 参数）的维护。
type ApisController struct {
	*ControllerImpl[*services.AdminApiService]
}

func NewApisController() *ApisController {
	return &ApisController{ControllerImpl: NewAdminController[*services.AdminApiService](services.NewAdminApiService())}
}

func (a *ApisController) Index(ctx http.Context) http.Response {
	if a.ActionOfGetData(ctx) {
		return a.ControllerImpl.Service.List(ctx)
	}

	return a.SuccessData(ctx, a.list(ctx))
}

func (a *ApisController) Show(ctx http.Context) http.Response {
	return a.ControllerImpl.Service.Show(ctx)
}

func (a *ApisController) Store(ctx http.Context) http.Response {
	return a.ControllerImpl.Service.Store(ctx)
}

func (a *ApisController) Update(ctx http.Context) http.Response {
	return a.ControllerImpl.Service.Update(ctx)
}

func (a *ApisController) Destroy(ctx http.Context) http.Response {
	return a.ControllerImpl.Service.Destroy(ctx)
}

// Templates 返回内置模板下拉选项。
func (a *ApisController) Templates(ctx http.Context) http.Response {
	return a.SuccessData(ctx, a.ControllerImpl.Service.Templates())
}

// ArgsSchema 返回指定模板的参数 schema。
func (a *ApisController) ArgsSchema(ctx http.Context) http.Response {
	name := ctx.Request().Input("template")
	for _, t := range a.ControllerImpl.Service.Templates() {
		if t["value"] == name {
			return a.SuccessData(ctx, t)
		}
	}

	return a.FailMsg(ctx, tools.AdminLang(ctx, "apis.template_format_error"))
}

func (a *ApisController) list(ctx http.Context) *renderers.Page {
	header := []any{a.CreateButton(ctx, a.form(ctx), true, "lg", "", "")}
	header = append(header, a.BaseHeaderToolBar()...)

	crud := a.BaseCRUD(ctx).
		FilterTogglable(true).
		HeaderToolbar(header).
		Filter(a.BaseFilter().Body([]any{
			gamis.TextControl().Name("title").Label(tools.AdminLang(ctx, "apis.title")).Size("md").Clearable(true),
			gamis.TextControl().Name("path").Label(tools.AdminLang(ctx, "apis.path")).Size("md").Clearable(true),
		})).
		FooterToolbar([]string{"statistics"}).
		BulkActions([]any{a.BulkDeleteButton(ctx)}).
		Columns([]any{
			gamis.TableColumn().Name("id").Label("ID").Sortable(true),
			gamis.TableColumn().Name("title").Label(tools.AdminLang(ctx, "apis.title")),
			gamis.TableColumn().Name("path").Label(tools.AdminLang(ctx, "apis.path")),
			gamis.TableColumn().Name("template_title").Label(tools.AdminLang(ctx, "apis.template")),
			gamis.TableColumn().Name("enabled").Label(tools.AdminLang(ctx, "apis.enabled")).
				QuickEdit(map[string]any{"type": "switch", "mode": "inline", "saveImmediately": true}),
			gamis.TableColumn().Name("updated_at").Label(tools.AdminLang(ctx, "updated_at")).Type("datetime").Sortable(true),
			a.RowActions(ctx, a.form(ctx), []any{
				a.RowEditButton(ctx, a.form(ctx), true, "lg", "", ""),
				a.RowDeleteButton(ctx, ""),
			}, "lg"),
		})

	return a.BaseList(crud)
}

func (a *ApisController) form(ctx http.Context) *renderers.Form {
	templates := a.ControllerImpl.Service.Templates()
	options := make([]map[string]any, 0, len(templates))
	for _, t := range templates {
		options = append(options, map[string]any{
			"label": t["label"],
			"value": t["value"],
		})
	}

	return a.BaseForm(ctx, false).Body([]any{
		gamis.TextControl().Name("title").Label(tools.AdminLang(ctx, "apis.title")).Required(true),
		gamis.TextControl().Name("path").Label(tools.AdminLang(ctx, "apis.path")).Required(true).
			Description(tools.AdminLang(ctx, "apis.path_desc")).
			Placeholder("eg: /admin/demo/list"),
		gamis.SelectControl().Name("template").Label(tools.AdminLang(ctx, "apis.template")).
			Options(options).Required(true),
		gamis.SwitchControl().Name("enabled").Label(tools.AdminLang(ctx, "apis.enabled")).
			OnText(tools.AdminLang(ctx, "admin.yes")).OffText(tools.AdminLang(ctx, "admin.no")).Value(1),
		// amis 的 json 控件没有对应的 gamis 渲染器，直接给 schema 字面量
		map[string]any{
			"type":   "json",
			"name":   "args",
			"label":  tools.AdminLang(ctx, "apis.args"),
			"height": 300,
		},
	})
}
