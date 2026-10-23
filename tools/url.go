package tools

import (
	"strings"

	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
)

// AppUrl 返回不带结尾斜杠的 APP_URL，未配置时为空字符串。
func AppUrl() string {
	return strings.TrimSuffix(cast.ToString(facades.Config().Env("APP_URL")), "/")
}

// AdminApiPrefix 返回后台 API 前缀，不含前导和结尾斜杠。
func AdminApiPrefix() string {
	return strings.Trim(cast.ToString(facades.Config().Get("admin.route.prefix", "admin-api")), "/")
}

// AdminPagePath 返回后台页面挂载路径，以 / 开头且不含结尾斜杠。
func AdminPagePath() string {
	return normalizeMountPath(cast.ToString(facades.Config().Get("admin.route.path", "/admin")), "/admin")
}

// AdminAssetsPath 返回后台静态资源挂载路径，以 / 开头且不含结尾斜杠。
func AdminAssetsPath() string {
	return normalizeMountPath(cast.ToString(facades.Config().Get("admin.route.assets_path", "/admin-assets")), "/admin-assets")
}

// AuthGuard 返回后台使用的 auth guard 名称，默认 admin。
//
// 登录、注销与中间件校验必须统一走这个 guard，否则签发的 token
// 会落到宿主的默认 guard 上，导致带 token 的请求被判为未认证。
func AuthGuard() string {
	guard := strings.TrimSpace(cast.ToString(facades.Config().Get("admin.auth.guard", "admin")))
	if guard == "" {
		return "admin"
	}

	return guard
}

func normalizeMountPath(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "/" {
		return fallback
	}
	if !strings.HasPrefix(value, "/") {
		value = "/" + value
	}

	return strings.TrimSuffix(value, "/")
}

func GetAdmin(str string) string {
	if len(str) == 0 || str[0] != '/' {
		str = "/" + str
	}
	base := AppUrl()
	prefix := "/" + AdminApiPrefix()
	if strings.HasSuffix(base, prefix) || strings.Contains(base, prefix+"/") {
		return base + str
	}

	return base + prefix + str
}

func GetAdminNil() string {
	base := AppUrl()
	prefix := "/" + AdminApiPrefix()
	if strings.HasSuffix(base, prefix) || strings.Contains(base, prefix+"/") {
		return base
	}

	return base + prefix
}

func GetApi(str string) string {
	if len(str) == 0 || str[0] != '/' {
		str = "/" + str
	}

	return AppUrl() + "/api" + str
}

func GetApiNil() string {
	return AppUrl() + "/api"
}

func Url(str string) string {
	return AppUrl() + str
}

func UrlNil() string {
	return AppUrl()
}

// GetImageUrl 上传图片地址
func GetImageUrl() string {
	return AppUrl() + "/" + AdminApiPrefix() + "/upload_image"
}

// GetFileUrl 上传文件地址
func GetFileUrl() string {
	return AppUrl() + "/" + AdminApiPrefix() + "/upload_file"
}
