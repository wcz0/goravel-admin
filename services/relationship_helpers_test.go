package services

import (
	"strings"
	"testing"

	"github.com/wcz0/goravel-admin/models/admin"
)

// TestPrefixDigit 标识符不能以数字开头。
func TestPrefixDigit(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Users", "Users"},
		{"2fa", "R2fa"},
		{"9", "R9"},
		{"", ""},
		{"_x", "_x"},
	}

	for _, c := range cases {
		if got := prefixDigit(c.in); got != c.want {
			t.Errorf("prefixDigit(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestLastSegment 取 FQCN 最后一段，兼容反斜杠与正斜杠。
func TestLastSegment(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"App\\Models\\User", "User"},
		{"App/Models/User", "User"},
		{"User", "User"},
		{"  App\\Models\\User  ", "User"},
		{"", ""},
		{"App\\Models\\", ""},
	}

	for _, c := range cases {
		if got := lastSegment(c.in); got != c.want {
			t.Errorf("lastSegment(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestFirstArg 按顺序取第一个非空字符串参数。
func TestFirstArg(t *testing.T) {
	params := map[string]any{
		"a":     "",
		"b":     "value-b",
		"c":     "value-c",
		"num":   42,
		"nil":   nil,
		"empty": "   ",
	}

	if got := firstArg(params, "a", "b", "c"); got != "value-b" {
		t.Errorf("应取第一个非空值 value-b，实际 %q", got)
	}
	// 非字符串类型一律跳过
	if got := firstArg(params, "num", "b"); got != "value-b" {
		t.Errorf("应跳过非字符串类型的 num，实际 %q", got)
	}
	// 全是空格也算非空（保持原行为：不做 TrimSpace）
	if got := firstArg(params, "empty", "c"); got != "   " {
		t.Errorf("空格值应被采用，实际 %q", got)
	}
	if got := firstArg(params, "missing", "nil"); got != "" {
		t.Errorf("都取不到时应返回空串，实际 %q", got)
	}
	if got := firstArg(nil, "a"); got != "" {
		t.Errorf("nil map 应返回空串，实际 %q", got)
	}
}

// TestOrDefault 空值回落。
func TestOrDefault(t *testing.T) {
	if got := orDefault("", "fallback"); got != "fallback" {
		t.Errorf("空值应回落，实际 %q", got)
	}
	if got := orDefault("value", "fallback"); got != "value" {
		t.Errorf("非空值应保留，实际 %q", got)
	}
	// 空格不算空（与 optionalString 的行为不同，锁定避免误改）
	if got := orDefault("   ", "fallback"); got != "   " {
		t.Errorf("空格应视为非空，实际 %q", got)
	}
}

// TestCastToString nil 安全解引用。
func TestCastToString(t *testing.T) {
	if got := castToString(nil); got != "" {
		t.Errorf("nil 应返回空串，实际 %q", got)
	}
	value := "hello"
	if got := castToString(&value); got != "hello" {
		t.Errorf("应返回指针内容，实际 %q", got)
	}
	empty := ""
	if got := castToString(&empty); got != "" {
		t.Errorf("空串指针应返回空串，实际 %q", got)
	}
}

// TestOptionalString 空字符串转 nil，避免写入无意义的空值。
func TestOptionalString(t *testing.T) {
	for _, empty := range []string{"", "   ", "\t", "\n"} {
		if got := optionalString(empty); got != nil {
			t.Errorf("空白值 %q 应转成 nil，实际 %q", empty, *got)
		}
	}

	got := optionalString("value")
	if got == nil {
		t.Fatal("非空值不应转成 nil")
	}
	if *got != "value" {
		t.Errorf("应返回原值，实际 %q", *got)
	}
}

// TestRelationshipTypes 必须覆盖 owl-admin 的 9 种关联，且 label/value 一致。
func TestRelationshipTypes(t *testing.T) {
	types := relationshipTypes()
	if len(types) != 9 {
		t.Errorf("应支持 9 种关联，实际 %d 种", len(types))
	}

	seen := make(map[string]bool, len(types))
	for _, item := range types {
		value := item["value"]
		if value == "" {
			t.Error("关联类型 value 不能为空")
		}
		if item["label"] != value {
			t.Errorf("label 与 value 应一致，实际 %q / %q", item["label"], value)
		}
		if seen[value] {
			t.Errorf("关联类型重复: %q", value)
		}
		seen[value] = true
	}

	for _, required := range []string{"hasOne", "hasMany", "belongsTo", "belongsToMany"} {
		if !seen[required] {
			t.Errorf("缺少必备关联类型 %q", required)
		}
	}
}

// TestPreviewCodeFallsBackToPlaceholder 参数缺失时仍要产出可编译的骨架，不能出现空标识符。
func TestPreviewCodeFallsBackToPlaceholder(t *testing.T) {
	code := previewCode(admin.AdminRelationship{
		Model: "",
		Title: "",
		Type:  "hasMany",
		Args:  nil,
	})

	for _, bad := range []string{"* (", "[] (", "* )"} {
		if strings.Contains(code, bad) {
			t.Errorf("缺少参数时生成了空标识符 %q，输出：\n%s", bad, code)
		}
	}
	if !strings.Contains(code, "Related") {
		t.Errorf("应使用 Related 作为占位模型名，输出：\n%s", code)
	}
}
