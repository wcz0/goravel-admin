package admin

import (
	"goravel/app/http/controllers"
	"github.com/goravel/framework/contracts/http"


)

type IndexController struct {
	*controllers.Controller
}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (i *IndexController) SaveSetting(ctx http.Context) http.Response  {
	return i.Success(ctx)

}

func (i *IndexController) GetSetting(ctx http.Context) http.Response {
	return i.Success(ctx)
}

func (i *IndexController) NoContext() {

}

func (i *IndexController) ImageUpload() {

}

func (i *IndexController) FileUpload() {

}

func (i *IndexController) RichFileUpload() {

}

func (i *IndexController) GetUserSetting() {

}

func (i *IndexController) PutUserSetting() {

}

func (i *IndexController) GetCurrentUser() {

}

func (i *IndexController) SearchIcon() {

}

func (i *IndexController) GetDashBoard() {

}