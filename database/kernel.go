package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"

	"goravel/database/migrations"
	"goravel/database/seeders"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000002CreateCasbinRulesTable{},
		&migrations.M20210101000002CreateJobsTable{},
		&migrations.M20240417000001CreateAdminUsersTable{},
		&migrations.M20240417000002CreateUsersTable{},
		&migrations.M20240418000001CreateAdminSettingsTable{},
		&migrations.M20240418000002CreateAdminMenusTable{},
		&migrations.M20240418000003CreateAdminPermissionMenuTable{},
		&migrations.M20240418000004CreateAdminPermissionsTable{},
		&migrations.M20240418000005CreateAdminExtensionsTable{},
		&migrations.M20240419000001CreateAdminRolesTable{},
		&migrations.M20240423000001CreateCasbinRulesSecondTable{},
		&migrations.M20240418000001CreateAdminSettingsTable{},
	}
}

func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
	}
}
