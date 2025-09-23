package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240418000002CreateAdminMenusTable struct{}

// Signature The unique signature for the migration.
func (r *M20240418000002CreateAdminMenusTable) Signature() string {
	return "20240418000002_create_admin_menus_table"
}

// Up Run the migrations.
func (r *M20240418000002CreateAdminMenusTable) Up() error {
	return facades.Schema().Create("admin_menus", func(table schema.Blueprint) {
		table.Primary("id")
		table.UnsignedInteger("id").AutoIncrement()
		table.Integer("parent_id").Nullable().Default(0)
		table.Integer("order").Nullable().Default(0)
		table.String("title", 100).Nullable().Default("")
		table.String("icon", 100).Nullable()
		table.String("uri", 255).Nullable()
		table.TinyInteger("url_type").Nullable().Default(0)
		table.TinyInteger("visible").Nullable().Default(0)
		table.TinyInteger("is_home").Nullable().Default(0)
		table.String("component", 255).Nullable()
		table.TinyInteger("is_full").Nullable().Default(0)
		table.String("extension", 255).Nullable()
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (r *M20240418000002CreateAdminMenusTable) Down() error {
	return facades.Schema().DropIfExists("admin_menus")
}