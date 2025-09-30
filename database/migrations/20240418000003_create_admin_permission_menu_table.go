package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240418000003CreateAdminPermissionMenuTable struct{}

// Signature The unique signature for the migration.
func (r *M20240418000003CreateAdminPermissionMenuTable) Signature() string {
	return "20240418000003_create_admin_permission_menu_table"
}

// Up Run the migrations.
func (r *M20240418000003CreateAdminPermissionMenuTable) Up() error {
	return facades.Schema().Create("admin_permission_menu", func(table schema.Blueprint) {
		table.Integer("permission_id").Comment("权限ID")
		table.Integer("menu_id").Comment("菜单ID")
		table.Timestamps()
		table.Index("permission_id", "menu_id")
	})
}

// Down Reverse the migrations.
func (r *M20240418000003CreateAdminPermissionMenuTable) Down() error {
	return facades.Schema().DropIfExists("admin_permission_menu")
}