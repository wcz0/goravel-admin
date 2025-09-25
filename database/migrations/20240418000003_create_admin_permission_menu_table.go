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
		table.Primary("id")
		table.UnsignedInteger("id").AutoIncrement()
		table.Integer("permission").Nullable()
		table.Integer("menu").Nullable()
		table.DateTime("created_at").Nullable()
		table.DateTime("updated_at").Nullable()
	})
}

// Down Reverse the migrations.
func (r *M20240418000003CreateAdminPermissionMenuTable) Down() error {
	return facades.Schema().DropIfExists("admin_permission_menu")
}