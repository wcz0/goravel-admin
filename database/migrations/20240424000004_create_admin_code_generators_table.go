package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240424000004CreateAdminCodeGeneratorsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240424000004CreateAdminCodeGeneratorsTable) Signature() string {
	return "20240424000004_create_admin_code_generators_table"
}

// Up Run the migrations.
func (r *M20240424000004CreateAdminCodeGeneratorsTable) Up() error {
	return facades.Schema().Create("admin_code_generators", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("title", 255).Default("").Comment("名称")
		table.String("table_name", 255).Default("").Comment("表名")
		table.String("primary_key", 255).Default("id").Comment("主键名")
		table.String("model_name", 255).Default("").Comment("模型名")
		table.String("controller_name", 255).Default("").Comment("控制器名")
		table.String("service_name", 255).Default("").Comment("服务名")
		table.LongText("columns").Comment("字段信息")
		table.TinyInteger("need_timestamps").Default(0).Comment("是否需要时间戳")
		table.TinyInteger("soft_delete").Default(0).Comment("是否需要软删除")
		table.Text("needs").Nullable().Comment("需要生成的代码")
		table.Text("menu_info").Nullable().Comment("菜单信息")
		table.Text("page_info").Nullable().Comment("页面信息")
		table.Timestamps()
	})
}

// Down Reverse the migrations.
func (r *M20240424000004CreateAdminCodeGeneratorsTable) Down() error {
	return facades.Schema().DropIfExists("admin_code_generators")
}