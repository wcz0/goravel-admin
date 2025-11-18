package admin

import (
	"goravel/app/models/admin"
	"goravel/app/services"
	"goravel/app/tools"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/wcz0/gamis"
	"github.com/wcz0/gamis/renderers"
)

type RoleController struct {
	*ControllerImpl[*services.AdminRoleService]
}

func NewRoleController() *RoleController {
    return &RoleController{
        ControllerImpl: NewAdminController[*services.AdminRoleService](services.NewAdminRoleService()),
    }
}

func (r *RoleController) Index(ctx http.Context) http.Response {
	if r.ActionOfGetData(ctx) {
		return r.ControllerImpl.Service.List(ctx)
	}
	return r.SuccessData(ctx, r.list(ctx))
}

func (r *RoleController) list(ctx http.Context) *renderers.Page {
	HeaderToolbar := []any{
		r.CreateButton(ctx, r.Form(ctx), true, "md", "", ""),
	}
	HeaderToolbar = append(HeaderToolbar, r.BaseHeaderToolBar()...)

    crud := r.BaseCRUD(ctx).
        Set("id", "role-crud").
        HeaderToolbar(HeaderToolbar).
        FilterTogglable(true).
		Filter(
			r.BaseFilter().Body([]any{
				gamis.TextControl().
					Name("name").
					Label(tools.AdminLang(ctx, "admin_role.name")).
					Size("md").
					Clearable(true).
					Placeholder(tools.AdminLang(ctx, "admin_role.name")),
				gamis.TextControl().
					Name("slug").
					Label(tools.AdminLang(ctx, "admin_role.slug")).
					Size("md").
					Clearable(true).
					Placeholder(tools.AdminLang(ctx, "admin_role.slug")),
				gamis.DateRangeControl().
					Name("created_at").
					Label(tools.AdminLang(ctx, "created_at")).
					Format("YYYY-MM-DD").
					Placeholder(tools.AdminLang(ctx, "created_at")),
			}),
		).
		Columns([]any{
			gamis.TableColumn().Name("id").Label("ID").Sortable(true),
			gamis.TableColumn().Name("name").Label(tools.AdminLang(ctx, "admin_role.name")),
			gamis.TableColumn().Name("slug").Label(tools.AdminLang(ctx, "admin_role.slug")).Type("tag"),
			gamis.TableColumn().
				Name("created_at").
				Label(tools.AdminLang(ctx, "created_at")).
				Type("datetime").
				Sortable(true),
			gamis.TableColumn().
				Name("updated_at").
				Label(tools.AdminLang(ctx, "updated_at")).
				Type("datetime").
				Sortable(true),
			r.RowActions(ctx, r.Form(ctx), []any{
            r.setPermissionButton(ctx, r.setPermission(ctx)),
            r.RowEditButton(ctx, r.Form(ctx), true, "md", "", "").HiddenOn("${slug == 'administrator'}"),
            gamis.AjaxAction().
                Label(tools.AdminLang(ctx, "delete")).
                Level("link").
                ClassName("text-danger").
                Api(map[string]any{
                    "url":  r.Extra.QueryPath(ctx) + "/${id}",
                    "method": "delete",
                }).
                Set("reload", "role-crud").
                Set("confirmText", tools.AdminLang(ctx, "confirm_delete")).
                HiddenOn("${slug == 'administrator'}"),
        }, "md"),
    })

	return r.BaseList(crud).Css(map[string]any{
		".tree-full": map[string]any{
			"overflow": "hidden !important",
		},
		".cxd-TreeControl > .cxd-Tree": map[string]any{
			"height":     "100% !important",
			"max-height": "100% !important",
		},
	})
}

func (r *RoleController) Show(ctx http.Context) http.Response {
	return r.ControllerImpl.Service.Show(ctx)
}

func (r *RoleController) Store(ctx http.Context) http.Response {
    hasError, response := r.HandleValidationErrors(ctx, map[string]string{
        "name": "required|string|max_len:255",
        "slug": "required|string|max_len:255",
    }, map[string]string{
        "name.required": "角色名称不能为空",
        "name.max_len": "角色名称最多 255 位",
        "slug.required": "角色标识不能为空",
        "slug.max_len": "角色标识最多 255 位",
    })
    if hasError {
        return response
    }
    return r.ControllerImpl.Service.Store(ctx)
}

func (r *RoleController) Update(ctx http.Context) http.Response {
	// 验证输入
	validation, err := ctx.Request().Validate(map[string]string{
		"id":   "required|number",
		"name": "required|string|max:255",
		"slug": "required|string|max:255",
	})
	if err != nil {
		return r.FailMsg(ctx, err.Error())
	}
	if validation.Fails() {
		return r.FailMsg(ctx, validation.Errors().All())
	}

	// 检查角色是否存在
	id := ctx.Request().Route("id")
	var role admin.AdminRole
	if err := facades.Orm().Query().Find(&role, id); err != nil {
		return r.FailMsg(ctx, "角色不存在")
	}

	return r.ControllerImpl.Service.Update(ctx)
}

func (r *RoleController) Destroy(ctx http.Context) http.Response {
	// 验证输入
	validation, err := ctx.Request().Validate(map[string]string{
		"id": "required|number",
	})
	if err != nil {
		return r.FailMsg(ctx, err.Error())
	}
	if validation.Fails() {
		return r.FailMsg(ctx, validation.Errors().All())
	}

	return r.ControllerImpl.Service.Destroy(ctx)
}

// AssignPermissions 为角色分配权限
func (r *RoleController) AssignPermissions(ctx http.Context) http.Response {
	// 验证输入
	validation, err := ctx.Request().Validate(map[string]string{
		"role_id":      "required|number",
		"permission_ids": "required|array",
	})
	if err != nil {
		return r.FailMsg(ctx, err.Error())
	}
	if validation.Fails() {
		return r.FailMsg(ctx, validation.Errors().All())
	}

	roleId := ctx.Request().Input("role_id", "")
	var role admin.AdminRole
	if err := facades.Orm().Query().Find(&role, roleId); err != nil {
		return r.FailMsg(ctx, "角色不存在")
	}

	permissionIds := ctx.Request().Input("permission_ids", "")
	if permissionIds == "" {
		return r.FailMsg(ctx, "权限ID列表不能为空")
	}
	var permissions []*admin.AdminPermission
	if err := facades.Orm().Query().Find(&permissions, permissionIds); err != nil {
		return r.FailMsg(ctx, "权限不存在")
	}

	// 同步权限
	if err := role.SyncPermissions(permissions); err != nil {
		return r.FailMsg(ctx, "权限分配失败: "+err.Error())
	}

	return r.SuccessMsg(ctx, "权限分配成功")
}

// GetPermissions 获取角色权限
func (r *RoleController) GetPermissions(ctx http.Context) http.Response {
	roleId := ctx.Request().Input("role_id", "")
	if roleId == "" {
		return r.FailMsg(ctx, "角色ID不能为空")
	}

	var role admin.AdminRole
	if err := facades.Orm().Query().Find(&role, roleId); err != nil {
		return r.FailMsg(ctx, "角色不存在")
	}

	permissions := role.AllPermissions()
	return r.SuccessData(ctx, permissions)
}

// SetPermission 权限设置页面
func (r *RoleController) setPermission(ctx http.Context) *renderers.Drawer {
    var permissions []admin.AdminPermission
    _ = facades.Orm().Query().Order("parent_id asc, custom_order asc, id asc").Get(&permissions)
    // 构建树并转换为 TreeControl 选项
    // 转为指针切片以复用构建函数
    permPtrs := make([]*admin.AdminPermission, 0, len(permissions))
    for i := range permissions {
        permPtrs = append(permPtrs, &permissions[i])
    }
    permissionTree := buildPermissionTree(permPtrs)
    treeOptions := convertPermissionTreeToOptions(permissionTree)

	return gamis.Drawer().
		Title(tools.AdminLang(ctx, "admin_role.set_permissions")).
		Resizable(true).
		CloseOnOutside(true).
		CloseOnEsc(true).
		Body([]any{
			gamis.Form().
				Api(tools.GetAdmin("system/admin_roles/save_permissions")).
				InitApi(tools.GetAdmin("system/admin_roles/${id}")).
				Mode("normal").
				Data(map[string]any{"role_id": "${id}"}).
				OnEvent(map[string]any{"submitSucc": []any{map[string]any{"actionType": "reload"}, map[string]any{"actionType": "closeDrawer"}}}).
				Body([]any{
					gamis.TreeControl().
						Name("permissions").
						Label("").
						Multiple(true).
						HeightAuto(true).
						Options(treeOptions).
						Searchable(true).
						Cascade(true).
						JoinValues(false).
						ExtractValue(true).
						Size("full").
						ClassName("h-full b-none").
						InputClassName("h-full tree-full").
						LabelField("name").
						ValueField("id"),
				}),
		})
}

// setPermissionButton 设置权限按钮
func (r *RoleController) setPermissionButton(ctx http.Context, drawer *renderers.Drawer) *renderers.LinkAction {
    action := (*renderers.LinkAction)(gamis.DrawerAction().
        Label(tools.AdminLang(ctx, "admin_role.set_permissions")).
        Level("link").
        Drawer(drawer))
    action.Set("actionType", "drawer").Link("")
    return action
}

// Form 表单页面
func (r *RoleController) Form(ctx http.Context) *renderers.Form {
    return r.BaseForm(ctx, false).Body([]any{
        gamis.GroupControl().Body([]any{
            gamis.TextControl().
                Name("name").
                Label(tools.AdminLang(ctx, "admin_role.name")).
                Required(true),
            gamis.TextControl().
                Name("slug").
                Label(tools.AdminLang(ctx, "admin_role.slug")).
                Required(true).
                Description(tools.AdminLang(ctx, "admin_role.slug_tip")),
        }),
    })
}

// Detail 详情页面
func (r *RoleController) Detail(ctx http.Context) http.Response {
	roleId := ctx.Request().Input("id", "")
	if roleId == "" {
		page := gamis.Page().
			Title(tools.AdminLang(ctx, "admin_role.detail")).
			Body([]any{
				gamis.TextControl().Name("error").Value(tools.AdminLang(ctx, "admin_role.id_required")),
			})
		return r.SuccessData(ctx, page)
	}

	// 获取角色信息
	var role admin.AdminRole
	if err := facades.Orm().Query().Find(&role, roleId); err != nil {
		page := gamis.Page().
			Title(tools.AdminLang(ctx, "admin_role.detail")).
			Body([]any{
				gamis.TextControl().Name("error").Value(tools.AdminLang(ctx, "admin_role.not_found")),
			})
		return r.SuccessData(ctx, page)
	}

	// 获取角色的权限
	permissions := role.AllPermissions()

	// 构建权限数据
	var permissionData []any
	for _, permission := range permissions {
		permissionData = append(permissionData, map[string]any{
			"id":   permission.ID,
			"name": permission.Name,
			"slug": permission.Slug,
		})
	}

	page := gamis.Page().
		Title(tools.AdminLang(ctx, "admin_role.detail")).
		Body([]any{
			gamis.TextControl().Name("id").Value(roleId).Hidden(true),
			gamis.TextControl().Name("name").Label(tools.AdminLang(ctx, "admin_role.name")).Value(role.Name),
			gamis.TextControl().Name("slug").Label(tools.AdminLang(ctx, "admin_role.slug")).Value(role.Slug),
			gamis.TextControl().Name("created_at").Label(tools.AdminLang(ctx, "created_at")).Value(role.CreatedAt),
			gamis.TextControl().Name("updated_at").Label(tools.AdminLang(ctx, "updated_at")).Value(role.UpdatedAt),
			gamis.Divider().Title(tools.AdminLang(ctx, "admin_role.permissions")),
			gamis.Table().
				Title(tools.AdminLang(ctx, "admin_role.permission_list")).
				Data(permissionData).
				Columns([]any{
					gamis.TableColumn().Name("id").Label(tools.AdminLang(ctx, "admin_role.id")),
					gamis.TableColumn().Name("name").Label(tools.AdminLang(ctx, "admin_role.name")),
					gamis.TableColumn().Name("slug").Label(tools.AdminLang(ctx, "admin_role.slug")),
					gamis.TableColumn().Name("description").Label(tools.AdminLang(ctx, "admin_role.description")),
				}),
		})
	return r.SuccessData(ctx, page)
}

func (r *RoleController) Options(ctx http.Context) http.Response {
    var roles []admin.AdminRole
    if err := facades.Orm().Query().Get(&roles); err != nil {
        return r.FailMsg(ctx, err.Error())
    }
    items := make([]map[string]any, 0, len(roles))
    for _, role := range roles {
        items = append(items, map[string]any{"id": role.ID, "name": role.Name})
    }
    return r.SuccessData(ctx, map[string]any{"items": items})
}

// SavePermissions 保存权限设置
func (r *RoleController) SavePermissions(ctx http.Context) http.Response {
    return r.ControllerImpl.Service.AssignPermissions(ctx)
}

// buildPermissionTree 构建权限树形结构
func buildPermissionTree(permissions []*admin.AdminPermission) []map[string]any {
    parentMap := make(map[uint][]*admin.AdminPermission)

	// 按父级ID分组
	for _, permission := range permissions {
		parentMap[permission.ParentId] = append(parentMap[permission.ParentId], permission)
	}

	// 递归构建树
    var buildTree func(parentId uint) []map[string]any
    buildTree = func(parentId uint) []map[string]any {
        items := make([]map[string]any, 0)
		children := parentMap[parentId]

		for _, permission := range children {
            item := map[string]any{
                "label": permission.Name,
                "value": permission.ID,
                "children": buildTree(permission.ID),
            }
            items = append(items, item)
        }

		return items
	}

    return buildTree(0)
}

// convertPermissionTreeToOptions 将权限树数据转换为TreeControl选项格式
func convertPermissionTreeToOptions(treeData []map[string]any) []map[string]any {
    var options []map[string]any

    for _, node := range treeData {
        option := map[string]any{
            "name": node["label"],
            "id":   node["value"],
        }

        if children, exists := node["children"]; exists {
            if childrenSlice, ok := children.([]map[string]any); ok && len(childrenSlice) > 0 {
                option["children"] = convertPermissionTreeToOptions(childrenSlice)
            }
        }

        options = append(options, option)
    }

    return options
}