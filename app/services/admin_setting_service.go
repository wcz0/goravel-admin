package services

import (
	"goravel/app/models"

	"github.com/goravel/framework/facades"
)

type AdminSettingService struct {
	*Service
	cacheKey string
}

func NewAdminSettingService() *AdminSettingService {
	return &AdminSettingService{
		Service: NewService(),
		cacheKey: "app_setting_",
	}
}

func (a *AdminSettingService) Get(key string, default_ any, fresh bool ) any {
	var adminSetting models.AdminSetting
	if fresh {
		value := facades.Orm().Query().Where("key", key).Select("values").First(&adminSetting)
		if value == nil {
			return default_
		}
	}
	value, err := facades.Cache().RememberForever(a.cacheKey + key, func() (any, error) {
		err := facades.Orm().Query().Where("key", key).Select("values").First(&adminSetting)
		if err!=nil {
			return nil, err
		}
		return adminSetting.Values, nil
	})
	if err != nil {
		return default_
	}
	if value != nil {
		return value
	}else {
		return default_
	}
}