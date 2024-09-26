package core

import (
	"errors"
	"goravel/app/models/admin"
	"goravel/app/response"
	"strings"

	"github.com/goravel/framework/auth"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type Permission struct {
	authExcept       []string
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

/**
* 检查用户是否登录
* @return bool true 未登录 false 已登录
*/
func (p *Permission) AuthIntercept(ctx http.Context) {
	config := facades.Config()
	if !config.GetBool("admin.auth.enable") {
		return
	}
	// configExcept := config.Get("admin.auth.except").([]string)
	// mergedExcept := append(p.authExcept, p.permissionExcept...)
	// mergedExcept = append(mergedExcept, configExcept...)
	// 白名单处理
	// isExcept := false
	// for _, except := range mergedExcept {
	// 	formattedPath := p.pathFormatting(except)
	// 	if except == formattedPath {
	// 		isExcept = true
	// 		break
	// 	}
	// }
	// 用户登录
	token := ctx.Request().Header("Authorization")
	payload, err := facades.Auth(ctx).Guard("admin").Parse(token)
	if err != nil {
		if errors.Is(err, auth.ErrorTokenExpired) {
			ctx.Request().AbortWithStatusJson(http.StatusOK, response.TokenExpired)
			return
		}else{
			ctx.Request().AbortWithStatusJson(http.StatusOK, response.Unauthorized)
			return
		}
	}
	var user admin.AdminUser
	if err := facades.Orm().Query().Where("id", payload.Key).First(&user); err != nil {
		ctx.Request().AbortWithStatusJson(http.StatusOK, response.Unauthorized)
		return
	}
	ctx.WithValue("user", user)
}

/**
	检查用户状态
*/
func (p *Permission) CheckUserStatus(ctx http.Context) {
	user := ctx.Value("user")
	if user == nil {
		return
	}
	if user.(admin.AdminUser).Enabled == 0 {
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
	// 判断是否为白名单
	configExcept := config.Get("admin.permission.except").([]string)
	excepted := append(p.permissionExcept, configExcept...)
	excepted = append(excepted, p.authExcept...)
	if config.GetBool("admin.show_development_tools") {
		excepted = append(excepted, "/dev_tools*")
	}
	if len(excepted) == 0 {
		return false
	}
	// 白名单处理
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
	// 判断是否为超级管理员
	user := ctx.Value("user").(admin.AdminUser)
	if user.IsAdministrator() {
		return false
	}
	allPermissions := user.AllPermissions()
	for _, permission := range allPermissions {
		if !permission.ShouldPassThrough(ctx) {
			return false
		}
	}
	return true
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


