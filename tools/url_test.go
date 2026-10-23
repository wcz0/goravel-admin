package tools

import "testing"

// TestNormalizeMountPath 覆盖后台挂载路径的边界：补前导斜杠、去结尾斜杠、空值回落。
func TestNormalizeMountPath(t *testing.T) {
	const fallback = "/admin"

	cases := []struct {
		name  string
		value string
		want  string
	}{
		{"标准值", "/admin", "/admin"},
		{"缺前导斜杠", "admin", "/admin"},
		{"多余结尾斜杠", "/admin/", "/admin"},
		{"首尾都有斜杠", "admin/", "/admin"},
		{"空值回落", "", fallback},
		{"纯斜杠回落", "/", fallback},
		{"仅空格回落", "   ", fallback},
		{"首尾空格", "  /manage  ", "/manage"},
		{"多级路径", "backoffice/v2", "/backoffice/v2"},
	}

	for _, c := range cases {
		if got := normalizeMountPath(c.value, fallback); got != c.want {
			t.Errorf("%s: normalizeMountPath(%q) = %q, want %q", c.name, c.value, got, c.want)
		}
	}
}

// TestNormalizeMountPathKeepsInnerSlash 确认内部斜杠不被误删。
func TestNormalizeMountPathKeepsInnerSlash(t *testing.T) {
	if got := normalizeMountPath("/a/b/c/", "/admin"); got != "/a/b/c" {
		t.Errorf("得到 %q，期望 /a/b/c", got)
	}
}
