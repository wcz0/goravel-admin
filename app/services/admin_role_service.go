package services

import (
	"goravel/app/models/admin"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type AdminRoleService struct {
	*AdminService[*admin.AdminRole]
}

func NewAdminRoleService() *AdminRoleService {
	return &AdminRoleService{
		AdminService: NewAdminService[*admin.AdminRole](admin.NewAdminRole()),
	}
}

// Store 创建角色
func (s *AdminRoleService) Store(ctx http.Context) http.Response {
	role := &admin.AdminRole{
		Name: ctx.Request().Input("name"),
		Slug: ctx.Request().Input("slug"),
	}

	// 检查角色名和标识是否已存在
	var existingRole admin.AdminRole
	if err := facades.Orm().Query().Where("name = ? OR slug = ?", role.Name, role.Slug).First(&existingRole); err == nil {
		return s.FailMsg(ctx, "角色名或标识已存在")
	}

	if err := facades.Orm().Query().Create(&role); err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	return s.Success(ctx, "角色创建成功")
}

// Update 更新角色
func (s *AdminRoleService) Update(ctx http.Context) http.Response {
	id := ctx.Request().InputInt("id")
	var role admin.AdminRole
	if err := facades.Orm().Query().Find(&role, id); err != nil {
		return s.FailMsg(ctx, "角色不存在")
	}

	role.Name = ctx.Request().Input("name")
	role.Slug = ctx.Request().Input("slug")

	// 检查角色名和标识是否已存在（排除当前角色）
	var existingRole admin.AdminRole
	if err := facades.Orm().Query().Where("id != ? AND (name = ? OR slug = ?)", id, role.Name, role.Slug).First(&existingRole); err == nil {
		return s.FailMsg(ctx, "角色名或标识已存在")
	}

	if err := facades.Orm().Query().Save(&role); err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	return s.Success(ctx, "角色更新成功")
}

// Destroy 删除角色
func (s *AdminRoleService) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().InputInt("id")
	var role admin.AdminRole
	if err := facades.Orm().Query().Find(&role, id); err != nil {
		return s.FailMsg(ctx, "角色不存在")
	}

	// 超级管理员角色不能删除
	if role.IsAdministrator() {
		return s.FailMsg(ctx, "超级管理员角色不能删除")
	}

	if _, err := facades.Orm().Query().Delete(&role); err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	return s.Success(ctx, "角色删除成功")
}

// List 获取角色列表
func (s *AdminRoleService) List(ctx http.Context) http.Response {
	var roles []admin.AdminRole
	if err := facades.Orm().Query().Order("id desc").Find(&roles); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessData(ctx, map[string]any{
		"items": roles,
	})
}

// Show 获取角色详情
func (s *AdminRoleService) Show(ctx http.Context) http.Response {
	id := ctx.Request().InputInt("id", 0)
	if id == 0 {
		return s.FailMsg(ctx, "角色ID不能为空")
	}

	var role admin.AdminRole
	if err := facades.Orm().Query().Find(&role, id); err != nil {
		return s.FailMsg(ctx, "角色不存在")
	}

	// 获取角色的权限
	permissions := role.AllPermissions()

	return s.SuccessData(ctx, map[string]any{
		"role":        role,
		"permissions": permissions,
	})
}

// PermissionOptions 获取权限选项（用于下拉框）
func (s *AdminRoleService) PermissionOptions(ctx http.Context) []map[string]any {
	var permissions []admin.AdminPermission
	if err := facades.Orm().Query().Order("parent_id asc, custom_order asc, id asc").Find(&permissions); err != nil {
		return []map[string]any{}
	}

	var options []map[string]any
	for _, permission := range permissions {
		options = append(options, map[string]any{
			"label": permission.Name,
			"value": permission.ID,
		})
	}
	return options
}

// SavePermissions 保存角色权限
func (s *AdminRoleService) SavePermissions(roleId string, permissions string) bool {
	if roleId == "" {
		return false
	}

	var role admin.AdminRole
	if err := facades.Orm().Query().Find(&role, roleId); err != nil {
		return false
	}

	var permissionList []*admin.AdminPermission
	if permissions != "" {
		if err := facades.Orm().Query().Find(&permissionList, permissions); err != nil {
			return false
		}
	}

	// 同步权限
	if err := role.SyncPermissions(permissionList); err != nil {
		return false
	}

	return true
}
