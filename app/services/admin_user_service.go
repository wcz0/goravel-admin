package services

import (
	"goravel/app/models/admin"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type AdminUserService struct {
	*AdminService[*admin.AdminUser]
}

func NewAdminUserService() *AdminUserService {
	return &AdminUserService{
		AdminService: NewAdminService[*admin.AdminUser](admin.NewAdminUser()),
	}
}

func (s *AdminUserService) Login(ctx http.Context) http.Response {
	var adminUser admin.AdminUser
	if err := facades.Orm().Query().Where("username", ctx.Request().Input("username")).
		Where("enabled", admin.Enabled_ON).
		First(&adminUser); err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	if !facades.Hash().Check(ctx.Request().Input("password"), adminUser.Password) {
		return s.FailMsg(ctx, "Password error.")
	}
	token, err := facades.Auth(ctx).Login(&adminUser)
	if err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	return s.SuccessData(ctx, map[string]any{
		"token": token,
		"message": "登录成功",
	})
}

func (a *AdminUserService) List(ctx http.Context) http.Response {
	query := facades.Orm().Query().Select("id", "username", "name", "avatar", "enabled", "created_at")

	if keyword := ctx.Request().Input("keyword"); keyword != "" {
		query.Where("username", "like", "%"+keyword+"%").
			OrWhere("name", "like", "%"+keyword+"%")
	}

	var users []admin.AdminUser
	if err := query.Get(&users); err != nil {
		return a.FailMsg(ctx, err.Error())
	}

	return a.SuccessData(ctx, map[string]any{
		"items": users,
	})
}

func (a *AdminUserService) Export(ctx http.Context) http.Response {
	return a.SuccessData(ctx, map[string]any{
		"items": []any{},
	})
}

func (a *AdminUserService) RoleOptions(ctx http.Context) []map[string]any {
	var roles []admin.AdminRole
	query := facades.Orm().Query()

	// 获取当前用户
	currentUser, ok := ctx.Value("admin_user").(*admin.AdminUser)
	if !ok || currentUser == nil {
		return []map[string]any{}
	}

	// 如果不是超级管理员，则不能分配超级管理员角色
	if !currentUser.IsAdministrator() {
		query.Where("slug", "!=", "administrator")
	}

	if err := query.Get(&roles); err != nil {
		return []map[string]any{}
	}

	options := make([]map[string]any, 0)
	for _, role := range roles {
		options = append(options, map[string]any{
			"id":   role.ID,
			"name": role.Name,
		})
	}

	return options
}

func (a *AdminUserService) QuickEdit(ctx http.Context) http.Response {
	return a.Success(ctx, "快速编辑成功")
}

func (a *AdminUserService) QuickEditItem(ctx http.Context) http.Response {
	return a.Success(ctx, "快速编辑项成功")
}

func (a *AdminUserService) Update(ctx http.Context) http.Response {
	return a.Success(ctx, "更新成功")
}

func (a *AdminUserService) Store(ctx http.Context) http.Response {
	// 获取请求参数
	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")
	name := ctx.Request().Input("name")
	avatar := ctx.Request().Input("avatar")
	enabled := ctx.Request().InputInt("enabled", 1) // 默认启用
	roles := ctx.Request().InputArray("roles") // 角色ID数组

	// 检查用户名是否已存在
	var existingUser admin.AdminUser
	if err := facades.Orm().Query().Where("username", username).First(&existingUser); err != nil {
		return a.FailMsg(ctx, err.Error())
	}
	return a.Success(ctx, "创建成功")
}
