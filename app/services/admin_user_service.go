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
	if err := facades.Orm().Query().Where("username", ctx.Request().Input("username")).First(&adminUser); err != nil {
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
