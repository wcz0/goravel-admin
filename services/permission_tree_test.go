package services

import (
	"fmt"
	"testing"

	"github.com/goravel/framework/database/orm"

	"github.com/wcz0/goravel-admin/models"
	"github.com/wcz0/goravel-admin/models/admin"
)

// permissionFixture 构造一棵三层的权限树：
//
//	系统(1) ─┬─ 角色管理(3)
//	         └─ 用户管理(2) ── 用户列表(4)
//	日志(5)
func permissionFixture() []admin.AdminPermission {
	return []admin.AdminPermission{
		{
			Name: "系统", Slug: "system",
			HttpMethod:  models.StringSlice{"GET", "POST"},
			HttpPath:    models.StringSlice{"/admin/system/*"},
			CustomOrder: 0,
			Model:       orm.Model{ID: 1},
		},
		{
			Name: "用户管理", Slug: "system.user",
			CustomOrder: 2,
			ParentId:    1,
			Model:       orm.Model{ID: 2},
		},
		{
			Name: "角色管理", Slug: "system.role",
			CustomOrder: 1,
			ParentId:    1,
			Model:       orm.Model{ID: 3},
		},
		{
			Name: "用户列表", Slug: "system.user.list",
			CustomOrder: 0,
			ParentId:    2,
			Model:       orm.Model{ID: 4},
		},
		{
			Name: "日志", Slug: "log",
			CustomOrder: 1,
			Model:       orm.Model{ID: 5},
		},
	}
}

// collectIDs 取出某一层节点的 id 序列，用于断言顺序。
func collectIDs(nodes []map[string]any) []uint {
	ids := make([]uint, 0, len(nodes))
	for _, node := range nodes {
		id, _ := node["id"].(uint)
		ids = append(ids, id)
	}

	return ids
}

// childIDs 取节点 children 的 id 序列。
func childIDs(t *testing.T, node map[string]any) []uint {
	t.Helper()

	children, ok := node["children"].([]map[string]any)
	if !ok {
		t.Fatalf("children 类型不是 []map[string]any: %T", node["children"])
	}

	return collectIDs(children)
}

// TestBuildPermissionTreeStructure 校验层级挂载正确。
func TestBuildPermissionTreeStructure(t *testing.T) {
	service := &AdminPermissionService{}
	tree := service.buildPermissionTree(permissionFixture())

	if len(tree) != 2 {
		t.Fatalf("应有 2 个根节点，实际 %d 个", len(tree))
	}

	roots := map[uint]map[string]any{}
	for _, node := range tree {
		roots[node["id"].(uint)] = node
	}

	system, ok := roots[1]
	if !ok {
		t.Fatal("缺少根节点 1（系统）")
	}
	if _, ok := roots[5]; !ok {
		t.Fatal("缺少根节点 5（日志）")
	}

	if got := childIDs(t, system); len(got) != 2 {
		t.Fatalf("系统下应有 2 个子节点，实际 %v", got)
	}

	// 孙节点应挂在 2（用户管理）下
	var userNode map[string]any
	for _, child := range system["children"].([]map[string]any) {
		if child["id"].(uint) == 2 {
			userNode = child
		}
	}
	if userNode == nil {
		t.Fatal("找不到用户管理节点")
	}
	if got := childIDs(t, userNode); len(got) != 1 || got[0] != 4 {
		t.Fatalf("用户管理下应只有用户列表(4)，实际 %v", got)
	}

	// 叶子节点 children 应为空数组，前端可直接渲染
	if got := childIDs(t, roots[5]); len(got) != 0 {
		t.Errorf("日志应是叶子节点，实际有子节点 %v", got)
	}
}

// TestBuildPermissionTreeMethodFallback 空 http 方法显示 ANY，空路径显示 /。
func TestBuildPermissionTreeMethodFallback(t *testing.T) {
	service := &AdminPermissionService{}
	tree := service.buildPermissionTree(permissionFixture())

	var logNode map[string]any
	for _, node := range tree {
		if node["id"].(uint) == 5 {
			logNode = node
		}
	}
	if logNode == nil {
		t.Fatal("找不到日志节点")
	}

	if logNode["http_method"] != "ANY" {
		t.Errorf("空方法应显示 ANY，实际 %v", logNode["http_method"])
	}
	if logNode["http_path"] != "/" {
		t.Errorf("空路径应显示 /，实际 %v", logNode["http_path"])
	}

	var systemNode map[string]any
	for _, node := range tree {
		if node["id"].(uint) == 1 {
			systemNode = node
		}
	}
	if systemNode["http_method"] != "GET, POST" {
		t.Errorf("多方法应逗号连接，实际 %v", systemNode["http_method"])
	}
}

// TestBuildPermissionTreeOrderIsStable 同一份数据多次构建，子节点顺序必须一致。
//
// 原实现用 map 遍历子节点，Go 的 map 迭代顺序是随机的，会导致前端权限树的
// 顺序每次刷新都变。
func TestBuildPermissionTreeOrderIsStable(t *testing.T) {
	service := &AdminPermissionService{}
	fixture := permissionFixture()

	var first []string
	for round := 0; round < 50; round++ {
		tree := service.buildPermissionTree(fixture)

		var signature []string
		var walk func(nodes []map[string]any)
		walk = func(nodes []map[string]any) {
			for _, node := range nodes {
				signature = append(signature, fmt.Sprint(node["id"]))
				walk(node["children"].([]map[string]any))
			}
		}
		walk(tree)

		if round == 0 {
			first = signature
			continue
		}
		if fmt.Sprint(signature) != fmt.Sprint(first) {
			t.Fatalf("第 %d 次构建的顺序与首次不一致:\n首次 %v\n本次 %v", round, first, signature)
		}
	}
}

// TestBuildPermissionTreeKeepsCustomOrder 子节点按 custom_order 升序排列。
func TestBuildPermissionTreeKeepsCustomOrder(t *testing.T) {
	service := &AdminPermissionService{}
	tree := service.buildPermissionTree(permissionFixture())

	var system map[string]any
	for _, node := range tree {
		if node["id"].(uint) == 1 {
			system = node
		}
	}
	if system == nil {
		t.Fatal("找不到根节点 1")
	}

	// 角色管理(3) 的 custom_order=1 小于 用户管理(2) 的 2
	got := childIDs(t, system)
	want := []uint{3, 2}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("子节点应按 custom_order 排列为 %v，实际 %v", want, got)
	}
}
