package tools

import "testing"

// TestMap2options 验证 map 到 amis options 的转换：key 作 value、value 作 label。
func TestMap2options(t *testing.T) {
	input := map[string]string{
		"enabled":  "启用",
		"disabled": "禁用",
	}

	options := Map2options(input)
	if len(options) != len(input) {
		t.Fatalf("应返回 %d 个选项，实际 %d 个", len(input), len(options))
	}

	// map 遍历顺序随机，按 value 反查
	got := make(map[string]string, len(options))
	for _, option := range options {
		if _, ok := option["label"]; !ok {
			t.Fatalf("选项缺少 label: %v", option)
		}
		if _, ok := option["value"]; !ok {
			t.Fatalf("选项缺少 value: %v", option)
		}
		got[option["value"]] = option["label"]
	}

	for key, label := range input {
		if got[key] != label {
			t.Errorf("value=%q 的 label 应为 %q，实际 %q", key, label, got[key])
		}
	}
}

// TestMap2optionsEmpty 空 map 返回 nil，调用方需要按 nil 安全处理。
func TestMap2optionsEmpty(t *testing.T) {
	if got := Map2options(nil); len(got) != 0 {
		t.Errorf("空 map 应返回空结果，实际 %v", got)
	}
}
