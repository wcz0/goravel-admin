package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers/admin"
	"goravel/app/http/middleware"
)

func Admin() {
	index := admin.NewIndexController()
	auth := admin.NewAuthController()
	user := admin.NewUserController()
	home := admin.NewHomeController()

	router := facades.Route()

	// admin-api
	router.Prefix("admin-api").Middleware(middleware.AdminLang()).Group(func(router route.Router) {
		router.Get("login", auth.LoginPage)
		router.Post("login", auth.Login)
		router.Get("logout", auth.Logout)
		router.Post("register", auth.Register)
		router.Post("_settings", index.SaveSetting)
		router.Get("_settings", index.GetSetting)
		router.Get("no-content", index.NoContext)
		router.Get("_download_export", index.DownloadExport)

		router.Middleware(middleware.Authenticate()).Group(func(router route.Router) {
			router.Post("upload_image", index.ImageUpload)
			router.Post("upload_file", index.FileUpload)
			router.Post("upload_rich", index.RichFileUpload)
			router.Get("menus", index.GetMenus)
			router.Get("current-user", auth.CurrentUser)
			router.Get("user_setting", user.GetUserSetting)
			router.Put("user_setting", user.PutUserSetting)

			// router.Resource("dashboard", admin.NewHomeController())

			router.Get("dashboard", home.Index)
			router.Middleware(middleware.Permission()).Group(func(router route.Router) {
				router.Prefix("system").Group(func(router route.Router) {
					router.Get("users", user.Index)
					router.Resource("admin_users", admin.NewUserController())
					router.Resource("admin_roles", admin.NewRoleController())
					router.Get("admin_roles/options", admin.NewRoleController().Options)
                    router.Resource("permissions", admin.NewPermissionController())
                    router.Get("permissions/options", admin.NewPermissionController().Options)
                    router.Get("permissions/tree_options", admin.NewPermissionController().TreeOptions)
                    router.Resource("admin_permissions", admin.NewPermissionController())
                    router.Post("_admin_permissions_auto_generate", admin.NewPermissionController().AutoGenerate)
                    router.Resource("menus", admin.NewMenuController())
                    router.Resource("admin_menus", admin.NewMenuController())
                    router.Post("admin_menus/save_order", admin.NewMenuController().SaveOrder)
					router.Post("admin_roles/save_permissions", admin.NewRoleController().SavePermissions)
				})
			})
		})
	})
}
