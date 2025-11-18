package services

import (
	"encoding/json"
	"errors"
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

/**
 * 获取配置
 * @param key 键
 * @param default_ 默认值
 * @param fresh 是否强制从数据库获取
 * @return string
 */
func (a *AdminSettingService) Get(key string, default_ any, fresh bool) any {
	var adminSetting admin.AdminSetting
	if fresh {
		err := facades.Orm().Query().Where("key", key).Select("values").First(&adminSetting)
		if err != nil {
			return default_
		}
		var result any
		if err := json.Unmarshal([]byte(adminSetting.Values), &result); err != nil {
			return default_
		}
		return result
	}

    var value any
    var err error
    func() {
        defer func() {
            if r := recover(); r != nil {
                err = errors.New("cache unavailable")
            }
        }()
        value, err = facades.Cache().RememberForever(a.cacheKey+key, func() (any, error) {
            queryErr := facades.Orm().Query().Where("key", key).Select("values").First(&adminSetting)
            if queryErr != nil {
                return nil, queryErr
            }
            if err != nil {
                return nil, err
            }
            return adminSetting.Values, nil
        })
    }()

	if err != nil {
		return default_
	}

    if value != nil {
        var result any
        if err := json.Unmarshal([]byte(value.(string)), &result); err != nil {
            return default_
        }
        return result
    } else {
        return default_
    }
}

func (a *AdminSettingService) SetMany(array map[string]any) error {
	for key, value := range array {
		// Marshal the value to a json string
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return err
		}

		var adminSetting admin.AdminSetting
		// 先查找是否存在
		err = facades.Orm().Query().Where("key", key).
			UpdateOrCreate(&adminSetting, admin.AdminSetting{Key: key}, admin.AdminSetting{Key: key, Values: string(jsonBytes)})
		if err != nil {
			return err
		}
		// 更新缓存
    func() {
        defer func() { _ = recover() }()
        _ = facades.Cache().Forever(a.cacheKey+key, string(jsonBytes))
    }()
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
    func() { defer func() { _ = recover() }(); facades.Cache().Forget(a.getCacheKey(key)) }()
}

func (a *AdminSettingService) getCacheKey(key string) string {
	return a.cacheKey + key
}
