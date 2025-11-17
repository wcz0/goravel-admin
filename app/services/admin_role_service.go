package services

import (
	"goravel/app/models/admin"
<<<<<<< HEAD
=======
	"strconv"
	"strings"
>>>>>>> 08c77dc3ed68fd34ac5aa196c797580ff3c72dcb

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type AdminRoleService struct {
	*AdminService[*admin.AdminRole]
}

func NewAdminRoleService() *AdminRoleService {
	return &AdminRoleService{
<<<<<<< HEAD
		AdminService: NewAdminService[*admin.AdminRole](admin.NewAdminRole()),
=======
		AdminService: NewAdminService[*admin.AdminRole](&admin.AdminRole{}),
>>>>>>> 08c77dc3ed68fd34ac5aa196c797580ff3c72dcb
	}
}

// Store 创建角色
<<<<<<< HEAD
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
=======
func (a *AdminRoleService) Store(ctx http.Context) http.Response {
	// 获取请求参数
	name := ctx.Request().Input("name")
	slug := ctx.Request().Input("slug")
	permissionIds := ctx.Request().InputArray("permissions") // 权限ID数组

	// 参数验证
	if name == "" {
		return a.FailMsg(ctx, "角色名称不能为空")
	}
	if slug == "" {
		return a.FailMsg(ctx, "角色标识不能为空")
	}

	// 检查角色名称是否已存在
	var existingRole admin.AdminRole
	if err := facades.Orm().Query().Where("name", name).First(&existingRole); err == nil {
		return a.FailMsg(ctx, "角色名称已存在")
	}

	// 检查角色标识是否已存在
	if err := facades.Orm().Query().Where("slug", slug).First(&existingRole); err == nil {
		return a.FailMsg(ctx, "角色标识已存在")
	}

	// 创建角色
	role := &admin.AdminRole{
		Name: name,
		Slug: slug,
	}

	// 开始事务
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return a.FailMsg(ctx, "开始事务失败: "+err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存角色
	if err := tx.Create(&role); err != nil {
		tx.Rollback()
		return a.FailMsg(ctx, "创建角色失败: "+err.Error())
	}

	// 分配权限
	if len(permissionIds) > 0 {
		var permissionIdList []interface{}
		for _, permissionId := range permissionIds {
			if permissionId != "" {
				permissionIdList = append(permissionIdList, permissionId)
			}
		}

		if len(permissionIdList) > 0 {
			var validPermissions []admin.AdminPermission
			if err := tx.WhereIn("id", permissionIdList).Find(&validPermissions); err != nil {
				tx.Rollback()
				return a.FailMsg(ctx, "获取权限失败: "+err.Error())
			}

			// 使用模型关联方法直接关联权限
			if err := tx.Model(&role).Association("AdminPermissions").Append(&validPermissions); err != nil {
				tx.Rollback()
				return a.FailMsg(ctx, "分配权限失败: "+err.Error())
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return a.FailMsg(ctx, "提交事务失败: "+err.Error())
	}

	return a.SuccessData(ctx, map[string]interface{}{
		"id":   role.ID,
		"name": role.Name,
		"slug": role.Slug,
	})
}

// Update 更新角色
func (a *AdminRoleService) Update(ctx http.Context) http.Response {
	// 获取角色ID
	id := ctx.Request().Input("id")
	if id == "" {
		return a.FailMsg(ctx, "角色ID不能为空")
	}

	// 获取请求参数
	name := ctx.Request().Input("name")
	slug := ctx.Request().Input("slug")
	permissionIds := ctx.Request().InputArray("permissions") // 权限ID数组

	// 查找角色
	var role admin.AdminRole
	if err := facades.Orm().Query().Where("id", id).First(&role); err != nil {
		return a.FailMsg(ctx, "角色不存在")
	}

	// 检查是否为超级管理员角色
	if role.Slug == "administrator" {
		return a.FailMsg(ctx, "超级管理员角色不能修改")
	}

	// 检查角色名称是否被其他角色使用
	if name != "" && name != role.Name {
		var existingRole admin.AdminRole
		if err := facades.Orm().Query().Where("name", name).Where("id", "!=", id).First(&existingRole); err == nil {
			return a.FailMsg(ctx, "角色名称已存在")
		}
	}

	// 检查角色标识是否被其他角色使用
	if slug != "" && slug != role.Slug {
		var existingRole admin.AdminRole
		if err := facades.Orm().Query().Where("slug", slug).Where("id", "!=", id).First(&existingRole); err == nil {
			return a.FailMsg(ctx, "角色标识已存在")
		}
	}

	// 开始事务
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return a.FailMsg(ctx, "开始事务失败: "+err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新角色基本信息
	if name != "" {
		role.Name = name
	}
	if slug != "" {
		role.Slug = slug
	}

	// 保存角色信息
	if err := tx.Save(&role); err != nil {
		tx.Rollback()
		return a.FailMsg(ctx, "更新角色失败: "+err.Error())
	}

	// 更新权限关联
	if len(permissionIds) > 0 {
		// 先删除现有的权限关联
		if _, err := tx.Where("role_id", role.ID).Delete(&admin.AdminRolePermission{}); err != nil {
			tx.Rollback()
			return a.FailMsg(ctx, "删除原有权限关联失败: "+err.Error())
		}

		// 处理新的权限ID
		var permissionIdList []interface{}
		for _, permissionId := range permissionIds {
			if permissionId != "" {
				permissionIdList = append(permissionIdList, permissionId)
			}
		}

		// 获取权限并重新关联
		if len(permissionIdList) > 0 {
			var permissions []admin.AdminPermission
			if err := tx.WhereIn("id", permissionIdList).Find(&permissions); err != nil {
				tx.Rollback()
				return a.FailMsg(ctx, "获取权限失败: "+err.Error())
			}

			// 使用中间表模型直接插入关联关系
			for _, permission := range permissions {
				rolePermission := &admin.AdminRolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
				}
				if err := tx.Create(rolePermission); err != nil {
					tx.Rollback()
					return a.FailMsg(ctx, "分配权限失败: "+err.Error())
				}
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return a.FailMsg(ctx, "提交事务失败: "+err.Error())
	}

	return a.SuccessData(ctx, map[string]interface{}{
		"id":   role.ID,
		"name": role.Name,
		"slug": role.Slug,
	})
}

// Destroy 删除角色
func (a *AdminRoleService) Destroy(ctx http.Context) http.Response {
	// 获取角色ID
	id := ctx.Request().InputInt("id")
	if id == 0 {
		return a.FailMsg(ctx, "角色ID不能为空")
	}

	// 查找角色
	var role admin.AdminRole
	if err := facades.Orm().Query().Where("id", id).First(&role); err != nil {
		return a.FailMsg(ctx, "角色不存在")
	}

	// 检查是否为超级管理员角色
	if role.Slug == "administrator" {
		return a.FailMsg(ctx, "超级管理员角色不能删除")
	}

	// 检查是否有用户使用该角色
	userCount, err := facades.Orm().Query().Model(&admin.AdminRoleUser{}).Where("role_id", id).Count()
	if err != nil {
		return a.FailMsg(ctx, "检查角色使用情况失败: "+err.Error())
	}
	if userCount > 0 {
		return a.FailMsg(ctx, "该角色正在被用户使用，无法删除")
	}

	// 开始事务
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return a.FailMsg(ctx, "开始事务失败: "+err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除角色权限关联
	if _, err := tx.Where("role_id", id).Delete(&admin.AdminRolePermission{}); err != nil {
		tx.Rollback()
		return a.FailMsg(ctx, "删除角色权限关联失败: "+err.Error())
	}

	// 删除角色
	if _, err := tx.Delete(&role); err != nil {
		tx.Rollback()
		return a.FailMsg(ctx, "删除角色失败: "+err.Error())
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return a.FailMsg(ctx, "提交事务失败: "+err.Error())
	}

	return a.SuccessMsg(ctx, "角色删除成功")
}

// List 获取角色列表
func (a *AdminRoleService) List(ctx http.Context) http.Response {
	page := ctx.Request().InputInt("page", 1)
	perPage := ctx.Request().InputInt("perPage", 15)
	name := ctx.Request().Input("name")

	query := facades.Orm().Query().Model(&admin.AdminRole{})

	// 搜索条件
	if name != "" {
		query.Where("name", "like", "%"+name+"%")
	}

	// 分页查询
	var roles []admin.AdminRole
	var total int64

	total, err := query.Count()
	if err != nil {
		return a.FailMsg(ctx, "获取角色总数失败: "+err.Error())
	}

	if err := query.Offset((page-1)*perPage).Limit(perPage).Get(&roles); err != nil {
		return a.FailMsg(ctx, "获取角色列表失败: "+err.Error())
	}

	// 构建返回数据
	items := make([]map[string]interface{}, 0)
	for _, role := range roles {
		items = append(items, map[string]interface{}{
			"id":         role.ID,
			"name":       role.Name,
			"slug":       role.Slug,
			"created_at": role.CreatedAt,
			"updated_at": role.UpdatedAt,
		})
	}

	return a.SuccessData(ctx, map[string]interface{}{
		"items": items,
		"total": total,
		"page":  page,
		"perPage": perPage,
	})
}

// PermissionOptions 获取权限选项
func (a *AdminRoleService) PermissionOptions(ctx http.Context) []map[string]any {
	var permissions []admin.AdminPermission
	if err := facades.Orm().Query().Get(&permissions); err != nil {
		return []map[string]any{}
	}

	options := make([]map[string]any, 0)
	for _, permission := range permissions {
		options = append(options, map[string]any{
			"id":   permission.ID,
			"name": permission.Name,
		})
	}

	return options
}

// Show 获取角色详情
func (a *AdminRoleService) Show(ctx http.Context) http.Response {
	id := ctx.Request().Input("id")
	if id == "" {
		return a.FailMsg(ctx, "角色ID不能为空")
	}

	var role admin.AdminRole
	if err := facades.Orm().Query().With("AdminPermissions").Where("id", id).First(&role); err != nil {
		return a.FailMsg(ctx, "角色不存在")
	}

	// 获取角色的权限ID列表
	permissionIds := make([]uint, 0)
	for _, permission := range role.AdminPermissions {
		permissionIds = append(permissionIds, permission.ID)
	}

	return a.SuccessData(ctx, map[string]interface{}{
		"id":          role.ID,
		"name":        role.Name,
		"slug":        role.Slug,
		"permissions": permissionIds,
		"created_at":  role.CreatedAt,
		"updated_at":  role.UpdatedAt,
	})
}

// AssignPermissions 为角色分配权限
func (a *AdminRoleService) AssignPermissions(ctx http.Context) http.Response {
	roleId := ctx.Request().InputInt("role_id")
	permissionIds := ctx.Request().Input("permission_ids")
	
	// 验证角色是否存在
	var role admin.AdminRole
	if err := facades.Orm().Query().Where("id", roleId).First(&role); err != nil {
		return a.FailMsg(ctx, "角色不存在")
	}
	
	// 开启事务
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return a.FailMsg(ctx, "开启事务失败: "+err.Error())
	}
	defer tx.Rollback()
	
	// 删除原有权限关联
	if _, err := tx.Where("role_id", roleId).Delete(&admin.AdminRolePermission{}); err != nil {
		return a.FailMsg(ctx, "删除原有权限关联失败: "+err.Error())
	}
	
	// 添加新的权限关联
	if permissionIds != "" {
		var permissions []interface{}
		for _, permissionId := range strings.Split(permissionIds, ",") {
			if pid, err := strconv.Atoi(strings.TrimSpace(permissionId)); err == nil && pid > 0 {
				permissions = append(permissions, &admin.AdminRolePermission{
					RoleID:       uint(roleId),
					PermissionID: uint(pid),
				})
			}
		}
		
		if len(permissions) > 0 {
			if err := tx.Create(&permissions); err != nil {
				return a.FailMsg(ctx, "分配权限失败: "+err.Error())
			}
		}
	}
	
	// 提交事务
	if err := tx.Commit(); err != nil {
		return a.FailMsg(ctx, "提交事务失败: "+err.Error())
	}
	
	return a.Success(ctx)
}

// GetRolePermissions 获取角色的权限列表
func (a *AdminRoleService) GetRolePermissions(ctx http.Context) http.Response {
	roleId := ctx.Request().InputInt("role_id")
	
	// 验证角色是否存在
	var role admin.AdminRole
	if err := facades.Orm().Query().Where("id", roleId).First(&role); err != nil {
		return a.FailMsg(ctx, "角色不存在")
	}
	
	// 获取角色的权限ID列表
	var rolePermissions []admin.AdminRolePermission
	if err := facades.Orm().Query().Where("role_id", roleId).Find(&rolePermissions); err != nil {
		return a.FailMsg(ctx, "获取角色权限失败: "+err.Error())
	}
	
	var permissionIds []uint
	for _, rp := range rolePermissions {
		permissionIds = append(permissionIds, rp.PermissionID)
	}
	
	// 获取权限详情
	var permissions []admin.AdminPermission
	if len(permissionIds) > 0 {
		// 转换为 []interface{} 类型
		var ids []interface{}
		for _, id := range permissionIds {
			ids = append(ids, id)
		}
		if err := facades.Orm().Query().WhereIn("id", ids).Find(&permissions); err != nil {
			return a.FailMsg(ctx, "获取权限详情失败: "+err.Error())
		}
	}
	
	return a.SuccessData(ctx, map[string]interface{}{
		"role":        role,
		"permissions": permissions,
		"permission_ids": permissionIds,
	})
}

// RemovePermission 移除角色的特定权限
func (a *AdminRoleService) RemovePermission(ctx http.Context) http.Response {
	roleId := ctx.Request().InputInt("role_id")
	permissionId := ctx.Request().InputInt("permission_id")
	
	// 验证参数
	if roleId <= 0 || permissionId <= 0 {
		return a.FailMsg(ctx, "参数无效")
	}
	
	// 删除权限关联
	if _, err := facades.Orm().Query().Where("role_id", roleId).Where("permission_id", permissionId).Delete(&admin.AdminRolePermission{}); err != nil {
		return a.FailMsg(ctx, "移除权限失败: "+err.Error())
	}
	
	return a.Success(ctx)
}
>>>>>>> 08c77dc3ed68fd34ac5aa196c797580ff3c72dcb
