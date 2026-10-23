package services

import (
	"reflect"
	"testing"
)

// TestParseIDList 解析逗号分隔的 id 列表，跳过空值与非法值。
func TestParseIDList(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []any
	}{
		{"单个 id", "1", []any{1}},
		{"多个 id", "1,2,3", []any{1, 2, 3}},
		{"带空格", " 1 , 2 ,3 ", []any{1, 2, 3}},
		{"空串", "", []any{}},
		{"只有逗号", ",,,", []any{}},
		{"含非法值", "1,abc,3", []any{1, 3}},
		{"含 0 与负数", "0,-1,2", []any{2}},
		{"重复 id", "1,1,2", []any{1, 1, 2}},
		{"尾随逗号", "1,2,", []any{1, 2}},
	}

	for _, c := range cases {
		got := parseIDList(c.in)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: parseIDList(%q) = %#v, want %#v", c.name, c.in, got, c.want)
		}
	}
}

// TestParseIDListReturnsAny 结果必须是 []any，否则传给 WhereIn 会类型不匹配。
func TestParseIDListReturnsAny(t *testing.T) {
	got := parseIDList("7")
	if len(got) != 1 {
		t.Fatalf("应解析出 1 个 id，实际 %d 个", len(got))
	}
	if _, ok := got[0].(int); !ok {
		t.Errorf("元素类型应为 int，实际 %T", got[0])
	}
}
