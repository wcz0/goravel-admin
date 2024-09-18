package tools

import (
	"github.com/goravel/framework/facades"
)

func GetAdmin(str string) string {
	config := facades.Config()
	return config.Env("APP_URL").(string) + "/" + config.GetString("admin.route.prefix") + str
}

func GetAdminNil() string {
	config := facades.Config()
	return config.Env("APP_URL").(string) + "/" + config.GetString("admin.route.prefix")
}

func GetApi(str string) string {
	config := facades.Config()
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
