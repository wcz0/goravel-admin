package main

// Stubs 提供安装时写入宿主的文件模板。
type Stubs struct{}

// Config 返回 config/admin.go 的覆盖模板。
func (s Stubs) Config() string {
	return `package config

// goravel-admin 的默认配置由扩展包提供，见 github.com/wcz0/goravel-admin/config。
//
// 需要覆盖时，在下面补一个 init 函数，只写想改动的项即可，
// 扩展包会在 ServiceProvider.Register 阶段把缺失的默认值补齐：
//
//	import "github.com/goravel/framework/facades"
//
//	func init() {
//		facades.Config().Add("admin", map[string]any{
//			"name": "My Admin",
//			"layout": map[string]any{
//				"title": "%title% | My Admin",
//			},
//		})
//	}
//
// 也可以只通过环境变量覆盖，例如 ADMIN_ROUTE_PREFIX、ADMIN_SITE_TITLE。
`
}
