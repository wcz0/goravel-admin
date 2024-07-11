package seeders

import (
	"goravel/app/models"

	"github.com/goravel/framework/facades"
)

type AdminDataSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *AdminDataSeeder) Signature() string {
	return "AdminDataSeeder"
}

// Run executes the seeder logic.
func (s *AdminDataSeeder) Run() error {
	// 创建初始用户

	hashedPassword, err := facades.Hash().Make("admin")
	if err != nil {
		return err
	}

	if err := facades.Orm().Query().Create(&models.AdminUser{
		Username:      "admin",
		Password:      hashedPassword,
		Name:          "Admin",
		Avatar:        "",
		RememberToken: "",
	}); err != nil {
		return err
	}

	// 创建初始角色
	if err := facades.Orm().Query().Create(&models.AdminRole{
		Name:  "Admin",
		Value: "root",
	}); err != nil {
		return err
	}

	// 用户分配角色

	// 创建初始权限
	if err := facades.Orm().Query().Create(&[]models.AdminPermission{
		{
			Name: "首页",
			Value: "/dashboard",
		},
		{
			Name: "系统",
			Value: "/system",
		},
		{
			Name: "后台用户",
			Value: "/admin_user",
		},
		{
			Name: "后台角色",
			Value: "/admin_role",
		},
		{
			Name: "后台权限",
			Value: "/admin_permission",
		},
		{
			Name: "后台菜单",
			Value: "/admin_menu",
		},

	}); err != nil {
		return err
	}

	// 角色分配权限

	// 创建初始菜单
	//
	

	// 权限绑定菜单

	return nil
}
