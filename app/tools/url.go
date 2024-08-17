package tools

import "github.com/goravel/framework/facades"

var url string
var adminUrl string

func init() {
	config := facades.Config()
	url = config.Env("APP_URL").(string)
	adminUrl = config.GetString("admin.route.prefix")
}

func GetAdmin(str string) string {
	return url + adminUrl + str
}

func GetAdminNil() string {
	return url + adminUrl
}

func GetApi(str string) string {
	return url + "/api" + str
}

func GetApiNil() string {
	return url + "/api"
}

func GetUrl(str string) string {
	return url + str
}

func GetUrlNil() string {
	return url
}

/**
 * 上传图片地址
 */
func GetImageUrl() string {
	return url + adminUrl + "/upload_image"
}

/**
 * 上传文件地址
 */
func GetFileUrl() string {
	return url + adminUrl + "/upload_file"
}
