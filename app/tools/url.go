package tools

import (
    "strings"
    "github.com/goravel/framework/facades"
)

func GetAdmin(str string) string {
    config := facades.Config()
    if len(str) == 0 || str[0] != '/' {
        str = "/" + str
    }
    base := config.Env("APP_URL").(string)
    prefix := "/" + config.GetString("admin.route.prefix")
    if strings.HasSuffix(base, prefix) || strings.Contains(base, prefix+"/") {
        return base + str
    }
    return base + prefix + str
}

func GetAdminNil() string {
    config := facades.Config()
    base := config.Env("APP_URL").(string)
    prefix := "/" + config.GetString("admin.route.prefix")
    if strings.HasSuffix(base, prefix) || strings.Contains(base, prefix+"/") {
        return base
    }
    return base + prefix
}

func GetApi(str string) string {
    config := facades.Config()
    if len(str) == 0 || str[0] != '/' {
        str = "/" + str
    }
    return config.Env("APP_URL").(string) + "/api" + str
}

func GetApiNil() string {
	config := facades.Config()
	return config.Env("APP_URL").(string) + "/api"
}

func Url(str string) string {
	config := facades.Config()
	return config.Env("APP_URL").(string) + str
}

func UrlNil() string {
	config := facades.Config()
	return config.Env("APP_URL").(string)
}

/**
 * 上传图片地址
 */
func GetImageUrl() string {
	config := facades.Config()
	return config.Env("APP_URL").(string) + "/" + config.GetString("admin.route.prefix") + "/upload_image"
}

/**
 * 上传文件地址
 */
func GetFileUrl() string {
	config := facades.Config()
	return config.Env("APP_URL").(string) + "/" + config.GetString("admin.route.prefix") + "/upload_file"
}
