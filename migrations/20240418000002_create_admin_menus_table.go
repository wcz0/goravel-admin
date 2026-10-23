package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240418000002CreateAdminMenusTable struct{}

// Signature The unique signature for the migration.
func (r *M20240418000002CreateAdminMenusTable) Signature() string {
	return "20240418000002_create_admin_menus_table"
}

// Up Run the migrations.
func (r *M20240418000002CreateAdminMenusTable) Up() error {
	return facades.Schema().Create("admin_menus", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.UnsignedInteger("parent_id").Default(0).Comment("父级菜单ID")
		table.Integer("custom_order").Default(0).Comment("排序")
		table.String("title", 100).Comment("菜单名称")
		table.String("icon", 100).Nullable().Comment("菜单图标")
		table.String("url", 255).Nullable().Comment("菜单路由")
		table.TinyInteger("url_type").Default(1).Comment("路由类型(1:路由,2:外链,3:iframe)")
		table.TinyInteger("visible").Default(1).Comment("是否可见")
		table.TinyInteger("is_home").Default(0).Comment("是否为首页")
		table.TinyInteger("keep_alive").Nullable().Comment("页面缓存")
		table.String("iframe_url", 255).Nullable().Comment("iframe_url")
		table.String("component", 255).Nullable().Comment("菜单组件")
		table.TinyInteger("is_full").Default(0).Comment("是否是完整页面")
		table.String("extension", 255).Nullable().Comment("扩展")
		table.Timestamps()
	})
}

// Down Reverse the migrations.
func (r *M20240418000002CreateAdminMenusTable) Down() error {
	return facades.Schema().DropIfExists("admin_menus")
}