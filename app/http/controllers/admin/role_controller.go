package admin

import "github.com/goravel/framework/contracts/http"

type RoleController struct {
}

func NewRoleController() *RoleController {
	return &RoleController{}
}

func (r *RoleController) Index(ctx http.Context) http.Response {

	return nil
}

func (r *RoleController) Show(ctx http.Context) http.Response {
	return nil
}

func (r *RoleController) Store(ctx http.Context) http.Response {
	return nil
}

func (r *RoleController) Update(ctx http.Context) http.Response {

	return nil
}

func (r *RoleController) Destroy(ctx http.Context) http.Response {

	return nil
}