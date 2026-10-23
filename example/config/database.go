package config

import (
	"github.com/goravel/framework/contracts/database/driver"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/path"
	sqlitefacades "github.com/goravel/sqlite/facades"
)

func init() {
	config := facades.Config()
	config.Add("database", map[string]any{
		// Default database connection name
		"default": config.Env("DB_CONNECTION", "sqlite"),

		// Database connections
		"connections": map[string]any{
			"sqlite": map[string]any{
				"driver":   "sqlite",
				"database": path.Database("goravel-admin.db"),
				"via": func() (driver.Driver, error) {
					return sqlitefacades.Sqlite("sqlite")
				},
			},
		},

		// Migration Repository Table
		"migrations": map[string]any{
			"table": "migrations",
		},
	})
}
