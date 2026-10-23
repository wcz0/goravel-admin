package seeders

import (
	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/facades"

	adminseeders "github.com/wcz0/goravel-admin/seeders"
)

type DatabaseSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *DatabaseSeeder) Signature() string {
	return "DatabaseSeeder"
}

// Run executes the seeder logic.
func (s *DatabaseSeeder) Run() error {
	return facades.Seeder().CallOnce([]seeder.Seeder{
		&adminseeders.AdminSeeder{},
	})
}
