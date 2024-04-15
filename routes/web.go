package routes

import (
	"github.com/goravel/framework/facades"
)

func Web() {
	facades.Route().Static("/admin", "./public/admin")
	facades.Route().StaticFile("/admin", "./public/admin/index.html")

	facades.Route().Static("/", "./public")
	facades.Route().StaticFile("/", "./public/index.html")
}
