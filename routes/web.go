package routes

import (
	"net/http"
	"github.com/goravel/framework/facades"
)

func Web() {
	facades.Route().StaticFile("/admin", "./public/admin/index.html")
	facades.Route().StaticFS("/admin", http.Dir("./public/admin/"))
	facades.Route().StaticFS("/admin-assets/assets", http.Dir("./public/admin/assets/"))
	facades.Route().StaticFS("/admin-assets/scripts", http.Dir("./public/admin/scripts/"))
}
