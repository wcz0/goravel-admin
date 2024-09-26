package services

import (
	"goravel/app/models/admin"

	"github.com/goravel/framework/facades"
)

type AdminSettingService struct {
	*Service
	cacheKey string
}

func NewAdminSettingService() *AdminSettingService {
	return &AdminSettingService{
		Service:  NewService(),
		cacheKey: "app_setting_",
	}
}

func (a *AdminSettingService) Get(key string, default_ any, fresh bool) any {
	var adminSetting admin.AdminSetting
	if fresh {
		value := facades.Orm().Query().Where("key", key).Select("values").First(&adminSetting)
		if value == nil {
			return default_
		}
	}
	value, err := facades.Cache().RememberForever(a.cacheKey+key, func() (any, error) {
		err := facades.Orm().Query().Where("key", key).Select("values").First(&adminSetting)
		if err != nil {
			return nil, err
		}
		return adminSetting.Values, nil
	})
	if err != nil {
		return default_
	}
	if value != nil {
		return value
	} else {
		return default_
	}
}

func (a *AdminSettingService) SetMany(array map[string]any) bool {
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return false
	}
	for key, value := range array {
		if !a.Set(key, value.(string)) {
			if err = tx.Rollback(); err != nil {
				return false
			}
			return false
		}
	}
	if err = tx.Commit(); err != nil {
		return false
	}
	return true
}

func (a *AdminSettingService) Set(key string, value string) bool {
	var adminSetting admin.AdminSetting
	if err := facades.Orm().Query().Where("key", key).FirstOrNew(&adminSetting, admin.AdminSetting{Key: key, Values: value}); err != nil {
		return false
	}
	if err := facades.Orm().Query().Save(&adminSetting); err != nil {
		return false
	}
	a.clearCache(key)
	return true
}

func (a *AdminSettingService) clearCache(key string) {
	facades.Cache().Forget(a.getCacheKey(key))
}

func (a *AdminSettingService) getCacheKey(key string) string {
	return a.cacheKey + key
}
