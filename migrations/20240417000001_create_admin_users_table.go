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
	return facades.Schema().Create("admin_users", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("username", 120).Comment("用户名")
		table.String("password", 80).Comment("密码")
		table.TinyInteger("enabled").Default(1).Comment("是否启用")
		table.String("name", 255).Default("").Comment("姓名")
		table.String("avatar", 255).Nullable().Comment("头像")
		table.String("remember_token", 100).Nullable().Comment("记住登录令牌")
		table.Timestamps()
		table.Unique("username")
	})
}

// Down Reverse the migrations.
func (r *M20240417000001CreateAdminUsersTable) Down() error {
	return facades.Schema().DropIfExists("admin_users")
}