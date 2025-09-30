package services

import (
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

func (s *AdminPermissionService) Store(ctx http.Context) http.Response {
	// 获取HTTP方法数组
	var httpMethods []string
	if methods := ctx.Request().Input("http_method"); methods != "" {
		httpMethods = strings.Split(methods, ",")
		for i, method := range httpMethods {
			httpMethods[i] = strings.TrimSpace(method)
		}
	}

	// 获取HTTP路径数组
	var httpPaths []string
	if paths := ctx.Request().Input("http_path"); paths != "" {
		httpPaths = strings.Split(paths, ",")
		for i, path := range httpPaths {
			httpPaths[i] = strings.TrimSpace(path)
		}
	}

	permission := &admin.AdminPermission{
		ParentId:    uint32(ctx.Request().InputInt("parent_id", 0)),
		Name:        ctx.Request().Input("name"),
		Slug:        ctx.Request().Input("slug"),
		HttpMethod:  httpMethods,
		HttpPath:    httpPaths,
		CustomOrder: ctx.Request().InputInt("custom_order", 0),
	}
	if err := facades.Orm().Query().Create(&permission); err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	return s.Success(ctx)
}

func (s *AdminPermissionService) Update(ctx http.Context) http.Response {
	var permission admin.AdminPermission
	if err := facades.Orm().Query().Where("id", ctx.Request().InputInt("id")).First(&permission); err != nil {
		return s.FailMsg(ctx, "Permission not found.")
	}

	// 更新基本字段
	permission.ParentId = uint32(ctx.Request().InputInt("parent_id", 0))
	permission.Name = ctx.Request().Input("name")
	permission.Slug = ctx.Request().Input("slug")
	permission.CustomOrder = ctx.Request().InputInt("custom_order", 0)

	// 更新HTTP方法数组
	if methods := ctx.Request().Input("http_method"); methods != "" {
		httpMethods := strings.Split(methods, ",")
		for i, method := range httpMethods {
			httpMethods[i] = strings.TrimSpace(method)
		}
		permission.HttpMethod = httpMethods
	} else {
		permission.HttpMethod = []string{}
	}

	// 更新HTTP路径数组
	if paths := ctx.Request().Input("http_path"); paths != "" {
		httpPaths := strings.Split(paths, ",")
		for i, path := range httpPaths {
			httpPaths[i] = strings.TrimSpace(path)
		}
		permission.HttpPath = httpPaths
	} else {
		permission.HttpPath = []string{}
	}

	if err := facades.Orm().Query().Save(&permission); err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	return s.Success(ctx)
}

func (s *AdminPermissionService) Show(ctx http.Context) http.Response {
	var permission admin.AdminPermission
	id := ctx.Request().InputInt("id")
	if err := facades.Orm().Query().Where("id", id).First(&permission); err != nil {
		return s.FailMsg(ctx, "权限不存在")
	}
	return s.SuccessData(ctx, permission)
}

func (s *AdminPermissionService) Destroy(ctx http.Context) http.Response {
	var permission admin.AdminPermission
	if err := facades.Orm().Query().Where("id", ctx.Request().InputInt("id")).First(&permission); err != nil {
		return s.FailMsg(ctx, "Permission not found.")
	}
	if _, err := facades.Orm().Query().Delete(&permission); err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	return s.Success(ctx)
}

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
        return s.FailMsg(ctx, err.Error())
    }

    // 构建树状结构
    treeData := s.buildPermissionTree(permissions)

    return s.SuccessData(ctx, treeData)
}

// PermissionOptions 获取权限选项
// PermissionOptions 获取权限选项列表，用于父级权限选择
func (s *AdminPermissionService) PermissionOptions(ctx http.Context) []map[string]any {
	var permissions []admin.AdminPermission
	if err := facades.Orm().Query().Order("custom_order asc").Order("created_at desc").Get(&permissions); err != nil {
		return []map[string]any{}
	}

	// 构建层级结构的选项
	options := make([]map[string]any, 0)

	// 添加根选项
	options = append(options, map[string]any{
		"id":   0,
		"name": "根权限",
	})

	// 构建权限映射
	permissionMap := make(map[uint]admin.AdminPermission)
	for _, permission := range permissions {
		permissionMap[permission.ID] = permission
	}

	// 递归构建选项
	s.buildPermissionOptions(&options, permissions, permissionMap, 0, "")

	return options
}

// buildPermissionOptions 递归构建权限选项
func (s *AdminPermissionService) buildPermissionOptions(options *[]map[string]any, permissions []admin.AdminPermission, permissionMap map[uint]admin.AdminPermission, parentId uint32, prefix string) {
	for _, permission := range permissions {
		if permission.ParentId == parentId {
			*options = append(*options, map[string]any{
				"id":   permission.ID,
				"name": prefix + permission.Name,
			})

			// 递归处理子权限
			s.buildPermissionOptions(options, permissions, permissionMap, uint32(permission.ID), prefix+"├─ ")
		}
	}
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
        if childPermission.ParentId == uint32(permission.ID) {
            childNode := s.buildPermissionNode(childPermission, permissionMap)
            node["children"] = append(node["children"].([]map[string]any), childNode)
        }
    }

    return node
}