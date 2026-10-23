package services

import (
	"errors"
	"regexp"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/goravel-admin/models/admin"
)

type AdminExtensionService struct {
	*Service
}

func NewAdminExtensionService() *AdminExtensionService {
	return &AdminExtensionService{Service: NewService()}
}

// extensionNamePattern 限制扩展名形态为 vendor/name。
var extensionNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+$`)

// EnabledNames 返回已启用的扩展名，供 _settings 下发给前端。
func (s *AdminExtensionService) EnabledNames() []string {
	var extensions []admin.AdminExtension
	if err := facades.Orm().Query().
		Where("is_enabled", admin.Extension_Enabled_ON).
		Find(&extensions); err != nil {
		return []string{}
	}

	names := make([]string, 0, len(extensions))
	for _, e := range extensions {
		names = append(names, e.Name)
	}

	return names
}

// List 返回扩展列表。
func (s *AdminExtensionService) List(ctx http.Context) http.Response {
	var extensions []admin.AdminExtension
	if err := facades.Orm().Query().Order("id asc").Find(&extensions); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	items := make([]map[string]any, 0, len(extensions))
	for _, e := range extensions {
		items = append(items, map[string]any{
			"id":         e.ID,
			"name":       e.Name,
			"is_enabled": e.IsEnabled,
			"created_at": e.CreatedAt,
			"updated_at": e.UpdatedAt,
		})
	}

	return s.SuccessData(ctx, map[string]any{"items": items})
}

// Store 登记一个扩展。
//
// Go 版无法在运行时加载扩展代码，这里只负责把扩展名写入名单，
// 扩展本体需要在宿主编译期注册。
func (s *AdminExtensionService) Store(ctx http.Context) http.Response {
	name := ctx.Request().Input("name")
	if !extensionNamePattern.MatchString(name) {
		return s.FailMsg(ctx, "扩展名称格式应为 vendor/name")
	}

	total, err := facades.Orm().Query().Model(&admin.AdminExtension{}).Where("name", name).Count()
	if err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	if total > 0 {
		return s.FailMsg(ctx, "该扩展已存在")
	}

	extension := admin.AdminExtension{Name: name, IsEnabled: admin.Extension_Enabled_OFF}
	if err := facades.Orm().Query().Create(&extension); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "创建成功")
}

// Enable 启用或停用扩展。
func (s *AdminExtensionService) Enable(ctx http.Context) http.Response {
	return s.toggle(ctx, admin.Extension_Enabled_ON)
}

// Disable 停用扩展。
func (s *AdminExtensionService) Disable(ctx http.Context) http.Response {
	return s.toggle(ctx, admin.Extension_Enabled_OFF)
}

func (s *AdminExtensionService) toggle(ctx http.Context, enabled uint8) http.Response {
	var extension admin.AdminExtension
	if err := facades.Orm().Query().Find(&extension, ctx.Request().InputInt("id")); err != nil {
		return s.FailMsg(ctx, "扩展不存在")
	}

	extension.IsEnabled = enabled
	if err := facades.Orm().Query().Save(&extension); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "操作成功")
}

// Destroy 从名单里移除扩展。
func (s *AdminExtensionService) Destroy(ctx http.Context) http.Response {
	if _, err := facades.Orm().Query().
		WhereIn("id", parseIDList(ctx.Request().Input("id"))).
		Delete(&admin.AdminExtension{}); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "删除成功")
}

// ErrExtensionUnsupported 说明扩展安装能力在 Go 版不可用。
var ErrExtensionUnsupported = errors.New("Go 版不支持在运行时安装扩展，请在宿主编译期注册")
