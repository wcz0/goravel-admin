package icons

import (
	"strings"
	"testing"
)

// TestCountLoaded 图标快照必须完整加载（iconify 全量约 22 万条）。
func TestCountLoaded(t *testing.T) {
	count := Count()
	if count < 200000 {
		t.Fatalf("图标数量异常偏少: %d，embed 数据可能不完整", count)
	}
}

// TestSearchFallsBackToHome 空关键字与纯空白都回落到 home，与前端默认行为一致。
func TestSearchFallsBackToHome(t *testing.T) {
	expected := Search("home")
	if len(expected) == 0 {
		t.Fatal("home 关键词应能搜到结果")
	}

	for _, empty := range []string{"", "   ", "\t\n"} {
		got := Search(empty)
		if len(got) != len(expected) {
			t.Errorf("空关键字 %q 结果数量 %d，期望与 home 一致的 %d", empty, len(got), len(expected))
		}
	}
}

// TestSearchIsCaseInsensitive 关键字大小写不敏感。
func TestSearchIsCaseInsensitive(t *testing.T) {
	lower := Search("arrow-left")
	upper := Search("ARROW-LEFT")

	if len(lower) == 0 {
		t.Fatal("arrow-left 应能搜到结果")
	}
	if len(lower) != len(upper) {
		t.Errorf("大小写结果数量不一致: %d vs %d", len(lower), len(upper))
	}
}

// TestSearchMatchesSubstring 返回的每一条都必须真的包含关键字。
func TestSearchMatchesSubstring(t *testing.T) {
	const keyword = "arrow-left"

	for _, name := range Search(keyword) {
		if !strings.Contains(strings.ToLower(name), keyword) {
			t.Errorf("结果 %q 不含关键字 %q", name, keyword)
		}
	}
}

// TestSearchRespectsLimit 单次最多返回 1000 条，避免把整个集合塞给前端。
func TestSearchRespectsLimit(t *testing.T) {
	// "a" 这种高频字母一定会匹配到远超上限的图标
	got := Search("a")
	if len(got) > maxResults {
		t.Fatalf("返回 %d 条，超过上限 %d", len(got), maxResults)
	}
	if len(got) != maxResults {
		t.Errorf("\"a\" 应命中满额 %d 条，实际 %d 条", maxResults, len(got))
	}
}

// TestSearchNoMatch 无匹配时返回空切片（而不是 nil），前端可直接遍历。
func TestSearchNoMatch(t *testing.T) {
	got := Search("this-icon-definitely-does-not-exist-zzz")
	if got == nil {
		t.Fatal("无匹配时应返回空切片而不是 nil")
	}
	if len(got) != 0 {
		t.Errorf("无匹配时应返回 0 条，实际 %d 条", len(got))
	}
}
