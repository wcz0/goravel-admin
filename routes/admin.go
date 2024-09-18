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
	// menu := admin.NewMenuController()
	user := admin.NewUserController()

	router := facades.Route()

	// router.Get("test", index.Test)

	// admin-api
	router.Prefix("admin-api").Group(func(router route.Router) {
		router.Get("login", auth.LoginPage)
		router.Post("login", auth.Login)
		router.Post("register", auth.Register)
		router.Post("_settings", index.SaveSetting)
		router.Get("_settings", index.GetSetting)
		router.Get("no-content", index.NoContext)
		router.Get("_download_export", index.DownloadExport)

		router.Middleware(middleware.Authenticate()).Group(func(router route.Router) {
			router.Post("upload_image", index.ImageUpload)
			router.Post("upload_file", index.FileUpload)
			router.Get("menus", index.GetMenus)
			router.Get("current-user", auth.CurrentUser)
			router.Post("upload_rich", index.RichFileUpload)
			router.Get("user_setting", user.GetUserSetting)
			router.Put("user_setting", user.PutUserSetting)

			router.Middleware(middleware.Permission()).Group(func(router route.Router) {

				router.Prefix("system").Group(func(router route.Router) {
					router.Resource("admin_users", admin.NewUserController())
					router.Resource("roles", admin.NewRoleController())
					router.Resource("permissions", admin.NewPermissionController())
					router.Resource("menus", admin.NewMenuController())


				})
			})
		})


	})

}
