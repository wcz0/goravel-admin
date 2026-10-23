package migrations

import "github.com/goravel/framework/contracts/database/schema"

// All 返回 goravel-admin 需要执行的全部迁移。
//
// 迁移由 admin.ServiceProvider 在 Boot 阶段通过 facades.Schema().Register 注册，
// 宿主不需要修改 bootstrap/migrations.go，也不会与宿主自己的 migrations 包重名。
func All() []schema.Migration {
	return []schema.Migration{
		&M20240417000001CreateAdminUsersTable{},
		&M20240418000001CreateAdminSettingsTable{},
		&M20240418000002CreateAdminMenusTable{},
		&M20240418000003CreateAdminPermissionMenuTable{},
		&M20240418000004CreateAdminPermissionsTable{},
		&M20240418000005CreateAdminExtensionsTable{},
		&M20240419000001CreateAdminRolesTable{},
		&M20240424000001CreateAdminRoleUsersTable{},
		&M20240424000002CreateAdminRolePermissionsTable{},
		&M20240424000003CreateAdminApisTable{},
		&M20240424000005CreateAdminPagesTable{},
		&M20240424000006CreateAdminRelationshipsTable{},
	}
}
