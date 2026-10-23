package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240424000005CreateAdminPagesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240424000005CreateAdminPagesTable) Signature() string {
	return "20240424000005_create_admin_pages_table"
}

// Up Run the migrations.
func (r *M20240424000005CreateAdminPagesTable) Up() error {
	return facades.Schema().Create("admin_pages", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("title", 255).Comment("页面名称")
		table.String("sign", 255).Comment("页面标识")
		table.LongText("schema").Comment("页面结构")
		table.Timestamps()
	})
}

// Down Reverse the migrations.
func (r *M20240424000005CreateAdminPagesTable) Down() error {
	return facades.Schema().DropIfExists("admin_pages")
}