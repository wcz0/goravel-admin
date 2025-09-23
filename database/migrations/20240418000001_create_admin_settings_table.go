package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240418000001CreateAdminSettingsTable struct{}

// Signature The unique signature for the migration.
func (r *M20240418000001CreateAdminSettingsTable) Signature() string {
	return "20240418000001_create_admin_settings_table"
}

// Up Run the migrations.
func (r *M20240418000001CreateAdminSettingsTable) Up() error {
	return facades.Schema().Create("admin_settings", func(table schema.Blueprint) {
		table.String("key", 255).Nullable().Default("")
		table.LongText("values").Nullable()
		table.Timestamps()
	})
}

// Down Reverse the migrations.
func (r *M20240418000001CreateAdminSettingsTable) Down() error {
	return facades.Schema().DropIfExists("admin_settings")
}