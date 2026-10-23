package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/wcz0/gamis/v6"
	"github.com/wcz0/gamis/v6/renderers"

	"github.com/wcz0/goravel-admin/services"
	"github.com/wcz0/goravel-admin/tools"
)

// RelationshipsController 负责「开发者工具 - 动态关联」。
//
// Goravel 的 ORM 关联必须静态声明在模型上，所以这里不改运行时行为，
// 只负责登记关联元数据并给出可复制的 struct 方法预览。
type RelationshipsController struct {
	*ControllerImpl[*services.AdminRelationshipService]
}

func NewRelationshipsController() *RelationshipsController {
	return &RelationshipsController{ControllerImpl: NewAdminController[*services.AdminRelationshipService](services.NewAdminRelationshipService())}
}

func (r *RelationshipsController) Index(ctx http.Context) http.Response {
	if r.ActionOfGetData(ctx) {
		return r.ControllerImpl.Service.List(ctx)
	}

	return r.SuccessData(ctx, r.list(ctx))
}

func (r *RelationshipsController) Show(ctx http.Context) http.Response {
	return r.ControllerImpl.Service.Show(ctx)
}

func (r *RelationshipsController) Store(ctx http.Context) http.Response {
	return r.ControllerImpl.Service.Store(ctx)
}

func (r *RelationshipsController) Update(ctx http.Context) http.Response {
	return r.ControllerImpl.Service.Update(ctx)
}

func (r *RelationshipsController) Destroy(ctx http.Context) http.Response {
	return r.ControllerImpl.Service.Destroy(ctx)
}

// Types 返回关联类型下拉。
func (r *RelationshipsController) Types(ctx http.Context) http.Response {
	return r.SuccessData(ctx, r.ControllerImpl.Service.Types())
}

func (r *RelationshipsController) list(ctx http.Context) *renderers.Page {
	header := []any{r.CreateButton(ctx, r.form(ctx), true, "lg", "", "")}
	header = append(header, r.BaseHeaderToolBar()...)

	crud := r.BaseCRUD(ctx).
		FilterTogglable(true).
		HeaderToolbar(header).
		Filter(r.BaseFilter().Body([]any{
			gamis.TextControl().Name("model").Label(tools.AdminLang(ctx, "relationships.model")).Size("md").Clearable(true),
			gamis.TextControl().Name("title").Label(tools.AdminLang(ctx, "relationships.title")).Size("md").Clearable(true),
		})).
		FooterToolbar([]string{"statistics"}).
		BulkActions([]any{r.BulkDeleteButton(ctx)}).
		Columns([]any{
			gamis.TableColumn().Name("id").Label("ID").Sortable(true),
			gamis.TableColumn().Name("model").Label(tools.AdminLang(ctx, "relationships.model")),
			gamis.TableColumn().Name("title").Label(tools.AdminLang(ctx, "relationships.title")),
			gamis.TableColumn().Name("type").Label(tools.AdminLang(ctx, "relationships.type")),
			gamis.TableColumn().Name("remark").Label(tools.AdminLang(ctx, "relationships.remark")),
			gamis.TableColumn().Name("updated_at").Label(tools.AdminLang(ctx, "updated_at")).Type("datetime").Sortable(true),
			r.RowActions(ctx, r.form(ctx), []any{
				r.RowEditButton(ctx, r.form(ctx), true, "lg", "", ""),
				r.RowDeleteButton(ctx, ""),
			}, "lg"),
		})

	return r.BaseList(crud)
}

func (r *RelationshipsController) form(ctx http.Context) *renderers.Form {
	types := r.ControllerImpl.Service.Types()
	options := make([]map[string]string, 0, len(types))
	for _, t := range types {
		options = append(options, map[string]string{"label": t["label"], "value": t["value"]})
	}

	return r.BaseForm(ctx, false).Body([]any{
		gamis.TextControl().Name("model").Label(tools.AdminLang(ctx, "relationships.model")).Required(true).
			Description(tools.AdminLang(ctx, "relationships.model_desc")).
			Placeholder("eg: AdminUser"),
		gamis.TextControl().Name("title").Label(tools.AdminLang(ctx, "relationships.title")).Required(true).
			Placeholder("eg: AdminRoles"),
		gamis.SelectControl().Name("type").Label(tools.AdminLang(ctx, "relationships.type")).Options(options).Required(true),
		gamis.TextControl().Name("remark").Label(tools.AdminLang(ctx, "relationships.remark")),
		map[string]any{
			"type":   "json",
			"name":   "args",
			"label":  tools.AdminLang(ctx, "relationships.args"),
			"height": 260,
			"value":  map[string]any{"related_model": "", "foreign_key": "", "local_key": ""},
		},
	})
}
