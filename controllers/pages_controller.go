package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/wcz0/gamis/v6"
	"github.com/wcz0/gamis/v6/renderers"

	"github.com/wcz0/goravel-admin/services"
	"github.com/wcz0/goravel-admin/tools"
)

// PagesController 负责「开发者工具 - 页面管理」。
//
// 页面结构保存为 amis JSON，前端通过 page_schema?sign=xxx 拉取后直接渲染，
// 所以这里的表单只做 JSON 编辑，不需要为每个页面单独写后端。
type PagesController struct {
	*ControllerImpl[*services.AdminPageService]
}

func NewPagesController() *PagesController {
	return &PagesController{ControllerImpl: NewAdminController[*services.AdminPageService](services.NewAdminPageService())}
}

func (p *PagesController) Index(ctx http.Context) http.Response {
	if p.ActionOfGetData(ctx) {
		return p.ControllerImpl.Service.List(ctx)
	}

	return p.SuccessData(ctx, p.list(ctx))
}

func (p *PagesController) Show(ctx http.Context) http.Response {
	return p.ControllerImpl.Service.Show(ctx)
}

func (p *PagesController) Store(ctx http.Context) http.Response {
	return p.ControllerImpl.Service.Store(ctx)
}

func (p *PagesController) Update(ctx http.Context) http.Response {
	return p.ControllerImpl.Service.Update(ctx)
}

func (p *PagesController) Destroy(ctx http.Context) http.Response {
	return p.ControllerImpl.Service.Destroy(ctx)
}

// Options 供菜单表单里的「页面」下拉选择使用。
func (p *PagesController) Options(ctx http.Context) http.Response {
	return p.ControllerImpl.Service.Options(ctx)
}

func (p *PagesController) list(ctx http.Context) *renderers.Page {
	header := []any{p.CreateButton(ctx, p.form(ctx), true, "lg", "", "")}
	header = append(header, p.BaseHeaderToolBar()...)

	crud := p.BaseCRUD(ctx).
		FilterTogglable(true).
		HeaderToolbar(header).
		Filter(p.BaseFilter().Body([]any{
			gamis.TextControl().Name("title").Label(tools.AdminLang(ctx, "pages.title")).Size("md").Clearable(true),
			gamis.TextControl().Name("sign").Label(tools.AdminLang(ctx, "pages.sign")).Size("md").Clearable(true),
		})).
		FooterToolbar([]string{"statistics"}).
		BulkActions([]any{p.BulkDeleteButton(ctx)}).
		Columns([]any{
			gamis.TableColumn().Name("id").Label("ID").Sortable(true),
			gamis.TableColumn().Name("title").Label(tools.AdminLang(ctx, "pages.title")),
			gamis.TableColumn().Name("sign").Label(tools.AdminLang(ctx, "pages.sign")),
			gamis.TableColumn().Name("updated_at").Label(tools.AdminLang(ctx, "updated_at")).Type("datetime").Sortable(true),
			p.RowActions(ctx, p.form(ctx), []any{
				p.RowEditButton(ctx, p.form(ctx), true, "lg", "", ""),
				p.RowDeleteButton(ctx, ""),
			}, "lg"),
		})

	return p.BaseList(crud)
}

func (p *PagesController) form(ctx http.Context) *renderers.Form {
	// custom-amis-editor 是前端自定义组件，gamis 没有对应的渲染器，直接给 schema 字面量
	editor := map[string]any{
		"type":      "custom-amis-editor",
		"name":      "schema",
		"label":     "",
		"mode":      "normal",
		"className": "h-full",
	}

	return p.BaseForm(ctx, false).Body([]any{
		gamis.TextControl().Name("title").Label(tools.AdminLang(ctx, "pages.title")).Required(true),
		gamis.TextControl().Name("sign").Label(tools.AdminLang(ctx, "pages.sign")).Required(true).Description(tools.AdminLang(ctx, "pages.sign_desc")),
		gamis.SubFormControl().Name("page").Label(tools.AdminLang(ctx, "pages.page")).Required(true).Form(
			gamis.Form().ClassName("h-full").Set("size", "full").Title("").Body([]any{editor}),
		),
	})
}
