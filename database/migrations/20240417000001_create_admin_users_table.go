package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240417000001CreateAdminUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20240417000001CreateAdminUsersTable) Signature() string {
	return "20240417000001_create_admin_users_table"
}

// Up Run the migrations.
func (r *M20240417000001CreateAdminUsersTable) Up() error {
	return facades.Schema().Table("admin_users", func(table schema.Blueprint) {
		table.Primary("id")
		table.UnsignedInteger("id").AutoIncrement()
		table.String("username", 120).Nullable().Default("")
		table.String("password", 80).Nullable().Default("")
		table.String("name", 255).Nullable().Default("")
		table.String("avatar", 255).Nullable()
		table.String("remember_token", 100).Nullable()
		table.Timestamps()
		table.SoftDeletes()
		table.Unique("admin_users_username_unique", "username")
	})
}

// Down Reverse the migrations.
func (r *M20240417000001CreateAdminUsersTable) Down() error {
	return facades.Schema().DropIfExists("admin_users")
}