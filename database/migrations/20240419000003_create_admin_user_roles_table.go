package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240419000003CreateAdminUserRolesTable struct{}

// Signature The unique signature for the migration.
func (r *M20240419000003CreateAdminUserRolesTable) Signature() string {
	return "20240419000003_create_admin_user_roles_table"
}

// Up Run the migrations.
func (r *M20240419000003CreateAdminUserRolesTable) Up() error {
	return facades.Schema().Create("admin_user_roles", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.UnsignedBigInteger("user_id").Comment("用户ID")
		table.UnsignedBigInteger("role_id").Comment("角色ID")
		table.Timestamps()
		
		// 唯一约束：防止重复分配
		table.Unique("user_id", "role_id")
	})
}

// Down Reverse the migrations.
func (r *M20240419000003CreateAdminUserRolesTable) Down() error {
	return facades.Schema().DropIfExists("admin_user_roles")
}