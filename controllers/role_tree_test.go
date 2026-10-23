package controllers

import (
	"fmt"
	"testing"

	"github.com/goravel/framework/database/orm"

	"github.com/wcz0/goravel-admin/models/admin"
)

// rolePermissionFixture 构造两级权限树，顺序按输入顺序。
func rolePermissionFixture() []*admin.AdminPermission {
	return []*admin.AdminPermission{
		{Name: "系统", Model: orm.Model{ID: 1}},
		{Name: "用户管理", ParentId: 1, Model: orm.Model{ID: 2}},
		{Name: "角色管理", ParentId: 1, Model: orm.Model{ID: 3}},
		{Name: "用户列表", ParentId: 2, Model: orm.Model{ID: 4}},
		{Name: "日志", Model: orm.Model{ID: 5}},
	}
}

// TestBuildRolePermissionTree 结构与顺序都应跟随输入。
func TestBuildRolePermissionTree(t *testing.T) {
	tree := buildPermissionTree(rolePermissionFixture())

	if len(tree) != 2 {
		t.Fatalf("应有 2 个根节点，实际 %d 个", len(tree))
	}

	labels := []string{tree[0]["label"].(string), tree[1]["label"].(string)}
	if fmt.Sprint(labels) != fmt.Sprint([]string{"系统", "日志"}) {
		t.Fatalf("根节点顺序应跟随输入，实际 %v", labels)
	}

	systemChildren := tree[0]["children"].([]map[string]any)
	if len(systemChildren) != 2 {
		t.Fatalf("系统下应有 2 个子节点，实际 %d 个", len(systemChildren))
	}
	if systemChildren[0]["label"] != "用户管理" || systemChildren[1]["label"] != "角色管理" {
		t.Errorf("子节点顺序应跟随输入，实际 %v / %v", systemChildren[0]["label"], systemChildren[1]["label"])
	}

	userChildren := systemChildren[0]["children"].([]map[string]any)
	if len(userChildren) != 1 || userChildren[0]["label"] != "用户列表" {
		t.Errorf("用户管理下应只有用户列表，实际 %v", userChildren)
	}

	// 叶子节点 children 是空切片而不是 nil，amis tree 需要能直接取到数组
	if got := tree[1]["children"].([]map[string]any); len(got) != 0 {
		t.Errorf("日志应为叶子节点，实际 %v", got)
	}
}

// TestBuildRolePermissionTreeOrphan 父节点不存在的记录不会出现在结果里（当前行为）。
func TestBuildRolePermissionTreeOrphan(t *testing.T) {
	permissions := []*admin.AdminPermission{
		{Name: "孤儿", ParentId: 99, Model: orm.Model{ID: 10}},
	}

	if tree := buildPermissionTree(permissions); len(tree) != 0 {
		t.Errorf("父节点缺失的记录应被跳过，实际返回 %d 个根节点", len(tree))
	}
}

// TestConvertPermissionTreeToOptions label/value 转成 name/id，children 递归转换。
func TestConvertPermissionTreeToOptions(t *testing.T) {
	tree := buildPermissionTree(rolePermissionFixture())
	options := convertPermissionTreeToOptions(tree)

	if len(options) != 2 {
		t.Fatalf("应转换出 2 个根选项，实际 %d 个", len(options))
	}

	system := options[0]
	if system["name"] != "系统" {
		t.Errorf("name 应为 系统，实际 %v", system["name"])
	}
	if system["id"] != uint(1) {
		t.Errorf("id 应为 1，实际 %v", system["id"])
	}

	children, ok := system["children"].([]map[string]any)
	if !ok || len(children) != 2 {
		t.Fatalf("系统下应有 2 个子选项，实际 %#v", system["children"])
	}

	// 叶子节点不带 children 键，避免前端渲染出空可展开箭头
	logOption := options[1]
	if _, exists := logOption["children"]; exists {
		t.Errorf("叶子节点不应带 children 键，实际 %v", logOption["children"])
	}
}

// TestConvertPermissionTreeToOptionsEmpty 空树与 nil 都返回空切片。
func TestConvertPermissionTreeToOptionsEmpty(t *testing.T) {
	if got := convertPermissionTreeToOptions(nil); len(got) != 0 {
		t.Errorf("nil 输入应返回空结果，实际 %v", got)
	}
	if got := convertPermissionTreeToOptions([]map[string]any{}); len(got) != 0 {
		t.Errorf("空输入应返回空结果，实际 %v", got)
	}
}
