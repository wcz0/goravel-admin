package admin

import (
	"errors"
	"fmt"
	"goravel/app/response"
	"goravel/app/tools"
	"reflect"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/wcz0/gamis"
	"github.com/wcz0/gamis/renderers"
)

type Controller[T any] struct {
	Service T
	*response.ErrorHandler
	Extra Extra
}

type Extra struct {
	AdminPrefix string
	IsCreate    bool
	IsEdit      bool
}

func NewAdminController[T any](service T, extra ...Extra) *Controller[T] {
	var e Extra
	if len(extra) > 0 {
		e = extra[0]
	} else {
		e = Extra{
			AdminPrefix: facades.Config().GetString("admin.route.prefix"),
			IsCreate:    false,
			IsEdit:      false,
		}
	}
	a := &Controller[T]{
		ErrorHandler: response.NewErrorHandler(),
		Service:      service,
		Extra:        e,
	}
	return a
}

func (e *Extra) QueryPath(ctx http.Context) string {
	path := ctx.Request().Path()
	path = strings.TrimPrefix(path, e.AdminPrefix)
	return path
}

func (c *Controller[T]) ActionOfGetData(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "getData"
}

func (c *Controller[T]) ActionOfExport(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "export"
}

func (c *Controller[T]) ActionOfQuickEdit(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "quickEdit"
}

func (c *Controller[T]) ActionOfQuickEditItem(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "quickEditItem"
}

func (c *Controller[T]) BaseList(crud any) *renderers.Page {
	return c.BasePage().Body(crud)
}

/**
 * 获取基础页面
 */
func (c *Controller[T]) BasePage() *renderers.Page {
	return gamis.Page().ClassName("m:overflow-auto")
}

/**
 * 获取基础表单
 */
func (c *Controller[T]) BaseForm(ctx http.Context, back bool) *renderers.Form {
	path := ctx.Request().Path()
	path = strings.TrimPrefix(path, facades.Config().GetString("admin.route.prefix"))
	form := gamis.Form().PanelClassName("px-48 m:px-0").Title("").PromptPageLeave("")
	if back {
		form.OnEvent(map[string]any{
			"submitSucc": map[string]any{
				"action": []any{
					map[string]any{
						"actionType": "custom",
						"script":     "window.history.back()",
					},
					map[string]any{
						"actionType": "custom",
						"script":     fmt.Sprintf(`window.$owl.hasOwnProperty('closeTabByPath') && window.$owl.closeTabByPath('%s')`, path),
					},
				},
			},
		})
	}
	return form
}

func (c *Controller[T]) BaseDetail(ctx http.Context) *renderers.Form {
	api, _ := c.GetShowGetDataPath(ctx)
	return gamis.Form().
		PanelClassName("px-48 m:px-0").
		Title("").
		Mode("horizontal").
		Actions([]any{}).
		InitApi(api)
}

/**
 * 返回列表按钮
 */
func (a *Controller[T]) BackButton(ctx http.Context) *renderers.OtherAction {
	// return gamis.Button().Label("返回").ClassName("w-full")
	path := ctx.Request().Path()
	path = strings.TrimPrefix(path, a.Extra.AdminPrefix)
	script := fmt.Sprintf(`window.$owl.hasOwnProperty('closeTabByPath') && window.$owl.closeTabByPath('%s')`, path)
	return gamis.OtherAction().
		Label(tools.AdminLang(ctx, "admin.back")).
		Icon("fa-solid fa-chevron-left").
		Level("primary").
		OnClick("window.history.back();" + script)
}

func (c *Controller[T]) BaseFilter() *renderers.Form {
	return gamis.Form().PanelClassName("base-filter").Title("").Actions([]any{
		gamis.Button().Label("重置").ActionType("clear-and-submit"),
		gamis.Component("submit").Label("搜索").Level("primary"),
	})
}

func (c *Controller[T]) BaseHeaderToolBar() any {
	return []any{
		"bulkActions",
		gamis.Component("reload").Align("right"),
		gamis.Component("filter-toggler").Align("right"),
	}
}

// 获取列表数据
func (c *Controller[T]) GetListGetDataPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().Path() + "?_action=getData")
}

// 获取 新增 保存 的路径
func (c *Controller[T]) GetStorePath(ctx http.Context) string {
	return "post:" + tools.GetAdmin(ctx.Request().Path())
}

// 获取导出数据
func (c *Controller[T]) GetExportPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().FullUrl() + "?_action=export")
}

// 获取快速编辑数据
func (c *Controller[T]) GetQuickEditPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().FullUrl() + "?_action=quickEdit")
}

// 获取快速编辑项目数据
func (c *Controller[T]) GetQuickEditItemPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().FullUrl() + "?_action=quickEditItem")
}

// 删除路径 TODO: 未测试
// func (c *AdminController[T]Impl) GetDeletePath(ctx http.Context) string {
// 	return "delete:" + tools.GetAdmin(c.Controller.QueryPath()+"/"+ctx.Request().Input("id"))
// }

// 批量删除 TODO: 未测试
func (c *Controller[T]) GetBulkDeletePath(ctx http.Context) string {
	return "delete:" + tools.GetAdmin(c.Extra.QueryPath(ctx)+"/"+ctx.Request().Input("ids"))
}

// 获取编辑页面路径
func (c *Controller[T]) GetEditPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().Path() + "/" + ctx.Request().Input("id") + "/edit")
}

// 获取编辑数据
func (c *Controller[T]) GetEditGetDataPath(ctx http.Context) string {
	return c.GetEditPath(ctx) + "?_action=getData"
}

// 详情页面
func (c *Controller[T]) GetShowPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().Path() + "/" + ctx.Request().Input("id"))
}

// 编辑保存
func (c *Controller[T]) GetUpdatePath(ctx http.Context) string {
	return "put:" + tools.GetAdmin(ctx.Request().Path()+"/"+ctx.Request().Input("id"))
}

// 获取详情
// TODO: 未测试
func (c *Controller[T]) GetShowGetDataPath(ctx http.Context) (string, error) {
	path := c.Extra.QueryPath(ctx)
	// 获取字段值
	service := reflect.ValueOf(c.Service)
	if service.Kind() == reflect.Ptr {
		service = service.Elem()
	}
	model := service.FieldByName("Model")
	if !model.IsValid() {
		return "", errors.New("model not found")
	}
	primaryKey := model.Addr().MethodByName("PrimaryKey")
	if !primaryKey.IsValid() {
		return "", errors.New("primary key not found")
	}
	values := primaryKey.Call(nil)
	if len(values) == 0 || values[0].Kind() != reflect.String {
		return "", errors.New("primary key not found")
	}
	path = path + "/{" + values[0].String() + "}"
	return tools.GetAdmin(path + "?_action=getData"), nil
}

// 新增页面
func (c *Controller[T]) GetCreatePath(ctx http.Context) string {
	return "/"+ strings.Trim(c.Extra.QueryPath(ctx), "/") + "/create"
}

func (c *Controller[T]) BaseCRUD(ctx http.Context) *renderers.CRUDTable {
	crud := gamis.CRUDTable().PerPage(20).AffixHeader(false).FilterTogglable(true).FilterDefaultVisible(false).
		Api(c.GetListGetDataPath(ctx)).QuickSaveApi(c.GetQuickEditPath(ctx)).QuickSaveItemApi(c.GetQuickEditItemPath(ctx)).
		BulkActions([]any{
			c.BulkDeleteButton(ctx),
		}).PerPageAvailable([]int{10, 20, 50, 100}).FooterToolbar([]string{
		"switch-per-page",
		"statistics",
		"pagination",
	}).HeaderToolbar([]any{})
	return crud
}

// func (c *AdminController[T]Impl) GetListGetDataPath() string {
// 	return tools.GetAdmin(c.Controller.QueryPath() + "?_action=getData")
// }

/**
* 批量删除按钮
 */
func (c *Controller[T]) BulkDeleteButton(ctx http.Context) *renderers.DialogAction {
	return gamis.DialogAction().
		Label("删除").
		Icon("fa-solid fa-trash-can").
		Dialog(
			gamis.Dialog().Title("删除").
				ClassName("py-2").Actions([]any{
				gamis.Action().ActionType("cancel").Label(tools.AdminLang(ctx, "admin.cancel")),
				gamis.Action().ActionType("submit").Label(tools.AdminLang(ctx, "admin.delete")).Level("danger"),
			}).Body([]any{
				gamis.Form().WrapWithPanel(false).Api(c.GetBulkDeletePath(ctx)).Body([]any{
					gamis.Tpl().ClassName("py-2").Tpl(tools.AdminLang(ctx, "admin.confirm_delete")),
				}),
			}))
}

// 创建按钮
func (c *Controller[T]) CreateButton(ctx http.Context, dialog bool, size string) *renderers.OtherAction {
	// if dialog {
	// 	form :=
	// }
	return nil
}

/**
 * 获取列表
 */
func (c *Controller[T]) GetListPath(ctx http.Context) string {
	path := c.Extra.QueryPath(ctx)
	return path
}
