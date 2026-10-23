package services

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestNormalizePath 接口路径必须归一成以 / 开头、无空白。
func TestNormalizePath(t *testing.T) {
	cases := []struct {
		name       string
		in         string
		want       string
		wantErr    bool
		errKeyword string
	}{
		{name: "补前导斜杠", in: "users/list", want: "/users/list"},
		{name: "已带斜杠", in: "/users/list", want: "/users/list"},
		{name: "单斜杠", in: "/", want: "/"},
		{name: "含空格", in: "/users list", wantErr: true, errKeyword: "空白"},
		{name: "含制表符", in: "/users\tlist", wantErr: true, errKeyword: "空白"},
		{name: "含换行", in: "/users\nlist", wantErr: true, errKeyword: "空白"},
		{name: "空字符串", in: "", wantErr: true, errKeyword: "不可为空"},
	}

	for _, c := range cases {
		got, err := normalizePath(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("%s: 期望报错，实际得到 %q", c.name, got)
				continue
			}
			if c.errKeyword != "" && !strings.Contains(err.Error(), c.errKeyword) {
				t.Errorf("%s: 错误信息应包含 %q，实际 %q", c.name, c.errKeyword, err.Error())
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: 不期望报错，实际 %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s: 得到 %q，期望 %q", c.name, got, c.want)
		}
	}
}

// TestTemplateTitle 已知模板返回中文标题，未知模板原样返回。
func TestTemplateTitle(t *testing.T) {
	if got := templateTitle("list"); got != "列表" {
		t.Errorf("list 的标题应为 列表，实际 %q", got)
	}
	if got := templateTitle("save_settings"); got != "保存配置" {
		t.Errorf("save_settings 的标题应为 保存配置，实际 %q", got)
	}
	if got := templateTitle("unknown_template"); got != "unknown_template" {
		t.Errorf("未知模板应回落为原名，实际 %q", got)
	}
	if got := templateTitle(""); got != "" {
		t.Errorf("空模板名应返回空串，实际 %q", got)
	}
}

// TestMarshalArgs nil 与不可序列化的值都返回 nil，不panic。
func TestMarshalArgs(t *testing.T) {
	if got := marshalArgs(nil); got != nil {
		t.Errorf("nil 应返回 nil，实际 %v", *got)
	}

	// chan 无法 JSON 序列化
	if got := marshalArgs(make(chan int)); got != nil {
		t.Errorf("不可序列化的值应返回 nil，实际 %v", *got)
	}

	got := marshalArgs(map[string]any{"model": "admin_users"})
	if got == nil {
		t.Fatal("map 应能序列化")
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(*got), &decoded); err != nil {
		t.Fatalf("序列化结果不是合法 JSON: %v", err)
	}
	if decoded["model"] != "admin_users" {
		t.Errorf("序列化结果不匹配: %v", decoded)
	}
}

// TestApiTemplatesIntegrity 内置模板表是前后端共用的契约，不允许重名或空方法。
func TestApiTemplatesIntegrity(t *testing.T) {
	templates := apiTemplates()
	if len(templates) != 10 {
		t.Errorf("内置模板应为 10 个，实际 %d 个", len(templates))
	}

	validMethods := map[string]bool{"get": true, "post": true, "put": true, "delete": true}
	seen := make(map[string]bool, len(templates))

	for _, tpl := range templates {
		if tpl.Name == "" {
			t.Error("模板 name 不能为空")
		}
		if seen[tpl.Name] {
			t.Errorf("模板 name 重复: %q", tpl.Name)
		}
		seen[tpl.Name] = true

		if tpl.Title == "" {
			t.Errorf("模板 %q 缺少标题", tpl.Name)
		}
		if !validMethods[tpl.Method] {
			t.Errorf("模板 %q 的 method %q 不合法", tpl.Name, tpl.Method)
		}
		if tpl.Args == nil {
			t.Errorf("模板 %q 缺少参数定义", tpl.Name)
		}
	}

	// 前端新增表单依赖这几个模板存在
	for _, required := range []string{"list", "detail", "create", "update", "delete"} {
		if !seen[required] {
			t.Errorf("缺少必备模板 %q", required)
		}
	}
}
