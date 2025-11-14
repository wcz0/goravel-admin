package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240418000004CreateAdminPermissionsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240418000004CreateAdminPermissionsTable) Signature() string {
	return "20240418000004_create_admin_permissions_table"
}

// Up Run the migrations.
func (r *M20240418000004CreateAdminPermissionsTable) Up() error {
	return facades.Schema().Create("admin_permissions", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.UnsignedBigInteger("parent_id").Default(0).Comment("父级权限ID")
		table.String("name", 255).Comment("权限名称")
		table.String("value", 255).Comment("权限值")
		table.Json("http_method").Comment("请求方法(JSON数组)")
		table.Json("http_path").Comment("请求路径(JSON数组)")
		table.Integer("custom_order").Default(0).Comment("排序")
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (r *M20240418000004CreateAdminPermissionsTable) Down() error {
	return facades.Schema().DropIfExists("admin_permissions")
}