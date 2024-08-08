package controllers

import (
	"goravel/app/response"
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"
	"github.com/wcz0/gamis"
	"github.com/wcz0/gamis/renderers"
)

type AdminController interface {
	QueryPath() string
	AdminPrefix() string
	isCreate() bool
	isEdit() bool
}

type Controller struct {
	*response.ErrorHandler
	Controller AdminController
}

func NewController() *Controller {
	return &Controller{
		ErrorHandler: response.NewErrorHandler(),
		Controller:   nil,
	}
}

func (c *Controller) ActionOfGetData(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "getData"
}

func (c *Controller) ActionOfExport(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "export"
}

func (c *Controller) ActionOfQuickEdit(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "quickEdit"
}

func (c *Controller) ActionOfQuickEditItem(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "quickEditItem"
}

func (c *Controller) BaseList(crud any) *renderers.Page {
	return c.BasePage().Body(crud)
}

func (c *Controller) BasePage() *renderers.Page {
	return gamis.Page().ClassName("m:overflow-auto")
}

// func (c. *Controller) BaseForm() *renderers.Form {
// 	return gamis.Form().PanelClassName("border-none").Id("form").Title("").Api(c.GetApi()).InitApi("/no-content").Body([]any{})
// }

func (c *Controller) BackButton() *renderers.OtherAction {
	// return gamis.Button().Label("返回").ClassName("w-full")

	return gamis.OtherAction().Label("返回").Icon("fa-solid fa-chevron-left").Level("primary").OnClick("window.history.back()")
}

func (c *Controller) BaseFilter() *renderers.Form {
	return gamis.Form().PanelClassName("base-filter").Title("").Actions([]any{
		gamis.Button().Label("重置").ActionType("clear-and-submit"),
		gamis.Component("submit").Label("搜索").Level("primary"),
	})
}

func (c *Controller) BaseHeaderToolBar() any {
	return []any{
		"bulkActions",
		gamis.Component("reload").Align("right"),
		gamis.Component("filter-toggler").Align("right"),
	}
}

// 获取列表数据
func (c *Controller) GetListGetDataPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().Path() + "?_action=getData")
}

// 获取 新增 保存 的路径
func (c *Controller) GetStorePath(ctx http.Context) string {
	return "post:" + tools.GetAdmin(ctx.Request().Path())
}

// 获取导出数据
func (c *Controller) GetExportPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().FullUrl() + "?_action=export")
}

// 获取快速编辑数据
func (c *Controller) GetQuickEditPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().FullUrl() + "?_action=quickEdit")
}

// 获取快速编辑项目数据
func (c *Controller) GetQuickEditItemPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().FullUrl() + "?_action=quickEditItem")
}

// 删除路径 TODO: 未测试
func (c *Controller) GetDeletePath(ctx http.Context) string {
	return "delete:" + tools.GetAdmin(c.Controller.QueryPath() + "/" + ctx.Request().Input("id"))
}

// 批量删除 TODO: 未测试
func (c *Controller) GetBulkDeletePath(ctx http.Context) string {
	return "delete:" + tools.GetAdmin(c.Controller.QueryPath() + "/" + ctx.Request().Input("ids"))
}

// 获取编辑页面路径
func (c *Controller) GetEditPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().Path() + "/" + ctx.Request().Input("id") + "/edit")
}

// 获取编辑数据
func (c *Controller) GetEditGetDataPath(ctx http.Context) string {
	return c.GetEditPath(ctx) + "?_action=getData"
}

// 详情页面
func (c *Controller) GetShowPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().Path() + "/" + ctx.Request().Input("id"))
}

// 编辑保存
func (c *Controller) GetUpdatePath(ctx http.Context) string {
	return "put:" + tools.GetAdmin(ctx.Request().Path() + "/" + ctx.Request().Input("id"))
}

// 获取详情
// TODO: 未测试
func (c *Controller) GetShowGetDataPath(ctx http.Context, id string) string {
	return tools.GetAdmin(ctx.Request().FullUrl() + "/" + id + "?_action=getData")
}

// 新增页面
func (c *Controller) GetCreatePath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().FullUrl() + "/create")
}

func (c *Controller) BaseCRUD(ctx http.Context) *renderers.CRUDTable {
	crud := gamis.CRUDTable().PerPage(20).AffixHeader(false).FilterTogglable(true).FilterDefaultVisible(false).
		Api(c.GetListGetDataPath(ctx)).QuickSaveApi(c.GetQuickEditPath(ctx)).QuickSaveItemApi(c.GetQuickEditItemPath(ctx)).
		BulkActions([]any{
			c.BulkDeleteButton(ctx),
		}).PerPageAvailable([]int{10, 20, 50, 100}).FooterToolbar([]string{
			"switch-per-page",
			"statistics",
			"pagination",
		}).HeaderToolbar([]any{

		})
	return crud
}

// func (c *Controller) GetListGetDataPath() string {
// 	return tools.GetAdmin(c.Controller.QueryPath() + "?_action=getData")
// }

// 批量删除按钮
func (c *Controller) BulkDeleteButton(ctx http.Context) *renderers.DialogAction {
	return gamis.DialogAction().Label("删除").Icon("fa-solid fa-trash-can").Dialog(gamis.Dialog().Title("删除").
		ClassName("py-2").Actions([]any{
			gamis.Action().ActionType("cancel").Label("取消"),
			gamis.Action().ActionType("submit").Label("确定").Level("danger"),
		}).Body([]any{
			gamis.Form().WrapWithPanel(false).Api(c.GetBulkDeletePath(ctx)).Body([]any{
				gamis.Tpl().ClassName("py-2").Tpl("确定要删除选中的数据吗?"),
			}),
		}))
}

// 创建按钮
func (c *Controller) CreateButton(ctx http.Context, dialog bool, size string) *renderers.OtherAction {
	// if dialog {
	// 	form :=
	// }
	return nil
}