package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/wcz0/gamis/v6"
	"github.com/wcz0/gamis/v6/renderers"

	"github.com/wcz0/goravel-admin/enums"
	"github.com/wcz0/goravel-admin/response"
	"github.com/wcz0/goravel-admin/tools"
)

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

// HandleValidationErrors 统一处理验证错误
func (c *ControllerImpl[T]) HandleValidationErrors(ctx http.Context, rules map[string]string, messages map[string]string) (bool, http.Response) {
	return handleValidationErrors(ctx, c.ErrorHandler, rules, messages)
}

// parseInt 辅助函数
func parseInt(s string) int {
	if i, err := fmt.Sscanf(s, "%d", new(int)); err == nil && i == 1 {
		var result int
		fmt.Sscanf(s, "%d", &result)
		return result
	}
	return 0
}

// 获取基础url
func (e *Extra) QueryPath(ctx http.Context) string {
	path := ctx.Request().Path()
	if strings.HasPrefix(path, "/"+e.AdminPrefix) {
		path = strings.TrimPrefix(path, "/"+e.AdminPrefix)
	} else {
		path = strings.TrimPrefix(path, e.AdminPrefix)
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
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
	return c.Extra.QueryPath(ctx) + "?_action=getData"
}

// 获取导出数据
func (c *ControllerImpl[T]) GetExportPath(ctx http.Context) string {
	return c.Extra.QueryPath(ctx) + "?_action=export"
}

// 删除路径 ?
func (c *ControllerImpl[T]) GetDeletePath(ctx http.Context, primaryKey ...string) string {
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	return "delete:" + c.Extra.QueryPath(ctx) + "/${" + key + "}"
}

// 批量删除 ?
func (c *ControllerImpl[T]) GetBulkDeletePath(ctx http.Context) string {
	return "delete:" + c.Extra.QueryPath(ctx) + "/${ids}"
}

// 获取编辑页面路径 ?
func (c *ControllerImpl[T]) GetEditPath(ctx http.Context, primaryKey ...string) string {
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	return "/" + strings.Trim(c.Extra.QueryPath(ctx), "/") + "/${" + key + "}/edit"
}

// 获取编辑数据 ?
func (c *ControllerImpl[T]) GetEditGetDataPath(ctx http.Context, primaryKey ...string) string {
	if api, err := c.GetShowGetDataPath(ctx); err == nil {
		return api
	}
	// fallback: old behavior
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	path := c.Extra.QueryPath(ctx)
	paths := strings.Split(path, "/")
	last := paths[len(paths)-1]
	if last == "edit" {
		path = "/${" + key + "}/edit"
	}
	return path + "?_action=getData"
}

// 详情页面 ?
func (c *ControllerImpl[T]) GetShowPath(ctx http.Context, primaryKey ...string) string {
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	return "/" + strings.Trim(c.Extra.QueryPath(ctx), "/") + "/${" + key + "}"
}

// 编辑保存 ?
func (c *ControllerImpl[T]) GetUpdatePath(ctx http.Context, primaryKey ...string) string {
	key := "id"
	if len(primaryKey) > 0 {
		key = primaryKey[0]
	}
	path := c.Extra.QueryPath(ctx)
	path = strings.TrimSuffix(path, "/edit")
	return "put:" + path + "/${" + key + "}"
}

// 获取快速编辑数据
func (c *ControllerImpl[T]) GetQuickEditPath(ctx http.Context) string {
	return ctx.Request().FullUrl() + "?_action=quickEdit"
}

// 获取快速编辑项目数据
func (c *ControllerImpl[T]) GetQuickEditItemPath(ctx http.Context) string {
	return ctx.Request().FullUrl() + "?_action=quickEditItem"
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
	return "post:" + c.Extra.QueryPath(ctx)
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
				gamis.Form().WrapWithPanel(false).
					Api(c.GetBulkDeletePath(ctx)).
					Set("data", map[string]any{"ids": "${selectedItems|map:id|join:,}"}).
					OnEvent(map[string]any{"submitSucc": []any{map[string]any{"actionType": "reload"}}}).
					Body([]any{
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
		onEvent := map[string]any{
			"submitSucc": []any{
				map[string]any{"actionType": "reload"},
			},
		}
		form = form.Api(c.GetStorePath(ctx)).OnEvent(onEvent)
		if _type == "drawer" {
			action = (*renderers.LinkAction)(gamis.DrawerAction().Drawer(
				gamis.Drawer().Title(title).Body(form).Size(size),
			))
			action.Set("actionType", "drawer").Link("")
			action.Set("onEvent", map[string]any{"submitSucc": []any{map[string]any{"actionType": "closeDrawer"}}})
		} else {
			action = (*renderers.LinkAction)(gamis.DialogAction().Dialog(
				gamis.Dialog().Title(title).Body(form).Size(size).Actions([]any{
					gamis.Action().ActionType("cancel").Label(tools.AdminLang(ctx, "cancel")),
					gamis.Action().ActionType("submit").Label(tools.AdminLang(ctx, "create")).Level("primary"),
				}),
			))
			action.Set("actionType", "dialog").Link("")
			action.Set("onEvent", map[string]any{"submitSucc": []any{map[string]any{"actionType": "closeDialog"}}})
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
		onEvent := map[string]any{
			"submitSucc": []any{
				map[string]any{"actionType": "reload"},
			},
		}
		if _type == "drawer" {
			onEvent["submitSucc"] = append(onEvent["submitSucc"].([]any), map[string]any{"actionType": "closeDrawer"})
		} else {
			onEvent["submitSucc"] = append(onEvent["submitSucc"].([]any), map[string]any{"actionType": "closeDialog"})
		}
		form = form.InitApi(c.GetEditGetDataPath(ctx)).Api(c.GetUpdatePath(ctx)).Redirect("").OnEvent(onEvent)
		if _type == "drawer" {
			action = (*renderers.LinkAction)(gamis.DrawerAction().Drawer(
				gamis.Drawer().Title(title).Body(form).Size(size),
			))
			action.Set("actionType", "drawer").Link("")
			action.Set("onEvent", map[string]any{"submitSucc": []any{map[string]any{"actionType": "closeDrawer"}}})
		} else {
			action = (*renderers.LinkAction)(gamis.DialogAction().Dialog(
				gamis.Dialog().Title(title).Body(form).Size(size).Actions([]any{
					gamis.Action().ActionType("cancel").Label(tools.AdminLang(ctx, "cancel")),
					gamis.Action().ActionType("submit").Label(tools.AdminLang(ctx, "save")).Level("primary"),
				}),
			))
			action.Set("actionType", "dialog").Link("")
			action.Set("onEvent", map[string]any{"submitSucc": []any{map[string]any{"actionType": "closeDialog"}}})
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
			gamis.Form().WrapWithPanel(false).
				Api(c.GetDeletePath(ctx)).
				Set("data", map[string]any{"id": "${id}"}).
				OnEvent(map[string]any{"submitSucc": []any{map[string]any{"actionType": "reload"}, map[string]any{"actionType": "closeDialog"}}}).
				Body([]any{
					gamis.Tpl().ClassName("py-2").Tpl(tools.AdminLang(ctx, "confirm_delete")),
				}),
		}),
	)
	action.Set("onEvent", map[string]any{"submitSucc": []any{map[string]any{"actionType": "reload"}, map[string]any{"actionType": "closeDialog"}}})
	action.Set("actionType", "dialog")
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
func (c *ControllerImpl[T]) BaseFilterConditionBuilder(ctx http.Context) *renderers.ConditionBuilderControl {
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
//
// 结构与 owl-admin 保持一致：一个带 loading 遮罩的下拉按钮，点「导出」时
// 请求当前列表接口（_action=export）拿到签名后的下载信息，再打开下载地址。
func (c *ControllerImpl[T]) ExportAction(ctx http.Context, disableSelectedItem bool) *renderers.Service {
	downloadPath := tools.GetAdmin("/_download_export")
	exportPath := c.GetExportPath(ctx)

	// doAction 是三个按钮共用的收尾动作：显示 loading → 请求导出 → 关 loading → 打开下载地址
	doAction := fmt.Sprintf(`doAction([{actionType:'setValue',componentId:'export-action',args:{value:{showExportLoading:true}}},{actionType:'ajax',args:{api:{url:url.toString(),method:'get'}}},{actionType:'setValue',componentId:'export-action',args:{value:{showExportLoading:false}}},{actionType:'custom',expression:'${event.data.responseResult.responseStatus===0}',script:'let res=event.data.responseResult.responseData;window.open(\'%s?path=\'+encodeURIComponent(res.path)+\'&expires=\'+encodeURIComponent(res.expires)+\'&signature=\'+encodeURIComponent(res.signature))'}])`, downloadPath)

	event := func(script string) map[string]any {
		return map[string]any{
			"click": map[string]any{
				"actions": []any{
					map[string]any{"actionType": "custom", "script": script},
				},
			},
		}
	}

	noPageData := tools.AdminLang(ctx, "export.page_no_data")
	noSelectedData := tools.AdminLang(ctx, "export.selected_rows_no_data")

	allScript := fmt.Sprintf(`let data=event.data;let params=Object.keys(data).filter(key=>key!=='page' && key!=='__super').reduce((obj,key)=>{obj[key]=data[key];return obj;},{});let url=new URL('%s',window.location.origin);Object.keys(params).forEach(key=>url.searchParams.append(key,(typeof params[key] == 'string' ? params[key] : JSON.stringify(params[key]))));%s`, exportPath, doAction)
	pageScript := fmt.Sprintf(`let ids=event.data.items.map(item=>item.id);if(ids.length===0){return doAction({actionType:'toast',args:{msgType:'warning',msg:'%s'}})};let url=new URL('%s',window.location.origin);url.searchParams.append('_ids',ids.join(','));%s`, noPageData, exportPath, doAction)
	selectedScript := fmt.Sprintf(`let ids=event.data.selectedItems.map(item=>item.id);if(ids.length===0){return doAction({actionType:'toast',args:{msgType:'warning',msg:'%s'}})};let url=new URL('%s',window.location.origin);url.searchParams.append('_ids',ids.join(','));%s`, noSelectedData, exportPath, doAction)

	buttons := []any{
		gamis.Button().Label(tools.AdminLang(ctx, "export.all")).OnEvent(event(allScript)),
		gamis.Button().Label(tools.AdminLang(ctx, "export.page")).OnEvent(event(pageScript)),
	}
	if !disableSelectedItem {
		buttons = append(buttons, gamis.Button().Label(tools.AdminLang(ctx, "export.selected_rows")).OnEvent(event(selectedScript)))
	}

	return gamis.Service().
		Id("export-action").
		Set("align", "right").
		Set("data", map[string]any{
			"showExportLoading": false,
		}).
		Body([]any{
			gamis.Spinner().Set("showOn", "${showExportLoading}").Set("overlay", true).Body([]any{
				gamis.DropdownButton().
					Label(tools.AdminLang(ctx, "export.title")).
					Set("icon", "fa-solid fa-download").
					Buttons(buttons).
					Set("closeOnClick", true).
					Align("right"),
			}),
		})
}

// ExportRows 是控制器的通用导出入口，子类提供表头与数据行即可。
func (c *ControllerImpl[T]) ExportRows(ctx http.Context, sheet string, headers []string, rows [][]any) http.Response {
	disk := facades.Config().GetString("admin.upload.disk")
	result, err := tools.BuildExport(disk, sheet, headers, rows)
	if err != nil {
		return c.FailMsg(ctx, err.Error())
	}

	return c.SuccessData(ctx, map[string]string{
		"path":      result.Path,
		"expires":   result.Expires,
		"signature": result.Signature,
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
	return c.Extra.QueryPath(ctx) + "/upload_file"
}

func (c *ControllerImpl[T]) UploadFile(ctx http.Context) http.Response {
	return c.upload(ctx, "file")
}

func (c *ControllerImpl[T]) UploadRichPath(ctx http.Context) string {
	return c.Extra.QueryPath(ctx) + "/upload_rich"
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
					"errno":             1,
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
			"status": enums.StatusSuccess,
			"code":   enums.Success,
			"msg":    tools.AdminLang(ctx, "upload_file_success"),
			"data": map[string]string{
				"url": link,
			},
			"doNotDisplayToast": 0,
			"errno":             0,
		})
	}
	return ctx.Response().Success().Json(http.Json{
		"status": enums.StatusSuccess,
		"code":   enums.Success,
		"msg":    tools.AdminLang(ctx, "upload_file_success"),
		"data": map[string]string{
			"link": link,
		},
		"doNotDisplayToast": 0,
		"link":              link,
	})
}

// 上传文件
func (c *ControllerImpl[T]) upload(ctx http.Context, _type string) http.Response {
	file, err := ctx.Request().File("file")
	if err != nil {
		return c.FailMsg(ctx, tools.AdminLang(ctx, "upload_file_error"))
	}
	config := facades.Config()
	// 目录配置形如 admin.upload.directory.image，需要按类型拼出完整键名
	directory := config.GetString("admin.upload.directory." + _type)
	path, err := file.Disk(config.GetString("admin.upload.disk")).Store(directory)
	if err != nil {
		return c.FailMsg(ctx, tools.AdminLang(ctx, "upload_file_error"))
	}
	url := facades.Storage().Disk(config.GetString("admin.upload.disk")).Url(path)
	return c.SuccessData(ctx, map[string]string{
		"value": url,
	})
}

// chunkTTL 是分块上传会话的存活时间。
const chunkTTL = 10 * time.Minute

// chunkCacheKey 是分块上传会话的缓存键前缀。
const chunkCacheKey = "admin_chunk_"

// chunkDir 是分片临时目录。
const chunkDir = "chunk"

// ChunkUploadStart 开启一次分块上传，返回会话 id。
func (c *ControllerImpl[T]) ChunkUploadStart(ctx http.Context) http.Response {
	uploadID := uuid.New().String()
	if err := facades.Cache().Put(chunkCacheKey+uploadID, true, chunkTTL); err != nil {
		return c.FailMsg(ctx, tools.AdminLang(ctx, "upload_file_error"))
	}

	disk := c.uploadDisk()
	if err := facades.Storage().Disk(disk).MakeDirectory(chunkDir + "/" + uploadID); err != nil {
		return c.FailMsg(ctx, tools.AdminLang(ctx, "upload_file_error"))
	}

	return c.SuccessData(ctx, map[string]string{"uploadId": uploadID})
}

// ChunkUpload 接收单个分片，返回该分片的 md5 供合并时校验。
func (c *ControllerImpl[T]) ChunkUpload(ctx http.Context) http.Response {
	uploadID := ctx.Request().Input("uploadId")
	partNumber := ctx.Request().InputInt("partNumber", -1)
	if !validChunkID(uploadID) || partNumber < 0 {
		return c.FailMsg(ctx, "分片参数不合法")
	}
	if facades.Cache().Get(chunkCacheKey+uploadID) == nil {
		return c.FailMsg(ctx, "上传会话已过期")
	}

	file, err := ctx.Request().File("file")
	if err != nil {
		return c.FailMsg(ctx, tools.AdminLang(ctx, "upload_file_error"))
	}

	disk := c.uploadDisk()
	storage := facades.Storage().Disk(disk)
	path := fmt.Sprintf("%s/%s/%d", chunkDir, uploadID, partNumber)

	// 直接用源文件字节落盘，避免 StoreAs 自动补扩展名导致分片名不可预期
	content, err := os.ReadFile(file.File())
	if err != nil {
		return c.FailMsg(ctx, tools.AdminLang(ctx, "upload_file_error"))
	}
	if err := storage.Put(path, string(content)); err != nil {
		return c.FailMsg(ctx, tools.AdminLang(ctx, "upload_file_error"))
	}

	sum := md5.Sum(content)

	return c.SuccessData(ctx, map[string]string{"eTag": hex.EncodeToString(sum[:])})
}

// ChunkUploadFinish 合并分片并落盘。
//
// partList 按前端给的顺序合并，每片都会与之前返回的 eTag 对比；任一不匹配
// 直接失败并清理临时目录，避免半成品污染上传目录。
func (c *ControllerImpl[T]) ChunkUploadFinish(ctx http.Context) http.Response {
	uploadID := ctx.Request().Input("uploadId")
	if !validChunkID(uploadID) || facades.Cache().Get(chunkCacheKey+uploadID) == nil {
		return c.FailMsg(ctx, "上传会话已过期")
	}

	filename := ctx.Request().Input("filename")
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	if ext == "" || !safeExtension(ext) {
		return c.FailMsg(ctx, "文件类型不允许")
	}

	target := c.resolveChunkTarget(ctx.Request().Input("t"))
	disk := c.uploadDisk()
	storage := facades.Storage().Disk(disk)

	parts, ok := ctx.Request().All()["partList"].([]any)
	if !ok || len(parts) == 0 {
		return c.FailMsg(ctx, "分片列表不可为空")
	}

	var builder strings.Builder
	chunkPaths := make([]string, 0, len(parts))
	for _, raw := range parts {
		part, ok := raw.(map[string]any)
		if !ok {
			return c.failChunk(ctx, disk, uploadID, "分片参数不合法")
		}
		partNumber := cast.ToInt(part["partNumber"])
		eTag := cast.ToString(part["eTag"])
		path := fmt.Sprintf("%s/%s/%d", chunkDir, uploadID, partNumber)

		if !storage.Exists(path) {
			return c.failChunk(ctx, disk, uploadID, fmt.Sprintf("分片 %d 不存在", partNumber))
		}
		content, err := storage.Get(path)
		if err != nil {
			return c.failChunk(ctx, disk, uploadID, tools.AdminLang(ctx, "upload_file_error"))
		}
		sum := md5.Sum([]byte(content))
		if !strings.EqualFold(hex.EncodeToString(sum[:]), eTag) {
			return c.failChunk(ctx, disk, uploadID, fmt.Sprintf("分片 %d 校验失败", partNumber))
		}
		builder.WriteString(content)
		chunkPaths = append(chunkPaths, path)
	}

	finalPath := fmt.Sprintf("%s/%s.%s", target, uploadID, ext)
	if err := storage.Put(finalPath, builder.String()); err != nil {
		return c.failChunk(ctx, disk, uploadID, tools.AdminLang(ctx, "upload_file_error"))
	}

	// 合并成功后清理会话与临时分片
	c.cleanChunk(disk, uploadID, chunkPaths)
	facades.Cache().Forget(chunkCacheKey + uploadID)

	return c.SuccessMsgData(ctx, tools.AdminLang(ctx, "upload_file_success"), map[string]string{
		"value": storage.Url(finalPath),
		"path":  finalPath,
	})
}

// uploadDisk 返回上传使用的磁盘名。
func (c *ControllerImpl[T]) uploadDisk() string {
	disk := facades.Config().GetString("admin.upload.disk")
	if disk == "" {
		return "public"
	}

	return disk
}

// resolveChunkTarget 校验并归一化合并目标目录。
//
// 只允许配置里声明过的上传目录，防止通过 t 参数把文件写到任意位置。
func (c *ControllerImpl[T]) resolveChunkTarget(raw string) string {
	allowed := map[string]bool{"uploads": true}
	directories, ok := facades.Config().Get("admin.upload.directory").(map[string]any)
	if ok {
		for _, value := range directories {
			if dir := cast.ToString(value); dir != "" {
				allowed[dir] = true
			}
		}
	}

	target := strings.Trim(strings.TrimSpace(raw), "/")
	if target == "" || !allowed[target] {
		return "uploads"
	}

	return target
}

// failChunk 清理失败的会话并返回错误。
func (c *ControllerImpl[T]) failChunk(ctx http.Context, disk string, uploadID string, msg string) http.Response {
	c.cleanChunk(disk, uploadID, nil)
	facades.Cache().Forget(chunkCacheKey + uploadID)

	return c.FailMsg(ctx, msg)
}

// cleanChunk 删除分片文件与临时目录。
//
// Storage 的 Delete 只在给定路径是文件时生效，目录需要单独移除，
// 因此这里先删已知分片，再删目录本身（本地磁盘时用 os.Remove 兜底）。
func (c *ControllerImpl[T]) cleanChunk(disk string, uploadID string, paths []string) {
	storage := facades.Storage().Disk(disk)
	dir := chunkDir + "/" + uploadID

	for _, path := range paths {
		_ = storage.Delete(path)
	}
	_ = storage.Delete(dir)

	// 本地磁盘才有真实路径，远程驱动删不到就留给宿主的生命周期策略
	if local := storage.Path(dir); local != "" {
		_ = os.RemoveAll(local)
	}
}

// validChunkID 校验会话 id，挡掉路径穿越。
func validChunkID(uploadID string) bool {
	if uploadID == "" || len(uploadID) > 64 {
		return false
	}
	if strings.ContainsAny(uploadID, "../\\\x00") {
		return false
	}
	_, err := uuid.Parse(uploadID)

	return err == nil
}

// safeExtension 挡掉脚本类后缀，避免上传目录被当成可执行入口。
func safeExtension(ext string) bool {
	switch ext {
	case "php", "php3", "php4", "php5", "php7", "php8", "phtml", "phar",
		"sh", "py", "asp", "aspx", "jsp", "env", "exe", "bat", "cmd", "vbs", "htaccess":
		return false
	}

	return true
}
