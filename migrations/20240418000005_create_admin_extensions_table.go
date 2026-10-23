package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240418000005CreateAdminExtensionsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240418000005CreateAdminExtensionsTable) Signature() string {
	return "20240418000005_create_admin_extensions_table"
}

// Up Run the migrations.
func (r *M20240418000005CreateAdminExtensionsTable) Up() error {
	if !facades.Schema().HasTable("admin_extensions") {
		return facades.Schema().Create("admin_extensions", func(table schema.Blueprint) {
			table.BigIncrements("id")
			table.String("name", 100).Comment("扩展名称")
			table.TinyInteger("is_enabled").Default(0).Comment("是否启用")
			table.Timestamps()
			table.Unique("name")
		})
	}
	
	return nil
}

// Down Reverse the migrations.
func (r *M20240418000005CreateAdminExtensionsTable) Down() error {
	return facades.Schema().DropIfExists("admin_extensions")
}