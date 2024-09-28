package services

import (
	"goravel/app/models/admin"
	"github.com/duke-git/lancet/v2/convertor"


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

/**
 * 获取配置
 * @param key 键
 * @param default_ 默认值
 * @param fresh 是否强制从数据库获取
 * @return string
 */
func (a *AdminSettingService) Get(key string, default_ string, fresh bool) string {
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
		return value.(string)
	} else {
		return default_
	}
}

func (a *AdminSettingService) SetMany(array map[string]any) error {
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return err
	}
	for key, value := range array {
		str := convertor.ToString(value)
		if err := a.Set(key, str); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (a *AdminSettingService) Set(key string, value string) error {
	var adminSetting admin.AdminSetting
	if err := facades.Orm().Query().Where("key", key).FirstOrNew(&adminSetting, admin.AdminSetting{Key: key, Values: value}); err != nil {
		return err
	}
	if err := facades.Orm().Query().Save(&adminSetting); err != nil {
		return err
	}
	a.clearCache(key)
	return nil
}

func (a *AdminSettingService) clearCache(key string) {
	facades.Cache().Forget(a.getCacheKey(key))
}

func (a *AdminSettingService) getCacheKey(key string) string {
	return a.cacheKey + key
}
