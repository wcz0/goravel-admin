package routes

import (
	"net/http"

	"github.com/goravel/framework/facades"
)

func Web() {
	// router := facades.Route()
	// router.Static("/admin", "./admin")
	facades.Route().StaticFile("/admin", "./public/admin/index.html")

	// // facades.Route().Static("/", "./public")
	// // router.StaticFile("/", "./public/index.html")

	// router.Get("/admin", func(c http.Context) {
	//     // 重定向到 SPA 的入口页面
	//     c.Response().Redirect(200, "/admin/index.html")
	// })

	facades.Route().StaticFS("/admin-assets/assets", http.Dir("./public/admin/assets/"))
	facades.Route().StaticFS("/admin-assets/scripts", http.Dir("./public/admin/scripts/"))
	// facades.Route().StaticFile("static-file", "./public/logo.png")
	// facades.Route().StaticFS("/admin", http.Dir("./public/admin/"))
	// facades.Route().Get("/admin", func(c ghttp.Context) ghttp.Response {
	// 	return c.Response().File(path.Public("admin") + "/index.html")
	// })
}
