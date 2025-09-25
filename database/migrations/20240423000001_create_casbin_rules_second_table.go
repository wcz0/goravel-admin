package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240423000001CreateCasbinRulesSecondTable struct{}

// Signature The unique signature for the migration.
func (r *M20240423000001CreateCasbinRulesSecondTable) Signature() string {
	return "20240423000001_create_casbin_rules_second_table"
}

// Up Run the migrations.
func (r *M20240423000001CreateCasbinRulesSecondTable) Up() error {
	return facades.Schema().Create("casbin_rules_second", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("ptype", 100).Default(nil)
		table.String("v0", 100).Default(nil)
		table.String("v1", 100).Default(nil)
		table.String("v2", 100).Default(nil)
		table.String("v3", 100).Default(nil)
		table.String("v4", 100).Default(nil)
		table.String("v5", 100).Default(nil)
		table.Unique("ptype", "v0", "v1", "v2", "v3", "v4", "v5")
	})
}

// Down Reverse the migrations.
func (r *M20240423000001CreateCasbinRulesSecondTable) Down() error {
	return facades.Schema().DropIfExists("casbin_rules_second")
}
