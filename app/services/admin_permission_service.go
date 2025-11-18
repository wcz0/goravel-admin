package services

import (
	"encoding/json"
	"goravel/app/models"
	"goravel/app/models/admin"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type AdminPermissionService struct {
	*AdminService[*admin.AdminPermission]
}

func NewAdminPermissionService() *AdminPermissionService {
	return &AdminPermissionService{
		AdminService: NewAdminService[*admin.AdminPermission](admin.NewPermission()),
	}
}

// Store 创建权限
func (s *AdminPermissionService) Store(ctx http.Context) http.Response {
    // 获取HTTP方法/路径（优先使用数组）
    methodsArr := ctx.Request().InputArray("http_method")
    pathsArr := ctx.Request().InputArray("http_path")

	permission := &admin.AdminPermission{
		ParentId:    uint(ctx.Request().InputInt("parent_id", 0)),
		Name:        ctx.Request().Input("name"),
		Slug:        ctx.Request().Input("value"),
		CustomOrder: ctx.Request().InputInt("custom_order", 0),
	}

    // 处理HTTP方法：数组优先，其次尝试JSON字符串，最后默认
    if len(methodsArr) > 0 {
        var httpMethodSlice []string
        for _, m := range methodsArr {
            if m != "" {
                httpMethodSlice = append(httpMethodSlice, strings.TrimSpace(m))
            }
        }
        if len(httpMethodSlice) > 0 {
            permission.HttpMethod = models.StringSlice(httpMethodSlice)
        }
    }
    if len(permission.HttpMethod) == 0 {
        httpMethodStr := ctx.Request().Input("http_method", "")
        if httpMethodStr != "" {
            var httpMethodSlice []string
            if err := json.Unmarshal([]byte(httpMethodStr), &httpMethodSlice); err == nil {
                permission.HttpMethod = models.StringSlice(httpMethodSlice)
            } else {
                permission.HttpMethod = models.StringSlice([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"})
            }
        } else {
            permission.HttpMethod = models.StringSlice([]string{"GET"})
        }
    }

    // 处理HTTP路径：数组优先，其次尝试JSON字符串，最后默认
    if len(pathsArr) > 0 {
        var httpPathSlice []string
        for _, p := range pathsArr {
            if p != "" {
                httpPathSlice = append(httpPathSlice, strings.TrimSpace(p))
            }
        }
        if len(httpPathSlice) > 0 {
            permission.HttpPath = models.StringSlice(httpPathSlice)
        }
    }
    if len(permission.HttpPath) == 0 {
        httpPathStr := ctx.Request().Input("http_path", "")
        if httpPathStr != "" {
            var httpPathSlice []string
            if err := json.Unmarshal([]byte(httpPathStr), &httpPathSlice); err == nil {
                permission.HttpPath = models.StringSlice(httpPathSlice)
            } else {
                permission.HttpPath = models.StringSlice([]string{"/*"})
            }
        } else {
            permission.HttpPath = models.StringSlice([]string{"/*"})
        }
    }

	// 检查权限名和值是否已存在
	var existingPermission admin.AdminPermission
	if err := facades.Orm().Query().Where("name = ? OR slug = ?", permission.Name, permission.Slug).First(&existingPermission); err == nil {
		return s.Fail(ctx, "权限名或值已存在")
	}

	if err := facades.Orm().Query().Create(&permission); err != nil {
		return s.Fail(ctx, err.Error())
	}
	return s.Success(ctx, "权限创建成功")
}

// Update 更新权限
func (s *AdminPermissionService) Update(ctx http.Context) http.Response {
    id := ctx.Request().InputInt("id")
    var permission admin.AdminPermission
    if err := facades.Orm().Query().Find(&permission, id); err != nil {
        return s.Fail(ctx, "权限不存在")
    }

    permission.ParentId = uint(ctx.Request().InputInt("parent_id", 0))
    permission.Name = ctx.Request().Input("name")
    permission.Slug = ctx.Request().Input("value")
    permission.CustomOrder = ctx.Request().InputInt("custom_order", 0)

    // 处理HTTP方法：数组优先，其次尝试JSON字符串，最后默认
    methodsArr := ctx.Request().InputArray("http_method")
    if len(methodsArr) > 0 {
        var httpMethodSlice []string
        for _, m := range methodsArr {
            if m != "" {
                httpMethodSlice = append(httpMethodSlice, strings.TrimSpace(m))
            }
        }
        if len(httpMethodSlice) > 0 {
            permission.HttpMethod = models.StringSlice(httpMethodSlice)
        }
    }
    if len(permission.HttpMethod) == 0 {
        httpMethodStr := ctx.Request().Input("http_method", "")
        if httpMethodStr != "" {
            var httpMethodSlice []string
            if err := json.Unmarshal([]byte(httpMethodStr), &httpMethodSlice); err == nil {
                permission.HttpMethod = models.StringSlice(httpMethodSlice)
            } else {
                permission.HttpMethod = models.StringSlice([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"})
            }
        } else {
            permission.HttpMethod = models.StringSlice([]string{"GET"})
        }
    }

    // 处理HTTP路径：数组优先，其次尝试JSON字符串，最后默认
    pathsArr := ctx.Request().InputArray("http_path")
    if len(pathsArr) > 0 {
        var httpPathSlice []string
        for _, p := range pathsArr {
            if p != "" {
                httpPathSlice = append(httpPathSlice, strings.TrimSpace(p))
            }
        }
        if len(httpPathSlice) > 0 {
            permission.HttpPath = models.StringSlice(httpPathSlice)
        }
    }
    if len(permission.HttpPath) == 0 {
        httpPathStr := ctx.Request().Input("http_path", "")
        if httpPathStr != "" {
            var httpPathSlice []string
            if err := json.Unmarshal([]byte(httpPathStr), &httpPathSlice); err == nil {
                permission.HttpPath = models.StringSlice(httpPathSlice)
            } else {
                permission.HttpPath = models.StringSlice([]string{"/*"})
            }
        } else {
            permission.HttpPath = models.StringSlice([]string{"/*"})
        }
    }

	// 检查权限名和值是否已存在（排除当前权限）
	var existingPermission admin.AdminPermission
	if err := facades.Orm().Query().Where("id != ? AND (name = ? OR slug = ?)", id, permission.Name, permission.Slug).First(&existingPermission); err == nil {
		return s.Fail(ctx, "权限名或值已存在")
	}

	if err := facades.Orm().Query().Save(&permission); err != nil {
		return s.Fail(ctx, err.Error())
	}
	return s.Success(ctx, "权限更新成功")
}

// Destroy 删除权限
func (s *AdminPermissionService) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().InputInt("id")
	var permission admin.AdminPermission
	if err := facades.Orm().Query().Find(&permission, id); err != nil {
		return s.Fail(ctx, "权限不存在")
	}

	if _, err := facades.Orm().Query().Delete(&permission); err != nil {
		return s.Fail(ctx, err.Error())
	}
	return s.Success(ctx, "权限删除成功")
}

// List 获取权限列表（树形结构）
func (s *AdminPermissionService) List(ctx http.Context) http.Response {
	query := facades.Orm().Query()

	// 支持关键字搜索 name / slug
	if keyword := ctx.Request().Input("keyword"); keyword != "" {
		query.Where("name", "like", "%"+keyword+"%").
			OrWhere("slug", "like", "%"+keyword+"%")
	}

	// 排序：custom_order 升序，创建时间降序
	query.Order("custom_order asc").Order("created_at desc")

	var permissions []admin.AdminPermission
	if err := query.Get(&permissions); err != nil {
		return s.Fail(ctx, err.Error())
	}

	// 构建树状结构
	treeData := s.buildPermissionTree(permissions)

	return s.SuccessData(ctx, map[string]any{"items": treeData})
}

// PermissionOptions 获取权限选项
// PermissionOptions 获取权限选项列表，用于父级权限选择
func (s *AdminPermissionService) PermissionOptions(ctx http.Context) http.Response {
	var permissions []admin.AdminPermission
	query := facades.Orm().Query().Order("parent_id asc, custom_order asc, id asc")
	if name := ctx.Request().Input("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if slug := ctx.Request().Input("slug"); slug != "" {
		query = query.Where("slug LIKE ?", "%"+slug+"%")
	}
	if err := query.Get(&permissions); err != nil {
		return s.Fail(ctx, err.Error())
	}

	items := make([]map[string]any, 0, len(permissions))
	for _, p := range permissions {
		items = append(items, map[string]any{"id": p.ID, "name": p.Name})
	}
	return s.SuccessData(ctx, map[string]any{"items": items})
}

func (s *AdminPermissionService) PermissionTreeOptions(ctx http.Context) http.Response {
	var permissions []admin.AdminPermission
	query := facades.Orm().Query().Order("parent_id asc, custom_order asc, id asc")
	if name := ctx.Request().Input("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if slug := ctx.Request().Input("slug"); slug != "" {
		query = query.Where("slug LIKE ?", "%"+slug+"%")
	}
	if err := query.Get(&permissions); err != nil {
		return s.Fail(ctx, err.Error())
	}

	tree := s.buildPermissionTree(permissions)
	var toOptions func([]map[string]any) []map[string]any
	toOptions = func(nodes []map[string]any) []map[string]any {
		opts := make([]map[string]any, 0, len(nodes))
		for _, n := range nodes {
			opt := map[string]any{"label": n["name"], "value": n["id"]}
			if children, ok := n["children"].([]map[string]any); ok && len(children) > 0 {
				opt["children"] = toOptions(children)
			}
			opts = append(opts, opt)
		}
		return opts
	}
	return s.SuccessData(ctx, map[string]any{"items": toOptions(tree)})
}

// Show 获取权限详情
func (s *AdminPermissionService) Show(ctx http.Context) http.Response {
	id := ctx.Request().InputInt("id", 0)
	if id == 0 {
		return s.Fail(ctx, "权限ID不能为空")
	}

	var permission admin.AdminPermission
	if err := facades.Orm().Query().Find(&permission, id); err != nil {
		return s.Fail(ctx, "权限不存在")
	}

	// 获取权限的HTTP方法和路径
	httpMethods := permission.GetHttpMethodOptions()
	httpPaths := permission.GetHttpPaths()

	return s.SuccessData(ctx, map[string]any{
		"permission":   permission,
		"http_methods": httpMethods,
		"http_paths":   httpPaths,
	})
}

// QuickEdit 快速编辑
func (s *AdminPermissionService) QuickEdit(ctx http.Context) http.Response {
	// 实现快速编辑逻辑
	return s.Success(ctx, "快速编辑成功")
}

// QuickEditItem 快速编辑单个项目
func (s *AdminPermissionService) QuickEditItem(ctx http.Context) http.Response {
	// 实现快速编辑单个项目逻辑
	return s.Success(ctx, "快速编辑项成功")
}

// 树形结构构建辅助类型
type AdminPermissionTreeNode struct {
	admin.AdminPermission
	Children []*AdminPermissionTreeNode `json:"children,omitempty"`
}

// buildPermissionTree 构建权限树
func buildPermissionTree(permissions []admin.AdminPermission, parentID uint) []*AdminPermissionTreeNode {
	var tree []*AdminPermissionTreeNode
	for _, permission := range permissions {
		if permission.ParentId == parentID {
			node := &AdminPermissionTreeNode{
				AdminPermission: permission,
			}
			// 递归构建子节点
			children := buildPermissionTree(permissions, permission.ID)
			if len(children) > 0 {
				node.Children = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// buildPermissionTree 构建权限树状结构
func (s *AdminPermissionService) buildPermissionTree(permissions []admin.AdminPermission) []map[string]any {
	// 创建权限映射
	permissionMap := make(map[uint]admin.AdminPermission)
	for _, permission := range permissions {
		permissionMap[permission.ID] = permission
	}

	// 构建树状结构
	var tree []map[string]any

	for _, permission := range permissions {
		// 只处理顶级权限（parent_id = 0）
		if permission.ParentId == 0 {
			node := s.buildPermissionNode(permission, permissionMap)
			tree = append(tree, node)
		}
	}

	return tree
}

// buildPermissionNode 构建权限节点（包含子节点）
func (s *AdminPermissionService) buildPermissionNode(permission admin.AdminPermission, permissionMap map[uint]admin.AdminPermission) map[string]any {
	// 处理 HttpMethod：空数组显示为 "ANY"
	httpMethod := "ANY"
	if len(permission.HttpMethod) > 0 {
		httpMethod = strings.Join(permission.HttpMethod, ", ")
	}

	// 处理 HttpPath：空数组显示为 "/"
	httpPath := "/"
	if len(permission.HttpPath) > 0 {
		httpPath = strings.Join(permission.HttpPath, ", ")
	}

	node := map[string]any{
		"id":           permission.ID,
		"name":         permission.Name,
		"slug":         permission.Slug,
		"http_method":  httpMethod,
		"http_path":    httpPath,
		"custom_order": permission.CustomOrder,
		"parent_id":    permission.ParentId,
		"created_at":   permission.CreatedAt,
		"updated_at":   permission.UpdatedAt,
		"children":     []map[string]any{},
	}

	// 递归构建子节点
	for _, childPermission := range permissionMap {
		if childPermission.ParentId == permission.ID {
			childNode := s.buildPermissionNode(childPermission, permissionMap)
			node["children"] = append(node["children"].([]map[string]any), childNode)
		}
	}

	return node
}
