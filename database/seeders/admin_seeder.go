package seeders

import (
	"fmt"
	"time"

	"goravel/app/models"
	"goravel/app/models/admin"

	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
)

const SEED_COMMENT = "2025-11-12 00:10:29"

// AdminSeeder 管理员数据种子
type AdminSeeder struct{}

func (seeder *AdminSeeder) Signature() string {
	return "AdminSeeder"
}

// Run 执行种子
func (seeder *AdminSeeder) Run() error {
	return seeder.seedAdminData()
}

// seedAdminData 填充管理员数据
func (seeder *AdminSeeder) seedAdminData() error {
	// 1. 清理现有数据（如果需要）
	// 注意：这里不清理数据，只是插入新数据

	// 2. 插入用户数据
	users := []*admin.AdminUser{
		{
			ID:        1,
			Username:  "admin",
			Password:  "$2a$12$qbXT0QDJh5PYU1WrpeT3.ufPNkSc.YaoIAuTrSPUee/mvSvgquIZm",
			Enabled:   admin.Enabled_ON,
			Name:      "Administrator",
			CreatedAt: time.Date(2024, time.July, 24, 16, 51, 57, 0, time.UTC),
			UpdatedAt: time.Date(2024, time.July, 24, 16, 51, 57, 0, time.UTC),
		},
	}
	for _, user := range users {
		if err := facades.Orm().Query().Create(user); err != nil {
			return fmt.Errorf("创建用户失败: %w", err)
		}
	}

	// 3. 插入角色数据
	roles := []*admin.AdminRole{
		{
			Name:  "Administrator",
			Slug:  "administrator",
			Model: orm.Model{ID: 1},
		},
	}
	for _, role := range roles {
		if err := facades.Orm().Query().Create(role); err != nil {
			return fmt.Errorf("创建角色失败: %w", err)
		}
	}

	// 4. 插入权限数据
	permissions := []*admin.AdminPermission{
		{
			Name:        "首页",
			Slug:        "home",
			HttpPath:    models.StringSlice{"/home*"},
			CustomOrder: 0,
			ParentId:    0,
			Model:       orm.Model{ID: 1},
		},
		{
			Name:        "系统",
			Slug:        "system",
			CustomOrder: 0,
			ParentId:    0,
			Model:       orm.Model{ID: 2},
		},
		{
			Name:        "管理员",
			Slug:        "admin_users",
			HttpPath:    models.StringSlice{"/admin_users*"},
			CustomOrder: 0,
			ParentId:    2,
			Model:       orm.Model{ID: 3},
		},
		{
			Name:        "角色",
			Slug:        "roles",
			HttpPath:    models.StringSlice{"/roles*"},
			CustomOrder: 0,
			ParentId:    2,
			Model:       orm.Model{ID: 4},
		},
		{
			Name:        "权限",
			Slug:        "permissions",
			HttpPath:    models.StringSlice{"/permissions*"},
			CustomOrder: 0,
			ParentId:    2,
			Model:       orm.Model{ID: 5},
		},
		{
			Name:        "菜单",
			Slug:        "menus",
			HttpPath:    models.StringSlice{"/menus*"},
			CustomOrder: 0,
			ParentId:    2,
			Model:       orm.Model{ID: 6},
		},
		{
			Name:        "设置",
			Slug:        "settings",
			HttpPath:    models.StringSlice{"/settings*"},
			CustomOrder: 0,
			ParentId:    2,
			Model:       orm.Model{ID: 7},
		},
	}
	for _, permission := range permissions {
		if err := facades.Orm().Query().Create(permission); err != nil {
			return fmt.Errorf("创建权限失败: %w", err)
		}
	}

	// 5. 插入菜单数据
	menus := []*admin.AdminMenu{
		{
			ID:          1,
			ParentId:    0,
			CustomOrder: 0,
			Title:       "dashboard",
			Icon:        "mdi:chart-line",
			Url:         "/dashboard",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			IsHome:      admin.IS_HOME_ON,
			Component:   "",
			IsFull:      0,
			Extension:   "",
			CreatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
			UpdatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
		},
		{
			ID:          2,
			ParentId:    0,
			CustomOrder: 0,
			Title:       "admin_system",
			Icon:        "material-symbols:settings-outline",
			Url:         "/system",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			IsHome:      admin.IS_HOME_OFF,
			Component:   "",
			IsFull:      0,
			Extension:   "",
			CreatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
			UpdatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
		},
		{
			ID:          3,
			ParentId:    2,
			CustomOrder: 0,
			Title:       "admin_users",
			Icon:        "ph:user-gear",
			Url:         "/system/admin_users",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			IsHome:      admin.IS_HOME_OFF,
			Component:   "",
			IsFull:      0,
			Extension:   "",
			CreatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
			UpdatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
		},
		{
			ID:          4,
			ParentId:    2,
			CustomOrder: 0,
			Title:       "admin_roles",
			Icon:        "carbon:user-role",
			Url:         "/system/admin_roles",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			IsHome:      admin.IS_HOME_OFF,
			Component:   "",
			IsFull:      0,
			Extension:   "",
			CreatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
			UpdatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
		},
		{
			ID:          5,
			ParentId:    2,
			CustomOrder: 0,
			Title:       "admin_permission",
			Icon:        "fluent-mdl2:permissions",
			Url:         "/system/admin_permissions",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			IsHome:      admin.IS_HOME_OFF,
			Component:   "",
			IsFull:      0,
			Extension:   "",
			CreatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
			UpdatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
		},
		{
			ID:          6,
			ParentId:    2,
			CustomOrder: 0,
			Title:       "admin_menu",
			Icon:        "ant-design:menu-unfold-outlined",
			Url:         "/system/admin_menus",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			IsHome:      admin.IS_HOME_OFF,
			Component:   "",
			IsFull:      0,
			Extension:   "",
			CreatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
			UpdatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
		},
		{
			ID:          7,
			ParentId:    2,
			CustomOrder: 0,
			Title:       "admin_setting",
			Icon:        "akar-icons:settings-horizontal",
			Url:         "/system/settings",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			IsHome:      admin.IS_HOME_OFF,
			Component:   "",
			IsFull:      0,
			Extension:   "",
			CreatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
			UpdatedAt:   time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
		},
	}
	for _, menu := range menus {
		if err := facades.Orm().Query().Create(menu); err != nil {
			return fmt.Errorf("创建菜单失败: %w", err)
		}
	}

	// 6. 插入设置数据
	settings := []*admin.AdminSetting{
		{
			Key:    "admin_locale",
			Values: "\"zh_CN\"",
		},
		{
			Key:    "system_theme_setting",
			Values: "{\"accordionMenu\":false,\"animateInDuration\":600,\"animateInType\":\"alpha\",\"animateOutDuration\":600,\"animateOutType\":\"alpha\",\"breadcrumb\":true,\"darkTheme\":false,\"enableTab\":true,\"footer\":false,\"keepAlive\":false,\"layoutMode\":\"default\",\"loginTemplate\":\"default\",\"siderTheme\":\"light\",\"tabIcon\":true,\"themeColor\":\"#1677ff\",\"topTheme\":\"light\"}",
		},
	}
	for _, setting := range settings {
		if err := facades.Orm().Query().Create(setting); err != nil {
			return fmt.Errorf("创建设置失败: %w", err)
		}
	}

	// 7. 关联角色和用户
	roleUsers := []*admin.AdminRoleUser{
		{RoleID: 1, UserID: 1},
	}
	for _, ru := range roleUsers {
		if err := facades.Orm().Query().Create(ru); err != nil {
			return fmt.Errorf("关联角色用户失败: %w", err)
		}
	}

	// 8. 关联角色和权限
	rolePermissions := []*admin.AdminRolePermission{
		{RoleID: 1, PermissionID: 1},
		{RoleID: 1, PermissionID: 2},
		{RoleID: 1, PermissionID: 3},
		{RoleID: 1, PermissionID: 4},
		{RoleID: 1, PermissionID: 5},
		{RoleID: 1, PermissionID: 6},
		{RoleID: 1, PermissionID: 7},
	}
	for _, rp := range rolePermissions {
		if err := facades.Orm().Query().Create(rp); err != nil {
			return fmt.Errorf("关联角色权限失败: %w", err)
		}
	}

	// 9. 关联权限和菜单
	permissionMenus := []*admin.AdminPermissionMenu{
		{PermissionID: 1, MenuID: 1},
		{PermissionID: 2, MenuID: 2},
		{PermissionID: 3, MenuID: 3},
		{PermissionID: 2, MenuID: 3},
		{PermissionID: 4, MenuID: 4},
		{PermissionID: 2, MenuID: 4},
		{PermissionID: 5, MenuID: 5},
		{PermissionID: 2, MenuID: 5},
		{PermissionID: 6, MenuID: 6},
		{PermissionID: 2, MenuID: 6},
		{PermissionID: 7, MenuID: 7},
		{PermissionID: 2, MenuID: 7},
	}
	for _, pm := range permissionMenus {
		if err := facades.Orm().Query().Create(pm); err != nil {
			return fmt.Errorf("关联权限菜单失败: %w", err)
		}
	}

	return nil
}
