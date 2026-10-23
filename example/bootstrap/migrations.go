package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"goravel/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000002CreateJobsTable{},
		&migrations.M20240417000002CreateUsersTable{},
	}
}
