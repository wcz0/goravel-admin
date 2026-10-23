package models

import (
	"reflect"
	"testing"
)

// TestStringSliceScanJSON 标准 JSON 数组。
func TestStringSliceScanJSON(t *testing.T) {
	var s StringSlice
	if err := s.Scan([]byte(`["GET","POST"]`)); err != nil {
		t.Fatalf("扫描失败: %v", err)
	}
	if !reflect.DeepEqual([]string(s), []string{"GET", "POST"}) {
		t.Errorf("得到 %v", s)
	}

	// string 与 []byte 两种来源都要支持
	var fromString StringSlice
	if err := fromString.Scan(`["PUT"]`); err != nil {
		t.Fatalf("字符串来源扫描失败: %v", err)
	}
	if !reflect.DeepEqual([]string(fromString), []string{"PUT"}) {
		t.Errorf("得到 %v", fromString)
	}
}

// TestStringSliceScanNull NULL 应解析为 nil，而不是空切片。
func TestStringSliceScanNull(t *testing.T) {
	s := StringSlice{"stale"}
	if err := s.Scan(nil); err != nil {
		t.Fatalf("NULL 不应报错: %v", err)
	}
	if s != nil {
		t.Errorf("NULL 应得到 nil，实际 %#v", s)
	}
}

// TestStringSliceScanDoubleEncoded 兼容被二次编码的历史数据（PHP 侧序列化遗留）。
func TestStringSliceScanDoubleEncoded(t *testing.T) {
	var s StringSlice
	// 外层是 JSON 字符串，内容是 JSON 数组文本
	if err := s.Scan(`"[\"GET\",\"POST\"]"`); err != nil {
		t.Fatalf("双重编码扫描失败: %v", err)
	}
	if !reflect.DeepEqual([]string(s), []string{"GET", "POST"}) {
		t.Errorf("得到 %v", s)
	}
}

// TestStringSliceScanPHPSingleQuotes 兼容 PHP 风格的 ['GET','POST']。
func TestStringSliceScanPHPSingleQuotes(t *testing.T) {
	var s StringSlice
	if err := s.Scan([]byte(`['GET','POST']`)); err != nil {
		t.Fatalf("PHP 数组扫描失败: %v", err)
	}
	if !reflect.DeepEqual([]string(s), []string{"GET", "POST"}) {
		t.Errorf("得到 %v", s)
	}
}

// TestStringSliceScanCommaSeparated 兼容逗号分隔的原始值。
func TestStringSliceScanCommaSeparated(t *testing.T) {
	var s StringSlice
	if err := s.Scan("GET,POST"); err != nil {
		t.Fatalf("逗号分隔扫描失败: %v", err)
	}
	if !reflect.DeepEqual([]string(s), []string{"GET", "POST"}) {
		t.Errorf("得到 %v", s)
	}

	var single StringSlice
	if err := single.Scan("GET"); err != nil {
		t.Fatalf("单值扫描失败: %v", err)
	}
	if !reflect.DeepEqual([]string(single), []string{"GET"}) {
		t.Errorf("得到 %v", single)
	}
}

// TestStringSliceScanRejectsBadTypes 非字符串、非字节类型直接报错。
func TestStringSliceScanRejectsBadTypes(t *testing.T) {
	cases := []any{123, 1.5, true, []int{1, 2}}

	for _, value := range cases {
		var s StringSlice
		if err := s.Scan(value); err == nil {
			t.Errorf("类型 %T 应报错，实际得到 %#v", value, s)
		}
	}
}

// TestStringSliceScanEmpty 空字符串视为脏数据（NULL 才是空值的正解）。
func TestStringSliceScanEmpty(t *testing.T) {
	var s StringSlice
	if err := s.Scan(""); err == nil {
		t.Errorf("空字符串应报错，实际得到 %#v", s)
	}
}

// TestStringSliceValue 写库时序列化成 JSON 数组。
func TestStringSliceValue(t *testing.T) {
	value, err := StringSlice{"GET", "POST"}.Value()
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	raw, ok := value.([]byte)
	if !ok {
		t.Fatalf("应返回 []byte，实际 %T", value)
	}
	if string(raw) != `["GET","POST"]` {
		t.Errorf("得到 %s", raw)
	}
}

// TestStringSliceRoundTrip 序列化后再解析应还原。
func TestStringSliceRoundTrip(t *testing.T) {
	original := StringSlice{"GET", "PUT", "DELETE"}

	value, err := original.Value()
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	var parsed StringSlice
	if err := parsed.Scan(value); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}
	if !reflect.DeepEqual(original, parsed) {
		t.Errorf("往返后不一致: %v -> %v", original, parsed)
	}
}
