package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers/admin"
)

func Admin() {
	indexController := admin.NewIndexController()
	authController := admin.NewAuthController()

	// admin-api
	facades.Route().Prefix("admin-api").Group(func(router route.Router) {
		router.Get("login", authController.LoginPage)
		router.Post("login", authController.Login)
		router.Post("register", authController.Register)
		router.Post("_settings", indexController.SaveSetting)
		router.Get("_settings", indexController.GetSetting)
	})
}
