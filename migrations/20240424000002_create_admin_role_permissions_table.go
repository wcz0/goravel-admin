package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240424000002CreateAdminRolePermissionsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240424000002CreateAdminRolePermissionsTable) Signature() string {
	return "20240424000002_create_admin_role_permissions_table"
}

// Up Run the migrations.
func (r *M20240424000002CreateAdminRolePermissionsTable) Up() error {
	return facades.Schema().Create("admin_role_permissions", func(table schema.Blueprint) {
		table.Integer("role_id").Comment("角色ID")
		table.Integer("permission_id").Comment("权限ID")
		table.Timestamps()
		
		table.Index("role_id")
		table.Index("permission_id")
	})
}

// Down Reverse the migrations.
func (r *M20240424000002CreateAdminRolePermissionsTable) Down() error {
	return facades.Schema().DropIfExists("admin_role_permissions")
}