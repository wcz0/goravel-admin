package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20240417000002CreateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20240417000002CreateUsersTable) Signature() string {
	return "20240417000002_create_users_table"
}

// Up Run the migrations.
func (r *M20240417000002CreateUsersTable) Up() error {
	return facades.Schema().Create("users", func(table schema.Blueprint) {
		table.BigIncrements("id")
		table.String("username", 255).Nullable().Default("")
		table.String("password", 255).Nullable().Default("")
		table.String("name", 255).Nullable().Default("")
		table.String("real_name", 255).Nullable()
		table.String("card_id", 255).Nullable().Default("")
		table.String("mark", 255).Nullable().Default("")
		table.UnsignedInteger("parther_id").Nullable().Default(0)
		table.UnsignedInteger("group_id").Nullable().Default(0)
		table.String("nickname", 255).Nullable().Default("")
		table.String("avatar", 255).Nullable().Default("")
		table.String("phone", 255).Nullable().Default("")
		table.String("email", 255).Nullable().Default("")
		table.String("add_ip", 255).Nullable().Default("")
		table.DateTime("last_time").Nullable()
		table.String("last_ip", 255).Nullable().Default("")
		table.String("money").Nullable().Default("0")
		table.String("brokerage_price").Nullable().Default("0")
		table.UnsignedInteger("integral").Nullable().Default(0)
		table.String("exp").Nullable().Default("0")
		table.UnsignedInteger("sign_num").Nullable().Default(0)
		table.UnsignedTinyInteger("state").Nullable().Default(1)
		table.UnsignedTinyInteger("level").Nullable().Default(0)
		table.UnsignedInteger("agent_level").Nullable().Default(0)
		table.UnsignedTinyInteger("is_spread").Nullable().Default(0)
		table.UnsignedBigInteger("spread_uid").Nullable()
		table.DateTime("spread_time").Nullable()
		table.String("user_type", 255).Nullable().Default("")
		table.UnsignedTinyInteger("is_promoter").Nullable().Default(0)
		table.UnsignedInteger("pay_count").Nullable().Default(0)
		table.UnsignedInteger("spread_count").Nullable().Default(0)
		table.DateTime("clean_time").Nullable()
		table.String("address", 255).Nullable().Default("")
		table.UnsignedInteger("admin_id").Nullable().Default(0)
		table.UnsignedTinyInteger("login_type").Nullable().Default(0)
		table.String("record_phone", 255).Nullable().Default("")
		table.UnsignedTinyInteger("member_level").Nullable().Default(0)
		table.UnsignedTinyInteger("member_ever").Nullable().Default(0)
		table.DateTime("overdue_time").Nullable()
		table.UnsignedTinyInteger("division_type").Nullable().Default(0)
		table.UnsignedTinyInteger("division_status").Nullable().Default(0)
		table.UnsignedTinyInteger("is_division").Nullable().Default(0)
		table.UnsignedTinyInteger("is_agent").Nullable().Default(0)
		table.UnsignedTinyInteger("is_staff").Nullable().Default(0)
		table.UnsignedInteger("division_id").Nullable().Default(0)
		table.Integer("agent_id").Nullable().Default(0)
		table.Integer("staff_id").Nullable().Default(0)
		table.UnsignedTinyInteger("division_percent").Nullable().Default(0)
		table.DateTime("division_end_time").Nullable()
		table.DateTime("division_last_time").Nullable()
		table.UnsignedInteger("division_invite").Nullable().Default(0)
		table.Timestamps()
		table.SoftDeletes()
		table.Unique("users_username_unique", "username")
	})
}

// Down Reverse the migrations.
func (r *M20240417000002CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}