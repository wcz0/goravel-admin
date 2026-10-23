package migrations

import (
	"regexp"
	"testing"
)

// migrationSignaturePattern Goravel 迁移名带时间戳前缀，用于排序执行。
var migrationSignaturePattern = regexp.MustCompile(`^\d{14}_[a-z0-9_]+$`)

// TestAllSignaturesUnique 签名重复会让迁移被跳过，导致表结构缺列。
func TestAllSignaturesUnique(t *testing.T) {
	list := All()
	if len(list) == 0 {
		t.Fatal("迁移列表为空")
	}

	seen := make(map[string]bool, len(list))
	for _, migration := range list {
		signature := migration.Signature()
		if seen[signature] {
			t.Errorf("迁移签名重复: %q", signature)
		}
		seen[signature] = true
	}
}

// TestAllSignaturesWellFormed 签名要符合 Goravel 的「时间戳_名称」格式。
func TestAllSignaturesWellFormed(t *testing.T) {
	for _, migration := range All() {
		if !migrationSignaturePattern.MatchString(migration.Signature()) {
			t.Errorf("迁移签名格式不合法: %q", migration.Signature())
		}
	}
}

// TestAllCount 表数量是安装契约的一部分，改动需要同步文档与 doctor 检查。
func TestAllCount(t *testing.T) {
	if got := len(All()); got != 13 {
		t.Errorf("迁移数量为 %d，期望 13（新增表请同时更新 admin:doctor 与文档）", got)
	}
}

// TestAllCoversCoreTables 核心表必须都在迁移清单里。
func TestAllCoversCoreTables(t *testing.T) {
	required := []string{
		"create_admin_users_table",
		"create_admin_roles_table",
		"create_admin_permissions_table",
		"create_admin_menus_table",
		"create_admin_settings_table",
		"create_admin_extensions_table",
		"create_admin_pages_table",
		"create_admin_apis_table",
		"create_admin_relationships_table",
	}

	seen := make(map[string]bool, len(required))
	for _, migration := range All() {
		seen[migration.Signature()] = true
	}

	for _, name := range required {
		found := false
		for signature := range seen {
			if len(signature) > len(name) && signature[len(signature)-len(name):] == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("缺少迁移 %q", name)
		}
	}
}
