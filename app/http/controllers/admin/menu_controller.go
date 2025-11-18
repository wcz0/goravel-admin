package admin

import (
    "goravel/app/services"
    "goravel/app/tools"

    "github.com/goravel/framework/contracts/http"
    "github.com/wcz0/gamis"
    "github.com/wcz0/gamis/renderers"
)

type MenuController struct {
    *ControllerImpl[*services.AdminMenuService]
}

func NewMenuController() *MenuController {
    return &MenuController{ControllerImpl: NewAdminController[*services.AdminMenuService](services.NewAdminMenuService())}
}

func (m *MenuController) Index(ctx http.Context) http.Response {
    if m.ActionOfGetData(ctx) {
        return m.ControllerImpl.Service.List(ctx)
    }
    return m.SuccessData(ctx, m.list(ctx))
}

func (m *MenuController) Show(ctx http.Context) http.Response {
    hasError, resp := m.HandleValidationErrors(ctx, map[string]string{"id": "required|number"}, nil)
    if hasError {
        return resp
    }
    return m.ControllerImpl.Service.Show(ctx)
}

func (m *MenuController) Store(ctx http.Context) http.Response {
    if m.ActionOfQuickEdit(ctx) {
        return m.ControllerImpl.Service.QuickEdit(ctx)
    }
    if m.ActionOfQuickEditItem(ctx) {
        return m.ControllerImpl.Service.QuickEditItem(ctx)
    }
    hasError, resp := m.HandleValidationErrors(ctx, map[string]string{
        "title":       "required|string|max:50",
        "icon":        "string|max:50",
        "url":         "required|string|max:255",
        "parent_id":   "number",
        "custom_order": "number",
        "visible":     "in:0,1",
        "is_home":     "in:0,1",
        "component":   "string",
        "keep_alive":  "number",
        "iframe_url":  "string",
        "is_full":     "number",
    }, nil)
    if hasError {
        return resp
    }
    return m.ControllerImpl.Service.Store(ctx)
}

func (m *MenuController) Update(ctx http.Context) http.Response {
    hasError, resp := m.HandleValidationErrors(ctx, map[string]string{
        "id":          "required|number",
        "title":       "required|string|max:50",
        "icon":        "string|max:50",
        "url":         "required|string|max:255",
        "parent_id":   "number",
        "custom_order": "number",
        "visible":     "in:0,1",
        "is_home":     "in:0,1",
        "component":   "string",
        "keep_alive":  "number",
        "iframe_url":  "string",
        "is_full":     "number",
    }, nil)
    if hasError {
        return resp
    }
    return m.ControllerImpl.Service.Update(ctx)
}

func (m *MenuController) Destroy(ctx http.Context) http.Response {
    hasError, resp := m.HandleValidationErrors(ctx, map[string]string{"id": "required|number"}, nil)
    if hasError {
        return resp
    }
    return m.ControllerImpl.Service.Destroy(ctx)
}

func (m *MenuController) SaveOrder(ctx http.Context) http.Response {
    return m.ControllerImpl.Service.Reorder(ctx)
}

func (m *MenuController) list(ctx http.Context) *renderers.Page {
    refreshEvent := map[string]any{"actions": []any{map[string]any{"actionType": "custom", "script": "window.$owl.refreshRoutes()"}}}
    header := []any{m.CreateButton(ctx, m.form(ctx), true, "md", "", "")}
    header = append(header, m.BaseHeaderToolBar()...)
    crud := m.BaseCRUD(ctx).
        Set("id", "menu-crud").
        PerPage(999).
        Set("draggable", true).
        Set("saveOrderApi", map[string]any{"url": m.Extra.QueryPath(ctx) + "/save_order", "data": map[string]any{"ids": "${ids}"}}).
        LoadDataOnce(true).
        Set("syncLocation", false).
        HeaderToolbar(header).
        FilterTogglable(true).
        Filter(m.BaseFilter().Body([]any{
            gamis.TextControl().Name("title").Label(tools.AdminLang(ctx, "admin_menu.title")).Size("md").Clearable(true).Placeholder(tools.AdminLang(ctx, "admin_menu.title")),
            gamis.TextControl().Name("url").Label(tools.AdminLang(ctx, "admin_menu.url")).Size("md").Clearable(true).Placeholder(tools.AdminLang(ctx, "admin_menu.url")),
        })).
        FooterToolbar([]string{"statistics"}).
        BulkActions([]any{m.BulkDeleteButton(ctx)}).
        Columns([]any{
            gamis.TableColumn().Name("id").Label("ID").Sortable(true),
            gamis.TableColumn().Name("title").Label(tools.AdminLang(ctx, "admin_menu.title")),
            gamis.TableColumn().Name("icon").Label(tools.AdminLang(ctx, "admin_menu.icon")),
            gamis.TableColumn().Name("url").Label(tools.AdminLang(ctx, "admin_menu.url")),
            gamis.TableColumn().Name("custom_order").Label(tools.AdminLang(ctx, "order")).QuickEdit(gamis.NumberControl().Min(0).Set("saveImmediately", true)).Sortable(true),
            gamis.TableColumn().Name("visible").Label(tools.AdminLang(ctx, "admin_menu.visible")).QuickEdit(map[string]any{"type": "checkbox", "mode": "inline", "saveImmediately": true}),
            gamis.TableColumn().Name("is_home").Label(tools.AdminLang(ctx, "admin_menu.is_home")).QuickEdit(map[string]any{"type": "checkbox", "mode": "inline", "saveImmediately": true}),
            m.RowActions(ctx, m.form(ctx), []any{
                m.RowEditButton(ctx, m.form(ctx), true, "md", "", ""),
                gamis.AjaxAction().Label(tools.AdminLang(ctx, "delete")).Level("link").ClassName("text-danger").Api(map[string]any{"url": m.Extra.QueryPath(ctx) + "/${id}", "method": "delete"}).Set("confirmText", tools.AdminLang(ctx, "confirm_delete")).Set("reload", "menu-crud"),
            }, "md"),
        }).
        OnEvent(map[string]any{"quickSaveItemSucc": refreshEvent, "saveOrderSucc": refreshEvent})
    return m.BaseList(crud)
}

func (m *MenuController) form(ctx http.Context) *renderers.Form {
    return m.BaseForm(ctx, false).Mode("normal").Body([]any{
        gamis.TextControl().Name("title").Label(tools.AdminLang(ctx, "admin_menu.title")).Required(true),
        gamis.TextControl().Name("icon").Label(tools.AdminLang(ctx, "admin_menu.icon")),
        gamis.TreeSelectControl().Name("parent_id").Label(tools.AdminLang(ctx, "admin_menu.parent_id")).Set("id", "parent_select").LabelField("title").ValueField("id").Set("showIcon", false).Source(m.Extra.QueryPath(ctx) + "?_action=getData"),
        gamis.GroupControl().Body([]any{
            gamis.NumberControl().Name("custom_order").Label(tools.AdminLang(ctx, "order")).Required(true).Set("displayMode", "enhance").Description(tools.AdminLang(ctx, "admin.order_asc")).Min(0).Value(0),
            gamis.SwitchControl().Name("visible").Label(tools.AdminLang(ctx, "admin_menu.visible")).OnText(tools.AdminLang(ctx, "admin.yes")).OffText(tools.AdminLang(ctx, "admin.no")).Value(1),
        }),
        gamis.TextControl().Name("url").Label(tools.AdminLang(ctx, "admin_menu.url")).Required(true).ValidateOnChange(true).Set("validations", map[string]any{"matchRegexp": "^(http(s)?\\:\\/\\/)?(\\/)+"}).Set("validationErrors", map[string]any{"matchRegexp": tools.AdminLang(ctx, "admin.need_start_with_slash")}).Placeholder("eg: /admin_menus"),
        gamis.GroupControl().Body([]any{
            gamis.TextControl().Name("iframe_url").Label("IframeUrl").ValidateOnChange(true).Set("validations", map[string]any{"matchRegexp": "^(http(s)?\\:\\/\\/)?(\\/)+"}).Set("validationErrors", map[string]any{"matchRegexp": tools.AdminLang(ctx, "admin.need_start_with_slash")}).Placeholder("eg: https://www.qq.com"),
        }),
        gamis.FieldSetControl().Title(tools.AdminLang(ctx, "admin.more")).Collapsable(true).Collapsed(true).Body([]any{
            gamis.TextControl().Name("component").Label(tools.AdminLang(ctx, "admin_menu.component")).Description(tools.AdminLang(ctx, "admin_menu.component_desc")).Value("amis"),
            gamis.SwitchControl().Name("keep_alive").Label(tools.AdminLang(ctx, "admin_menu.keep_alive")).OnText(tools.AdminLang(ctx, "admin.yes")).OffText(tools.AdminLang(ctx, "admin.no")).Value(0),
            gamis.SwitchControl().Name("is_home").Label(tools.AdminLang(ctx, "admin_menu.is_home")).OnText(tools.AdminLang(ctx, "admin.yes")).OffText(tools.AdminLang(ctx, "admin.no")).Description(tools.AdminLang(ctx, "admin_menu.is_home_description")).Value(0),
            gamis.SwitchControl().Name("is_full").Label(tools.AdminLang(ctx, "admin_menu.is_full")).OnText(tools.AdminLang(ctx, "admin.yes")).OffText(tools.AdminLang(ctx, "admin.no")).Description(tools.AdminLang(ctx, "admin_menu.is_full_description")).Value(0),
        }),
    }).OnEvent(map[string]any{"inited": map[string]any{"actions": []any{map[string]any{"actionType": "setValue", "componentId": "parent_select", "args": map[string]any{"value": "${responseData.parent_id || \"\"}"}}}}})
}