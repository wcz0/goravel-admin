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

func (p *Permission) CheckUserStatus(ctx http.Context) bool {


func (p *Permission) PermissionIntercept(ctx http.Context) bool {
	return false
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