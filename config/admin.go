package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("admin", map[string]any{
		// App Name
		"name": "Goravel Admin",

		"logo":           config.Env("ADMIN_LOGO", "/admin/logo.png"),
		"default_avatar": "/admin/avatar.png",

		// TODO: 未实现
		// "auth": map[string]any{
		// 	"login_captcha": config.Env("ADMIN_LOGIN_CAPTCHA", false),
		// 	"enable":        true,
		// 	"model":         "admin_user",
		// 	"controller":    "admin",
		// 	"guard":         "admin_user",
		// 	"except":        []any{},
		// },

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
		// "show_development_tools": config.Env("ADMIN_SHOW_DEVELOPMENT_TOOLS", true),

		// 是否显示 [权限] 功能中的自动生成按钮
		// "show_auto_generate_permission_button": config.Env("ADMIN_SHOW_AUTO_GENERATE_PERMISSION_BUTTON", true),

		// 扩展
		// "extension": map[string]any{
		// 	"dir": base_path("extensions"),
		// },

		"layout": map[string]any{
			// 浏览器标题, 功能名称使用 %title% 代替
			"title": "%title% | GoravelAdmin",
			"header": map[string]any{
				// 是否显示 [刷新] 按钮
				"refresh": true,
				// 是否显示 [全屏] 按钮
				"full_screen": true,
				// 是否显示 [主题配置] 按钮
				"theme_config": true,
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
			"footer": "<a href='https://github.com/wcz0/goravel-admin' target='_blank'>Owl Admin</a>",
		},
	})
}
