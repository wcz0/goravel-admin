package services

import (
	"goravel/app/models/admin"

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
		Where("enabled", 1).
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
	return s.SuccessMsgData(ctx, "登录成功", map[string]string{
		"token": token,
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

	return a.SuccessData(ctx, users)
}

func (a *AdminUserService) Export(ctx http.Context) http.Response {
	return a.SuccessData(ctx, []any{})
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
