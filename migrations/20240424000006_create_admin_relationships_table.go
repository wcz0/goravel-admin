package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240424000006CreateAdminRelationshipsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240424000006CreateAdminRelationshipsTable) Signature() string {
	return "20240424000006_create_admin_relationships_table"
}

// Up Run the migrations.
func (r *M20240424000006CreateAdminRelationshipsTable) Up() error {
	return facades.Schema().Create("admin_relationships", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("model", 255).Comment("模型")
		table.String("title", 255).Comment("关联名称")
		table.String("type", 255).Comment("关联类型")
		table.String("remark", 255).Nullable().Comment("关联名称")
		table.Text("args").Nullable().Comment("关联参数")
		table.Text("extra").Nullable().Comment("额外参数")
		table.Timestamps()
	})
}

// Down Reverse the migrations.
func (r *M20240424000006CreateAdminRelationshipsTable) Down() error {
	return facades.Schema().DropIfExists("admin_relationships")
}