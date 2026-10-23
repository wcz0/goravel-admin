package services

import (
	"strings"
	"testing"

	"github.com/wcz0/goravel-admin/models/admin"
)

func TestGoIdent(t *testing.T) {
	cases := []struct {
		raw      string
		exported bool
		want     string
	}{
		{"admin_users", true, "AdminUsers"},
		{"Smoke Roles", true, "SmokeRoles"},
		{"App\\Models\\User", true, "AppModelsUser"},
		{"User", true, "User"},
		{"2fa_codes", true, "R2faCodes"},
		{"admin_users", false, "admin_users"},
		{"用户", true, ""},
		{"", true, ""},
	}

	for _, c := range cases {
		if got := goIdent(c.raw, c.exported); got != c.want {
			t.Errorf("goIdent(%q, %v) = %q, want %q", c.raw, c.exported, got, c.want)
		}
	}
}

func TestPreviewCodeBelongsToMany(t *testing.T) {
	args := `{"related":"admin_roles","table":"role_user","foreign_key":"user_id","related_key":"role_id"}`
	code := previewCode(admin.AdminRelationship{
		Model: "admin_users",
		Title: "Smoke Roles",
		Type:  "belongsToMany",
		Args:  &args,
	})

	for _, want := range []string{
		"func (m *AdminUsers) SmokeRoles() (items []AdminRoles, err error)",
		`Join("inner join role_user`,
		`role_user.role_id`,
		`role_user.user_id`,
	} {
		if !strings.Contains(code, want) {
			t.Errorf("previewCode 缺少 %q，实际输出：\n%s", want, code)
		}
	}
}

func TestPreviewCodeHasMany(t *testing.T) {
	args := `{"related":"admin_menus","foreign_key":"user_id"}`
	code := previewCode(admin.AdminRelationship{
		Model: "admin_users",
		Title: "Menus",
		Type:  "hasMany",
		Args:  &args,
	})

	if !strings.Contains(code, `Where("user_id", m.ID)`) {
		t.Errorf("hasMany 预览应生成 where 条件，实际输出：\n%s", code)
	}
}

func TestExtractSchema(t *testing.T) {
	cases := []struct {
		name    string
		input   map[string]any
		want    string
		wantErr bool
	}{
		{
			name:  "嵌套 page.schema（前端表单）",
			input: map[string]any{"page": map[string]any{"schema": map[string]any{"type": "page"}}},
			want:  `{"type":"page"}`,
		},
		{
			name:  "直接提交对象",
			input: map[string]any{"schema": map[string]any{"type": "page", "body": "hi"}},
			want:  `{"body":"hi","type":"page"}`,
		},
		{
			name:  "直接提交 JSON 文本",
			input: map[string]any{"schema": `{"type":"page"}`},
			want:  `{"type":"page"}`,
		},
		{
			name:    "非法 JSON 文本",
			input:   map[string]any{"schema": `{oops`},
			wantErr: true,
		},
		{
			name:    "空结构",
			input:   map[string]any{"schema": ""},
			wantErr: true,
		},
		{
			name:    "字段缺失",
			input:   map[string]any{"title": "x"},
			wantErr: true,
		},
	}

	for _, c := range cases {
		got, err := extractSchema(c.input)
		if c.wantErr {
			if err == nil {
				t.Errorf("%s: 期望报错，实际得到 %q", c.name, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: 不期望报错，实际得到 %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s: 得到 %q，期望 %q", c.name, got, c.want)
		}
	}
}

func TestFirstByIDRejectsZero(t *testing.T) {
	for _, id := range []int{0, -1} {
		if err := firstByID(id, &admin.AdminPage{}); err != ErrRecordNotFound {
			t.Errorf("firstByID(%d) 应返回 ErrRecordNotFound，实际 %v", id, err)
		}
	}
}
