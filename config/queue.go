package config

import (
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
	redisfacades "github.com/goravel/redis/facades"
)

func init() {
	config := facades.Config()
	config.Add("queue", map[string]any{
		// Default Queue Connection Name
		"default": config.Env("QUEUE_CONNECTION", "sync"),

		// Queue Connections
		//
		// Here you may configure the connection information for each server that is used by your application.
		// Drivers: "sync", "database", "custom"
		"connections": map[string]any{
			"sync": map[string]any{
				"driver": "sync",
			},
			"database": map[string]any{
				"driver":     "database",
				"connection": "mysql",
				"queue":      "default",
				"concurrent": 1,
			},
			"redis": map[string]any{
				"driver":     "custom",
				"connection": "default",
				"queue":      config.Env("REDIS_QUEUE", "default"),
				"via": func() (queue.Driver, error) {
					return redisfacades.Queue("redis")
				},
			},
		},

		"failed": map[string]any{
			"database": config.Env("DB_CONNECTION", "postgres"),
			"table":    "failed_jobs",
		},
	})
}
