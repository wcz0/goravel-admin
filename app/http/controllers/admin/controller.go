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

type Controller interface {
	Form(ctx http.Context) *renderers.Form
}

type ControllerImpl[T any] struct {
	Service T
	*response.ErrorHandler
	Extra Extra
}

// 基础属性
type Extra struct {
	AdminPrefix string
	IsCreate    bool
	IsEdit      bool
}

func NewAdminController[T any](service T, extra ...Extra) *ControllerImpl[T] {
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
	a := &ControllerImpl[T]{
		ErrorHandler: response.NewErrorHandler(),
		Service:      service,
		Extra:        e,
	}
	return a
}

// 获取基础url
func (e *Extra) QueryPath(ctx http.Context) string {
	path := ctx.Request().Path()
	path = strings.TrimPrefix(path, e.AdminPrefix)
	return path
}

func (c *ControllerImpl[T]) ActionOfGetData(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "getData"
}

func (c *ControllerImpl[T]) ActionOfExport(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "export"
}

func (c *ControllerImpl[T]) ActionOfQuickEdit(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "quickEdit"
}

func (c *ControllerImpl[T]) ActionOfQuickEditItem(ctx http.Context) bool {
	return ctx.Request().Input("_action") == "quickEditItem"
}


/**
 * QueryPathTrait
 */

 // 获取列表数据
func (c *ControllerImpl[T]) GetListGetDataPath(ctx http.Context) string {
	return tools.GetAdmin(c.Extra.QueryPath(ctx) + "?_action=getData")
}

// 获取导出数据
func (c *ControllerImpl[T]) GetExportPath(ctx http.Context) string {
	return tools.GetAdmin(c.Extra.QueryPath(ctx) + "?_action=export")
}

// 删除路径 ?
func (c *ControllerImpl[T]) GetDeletePath(ctx http.Context, primaryKey ...string) string {
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	return "delete:" + tools.GetAdmin(c.Extra.QueryPath(ctx)+"/${"+key+"}")
}

// 批量删除 ?
func (c *ControllerImpl[T]) GetBulkDeletePath(ctx http.Context) string {
	return "delete:" + tools.GetAdmin(c.Extra.QueryPath(ctx)+"/${ids}")
}

// 获取编辑页面路径 ?
func (c *ControllerImpl[T]) GetEditPath(ctx http.Context, primaryKey ...string) string {
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	return "/"+strings.Trim(c.Extra.QueryPath(ctx), "/")+"/${"+key+"}/edit"
}

// 获取编辑数据 ?
func (c *ControllerImpl[T]) GetEditGetDataPath(ctx http.Context, primaryKey ...string) string {
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	path := c.Extra.QueryPath(ctx)
	paths := strings.Split(path, "/")
	last := paths[len(paths)-1]
	if last == "edit" {
		path = "/${"+key+"}/edit"
	}
	return tools.GetAdmin(path + "?_action=getData")
}

// 详情页面 ?
func (c *ControllerImpl[T]) GetShowPath(ctx http.Context, primaryKey ...string) string {
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	return "/"+strings.Trim(c.Extra.QueryPath(ctx), "/")+"/${"+key+"}"
}

// 编辑保存 ?
func (c *ControllerImpl[T]) GetUpdatePath(ctx http.Context, primaryKey ...string) string {
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	path := c.Extra.QueryPath(ctx)
	paths := strings.Split(path, "/")
	last := paths[len(paths)-1]
	if last == "edit" {
		path = "/${"+key+"}/edit"
	}
	return "put:" + tools.GetAdmin(path)
}

// 获取快速编辑数据
func (c *ControllerImpl[T]) GetQuickEditPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().FullUrl() + "?_action=quickEdit")
}

// 获取快速编辑项目数据
func (c *ControllerImpl[T]) GetQuickEditItemPath(ctx http.Context) string {
	return tools.GetAdmin(ctx.Request().FullUrl() + "?_action=quickEditItem")
}

// 获取详情 ?
func (c *ControllerImpl[T]) GetShowGetDataPath(ctx http.Context) (string, error) {
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
func (c *ControllerImpl[T]) GetCreatePath(ctx http.Context) string {
	return "/" + strings.Trim(c.Extra.QueryPath(ctx), "/") + "/create"
}

// 获取 新增 保存 的路径
func (c *ControllerImpl[T]) GetStorePath(ctx http.Context) string {
	return "post:" + tools.GetAdmin(ctx.Request().Path())
}

/**
 * 获取列表
 */
 func (c *ControllerImpl[T]) GetListPath(ctx http.Context) string {
	path := c.Extra.QueryPath(ctx)
	return path
}

/**
 * ElementTrait
 */

/**
 * 获取基础页面
 */
func (c *ControllerImpl[T]) BasePage() *renderers.Page {
	return gamis.Page().ClassName("m:overflow-auto")
}

/**
 * 返回列表按钮
 */
func (a *ControllerImpl[T]) BackButton(ctx http.Context) *renderers.OtherAction {
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

/**
* 批量删除按钮
 */
func (c *ControllerImpl[T]) BulkDeleteButton(ctx http.Context) *renderers.DialogAction {
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
func (c *ControllerImpl[T]) CreateButton(ctx http.Context, form renderers.Form, dialog bool, size string, title string, _type string) *renderers.OtherAction {
	if title == "" {
		title = tools.AdminLang(ctx, "admin.create")
	}
	action := gamis.LinkAction().Link(c.GetCreatePath(ctx))
	if dialog {
		form = *form.Api()
	}

	return nil
}

// 行编辑按钮
func (c *ControllerImpl[T]) rowEditButton(ctx http.Context, form renderers.Form, dialog bool, size string, title string) *renderers.LinkAction {
	// todo
	if title == "" {
		title = tools.AdminLang(ctx, "admin.edit")
	}
	action := gamis.LinkAction().Link(c.GetEditPath(ctx))
	if dialog {

	}
	action = action.Label(title).Level("link")
	return action
}

// 行详情按钮
func (c *ControllerImpl[T]) rowShowButton(ctx http.Context, dialog bool, size string, title string) *renderers.DialogAction {
	return nil
}

// 行删除按钮
func (c *ControllerImpl[T]) rowDeleteButton(ctx http.Context, dialog bool, size string, title string) *renderers.DialogAction {
	return nil
}

// 行操作按钮
func (c *ControllerImpl[T]) rowActions(ctx http.Context, dialog bool, size string, title string) []any {
	return nil
}

// 基础筛选器
func (c *ControllerImpl[T]) BaseFilter() *renderers.Form {
	return gamis.Form().PanelClassName("base-filter").Title("").Actions([]any{
		gamis.Button().Label("重置").ActionType("clear-and-submit"),
		gamis.Component("submit").Label("搜索").Level("primary"),
	})
}

// 基础筛选 - 条件构造器
func (c *ControllerImpl[T]) baseFilterConditionBuilder(ctx http.Context) map[string]any {
	return nil
}

// 基础 CRUD
func (c *ControllerImpl[T]) BaseCRUD(ctx http.Context) *renderers.CRUDTable {
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

// 基础顶部工具栏
func (c *ControllerImpl[T]) BaseHeaderToolBar() any {
	return []any{
		"bulkActions",
		gamis.Component("reload").Align("right"),
		gamis.Component("filter-toggler").Align("right"),
	}
}

/**
 * 获取基础表单
 */
 func (c *ControllerImpl[T]) BaseForm(ctx http.Context, back bool) *renderers.Form {
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

// 基础详情
func (c *ControllerImpl[T]) BaseDetail(ctx http.Context) *renderers.Form {
	api, _ := c.GetShowGetDataPath(ctx)
	return gamis.Form().
		PanelClassName("px-48 m:px-0").
		Title("").
		Mode("horizontal").
		Actions([]any{}).
		InitApi(api)
}

// 基础列表 #
func (c *ControllerImpl[T]) BaseList(crud any) *renderers.Page {
	return gamis.Page().ClassName("m:overflow-auto").Body(crud)
}

// 导出按钮 ?
func (c *ControllerImpl[T]) ExportAction(ctx http.Context, disableSelectedItem bool) *renderers.Service {
	return gamis.Service().
		Id("export-action").
		Set("align", "right").
		Set("data", map[string]any{
			"showExportLoading": false,
		}).Body([]any{

		})
}