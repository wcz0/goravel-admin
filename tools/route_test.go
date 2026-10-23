package tools

import "testing"

// TestMatchPath 钉住权限匹配的路径口径。
//
// 权限表里的 pattern 只在前面补 /admin-api 前缀，随后要和 ctx.Request().Path()
// 比较；两者都必须带前导斜杠，否则通配符永远命中不了（曾经因此让非超管用户全站 403）。
func TestMatchPath(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		{"精确命中", "/admin-api/dashboard", "/admin-api/dashboard", true},
		{"尾部斜杠忽略", "/admin-api/dashboard/", "/admin-api/dashboard", true},
		{"通配符命中", "/admin-api/system/admin_roles*", "/admin-api/system/admin_roles", true},
		{"通配符命中子路径", "/admin-api/system/admin_roles*", "/admin-api/system/admin_roles/1", true},
		{"通配符不跨前缀", "/admin-api/system/admin_roles*", "/admin-api/dashboard", false},
		{"全局通配", "/admin-api*", "/admin-api/system/admin_users", true},
		{"路径更长但前缀不同", "/admin-api/dashboard", "/admin-api/dashboard-x", false},
		{"空 pattern 视作根", "", "/", true},
	}

	for _, c := range cases {
		if got := MatchPath(c.pattern, c.path); got != c.want {
			t.Errorf("%s: MatchPath(%q, %q) = %v, want %v", c.name, c.pattern, c.path, got, c.want)
		}
	}
}
