package admin

import (
	"goravel/app/models"
	"goravel/app/tools"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
)

// HTTP方法常量
var HttpMethods = []string{
	"GET",
	"POST",
	"PUT",
	"DELETE",
	"PATCH",
	"OPTIONS",
	"HEAD",
}

type AdminPermission struct {
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	HttpMethod  models.StringSlice `gorm:"type:json" json:"http_method"`
	HttpPath    models.StringSlice `gorm:"type:json" json:"http_path"`
	CustomOrder int                `json:"custom_order"`
	ParentId    uint32             `json:"parent_id"`
	orm.Model
	AdminRoles []*AdminRole `gorm:"many2many:admin_role_permissions;joinForeignKey:permission_id;joinReferences:role_id" json:"-"`
	AdminMenus []*AdminMenu `gorm:"many2many:admin_permission_menu;joinForeignKey:permission_id;joinReferences:menu_id" json:"-"`
}

func NewPermission () *AdminPermission {
	return &AdminPermission{}
}

func (a *AdminPermission) ShouldPassThrough(ctx http.Context) bool {
	// 未配置 http_path 时不放行，避免误授全局权限
	if len(a.HttpPath) == 0 {
		return false
	}
	
	routePrefix := strings.Trim(facades.Config().GetString("admin.route.prefix"), "/")
	
	// 构建匹配规则
	for _, path := range a.HttpPath {
		var methods []string
		var pathPattern string
		
		// 检查是否包含方法前缀 (如 "GET,POST:/api/users")
		if strings.Contains(path, ":") {
			parts := strings.SplitN(path, ":", 2)
			methodStr := strings.TrimSpace(parts[0])
			pathPattern = strings.TrimSpace(parts[1])
			
			// 解析方法列表
			if methodStr != "" {
				methods = strings.Split(methodStr, ",")
				for i, method := range methods {
					methods[i] = strings.TrimSpace(strings.ToUpper(method))
				}
			}
		} else {
			pathPattern = path
			// 如果路径中没有指定方法，使用权限配置的方法
			for _, method := range a.HttpMethod {
				methods = append(methods, strings.ToUpper(method))
			}
		}
		
		// 构建完整路径
		fullPath := pathPattern
		if routePrefix != "" && !strings.HasPrefix(pathPattern, "/"+routePrefix) {
			fullPath = "/" + routePrefix + pathPattern
		}
		
		// 检查路径匹配
		if a.matchRequest(fullPath, methods, ctx) {
			return true
		}
	}
	
	return false
}

// matchRequest 检查请求是否匹配指定的路径和方法
func (a *AdminPermission) matchRequest(pathPattern string, methods []string, ctx http.Context) bool {
	// 处理根路径
	if pathPattern == "/" {
		pathPattern = "/"
	} else {
		pathPattern = strings.Trim(pathPattern, "/")
	}
	
	// 检查路径匹配
	if !tools.RequestIs(ctx, pathPattern) {
		return false
	}
	
	// 如果没有指定方法或方法为空，则匹配所有方法
	if len(methods) == 0 {
		return true
	}
	
	// 检查方法匹配
	requestMethod := strings.ToUpper(ctx.Request().Method())
	for _, method := range methods {
		if method == requestMethod {
			return true
		}
	}
	
	return false
}

// GetHttpMethodOptions 获取HTTP方法选项
func (a *AdminPermission) GetHttpMethodOptions() []map[string]any {
	options := make([]map[string]any, 0, len(HttpMethods))
	for _, method := range HttpMethods {
		options = append(options, map[string]any{
			"label": method,
			"value": method,
		})
	}
	return options
}

// HasMethod 检查权限是否包含指定的HTTP方法
func (a *AdminPermission) HasMethod(method string) bool {
	if len(a.HttpMethod) == 0 {
		return true // 空数组表示所有方法
	}
	
	method = strings.ToUpper(method)
	for _, m := range a.HttpMethod {
		if strings.ToUpper(m) == method {
			return true
		}
	}
	return false
}

// HasPath 检查权限是否包含指定的路径
func (a *AdminPermission) HasPath(path string) bool {
	if len(a.HttpPath) == 0 {
		return false
	}
	
	for _, p := range a.HttpPath {
		if p == path {
			return true
		}
	}
	return false
}

// IsParentOf 检查当前权限是否是指定权限的父级
func (a *AdminPermission) IsParentOf(permission *AdminPermission) bool {
	return permission.ParentId == uint32(a.ID)
}

// IsChildOf 检查当前权限是否是指定权限的子级
func (a *AdminPermission) IsChildOf(permission *AdminPermission) bool {
	return a.ParentId == uint32(permission.ID)
}

func (a *AdminPermission) PrimaryKey() string {
	return "id"
}