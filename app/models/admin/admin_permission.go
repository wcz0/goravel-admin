package admin

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"goravel/app/tools"
	"strings"
	"time"
)

type AdminPermission struct {
	ID         uint
	Name       string           `gorm:"uniqueIndex:admin_permissions_name_index"`
	Slug       string           `gorm:"uniqueIndex:admin_permissions_slug_index"`
	HttpMethod models.StringSlice `gorm:"type:json;comment:请求方法"`
	HttpPath   models.StringSlice `gorm:"type:json;comment:请求路径"`
	ParentId   uint             `gorm:"default:0;comment:父级权限ID"`
	CustomOrder int             `gorm:"default:0;comment:排序"`
	CreatedAt time.Time
	UpdatedAt time.Time
	AdminRoles []*AdminRole `gorm:"many2many:admin_role_permissions;joinForeignKey:permission_id;joinReferences:role_id"`
	AdminMenus []*AdminMenu `gorm:"many2many:admin_permission_menu;joinForeignKey:permission_id;joinReferences:menu_id"`
}

// NewPermission 创建新的权限实例
func NewPermission() *AdminPermission {
	return &AdminPermission{}
}

// ShouldPassThrough 检查请求是否应该通过权限验证
func (a *AdminPermission) ShouldPassThrough(ctx http.Context) bool {
	// 如果没有配置路径和方法，则不通过（避免误授全局权限）
	if len(a.HttpPath) == 0 && len(a.HttpMethod) == 0 {
		return false
	}
	
	routePrefix := facades.Config().GetString("admin.route.prefix")
	
	// 遍历路径配置
	for _, path := range a.HttpPath {
		fullPath := strings.TrimSuffix(routePrefix, "/") + path
		
		// 检查路径匹配
		if tools.RequestIs(ctx, fullPath) {
			// 如果没有配置方法限制，任何方法都允许
			if len(a.HttpMethod) == 0 {
				return true
			}
			
			// 检查方法匹配
			for _, method := range a.HttpMethod {
				if tools.RequestMethodIs(ctx, method) {
					return true
				}
			}
		}
	}
	
	return false
}

// GetHttpMethods 获取HTTP方法列表
func (a *AdminPermission) GetHttpMethods() []string {
	if len(a.HttpMethod) == 0 {
		return []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	}
	return a.HttpMethod
}

// GetHttpPaths 获取HTTP路径列表
func (a *AdminPermission) GetHttpPaths() []string {
	if len(a.HttpPath) == 0 {
		return []string{"/*"}
	}
	return a.HttpPath
}

// HasMenu 检查权限是否关联菜单
func (a *AdminPermission) HasMenu() bool {
	return len(a.AdminMenus) > 0
}

// PrimaryKey 返回主键名称
func (a *AdminPermission) PrimaryKey() string {
	return "id"
}