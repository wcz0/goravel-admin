package bootstrap

import (
	contractsconsole "github.com/goravel/framework/contracts/console"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation"

	"goravel/app/console/commands"
	"goravel/config"
	"goravel/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithMigrations(Migrations).
		WithCommands(func() []contractsconsole.Command {
			return []contractsconsole.Command{
				&commands.Service{},
			}
		}).
		WithRouting(func() {
			routes.Web()
			routes.Api()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
