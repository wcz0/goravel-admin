package routes

import (
    "github.com/goravel/framework/contracts/route"
    "github.com/goravel/framework/facades"
    api "goravel/app/http/controllers/api"
)

func Api() {
    user := api.NewUserController()
    facades.Route().Prefix("api").Group(func(r route.Router) {
        r.Get("user", user.Show)
    })
}
