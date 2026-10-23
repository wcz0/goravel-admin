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
		table.BigIncrements("id")
		table.String("name", 255).Comment("角色名称")
		table.String("slug", 255).Comment("角色标识")
		
		// 唯一约束：name 和 slug 都应该是唯一的
		table.Unique("name")
		table.Unique("slug")
		
		table.Timestamps()
		table.SoftDeletes()
	})
}

// Down Reverse the migrations.
func (r *M20240419000001CreateAdminRolesTable) Down() error {
	return facades.Schema().DropIfExists("admin_roles")
}