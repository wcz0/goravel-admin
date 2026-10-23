// Package admin 把 goravel-admin 作为可安装的 Goravel 扩展包对外提供。
//
// 安装方式：
//
//	go run github.com/wcz0/goravel-admin/setup install
//
// 安装会注册 ServiceProvider 与后台数据表迁移，随后执行
//
//	go run . artisan migrate
//	go run . artisan admin:init
//
// 即可访问 /admin。
package admin

import (
	"io/fs"
	nethttp "net/http"
	"strings"

	contractsconfig "github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/foundation"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/path"
	"github.com/spf13/cast"

	"github.com/wcz0/goravel-admin/commands"
	adminconfig "github.com/wcz0/goravel-admin/config"
	"github.com/wcz0/goravel-admin/core"
	adminmigrations "github.com/wcz0/goravel-admin/migrations"
	"github.com/wcz0/goravel-admin/routes"
	"github.com/wcz0/goravel-admin/tools"
)

// Module 是扩展包的模块路径，同时也是 vendor:publish 的包名。
const Module = "github.com/wcz0/goravel-admin"

// ServiceProvider 注册并启动 goravel-admin。
type ServiceProvider struct{}

// Register 合并默认配置、注册菜单绑定与多语言资源。
func (r *ServiceProvider) Register(app foundation.Application) {
	// 注册菜单
	app.Instance("admin.menu", core.NewMenu())

	config := app.MakeConfig()

	// 深度合并默认配置，宿主已设置的值优先
	mergeConfig(config, adminconfig.New())

	// 确保 auth guard / provider 存在，后台才能开箱即用
	mergeAuthConfig(config)

	// 内嵌多语言文件作为文件加载器的兜底
	mergeLangFS(config)

	// 可选：把语言文件和前端资源发布到宿主目录，便于自定义或交给 Nginx 托管
	app.Publishes(Module, map[string]string{
		"lang": path.Lang(),
	}, "goravel-admin-lang")
	app.Publishes(Module, map[string]string{
		"assets/dist": path.Public("admin"),
	}, "goravel-admin-assets")
}

// Boot 注册后台路由、静态资源、迁移与命令。
func (r *ServiceProvider) Boot(app foundation.Application) {
	facades.Route().Group(func(router route.Router) {
		routes.Admin()
	})

	// 迁移由扩展包自注册。宿主不必改动 bootstrap/migrations.go，
	// 也避免了包内 migrations 包与宿主 migrations 包同名冲突。
	facades.Schema().Register(adminmigrations.All())

	r.registerAssets()

	app.MakeArtisan().Register([]console.Command{
		&commands.Admin{},
		&commands.DumpAdminData{},
		&commands.CreateUser{},
		&commands.ResetPassword{},
		&commands.Doctor{},
	})
}

// registerAssets 挂载内嵌的前端构建产物与后台入口页面。
func (r *ServiceProvider) registerAssets() {
	router := facades.Route()
	router.StaticFS(tools.AdminAssetsPath(), nethttp.FS(DistFS()))

	page := tools.AdminPagePath()
	router.Get(page, r.adminPage)
	router.Get(page+"/", r.adminPage)
}

// adminPage 输出后台入口 HTML，并把构建时的资源前缀替换成当前配置。
func (r *ServiceProvider) adminPage(ctx contractshttp.Context) contractshttp.Response {
	content, err := fs.ReadFile(DistFS(), "index.html")
	if err != nil {
		return ctx.Response().String(contractshttp.StatusInternalServerError, "goravel-admin: index.html not found")
	}

	html := string(content)
	html = strings.ReplaceAll(html, "/admin-assets", tools.AdminAssetsPath())
	html = strings.ReplaceAll(html, `window.$adminApiPrefix = '/admin-api'`,
		`window.$adminApiPrefix = '/`+tools.AdminApiPrefix()+`'`)

	return ctx.Response().
		Header("Cache-Control", "no-cache, no-store, must-revalidate").
		Data(contractshttp.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// mergeConfig 把默认配置深度合并到宿主的 admin 配置中，宿主值优先。
func mergeConfig(config contractsconfig.Config, defaults map[string]any) {
	config.Add("admin", deepMerge(defaults, cast.ToStringMap(config.Get("admin"))))
}

func deepMerge(defaults, overrides map[string]any) map[string]any {
	merged := make(map[string]any, len(defaults))
	for key, value := range defaults {
		merged[key] = value
	}

	for key, value := range overrides {
		defaultMap, defaultIsMap := merged[key].(map[string]any)
		overrideMap, overrideIsMap := value.(map[string]any)
		if defaultIsMap && overrideIsMap {
			merged[key] = deepMerge(defaultMap, overrideMap)
			continue
		}
		merged[key] = value
	}

	return merged
}

// mergeAuthConfig 保证 admin guard / provider 已配置，宿主可自行覆盖。
func mergeAuthConfig(config contractsconfig.Config) {
	auth := cast.ToStringMap(config.Get("auth"))

	guards := cast.ToStringMap(auth["guards"])
	if _, ok := guards["admin"]; !ok {
		guards["admin"] = map[string]any{
			"driver":   "jwt",
			"provider": "admin",
		}
	}
	auth["guards"] = guards

	providers := cast.ToStringMap(auth["providers"])
	if _, ok := providers["admin"]; !ok {
		providers["admin"] = map[string]any{
			"driver": "orm",
			"model":  "github.com/wcz0/goravel-admin/models/admin.AdminUser",
		}
	}
	auth["providers"] = providers

	config.Add("auth", auth)
}

// mergeLangFS 把内嵌语言文件接入翻译器。文件加载器（宿主 lang 目录）优先级更高。
func mergeLangFS(config contractsconfig.Config) {
	if existing, ok := config.Get("app.lang_fs").(fs.FS); ok && existing != nil {
		config.Add("app.lang_fs", multiFS{existing, LangFS()})
		return
	}

	config.Add("app.lang_fs", LangFS())
}

// multiFS 按顺序查找多个文件系统，靠前的优先。
type multiFS []fs.FS

func (m multiFS) Open(name string) (fs.File, error) {
	for _, fsys := range m {
		file, err := fsys.Open(name)
		if err == nil {
			return file, nil
		}
	}

	return nil, fs.ErrNotExist
}
