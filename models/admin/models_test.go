package admin

import (
	"testing"

	"github.com/goravel/framework/database/orm"

	"github.com/wcz0/goravel-admin/models"
)

// TestTableNames 中间表名必须与迁移、seeders 里写的一致，否则关联查询直接空结果。
func TestTableNames(t *testing.T) {
	cases := []struct {
		name  string
		got   string
		want  string
		model any
	}{
		{"权限菜单关联", AdminPermissionMenu{}.TableName(), "admin_permission_menu", AdminPermissionMenu{}},
		{"角色权限关联", AdminRolePermission{}.TableName(), "admin_role_permissions", AdminRolePermission{}},
		{"角色用户关联", AdminRoleUser{}.TableName(), "admin_role_users", AdminRoleUser{}},
	}

	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s: TableName() = %q, want %q", c.name, c.got, c.want)
		}
	}
}

// TestPrimaryKey 主键名统一为 id，通用服务依赖它拼查询。
func TestPrimaryKey(t *testing.T) {
	for name, model := range map[string]interface{ PrimaryKey() string }{
		"AdminApi":          &AdminApi{},
		"AdminExtension":    &AdminExtension{},
		"AdminPage":         &AdminPage{},
		"AdminPermission":   &AdminPermission{},
		"AdminRelationship": &AdminRelationship{},
		"AdminUser":         &AdminUser{},
	} {
		if got := model.PrimaryKey(); got != "id" {
			t.Errorf("%s.PrimaryKey() = %q, want id", name, got)
		}
	}
}

// TestAdminRoleIsAdministrator 超级管理员靠 slug 识别。
func TestAdminRoleIsAdministrator(t *testing.T) {
	cases := []struct {
		slug string
		want bool
	}{
		{"administrator", true},
		{"Administrator", false}, // 大小写不匹配视为普通角色
		{"admin", false},
		{"", false},
	}

	for _, c := range cases {
		role := &AdminRole{Slug: c.slug}
		if got := role.IsAdministrator(); got != c.want {
			t.Errorf("slug=%q 时 IsAdministrator() = %v, want %v", c.slug, got, c.want)
		}
	}
}

// TestPermissionHasMethod 空方法列表视为「所有方法」。
func TestPermissionHasMethod(t *testing.T) {
	all := &AdminPermission{}
	if !all.HasMethod("GET") {
		t.Error("空方法列表应匹配任何方法")
	}

	limited := &AdminPermission{HttpMethod: models.StringSlice{"GET", "post"}}
	for _, method := range []string{"GET", "get", "POST", "Post"} {
		if !limited.HasMethod(method) {
			t.Errorf("应匹配方法 %q", method)
		}
	}
	for _, method := range []string{"DELETE", "PUT", ""} {
		if limited.HasMethod(method) {
			t.Errorf("不应匹配方法 %q", method)
		}
	}
}

// TestPermissionHasPath 空路径列表不匹配任何路径（与 HasMethod 的宽松语义相反）。
func TestPermissionHasPath(t *testing.T) {
	empty := &AdminPermission{}
	if empty.HasPath("/admin/users") {
		t.Error("空路径列表不应匹配任何路径")
	}

	limited := &AdminPermission{HttpPath: models.StringSlice{"/admin/users", "/admin/roles"}}
	if !limited.HasPath("/admin/users") {
		t.Error("应匹配已配置路径")
	}
	if limited.HasPath("/admin/users/1") {
		t.Error("HasPath 是精确匹配，不应命中子路径")
	}
	if limited.HasPath("/ADMIN/USERS") {
		t.Error("路径匹配区分大小写")
	}
}

// TestPermissionGetHttpPaths 空路径回落成通配符，用于权限编码。
func TestPermissionGetHttpPaths(t *testing.T) {
	empty := &AdminPermission{}
	paths := empty.GetHttpPaths()
	if len(paths) != 1 || paths[0] != "/*" {
		t.Errorf("空路径应回落 [/*]，实际 %v", paths)
	}

	limited := &AdminPermission{HttpPath: models.StringSlice{"/a", "/b"}}
	if got := limited.GetHttpPaths(); len(got) != 2 || got[0] != "/a" || got[1] != "/b" {
		t.Errorf("应原样返回配置路径，实际 %v", got)
	}
}

// TestPermissionParentChild 父子关系判断只看 parent_id。
func TestPermissionParentChild(t *testing.T) {
	parent := &AdminPermission{Model: orm.Model{ID: 1}}
	child := &AdminPermission{Model: orm.Model{ID: 2}, ParentId: 1}

	if !parent.IsParentOf(child) {
		t.Error("1 应是 2 的父级")
	}
	if child.IsParentOf(parent) {
		t.Error("2 不可能是 1 的父级")
	}
	if !child.IsChildOf(parent) {
		t.Error("2 应是 1 的子级")
	}
	if parent.IsChildOf(child) {
		t.Error("1 不可能是 2 的子级")
	}
}

// TestPermissionHasMenu 关联菜单为空时不算有菜单。
func TestPermissionHasMenu(t *testing.T) {
	if (&AdminPermission{}).HasMenu() {
		t.Error("未关联菜单时 HasMenu() 应为 false")
	}

	withMenu := &AdminPermission{AdminMenus: []*AdminMenu{{}}}
	if !withMenu.HasMenu() {
		t.Error("有关联菜单时 HasMenu() 应为 true")
	}
}

// TestGetHttpMethodOptions 方法下拉选项必须覆盖 HttpMethods 全部常量。
func TestGetHttpMethodOptions(t *testing.T) {
	options := (&AdminPermission{}).GetHttpMethodOptions()
	if len(options) != len(HttpMethods) {
		t.Fatalf("选项数量 %d 与 HttpMethods %d 不一致", len(options), len(HttpMethods))
	}

	for i, method := range HttpMethods {
		if options[i]["label"] != method || options[i]["value"] != method {
			t.Errorf("第 %d 项应为 %q，实际 %v", i, method, options[i])
		}
	}
}
