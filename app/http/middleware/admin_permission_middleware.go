package middleware

import (
	"fmt"
	"goravel/app/models/admin"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// AdminPermissionMiddleware 权限验证中间件
type AdminPermissionMiddleware struct{}

// NewAdminPermissionMiddleware 创建权限中间件实例
func NewAdminPermissionMiddleware() *AdminPermissionMiddleware {
	return &AdminPermissionMiddleware{}
}

// Handle 处理权限验证请求
func (m *AdminPermissionMiddleware) Handle(ctx http.Context, next http.HandlerFunc) error {
	// 获取当前用户
	currentUser, ok := ctx.Value("admin_user").(*admin.AdminUser)
	if !ok || currentUser == nil {
		ctx.Response().Json(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return fmt.Errorf("unauthorized")
	}

	// 超级管理员跳过权限检查
	if currentUser.IsAdministrator() {
		next(ctx)
		return nil
	}

	// 获取请求的权限标识
	permissionValue := m.getPermissionValue(ctx)
	if permissionValue == "" {
		ctx.Response().Json(http.StatusForbidden, map[string]interface{}{
			"code":    403,
			"message": "缺少权限标识",
		})
		return fmt.Errorf("permission denied")
	}

	// 检查用户是否具有该权限
	hasPermission := m.checkUserPermission(currentUser, permissionValue)
	if !hasPermission {
		ctx.Response().Json(http.StatusForbidden, map[string]interface{}{
			"code":    403,
			"message": "没有权限访问该资源",
		})
		return fmt.Errorf("permission denied")
	}

	next(ctx)
	return nil
}

// getPermissionValue 获取当前请求的权限标识
func (m *AdminPermissionMiddleware) getPermissionValue(ctx http.Context) string {
	// 首先尝试从请求头获取权限标识
	if permission := ctx.Request().Header("X-Permission", ""); permission != "" {
		return permission
	}

	// 从请求参数获取权限标识
	if permission := ctx.Request().Input("permission", ""); permission != "" {
		return permission
	}

	// 根据路由和HTTP方法自动生成权限标识
	return m.generatePermissionFromRoute(ctx)
}

// generatePermissionFromRoute 根据路由和HTTP方法自动生成权限标识
func (m *AdminPermissionMiddleware) generatePermissionFromRoute(ctx http.Context) string {
	method := ctx.Request().Method()
	path := ctx.Request().Path()

	// 解析路径
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 3 {
		return ""
	}

	// 获取模块名
	module := pathParts[2]

	// 根据HTTP方法生成操作权限
	var operation string
	switch method {
	case http.MethodGet:
		operation = "index"
	case http.MethodPost:
		operation = "store"
	case http.MethodPut, http.MethodPatch:
		operation = "update"
	case http.MethodDelete:
		operation = "destroy"
	default:
		operation = "view"
	}

	return fmt.Sprintf("%s.%s", module, operation)
}

// checkUserPermission 检查用户是否具有指定权限
func (m *AdminPermissionMiddleware) checkUserPermission(user *admin.AdminUser, permissionValue string) bool {
	// 获取用户的所有权限
	permissions := m.getUserPermissions(user)
	
	// 检查是否直接拥有该权限
	for _, permission := range permissions {
		if permission.Slug == permissionValue {
			return true
		}
	}

	// 检查是否通过角色继承该权限
	return m.checkPermissionThroughRoles(user, permissionValue)
}

// getUserPermissions 获取用户的所有权限（包括角色继承的权限）
func (m *AdminPermissionMiddleware) getUserPermissions(user *admin.AdminUser) []admin.AdminPermission {
	var permissions []admin.AdminPermission

	// 获取用户直接关联的角色
	var roles []admin.AdminRole
	if err := facades.Orm().Query().Model(&user).Association("AdminRoles").Find(&roles); err != nil {
		return permissions
	}

	// 获取所有角色的权限
	for _, role := range roles {
		var rolePermissions []admin.AdminPermission
		if err := facades.Orm().Query().Model(&role).Association("AdminPermissions").Find(&rolePermissions); err != nil {
			continue
		}
		permissions = append(permissions, rolePermissions...)
	}

	return permissions
}

// checkPermissionThroughRoles 检查权限是否通过角色继承获得
func (m *AdminPermissionMiddleware) checkPermissionThroughRoles(user *admin.AdminUser, permissionValue string) bool {
	// 获取用户的所有角色
	var roles []admin.AdminRole
	if err := facades.Orm().Query().Model(&user).Association("AdminRoles").Find(&roles); err != nil {
		return false
	}

	// 检查每个角色是否具有该权限
	for _, role := range roles {
		var rolePermissions []admin.AdminPermission
		if err := facades.Orm().Query().Model(&role).Association("AdminPermissions").Find(&rolePermissions); err != nil {
			continue
		}

		for _, permission := range rolePermissions {
			if permission.Slug == permissionValue {
				return true
			}
		}
	}

	return false
}

// HasPermission 检查用户是否具有指定权限的辅助方法
func HasPermission(user *admin.AdminUser, permissionValue string) bool {
	middleware := NewAdminPermissionMiddleware()
	return middleware.checkUserPermission(user, permissionValue)
}

// GetUserPermissions 获取用户所有权限的辅助方法
func GetUserPermissions(user *admin.AdminUser) []admin.AdminPermission {
	middleware := NewAdminPermissionMiddleware()
	return middleware.getUserPermissions(user)
}