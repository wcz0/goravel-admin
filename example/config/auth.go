package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("auth", map[string]any{
		// Authentication Defaults
		"defaults": map[string]any{
			"guard": "user",
		},

		// Authentication Guards
		//
		// goravel-admin 会自动补充 admin guard（jwt + admin provider），
		// 如果宿主希望使用不同的驱动，在这里显式声明即可覆盖。
		"guards": map[string]any{
			"user": map[string]any{
				"driver":   "jwt",
				"provider": "user",
			},
		},

		"providers": map[string]any{
			"user": map[string]any{
				"driver": "orm",
				"model":  "goravel/app/models.User",
			},
		},
	})
}
