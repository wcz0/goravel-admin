package routes

import (
	"github.com/goravel/framework/facades"
)

func Web() {
	// 后台的页面与静态资源由 goravel-admin 的 ServiceProvider 自动挂载：
	//   /admin          后台入口
	//   /admin-assets   前端构建产物

	// public files
	facades.Route().Static("storage", "./storage/app/public")
}
