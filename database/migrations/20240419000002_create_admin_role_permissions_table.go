package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240419000002CreateAdminRolePermissionsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240419000002CreateAdminRolePermissionsTable) Signature() string {
	return "20240419000002_create_admin_role_permissions_table"
}

// Up Run the migrations.
func (r *M20240419000002CreateAdminRolePermissionsTable) Up() error {
	return facades.Schema().Create("admin_role_permissions", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.UnsignedBigInteger("role_id").Comment("角色ID")
		table.UnsignedBigInteger("permission_id").Comment("权限ID")
		table.Timestamps()
		
		// 唯一约束：防止重复分配
		table.Unique("role_id", "permission_id")
	})
}

// Down Reverse the migrations.
func (r *M20240419000002CreateAdminRolePermissionsTable) Down() error {
	return facades.Schema().DropIfExists("admin_role_permissions")
}