package admin

import (
	"goravel/app/http/controllers"

	"github.com/goravel/framework/contracts/http"
)

type UserController struct {
	*controllers.Controller
}

func NewUserController() *UserController {
	return &UserController{
		Controller: controllers.NewController(),
	}
}

func (u *UserController) PutUserSetting(c http.Context) http.Response {
	return u.Success(c)
}


func (r *UserController) Index(ctx http.Context) http.Response {
	return nil
}

func (r *UserController) Show(ctx http.Context) http.Response {
	return nil
}

func (r *UserController) Store(ctx http.Context) http.Response {
	return nil
}

func (r *UserController) Update(ctx http.Context) http.Response {
	return nil
}

func (r *UserController) Destroy(ctx http.Context) http.Response {
	return nil
}