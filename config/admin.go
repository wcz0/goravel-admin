package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("admin", map[string]any{
		// App Name
		"name": "Goravel Admin",

		"logo": config.Env("ADMIN_LOGO", "/admin/logo.png"),

		"default_avatar": "/admin/avatar.png",

		// TODO: 未实现
		// 应用路由
		"route": map[string]any{
			"prefix": config.Env("ADMIN_ROUTE_PREFIX", "admin-api"),
			"domain": config.Env("ADMIN_DOMAIN"),
			// "middleware": []string{
			// 	"admin",
			// },
			"without_extra_routes": []string{
				"/dashboard",
			},
		},

		"auth": map[string]any{
			"login_captcha":    config.Env("ADMIN_LOGIN_CAPTCHA", true),
			"enable":           config.Env("ADMIN_ENABLE_AUTH", true),
			"permission":       config.Env("ADMIN_ENABLE_PERMISSION", true),
			"token_expiration": config.Env("ADMIN_TOKEN_EXPIRATION", 86400),
			// TODO: 未实现
			// 	"model":         "admin_user",
			// 	"controller":    "admin",
			// 	"guard":         "admin_user",
			"except": []string{},
		},

		// TODO: 未实现
		// "upload": map[string]any{
		// 	"disk": "public",
		// 	"directory": map[string]any{
		// 		"image": "images",
		// 		"file":  "files",
		// 		"rich":  "rich",
		// 	},
		// },

		"https": config.Env("ADMIN_HTTPS", false),

		// TODO: 未实现
		// 是否显示 [开发者工具]
		"show_development_tools": config.Env("ADMIN_SHOW_DEVELOPMENT_TOOLS", false),

		// TODO: 未实现
		// 是否显示 [权限] 功能中的自动生成按钮
		"show_auto_generate_permission_button": config.Env("ADMIN_SHOW_AUTO_GENERATE_PERMISSION_BUTTON", false),

		// TODO: 未实现 扩展
		// "extension": map[string]any{
		// 	"dir": base_path("extensions"),
		// },

		"layout": map[string]any{
			// 浏览器标题, 功能名称使用 %title% 代替
			"title": config.Env("ADMIN_SITE_TITLE", "%title% | Goravel Admin"),
			"header": map[string]any{
				// 是否显示 [刷新] 按钮
				"refresh": config.Env("ADMIN_HEADER_REFRESH", true),
				// 是否显示 [暗黑模式] 按钮
				"dark": config.Env("ADMIN_HEADER_DARK", true),
				// 是否显示 [全屏] 按钮
				"full_screen": config.Env("ADMIN_HEADER_FULL_SCREEN", true),
				// 是否显示 [多语言] 按钮
				"locale_toggle": config.Env("ADMIN_HEADER_LOCALE_TOGGLE", true),
				// 是否显示 [主题配置] 按钮
				"theme_config": config.Env("ADMIN_HEADER_THEME_CONFIG", true),
			},
			// 多语言选项
			"locale_options": map[string]string{
				"en":    "English",
				"zh_CN": "简体中文",
			},
			/*
			 * keep_alive 页面缓存黑名单
			 *
			 * eg:
			 * 列表: /user
			 * 详情: /user/:id
			 * 编辑: /user/:id/edit
			 * 新增: /user/create
			 */
			"keep_alive_exclude": []any{},
			// 底部信息
			"footer": "<a href='https://github.com/wcz0/goravel-admin' target='_blank'>Goravel Admin</a>",
		},
	})
}
