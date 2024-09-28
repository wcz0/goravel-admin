package services

import (
	"goravel/app/models/admin"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type AdminPermissionService struct {
	*Service
}

func NewAdminPermissionService() *AdminPermissionService {
	return &AdminPermissionService{
		Service: NewService(),
	}
}

func (s *AdminPermissionService) Store(ctx http.Context) http.Response {
	permission := &admin.AdminPermission{
		ParentId: uint(ctx.Request().InputInt("parent_id", 0)),
		Name:     ctx.Request().Input("name"),
		// Value:    ctx.Request().Input("value"),
		// Method:   ctx.Request().Input("method"),
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
	permission.ParentId = uint(ctx.Request().InputInt("parent_id", 0))
	permission.Name = ctx.Request().Input("name")
	// permission.Value = ctx.Request().Input("value")
	// permission.Method = ctx.Request().Input("method")
	if err := facades.Orm().Query().Save(&permission); err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	return s.Success(ctx)
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
	var permissions []admin.AdminPermission
	if err := facades.Orm().Query().Get(&permissions); err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	return s.SuccessData(ctx, map[string]any{
		"items": buildTree(permissions, nil),
	})
}

type AdminPermissionTreeNode struct {
	admin.AdminPermission
	Children []*AdminPermissionTreeNode
}

func buildTree(nodes []admin.AdminPermission, parent *admin.AdminPermission) *AdminPermissionTreeNode {
	tree := &AdminPermissionTreeNode{}
	for _, node := range nodes {
		if node.ParentId == parent.ID {
			child := buildTree(nodes, &node)
			tree.Children = append(tree.Children, child)
		}
	}
	return tree
}
