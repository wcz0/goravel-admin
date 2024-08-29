package config

import (

	"github.com/goravel/framework/facades"
	gmodels "github.com/wcz0/goravel-authz/models"
)

func init() {
	config := facades.Config()
	config.Add("casbin", map[string]any{

		// Casbin default
		"default": "basic",

		"models": map[string]any{
			"basic": gmodels.NewRule(),
			// "second": models.NewAdminRule(),
		},
	})
}
