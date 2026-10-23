package controllers

import (
	"strings"
	"testing"
)

// TestParseInt 宽松解析：能解析出前导数字就返回，否则 0。
func TestParseInt(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"12", 12},
		{"0", 0},
		{"-3", -3},
		{"", 0},
		{"abc", 0},
		// Sscanf 会解析到第一个非数字字符为止，锁定该宽松行为
		{"12abc", 12},
		{" 42", 42},
	}

	for _, c := range cases {
		if got := parseInt(c.in); got != c.want {
			t.Errorf("parseInt(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

// TestValidChunkID 上传会话 id 必须是 UUID，且挡住路径穿越与超长输入。
func TestValidChunkID(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"标准 UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"大写 UUID", "550E8400-E29B-41D4-A716-446655440000", true},
		{"空", "", false},
		{"非 UUID", "abc", false},
		{"带上级目录", "../etc/passwd", false},
		{"带反斜杠", "550e8400-e29b-41d4-a716-446655440000\\x", false},
		{"带 NUL", "550e8400-e29b-41d4-a716-446655440000\x00", false},
		{"超长", strings.Repeat("a", 65), false},
		{"长度临界但非 UUID", strings.Repeat("a", 64), false},
	}

	for _, c := range cases {
		if got := validChunkID(c.in); got != c.want {
			t.Errorf("%s: validChunkID(%q) = %v, want %v", c.name, c.in, got, c.want)
		}
	}
}

// TestSafeExtensionBlocksScripts 脚本类后缀一律拒绝，含大写与前导点变体。
func TestSafeExtensionBlocksScripts(t *testing.T) {
	blocked := []string{
		"php", "php3", "php4", "php5", "php7", "php8", "phtml", "phar",
		"sh", "py", "asp", "aspx", "jsp", "env", "exe", "bat", "cmd", "vbs", "htaccess",
		// 大小写与点前缀变体：调用方现在虽然已经 ToLower，这里也不能漏
		"PHP", "Php", "EXE", "JSP", ".php", ".PHP", "BAT",
	}

	for _, ext := range blocked {
		if safeExtension(ext) {
			t.Errorf("后缀 %q 应被拒绝", ext)
		}
	}
}

// TestSafeExtensionAllowsMedia 常见文档、图片、压缩包应放行。
func TestSafeExtensionAllowsMedia(t *testing.T) {
	allowed := []string{
		"png", "jpg", "jpeg", "gif", "webp", "svg", "ico",
		"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx",
		"zip", "rar", "7z", "gz", "tar",
		"mp4", "mp3", "wav", "txt", "csv", "json", "md",
	}

	for _, ext := range allowed {
		if !safeExtension(ext) {
			t.Errorf("后缀 %q 应被放行", ext)
		}
	}
}

// TestNormalizeIntInput 表单里布尔/数字字符串要按 int 规则归一。
func TestNormalizeIntInput(t *testing.T) {
	rules := map[string]string{
		"enabled":   "required|int",
		"parent_id": "int",
		"name":      "required|max:50",
	}

	input := map[string]any{
		"enabled":   float64(1),
		"parent_id": "42",
		"name":      "123",
		"extra":     true,
	}
	normalizeIntInput(input, rules)

	if v, ok := input["enabled"].(int); !ok || v != 1 {
		t.Errorf("float64 应转成 int 1，实际 %#v", input["enabled"])
	}
	if v, ok := input["parent_id"].(int); !ok || v != 42 {
		t.Errorf("数字字符串应转成 int 42，实际 %#v", input["parent_id"])
	}
	// 非 int 规则不处理，否则会把纯数字名称当数字吞掉
	if v, ok := input["name"].(string); !ok || v != "123" {
		t.Errorf("非 int 规则字段不应被转换，实际 %#v", input["name"])
	}
	// 未声明规则的字段不动
	if v, ok := input["extra"].(bool); !ok || !v {
		t.Errorf("未声明规则的字段不应被转换，实际 %#v", input["extra"])
	}
}

// TestNormalizeIntInputBooleans 布尔与真假字符串统一成 1/0。
func TestNormalizeIntInputBooleans(t *testing.T) {
	rules := map[string]string{"a": "int", "b": "int", "c": "int", "d": "int", "e": "int"}
	input := map[string]any{
		"a": true,
		"b": false,
		"c": "true",
		"d": "FALSE",
		"e": "not-a-number",
	}
	normalizeIntInput(input, rules)

	want := map[string]int{"a": 1, "b": 0, "c": 1, "d": 0}
	for field, expected := range want {
		if got, ok := input[field].(int); !ok || got != expected {
			t.Errorf("%s 应转成 %d，实际 %#v", field, expected, input[field])
		}
	}
	// 无法转换的值保持原样，交给校验器报错
	if got, ok := input["e"].(string); !ok || got != "not-a-number" {
		t.Errorf("无法转换的值应保持原样，实际 %#v", input["e"])
	}
}
