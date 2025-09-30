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
	configExcept := config.Get("admin.auth.except").([]string)
	mergedExcept := append(p.authExcept, configExcept...)
	// 白名单处理
	for _, except := range mergedExcept {
		if p.pathMatches(ctx, except) {
			return
		}
	}
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
		if ctx.Request().Path() == p.pathFormatting(except) {
			isExcept = true
			break
		}
	}
	if isExcept {
		return false
	}
	// 判断是否为超级管理员
	userValue := ctx.Value("user")
	if userValue == nil {
		return true // 用户未登录，拒绝访问
	}

	user, ok := userValue.(admin.AdminUser)
	if !ok {
		return true // 类型断言失败，拒绝访问
	}

	if user.IsAdministrator() {
		return false // 超级管理员，允许访问
	}

	allPermissions := user.AllPermissions()
	if len(allPermissions) == 0 {
		return true // 没有权限，拒绝访问
	}

	// 检查是否有匹配的权限
	for _, permission := range allPermissions {
		if permission.ShouldPassThrough(ctx) {
			return false // 找到匹配权限，允许访问
		}
	}
	return true // 没有匹配的权限，拒绝访问
}


func (p *Permission) pathMatches(ctx http.Context, except string) bool {
	path := ctx.Request().Path()
	// 去掉前缀进行匹配
	prefix := "/" + strings.Trim(facades.Config().GetString("admin.route.prefix"), "/")
	path = strings.TrimPrefix(path, prefix)
	path = strings.Trim(path, "/")
	except = strings.Trim(except, "/")

	if except == "" {
		return path == ""
	}
	return path == except
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


