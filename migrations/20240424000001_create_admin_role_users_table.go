package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240424000001CreateAdminRoleUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20240424000001CreateAdminRoleUsersTable) Signature() string {
	return "20240424000001_create_admin_role_users_table"
}

// Up Run the migrations.
func (r *M20240424000001CreateAdminRoleUsersTable) Up() error {
	return facades.Schema().Create("admin_role_users", func(table schema.Blueprint) {
		table.Integer("role_id").Comment("角色ID")
		table.Integer("user_id").Comment("用户ID")
		table.Timestamps()
		table.Index("role_id")
		table.Index("user_id")
	})
}

// Down Reverse the migrations.
func (r *M20240424000001CreateAdminRoleUsersTable) Down() error {
	return facades.Schema().DropIfExists("admin_role_users")
}