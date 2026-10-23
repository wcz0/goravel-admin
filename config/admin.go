// Package config 提供 goravel-admin 的默认配置。
//
// 宿主应用无需复制任何文件即可获得这些默认值：admin.ServiceProvider 会在
// Register 阶段把它们深度合并进 app 的 "admin" 配置键。若需要覆盖，只要在宿主的
// config 目录里新增自己的 admin 配置（可执行 `artisan package:install` 生成骨架），
// 未显式设置的项仍会回落到这里的默认值。
package config

import (
	"github.com/goravel/framework/facades"
)

// New 返回 goravel-admin 的默认配置。
func New() map[string]any {
	config := facades.Config()

	return map[string]any{
		// 后台名称
		"name": config.Env("ADMIN_NAME", "Goravel Admin"),

		// 后台 Logo，由内嵌静态资源提供
		"logo": config.Env("ADMIN_LOGO", "/admin-assets/logo.png"),

		// 默认头像
		"default_avatar": config.Env("ADMIN_DEFAULT_AVATAR", "/admin-assets/default-avatar.png"),

		// 应用路由
		"route": map[string]any{
			// 后端 API 前缀
			"prefix": config.Env("ADMIN_ROUTE_PREFIX", "admin-api"),
			// 后台页面挂载路径
			"path": config.Env("ADMIN_ROUTE_PATH", "/admin"),
			// 静态资源挂载路径，需要与前端构建时的 base 保持一致
			"assets_path": config.Env("ADMIN_ASSETS_PATH", "/admin-assets"),
			// 不注册额外路由的白名单
			"without_extra_routes": []string{
				"/dashboard",
			},
		},

		// 认证
		"auth": map[string]any{
			"login_captcha":    config.Env("ADMIN_LOGIN_CAPTCHA", true),
			"enable":           config.Env("ADMIN_ENABLE_AUTH", true),
			"permission":       config.Env("ADMIN_ENABLE_PERMISSION", true),
			"token_expiration": config.Env("ADMIN_TOKEN_EXPIRATION", 86400),
			// 使用的 auth guard，需要在 config/auth.go 中存在
			"guard": config.Env("ADMIN_AUTH_GUARD", "admin"),
			// 只放登录相关接口；不要在这里放业务路由，否则会绕过登录校验
			"except": []string{},
		},

		// 权限
		"permission": map[string]any{
			"enable": config.Env("ADMIN_ENABLE_PERMISSION", true),
			"except": []string{},
		},

		// 上传
		"upload": map[string]any{
			"disk": config.Env("ADMIN_UPLOAD_DISK", "public"),
			"directory": map[string]any{
				"image": "images",
				"file":  "files",
				"rich":  "rich",
			},
		},

		// 是否强制 HTTPS
		"https": config.Env("ADMIN_HTTPS", false),

		// 是否显示 [开发者工具] 菜单
		"show_development_tools": config.Env("ADMIN_SHOW_DEVELOPMENT_TOOLS", false),

		// 是否显示 [权限] 功能中的自动生成按钮（前端预留，暂未实现）
		"show_auto_generate_permission_button": config.Env("ADMIN_SHOW_AUTO_GENERATE_PERMISSION_BUTTON", false),

		// 布局
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
	}
}
