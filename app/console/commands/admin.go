package commands

import (
	"fmt"
	"goravel/app/models/admin"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
}

// Signature 命令的名称和签名
func (receiver *Admin) Signature() string {
	return "admin:init"
}

// Description 命令的描述
func (receiver *Admin) Description() string {
	return "初始化admin database"
}

func (receiver *Service) Extend() command.Extend {
	return command.Extend{
		Category: "admin",
	}
}

// Handle 执行命令的逻辑
func (receiver *Admin) Handle(ctx console.Context) error {
	// 初始化数据库
	// if err := facades.Migration().Run(); err != nil {
	// 	return err
	// }

	// 使用ORM方式初始化数据库
	if err := receiver.initDatabase(); err != nil {
		return err
	}

	fmt.Println("数据库初始化成功")
	return nil
}

// initDatabase 初始化数据库
func (receiver *Admin) initDatabase() error {
	// 使用ORM方式插入数据

	// 1. 创建管理员角色
	adminRole := &admin.AdminRole{
		Name: "Administrator",
		Slug: "administrator",
	}
	if err := facades.Orm().Query().Create(adminRole); err != nil {
		return fmt.Errorf("创建管理员角色失败: %w", err)
	}

	// 2. 创建操作员角色
	operatorRole := &admin.AdminRole{
		Name: "Operator",
		Slug: "operator",
	}
	if err := facades.Orm().Query().Create(operatorRole); err != nil {
		return fmt.Errorf("创建操作员角色失败: %w", err)
	}

	// 3. 创建管理员用户
	adminUser := &admin.AdminUser{
		Username: "admin",
		Password: facades.Hash().Make("admin"),
		Enabled:  admin.Enabled_ON,
		Name:     "Administrator",
		Avatar:   "https://cdn.learnku.com/uploads/avatars/7265_1539760150.jpg!/both/100x100",
	}
	if err := facades.Orm().Query().Create(adminUser); err != nil {
		return fmt.Errorf("创建管理员用户失败: %w", err)
	}

	// 4. 关联管理员用户和管理员角色
	if err := facades.Orm().Query().Model(adminUser).Association("AdminRoles").Append(adminRole); err != nil {
		return fmt.Errorf("关联管理员用户和角色失败: %w", err)
	}

	// 5. 创建权限
	permissions := []*admin.AdminPermission{
		{
			Name:        "All permission",
			Slug:        "*",
			HttpMethod:  []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
			HttpPath:    []string{"*"},
			CustomOrder: 1,
		},
		{
			Name:        "Dashboard",
			Slug:        "dashboard",
			HttpMethod:  []string{"GET"},
			HttpPath:    []string{"/"},
			CustomOrder: 2,
		},
		{
			Name:        "Login",
			Slug:        "auth.login",
			HttpMethod:  []string{"GET", "POST"},
			HttpPath:    []string{"/auth/login", "/auth/logout"},
			CustomOrder: 3,
		},
		{
			Name:        "User setting",
			Slug:        "auth.setting",
			HttpMethod:  []string{"GET", "PUT"},
			HttpPath:    []string{"/auth/setting"},
			CustomOrder: 4,
		},
		{
			Name:        "Auth management",
			Slug:        "auth.management",
			HttpMethod:  []string{"GET", "POST"},
			HttpPath:    []string{"/auth/roles", "/auth/permissions", "/auth/menu", "/auth/logs"},
			CustomOrder: 5,
		},
	}

	for _, permission := range permissions {
		if err := facades.Orm().Query().Create(permission); err != nil {
			return fmt.Errorf("创建权限失败: %w", err)
		}
	}

	// 6. 关联角色和权限
	// 管理员角色拥有所有权限
	for _, permission := range permissions {
		if err := facades.Orm().Query().Model(adminRole).Association("AdminPermissions").Append(permission); err != nil {
			return fmt.Errorf("关联管理员角色和权限失败: %w", err)
		}
	}

	// 操作员角色只有基本权限
	operatorPermissions := permissions[1:] // 除了All permission外的权限
	for _, permission := range operatorPermissions {
		if err := facades.Orm().Query().Model(operatorRole).Association("AdminPermissions").Append(permission); err != nil {
			return fmt.Errorf("关联操作员角色和权限失败: %w", err)
		}
	}

	// 7. 创建菜单
	menus := []*admin.AdminMenu{
		{
			ParentId:   0,
			CustomOrder: 1,
			Title:      "Dashboard",
			Icon:       "mdi:chart-line",
			Url:        "/dashboard",
			UrlType:    admin.TYPE_ROUTE,
			Visible:    admin.Enabled_ON,
			IsHome:     admin.,
		},
		{
			ParentId:   0,
			CustomOrder: 2,
			Title:      "Admin",
			Icon:       "fa-tasks",
			Uri:        "",
		},
		{
			ParentId:   2,
			CustomOrder: 3,
			Title:      "Users",
			Icon:       "fa-users",
			Uri:        "auth/users",
		},
		{
			ParentId:   2,
			CustomOrder: 4,
			Title:      "Roles",
			Icon:       "fa-user",
			Uri:        "auth/roles",
		},
		{
			ParentId:   2,
			CustomOrder: 5,
			Title:      "Permission",
			Icon:       "fa-ban",
			Uri:        "auth/permissions",
		},
		{
			ParentId:   2,
			CustomOrder: 6,
			Title:      "Menu",
			Icon:       "fa-bars",
			Uri:        "auth/menu",
		},
		{
			ParentId:   2,
			CustomOrder: 7,
			Title:      "Operation log",
			Icon:       "fa-history",
			Uri:        "auth/logs",
		},
	}

	for _, menu := range menus {
		if err := facades.Orm().Query().Create(menu); err != nil {
			return fmt.Errorf("创建菜单失败: %w", err)
		}
	}

	// 8. 关联权限和菜单
	// Dashboard权限关联Dashboard菜单
	dashboardPermission := permissions[1]
	dashboardMenu := menus[0]
	if err := facades.Orm().Query().Model(dashboardPermission).Association("AdminMenus").Append(dashboardMenu); err != nil {
		return fmt.Errorf("关联Dashboard权限和菜单失败: %w", err)
	}

	// Auth management权限关联Admin下的所有菜单
	authManagementPermission := permissions[4]
	adminMenus := menus[1:]
	for _, menu := range adminMenus {
		if err := facades.Orm().Query().Model(authManagementPermission).Association("AdminMenus").Append(menu); err != nil {
			return fmt.Errorf("关联Auth management权限和菜单失败: %w", err)
		}
	}

	// 9. 创建系统设置
	settings := []*admin.AdminSetting{
		{
			Slug:  "site_title",
			Value: "Goravel Admin",
		},
		{
			Slug:  "site_logo",
			Value: "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"feather feather-feather\"><path d=\"M20.24 12.24a6 6 0 0 0-8.49-8.49L5 10.5V19h8.5z\"></path><line x1=\"16\" y1=\"8\" x2=\"2\" y2=\"22\"></line><line x1=\"17.5\" y1=\"15\" x2=\"9\" y2=\"15\"></line></svg>",
		},
		{
			Slug:  "site_logo_mini",
			Value: "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"feather feather-feather\"><path d=\"M20.24 12.24a6 6 0 0 0-8.49-8.49L5 10.5V19h8.5z\"></path><line x1=\"16\" y1=\"8\" x2=\"2\" y2=\"22\"></line><line x1=\"17.5\" y1=\"15\" x2=\"9\" y2=\"15\"></line></svg>",
		},
		{
			Slug:  "site_footer",
			Value: "© Goravel Admin",
		},
	}

	for _, setting := range settings {
		if err := facades.Orm().Query().Create(setting); err != nil {
			return fmt.Errorf("创建系统设置失败: %w", err)
		}
	}

	return nil
}
