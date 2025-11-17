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
				r.RowEditButton(ctx, r.Form(ctx), true, "md", "", ""),
				r.RowDeleteButton(ctx, ""),
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
	// 验证输入
	validation, err := ctx.Request().Validate(map[string]string{
		"name": "required|string|max:255",
		"slug": "required|string|max:255|unique:admin_roles,slug",
	})
	if err != nil {
		return r.FailMsg(ctx, err.Error())
	}
	if validation.Fails() {
		return r.FailMsg(ctx, validation.Errors().All())
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
	// 使用AdminPermissionService获取权限树结构
	permissionService := services.NewAdminPermissionService()
	permissionResponse := permissionService.List(ctx)

	// 提取权限树数据
	var permissionTree []map[string]any
	if permissionResponseData, ok := permissionResponse.Data().(map[string]any); ok {
		if items, ok := permissionResponseData["items"]; ok {
			if treeItems, ok := items.([]map[string]any); ok {
				permissionTree = treeItems
			}
		}
	}

	// 转换为TreeControl需要的格式
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
				Data(map[string]any{"id": "${id}"}).
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
	return (*renderers.LinkAction)(gamis.DrawerAction().
		Label(tools.AdminLang(ctx, "admin_role.set_permissions")).
		Level("link").
		Drawer(drawer))
}

// Form 表单页面
func (r *RoleController) Form(ctx http.Context) *renderers.Form {
	roleId := ctx.Request().Input("id", "")
	isEdit := roleId != ""

	var role admin.AdminRole
	if isEdit {
		if err := facades.Orm().Query().Find(&role, roleId); err != nil {
			// 返回一个空的表单，显示错误信息
			return gamis.Form().
				Title(tools.AdminLang(ctx, "admin_role.edit")).
				Body(gamis.Form().Body(gamis.TextControl().Label("Error").Value("角色不存在")))
		}
	}

	return gamis.Form().
		Title(tools.AdminLang(ctx, "admin_role."+map[bool]string{true: "edit", false: "create"}[isEdit])).
		Api(tools.GetAdmin("system/roles")).
		InitApi(tools.GetAdmin("system/roles/" + roleId + "/edit")).
		Mode("normal").
		Data(map[string]any{
			"id":   role.ID,
			"name": role.Name,
			"slug": role.Slug,
		}).
		Body([]any{
			gamis.GroupControl().Body([]any{
				gamis.TextControl().
					Name("name").
					Label(tools.AdminLang(ctx, "admin_role.name")).
					Required(true).
					Value(role.Name),
				gamis.TextControl().
					Name("slug").
					Label(tools.AdminLang(ctx, "admin_role.slug")).
					Required(true).
					Value(role.Slug).
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

// SavePermissions 保存权限设置
func (r *RoleController) SavePermissions(ctx http.Context) http.Response {
	roleId := ctx.Request().Input("id", "")
	permissions := ctx.Request().Input("permissions", "")

	result := r.Service.SavePermissions(roleId, permissions)
	if result {
		return r.SuccessMsg(ctx, tools.AdminLang(ctx, "admin.save"))
	} else {
		return r.FailMsg(ctx, tools.AdminLang(ctx, "admin.save_failed"))
	}
}

// buildPermissionTree 构建权限树形结构
func buildPermissionTree(permissions []*admin.AdminPermission) []map[string]interface{} {
	parentMap := make(map[uint][]*admin.AdminPermission)

	// 按父级ID分组
	for _, permission := range permissions {
		parentMap[permission.ParentId] = append(parentMap[permission.ParentId], permission)
	}

	// 递归构建树
	var buildTree func(parentId uint) []map[string]interface{}
	buildTree = func(parentId uint) []map[string]interface{} {
		items := make([]map[string]interface{}, 0)
		children := parentMap[parentId]

		for _, permission := range children {
			item := map[string]interface{}{
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
			"label": node["name"],
			"value": node["id"],
		}

		// 处理子节点
		if children, exists := node["children"]; exists {
			if childrenSlice, ok := children.([]map[string]any); ok && len(childrenSlice) > 0 {
				option["children"] = convertPermissionTreeToOptions(childrenSlice)
			}
		}

		options = append(options, option)
	}

	return options
}