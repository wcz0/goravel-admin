package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	admin "github.com/wcz0/goravel-admin/controllers"
	"github.com/wcz0/goravel-admin/middleware"
	"github.com/wcz0/goravel-admin/tools"
)

// Admin 注册后台 API 路由（默认前缀 /admin-api）。
func Admin() {
	index := admin.NewIndexController()
	auth := admin.NewAuthController()
	user := admin.NewUserController()
	home := admin.NewHomeController()
	role := admin.NewRoleController()
	permission := admin.NewPermissionController()
	menu := admin.NewMenuController()

	config := facades.Config()

	// 需要在登录后访问的路由
	protected := func(router route.Router) {
		router.Post("upload_image", index.ImageUpload)
		router.Post("upload_file", index.FileUpload)
		router.Post("upload_rich", index.RichFileUpload)
		router.Get("menus", index.GetMenus)
		router.Get("current-user", auth.CurrentUser)
		router.Get("user_setting", user.GetUserSetting)
		router.Put("user_setting", user.PutUserSetting)
		router.Get("dashboard", home.Index)

		// 需要在拥有权限后访问的路由
		authorized := func(router route.Router) {
			router.Prefix("system").Group(func(router route.Router) {
				router.Get("users", user.Index)
				router.Resource("admin_users", admin.NewUserController())
				router.Resource("admin_roles", role)
				router.Get("admin_roles/options", role.Options)
				router.Post("admin_roles/save_permissions", role.SavePermissions)
				router.Resource("permissions", permission)
				router.Get("permissions/options", permission.Options)
				router.Get("permissions/tree_options", permission.TreeOptions)
				router.Resource("admin_permissions", permission)
				router.Post("_admin_permissions_auto_generate", permission.AutoGenerate)
				router.Resource("menus", menu)
				router.Resource("admin_menus", menu)
				router.Post("admin_menus/save_order", menu.SaveOrder)
			})
		}

		if config.GetBool("admin.permission.enable", true) {
			router.Middleware(middleware.Permission()).Group(authorized)
			return
		}

		authorized(router)
	}

	router := facades.Route().Prefix(tools.AdminApiPrefix()).Middleware(middleware.AdminLang())
	router.Group(func(router route.Router) {
		router.Get("login", auth.LoginPage)
		router.Post("login", auth.Login)
		router.Get("logout", auth.Logout)
		router.Post("register", auth.Register)
		router.Post("_settings", index.SaveSetting)
		router.Get("_settings", index.GetSetting)
		router.Get("no-content", index.NoContext)
		router.Get("_download_export", index.DownloadExport)

		if config.GetBool("admin.auth.enable", true) {
			router.Middleware(middleware.Authenticate()).Group(protected)
			return
		}

		protected(router)
	})
}
