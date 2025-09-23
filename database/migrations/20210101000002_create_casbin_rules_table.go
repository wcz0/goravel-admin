package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20210101000002CreateCasbinRulesTable struct{}

// Signature The unique signature for the migration.
func (r *M20210101000002CreateCasbinRulesTable) Signature() string {
	return "20210101000002_create_casbin_rules_table"
}

// Up Run the migrations.
func (r *M20210101000002CreateCasbinRulesTable) Up() error {
	return facades.Schema().Create("casbin_rules", func(table schema.Blueprint) {
		table.Primary("id")
		table.UnsignedBigInteger("id").AutoIncrement()
		table.String("ptype", 100).Nullable()
		table.String("v0", 100).Nullable()
		table.String("v1", 100).Nullable()
		table.String("v2", 100).Nullable()
		table.String("v3", 100).Nullable()
		table.String("v4", 100).Nullable()
		table.String("v5", 100).Nullable()
		table.Unique("casbin_rules_unique", "ptype", "v0", "v1", "v2", "v3", "v4", "v5")
	})
}

// Down Reverse the migrations.
func (r *M20210101000002CreateCasbinRulesTable) Down() error {
	return facades.Schema().DropIfExists("casbin_rules")
}