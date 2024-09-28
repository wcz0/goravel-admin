package providers

import (
	"goravel/app/support/core"

	"github.com/goravel/framework/contracts/foundation"
)

type AdminServiceProvider struct {
}

func (receiver *AdminServiceProvider) Register(app foundation.Application) {
	// 注册菜单
	app.Instance("admin.menu", core.NewMenu())
}

func (receiver *AdminServiceProvider) Boot(app foundation.Application) {

}
