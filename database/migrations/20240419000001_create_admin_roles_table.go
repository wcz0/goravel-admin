package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240419000001CreateAdminRolesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240419000001CreateAdminRolesTable) Signature() string {
	return "20240419000001_create_admin_roles_table"
}

// Up Run the migrations.
func (r *M20240419000001CreateAdminRolesTable) Up() error {
	return facades.Schema().Create("admin_roles", func(table schema.Blueprint) {
		table.Primary("id")
		table.UnsignedInteger("id").AutoIncrement()
		table.String("name", 255).Nullable().Default("").Comment("角色名称")
		table.String("value", 255).Nullable().Default("").Comment("值")
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (r *M20240419000001CreateAdminRolesTable) Down() error {
	return facades.Schema().DropIfExists("admin_roles")
}