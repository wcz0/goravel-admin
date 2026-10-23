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
	pages := admin.NewPagesController()
	apis := admin.NewApisController()
	relationships := admin.NewRelationshipsController()
	extensions := admin.NewExtensionsController()

	config := facades.Config()

	// 需要在登录后访问的路由
	protected := func(router route.Router) {
		router.Post("upload_image", index.ImageUpload)
		router.Post("upload_file", index.FileUpload)
		router.Post("upload_rich", index.RichFileUpload)
		router.Any("upload_chunk_start", index.ChunkUploadStart)
		router.Any("upload_chunk", index.ChunkUpload)
		router.Any("upload_chunk_finish", index.ChunkUploadFinish)
		router.Get("menus", index.GetMenus)
		router.Get("current-user", auth.CurrentUser)
		router.Get("user_setting", user.GetUserSetting)
		router.Put("user_setting", user.PutUserSetting)
		router.Get("dashboard", home.Index)
		router.Get("page_schema", index.PageSchema)
		router.Get("_iconify_search", index.IconifySearch)

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

			// 开发者工具，仅在配置开启时注册，避免多暴露接口
			if config.GetBool("admin.show_development_tools") {
				router.Prefix("dev_tools").Group(func(router route.Router) {
					router.Resource("pages", pages)
					router.Get("pages/options", pages.Options)
					router.Resource("apis", apis)
					router.Get("api/templates", apis.Templates)
					router.Get("api/args_schema", apis.ArgsSchema)
					router.Resource("relationships", relationships)
					router.Get("relation/types", relationships.Types)
					router.Resource("extensions", extensions)
					router.Post("extensions/enable", extensions.Enable)
					router.Post("extensions/disable", extensions.Disable)
				})
			}
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
		router.Get("captcha", auth.Captcha)
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
