package admin

import (
	"errors"
	"fmt"
	"goravel/app/enums"
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
	return tools.Url(c.Extra.QueryPath(ctx) + "?_action=getData")
}

// 获取导出数据
func (c *ControllerImpl[T]) GetExportPath(ctx http.Context) string {
	return tools.Url(c.Extra.QueryPath(ctx) + "?_action=export")
}

// 删除路径 ?
func (c *ControllerImpl[T]) GetDeletePath(ctx http.Context, primaryKey ...string) string {
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	return "delete:" + tools.Url(c.Extra.QueryPath(ctx)+"/${"+key+"}")
}

// 批量删除 ?
func (c *ControllerImpl[T]) GetBulkDeletePath(ctx http.Context) string {
	return "delete:" + tools.Url(c.Extra.QueryPath(ctx)+"/${ids}")
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
	return tools.Url(path + "?_action=getData")
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
	return "put:" + tools.Url(path)
}

// 获取快速编辑数据
func (c *ControllerImpl[T]) GetQuickEditPath(ctx http.Context) string {
	return tools.Url(ctx.Request().FullUrl()) + "?_action=quickEdit"
}

// 获取快速编辑项目数据
func (c *ControllerImpl[T]) GetQuickEditItemPath(ctx http.Context) string {
	return tools.Url(ctx.Request().FullUrl()) + "?_action=quickEditItem"
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
	return tools.Url(path + "?_action=getData"), nil
}

// 新增页面
func (c *ControllerImpl[T]) GetCreatePath(ctx http.Context) string {
	return "/" + strings.Trim(c.Extra.QueryPath(ctx), "/") + "/create"
}

// 获取 新增 保存 的路径
func (c *ControllerImpl[T]) GetStorePath(ctx http.Context) string {
	return "post:" + tools.Url(c.Extra.QueryPath(ctx))
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
		Label(tools.AdminLang(ctx, "back")).
		Icon("fa-solid fa-chevron-left").
		Level("primary").
		OnClick("window.history.back();" + script)
}

/**
* 批量删除按钮
 */
func (c *ControllerImpl[T]) BulkDeleteButton(ctx http.Context) *renderers.DialogAction {
	return gamis.DialogAction().
		Label(tools.AdminLang(ctx, "delete")).
		Icon("fa-solid fa-trash-can").
		Dialog(
			gamis.Dialog().Title(tools.AdminLang(ctx, "delete")).
				ClassName("py-2").Actions([]any{
					gamis.Action().ActionType("cancel").Label(tools.AdminLang(ctx, "cancel")),
					gamis.Action().ActionType("submit").Label(tools.AdminLang(ctx, "delete")).Level("danger"),
			}).Body([]any{
				gamis.Form().WrapWithPanel(false).Api(c.GetBulkDeletePath(ctx)).Body([]any{
					gamis.Tpl().ClassName("py-2").Tpl(tools.AdminLang(ctx, "confirm_delete")),
				}),
			}),
		)
}




// 创建按钮
func (c *ControllerImpl[T]) CreateButton(ctx http.Context, form *renderers.Form, dialog bool, size string, title string, _type string) *renderers.LinkAction {
	if title == "" {
		title = tools.AdminLang(ctx, "create")
	}
	action := gamis.LinkAction().Link(c.GetCreatePath(ctx))

	if dialog {
		form = form.Api(c.GetStorePath(ctx)).OnEvent(map[string]any{})
		if _type == "drawer" {
			action = (*renderers.LinkAction)(gamis.DrawerAction().Drawer(
				gamis.Drawer().Title(title).Body(form).Size(size),
			))
		} else {
			action = (*renderers.LinkAction)(gamis.DialogAction().Dialog(
				gamis.Dialog().Title(title).Body(form).Size(size),
			))
		}
	}

	action.Label(title).Icon("fa fa-add").Level("primary")
	return action
}

// 行编辑按钮
func (c *ControllerImpl[T]) RowEditButton(ctx http.Context, form *renderers.Form, dialog bool, size string, title string, _type string) *renderers.LinkAction {
	if title == "" {
		title = tools.AdminLang(ctx, "edit")
	}
	action := gamis.LinkAction().Link(c.GetEditPath(ctx))
	if dialog {
		form = form.Api(c.GetUpdatePath(ctx)).Api(c.GetEditGetDataPath(ctx)).InitApi(c.GetEditGetDataPath(ctx)).Redirect("").OnEvent(map[string]any{})
		if _type == "drawer" {
			action = (*renderers.LinkAction)(gamis.DrawerAction().Drawer(
				gamis.Drawer().Title(title).Body(form).Size(size),
			))
		} else {
			action = (*renderers.LinkAction)(gamis.DialogAction().Dialog(
				gamis.Dialog().Title(title).Body(form).Size(size),
			))
		}
	}
	action = action.Label(title).Level("link")
	return action
}

// 行详情按钮
func (c *ControllerImpl[T]) RowShowButton(ctx http.Context, form *renderers.Form, dialog bool, size string, title string, _type string) *renderers.LinkAction {
	if title == "" {
		title = tools.AdminLang(ctx, "show")
	}
	action := gamis.LinkAction().Link(c.GetShowPath(ctx))
	if dialog {
		if _type == "drawer" {
			action = (*renderers.LinkAction)(gamis.DrawerAction().Drawer(
				gamis.Drawer().Title(title).Body(form).Size(size).Actions([]any{}).CloseOnEsc("").CloseOnOutside(""),
			))
		} else {
			action = (*renderers.LinkAction)(gamis.DialogAction().Dialog(
				gamis.Dialog().Title(title).Body(form).Size(size).Actions([]any{}).CloseOnEsc("").CloseOnOutside(""),
			))
		}
	}
	action = action.Label(title).Level("link")
	return action
}

// 行删除按钮
func (c *ControllerImpl[T]) RowDeleteButton(ctx http.Context, title string) *renderers.DialogAction {
	if title == "" {
		title = tools.AdminLang(ctx, "delete")
	}
	action := gamis.DialogAction().Label(title).Level("link").ClassName("text-danger").Dialog(
		gamis.Dialog().Title(title).ClassName("py-2").Actions([]any{
			gamis.Action().ActionType("cancel").Label(tools.AdminLang(ctx, "cancel")),
			gamis.Action().ActionType("submit").Label(tools.AdminLang(ctx, "delete")).Level("danger"),
		}).Body([]any{
			gamis.Form().WrapWithPanel(false).Api(c.GetDeletePath(ctx)).Body([]any{
				gamis.Tpl().ClassName("py-2").Tpl(tools.AdminLang(ctx, "confirm_delete")),
			}),
		}),
	)
	return action
}

// 行操作按钮
func (c *ControllerImpl[T]) RowActions(ctx http.Context, form *renderers.Form, dialog any, size string) *renderers.Operation {
	// 判断 dialog 是否为切片
	dialogValue := reflect.ValueOf(dialog)
	if dialogValue.Kind() == reflect.Slice {
		// 如果是切片，遍历处理
			return gamis.Operation().Label(tools.AdminLang(ctx, "actions")).Buttons(dialog)
	} else {
	// 添加删除按钮
	return gamis.Operation().Label(tools.AdminLang(ctx, "actions")).Buttons([]any{
		c.RowShowButton(ctx, form, false, size, "", ""),
		c.RowEditButton(ctx, form, false, size, "", ""),
		c.RowDeleteButton(ctx, ""),
	})
}

}

// 基础筛选器
func (c *ControllerImpl[T]) BaseFilter() *renderers.Form {
	return gamis.Form().PanelClassName("base-filter").Title("").Actions([]any{
		gamis.Button().Label("重置").ActionType("clear-and-submit"),
		gamis.Component("submit").Label("搜索").Level("primary"),
	})
}

// 基础筛选 - 条件构造器
func (c *ControllerImpl[T]) BaseFilterConditionBuilder(ctx http.Context) *renderers.ConditionBuilderControl{
	return gamis.ConditionBuilderControl().Name("filter_condition_builder")
}

// 基础 CRUD
func (c *ControllerImpl[T]) BaseCRUD(ctx http.Context) *renderers.CRUDTable {
	crud := gamis.CRUDTable().PerPage(20).AffixHeader(false).FilterTogglable(true).FilterDefaultVisible(false).
		Api(c.GetListGetDataPath(ctx)).
		QuickSaveApi(c.GetQuickEditPath(ctx)).
		QuickSaveItemApi(c.GetQuickEditItemPath(ctx)).
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
func (c *ControllerImpl[T]) BaseHeaderToolBar() []any {
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

// 图片上传路径
func (c *ControllerImpl[T]) UploadImagePath(ctx http.Context) string {
	return tools.GetAdmin("/upload_image")
}

func (c *ControllerImpl[T]) UploadImage(ctx http.Context) http.Response {
	return c.upload(ctx, "image")
}

func (c *ControllerImpl[T]) UploadFilePath(ctx http.Context) string {
	return c.Extra.QueryPath(ctx) +"/upload_file"
}


func (c *ControllerImpl[T]) UploadFile(ctx http.Context) http.Response {
	return c.upload(ctx, "file")
}


func (c *ControllerImpl[T]) UploadRichPath(ctx http.Context) string{
	return c.Extra.QueryPath(ctx) +"/upload_rich"
}

func (c *ControllerImpl[T]) UploadRich(ctx http.Context) http.Response {
	fromWangEditor := false
	file, err := ctx.Request().File("file")
	if err != nil {
		fromWangEditor = true
		file, err = ctx.Request().File("wangeditor-uploaded-image")
		if err != nil {
			file, err = ctx.Request().File("wangeditor-uploaded-video")
			if err != nil {
				return ctx.Response().Success().Json(http.Json{
					"status":            enums.StatusFailed,
					"code":              enums.Failed,
					"msg":               tools.AdminLang(ctx, "upload_file_error"),
					"data":              []any{},
					"doNotDisplayToast": 0,
					"errno": 1,
				})
			}
		}
	}
	config := facades.Config()
	path, err := file.Disk(config.GetString("admin.upload.disk")).Store(config.GetString("admin.upload.directory.rich"))
	if err != nil {
		return c.FailMsg(ctx, tools.AdminLang(ctx, "upload_file_error"))
	}
	link := facades.Storage().Disk(config.GetString("admin.upload.disk")).Url(path)
	if fromWangEditor {
		return ctx.Response().Success().Json(http.Json{
			"status":            enums.StatusSuccess,
			"code":              enums.Success,
			"msg":               tools.AdminLang(ctx, "upload_file_success"),
			"data":              map[string]string{
				"url": link,
			},
			"doNotDisplayToast": 0,
			"errno": 0,
		})
	}
	return ctx.Response().Success().Json(http.Json{
		"status":            enums.StatusSuccess,
		"code":              enums.Success,
		"msg":               tools.AdminLang(ctx, "upload_file_success"),
		"data":              map[string]string{
			"link": link,
		},
		"doNotDisplayToast": 0,
		"link": link,
	})
}

// 上传文件
func (c *ControllerImpl[T]) upload(ctx http.Context, _type string) http.Response {
	file, err := ctx.Request().File("file")
	if err != nil {
		return c.FailMsg(ctx, tools.AdminLang(ctx, "upload_file_error"))
	}
	config := facades.Config()
	path, err := file.Disk(config.GetString("admin.upload.disk")).Store(config.GetString("admin.upload.directory" + _type))
	if err != nil {
		return c.FailMsg(ctx, tools.AdminLang(ctx, "upload_file_error"))
	}
	return c.SuccessData(ctx, map[string]string{
		"value": path,
	})
}

func (c *ControllerImpl[T]) ChunkUploadStart(ctx http.Context) http.Response {
	return c.Success(ctx)
}

// TODO: 分片上传
func (c *ControllerImpl[T]) ChunkUpload(ctx http.Context) http.Response {
	return c.Success(ctx)
}

func (c *ControllerImpl[T]) ChunkUploadFinish(ctx http.Context) http.Response {
	return c.Success(ctx)
}