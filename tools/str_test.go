package tools

import (
	"os"
	"path/filepath"
	"testing"
)

// TestKebabToCamel 覆盖前端 sign（kebab-case）到 Go 命名的转换。
func TestKebabToCamel(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"admin-user", "adminUser"},
		{"admin-user-list", "adminUserList"},
		{"user", "user"},
		{"User", "user"},
		{"a-b-c", "aBC"},
		{"", ""},
		// 已无中划线时保持原样（除首段小写）
		{"admin_users", "admin_users"},
	}

	for _, c := range cases {
		if got := KebabToCamel(c.in); got != c.want {
			t.Errorf("KebabToCamel(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestPascalToCamel 覆盖大驼峰转小驼峰。
func TestPascalToCamel(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"AdminUser", "adminUser"},
		{"User", "user"},
		{"", ""},
		{"A", "a"},
		// 非 ASCII 首字符也要能处理（按 rune 而非 byte 截断）
		{"用户", "用户"},
	}

	for _, c := range cases {
		if got := PascalToCamel(c.in); got != c.want {
			t.Errorf("PascalToCamel(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestPascalToSnake 覆盖大驼峰转蛇形。
func TestPascalToSnake(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"AdminUser", "admin_user"},
		{"UserName", "user_name"},
		{"A", "a"},
		{"", ""},
		// 连续大写会逐字符插入下划线，这是已知行为，锁定避免回归后无人察觉
		{"HTTPServer", "h_t_t_p_server"},
	}

	for _, c := range cases {
		if got := PascalToSnake(c.in); got != c.want {
			t.Errorf("PascalToSnake(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestSnakeToPascal 覆盖蛇形转大驼峰。
func TestSnakeToPascal(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"admin_user", "AdminUser"},
		{"user", "User"},
		{"a_b_c", "ABC"},
		{"", ""},
		// 连续下划线会被 FieldsFunc 折叠
		{"admin__user", "AdminUser"},
	}

	for _, c := range cases {
		if got := SnakeToPascal(c.in); got != c.want {
			t.Errorf("SnakeToPascal(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestReadDir 只返回文件，不含子目录；目录不存在时返回空切片而不是 panic。
func TestReadDir(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"a.txt", "b.js"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatalf("写入文件失败: %v", err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatalf("创建子目录失败: %v", err)
	}

	files := ReadDir(dir)
	if len(files) != 2 {
		t.Fatalf("应返回 2 个文件，实际 %d 个: %v", len(files), files)
	}
	for _, name := range files {
		if name == "sub" {
			t.Error("ReadDir 不应返回子目录")
		}
	}

	if got := ReadDir(filepath.Join(dir, "not-exist")); len(got) != 0 {
		t.Errorf("目录不存在时应返回空切片，实际 %v", got)
	}
}
