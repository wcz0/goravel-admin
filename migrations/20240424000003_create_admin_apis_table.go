package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240424000003CreateAdminApisTable struct{}

// Signature The unique signature for the migration.
func (r *M20240424000003CreateAdminApisTable) Signature() string {
	return "20240424000003_create_admin_apis_table"
}

// Up Run the migrations.
func (r *M20240424000003CreateAdminApisTable) Up() error {
	return facades.Schema().Create("admin_apis", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("title", 255).Comment("接口名称")
		table.String("path", 255).Comment("接口路径")
		table.String("template", 255).Comment("接口模板")
		table.TinyInteger("enabled").Default(1).Comment("是否启用")
		table.LongText("args").Nullable().Comment("接口参数")
		table.Timestamps()
	})
}

// Down Reverse the migrations.
func (r *M20240424000003CreateAdminApisTable) Down() error {
	return facades.Schema().DropIfExists("admin_apis")
}