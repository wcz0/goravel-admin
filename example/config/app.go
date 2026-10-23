package config

import (
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/carbon"
)

// Boot Start all init methods of the current folder to bootstrap all config.
func Boot() {}

func init() {
	config := facades.Config()
	config.Add("app", map[string]any{
		// Application Name
		"name": config.Env("APP_NAME", "Goravel Admin"),

		// Application Environment
		"env": config.Env("APP_ENV", "production"),

		// Application Debug Mode
		"debug": config.Env("APP_DEBUG", false),

		// Application Timezone
		"timezone": carbon.Shanghai,

		// Application Locale Configuration
		"locale": config.Env("APP_LOCALE", "zh"),

		// Application Fallback Locale
		"fallback_locale": config.Env("APP_FALLBACK_LOCALE", "zh"),

		// Encryption Key
		//
		// 32 character string, otherwise these encrypted strings
		// will not be safe. Please do this before deploying an application!
		"key": config.Env("APP_KEY", ""),

		"maintenance": map[string]any{
			"driver": config.Env("APP_MAINTENANCE_DRIVER", "file"),
			"store":  config.Env("APP_MAINTENANCE_STORE", ""),
		},
	})
}
