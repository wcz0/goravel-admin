package admin

import "github.com/goravel/framework/contracts/http"

type MenuController struct {
}

func (m *MenuController) GetCurrentMenus() {

}

func NewMenuController() *MenuController {
	return &MenuController{}
}

func (m *MenuController) Index(ctx http.Context) http.Response {
	return nil
}

func (m *MenuController) Show(ctx http.Context) http.Response {
	return nil
}

func (m *MenuController) Store(ctx http.Context) http.Response {

	return nil
}

func (m *MenuController) Update(ctx http.Context) http.Response {

	return nil
}

func (m *MenuController) Destroy(ctx http.Context) http.Response {

	return nil
}