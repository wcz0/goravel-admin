package services

import (
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type AdminUserService struct {
	*Service
}

func NewAdminUserService() *AdminUserService {
	return &AdminUserService{
		Service: NewService(),
	}
}

func (s *AdminUserService) Login(ctx http.Context) http.Response {
	var adminUser models.AdminUser
	if err := facades.Orm().Query().Where("username", ctx.Request().Input("username")).First(&adminUser); err != nil {
		return s.MsgError(ctx, err.Error())
	}
	if !facades.Hash().Check(ctx.Request().Input("password"), adminUser.Password) {
		return s.MsgError(ctx, "Password error.")
	}
	token, err := facades.Auth(ctx).Login(&adminUser)
	if err != nil {
		return s.MsgError(ctx, err.Error())
	}
	return s.DataSuccess(ctx, map[string]string{
		"token": token,
	})
}
