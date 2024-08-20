package core

import (
	"goravel/app/models"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type Permission struct {
	authExcept []string
	permissionExcept []string
	// adminConfig map[string]any
	// routePrefix string
	// showDevTools bool
}

func NewPermission() *Permission {
	return &Permission{
		authExcept: []string{
			"login",
			"logout",
			"no-content",
			"_settings",
			"captcha",
			"_download_export",
		},
		permissionExcept: []string{
			"menus",
			"current-user",
			"user_setting",
			"login",
			"logout",
			"no-content",
			"_settings",
			"upload_image",
			"upload_file",
			"upload_rich",
			"captcha",
			"_download_export",
		},
	}
}

func (p *Permission) AuthIntercept(ctx http.Context) bool {
	config := facades.Config()
	if !config.GetBool("admin.auth.enable") {
		return false
	}
	configExcept := config.Get("admin.auth.except").([]string)
	mergedExcept := append(p.authExcept, p.permissionExcept...)
	mergedExcept = append(mergedExcept, configExcept...)
	isExcept := false
	for _, except := range mergedExcept {
		formattedPath := p.pathFormatting(except)
		if except == formattedPath {
			isExcept = true
			break
		}
	}
	var user models.AdminUser
	err := facades.Auth(ctx).Guard("admin").User(&user)
	if err != nil {
		return false
	}
	ctx.WithValue("user", user)
	return !isExcept
}

func (p *Permission) CheckUserStatus(ctx http.Context) {
	var user models.AdminUser
	err := facades.Auth(ctx).Guard("admin").User(&user)
	if err != nil {
		facades.Auth(ctx).Logout()
	}
	if user.Enabled == 0 {
		facades.Auth(ctx).Logout()
	}
}

/**
	权限拦截
*/
func (p *Permission) PermissionIntercept(ctx http.Context) bool {
	config := facades.Config()
	if !config.GetBool("admin.permission.enable") {
		return false
	}
	if ctx.Request().Path() == config.GetString("admin.route.prefix") {
		return false
	}
	configExcept := config.Get("admin.permission.except").([]string)
	excepted := append(p.permissionExcept, configExcept...)
	excepted = append(excepted, p.authExcept...)
	if config.GetBool("admin.show_development_tools") {
		excepted = append(excepted, "/dev_tools*")
	}
	if len(excepted) == 0 {
		return false
	}

	isExcept := false
	for _, except := range excepted {
		formattedPath := p.pathFormatting(except)
		if p.pathMatches(ctx, formattedPath) {
			isExcept = true
			break
		}
	}
	if isExcept {
		return false
	}
	user := ctx.Value("user")





	return false
}

func (p *Permission) pathMatches(ctx http.Context, except string) bool {
    path := ctx.Request().Path()
    if except == "/" {
        return path == except
    }
    return path == strings.Trim(except, "/")
}

func (p *Permission) pathFormatting(path string) string {
	prefix := "/" + strings.Trim(facades.Config().GetString("admin.route.prefix"), "/")
	if prefix == "/" {
		prefix = ""
	}
	path = strings.Trim(path, "/")
	if path == "" {
		path = "index"
	}
	return prefix + "/" + path
}

func (p *Permission) checkRoutePermission(ctx http.Context) bool {


	args := strings.Split(middleware[len(middlewarePrefix):], ",")

	method := args[0]
	args = args[1:]

	if !hasMethod(AdminPermissionModel(), method) {
		panic("Invalid permission method [" + method + "].")
	}

	callMethod(AdminPermissionModel(), method, args)

	return true
}