package tools

import "github.com/goravel/framework/facades"


var url string

func init() {
	url = facades.Config().Env("APP_URL", "http://localhost").(string)
}

func  GetAdmin(str string) string {
	return url + "/admin-api" + str
}

func GetAdminNil() string {
	return url + "/admin-api";
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
	return url + "/admin-api/upload_image"
}

/**
 * 上传文件地址
 */
func GetFileUrl() string {
	return url + "/admin-api/upload_file"
}
