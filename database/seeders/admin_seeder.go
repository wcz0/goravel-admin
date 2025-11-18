package seeders

import (
    "fmt"
    "time"

    "goravel/app/models"
    "goravel/app/models/admin"

    "github.com/goravel/framework/facades"
    "github.com/goravel/framework/database/orm"
)

const SEED_COMMENT = "2025-11-12 00:10:29"

// AdminSeeder 管理员数据种子
type AdminSeeder struct{}

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
ID: 1,
Username: "admin",
Password: "$2y$12$cgrmY8E0.UIQ6u3vMuxFt.qcFNg4oXTbg4kLxHAzh3HqyHVd/vyZC",
Enabled: 1,
Name: "Administrator",
Avatar: "",
RememberToken: "",
CreatedAt: time.Time{},
	UpdatedAt: time.Time{},
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
        Name: "Administrator",
        Slug: "administrator",
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
        Name: "首页",
        Slug: "home",
        HttpMethod: models.StringSlice{},
        HttpPath: models.StringSlice{"/home*"},
        CustomOrder: 0,
        Model: orm.Model{ID: 1},
    },
    {
        Name: "系统",
        Slug: "system",
        HttpMethod: models.StringSlice{},
        HttpPath: models.StringSlice{},
        CustomOrder: 0,
        Model: orm.Model{ID: 2},
    },
    {
        Name: "管理员",
        Slug: "admin_users",
        HttpMethod: models.StringSlice{},
        HttpPath: models.StringSlice{"/admin_users*"},
        CustomOrder: 0,
        Model: orm.Model{ID: 3},
    },
    {
        Name: "角色",
        Slug: "roles",
        HttpMethod: models.StringSlice{},
        HttpPath: models.StringSlice{"/roles*"},
        CustomOrder: 0,
        Model: orm.Model{ID: 4},
    },
    {
        Name: "权限",
        Slug: "permissions",
        HttpMethod: models.StringSlice{},
        HttpPath: models.StringSlice{"/permissions*"},
        CustomOrder: 0,
        Model: orm.Model{ID: 5},
    },
    {
        Name: "菜单",
        Slug: "menus",
        HttpMethod: models.StringSlice{},
        HttpPath: models.StringSlice{"/menus*"},
        CustomOrder: 0,
        Model: orm.Model{ID: 6},
    },
    {
        Name: "设置",
        Slug: "settings",
        HttpMethod: models.StringSlice{},
        HttpPath: models.StringSlice{"/settings*"},
        CustomOrder: 0,
        Model: orm.Model{ID: 7},
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
ID: 1,
ParentId: 0,
CustomOrder: 0,
Title: "dashboard",
Icon: "mdi:chart-line",
Url: "/dashboard",
UrlType: 1,
Visible: 1,
IsHome: 1,
Component: "",
IsFull: 0,
Extension: "",
CreatedAt: time.Time{},
UpdatedAt: time.Time{},
},
{
ID: 2,
ParentId: 0,
CustomOrder: 0,
Title: "admin_system",
Icon: "material-symbols:settings-outline",
Url: "/system",
UrlType: 1,
Visible: 1,
IsHome: 0,
Component: "",
IsFull: 0,
Extension: "",
CreatedAt: time.Time{},
UpdatedAt: time.Time{},
},
{
ID: 3,
ParentId: 2,
CustomOrder: 0,
Title: "admin_users",
Icon: "ph:user-gear",
Url: "/system/admin_users",
UrlType: 1,
Visible: 1,
IsHome: 0,
Component: "",
IsFull: 0,
Extension: "",
CreatedAt: time.Time{},
UpdatedAt: time.Time{},
},
{
ID: 4,
ParentId: 2,
CustomOrder: 0,
Title: "admin_roles",
Icon: "carbon:user-role",
Url: "/system/admin_roles",
UrlType: 1,
Visible: 1,
IsHome: 0,
Component: "",
IsFull: 0,
Extension: "",
CreatedAt: time.Time{},
UpdatedAt: time.Time{},
},
{
ID: 5,
ParentId: 2,
CustomOrder: 0,
Title: "admin_permission",
Icon: "fluent-mdl2:permissions",
Url: "/system/admin_permissions",
UrlType: 1,
Visible: 1,
IsHome: 0,
Component: "",
IsFull: 0,
Extension: "",
CreatedAt: time.Time{},
UpdatedAt: time.Time{},
},
{
ID: 6,
ParentId: 2,
CustomOrder: 0,
Title: "admin_menu",
Icon: "ant-design:menu-unfold-outlined",
Url: "/system/admin_menus",
UrlType: 1,
Visible: 1,
IsHome: 0,
Component: "",
IsFull: 0,
Extension: "",
CreatedAt: time.Time{},
UpdatedAt: time.Time{},
},
{
ID: 7,
ParentId: 2,
CustomOrder: 0,
Title: "admin_setting",
Icon: "akar-icons:settings-horizontal",
Url: "/system/settings",
UrlType: 1,
Visible: 1,
IsHome: 0,
Component: "",
IsFull: 0,
Extension: "",
CreatedAt: time.Time{},
UpdatedAt: time.Time{},
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
Key: "admin_locale",
Values: "\"en\"",
},
{
Key: "system_theme_setting",
Values: "{\"darkTheme\":false,\"footer\":false,\"breadcrumb\":true,\"themeColor\":\"#1677ff\",\"layoutMode\":\"default\",\"siderTheme\":\"light\",\"topTheme\":\"light\",\"animateInType\":\"alpha\",\"animateInDuration\":600,\"animateOutType\":\"alpha\",\"animateOutDuration\":600,\"loginTemplate\":\"default\",\"keepAlive\":false,\"enableTab\":false,\"tabIcon\":true,\"accordionMenu\":false}",
},
}
	for _, setting := range settings {
		if err := facades.Orm().Query().Create(setting); err != nil {
			return fmt.Errorf("创建设置失败: %w", err)
		}
	}

	// 7. 关联角色和用户
	roleUsers := []map[string]interface{}{
{
"role_id": 1,
"user_id": 1,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
}
	for _, ru := range roleUsers {
		// 通过原生SQL执行关联插入
		query := "INSERT INTO admin_role_users (role_id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?)"
		_, err := facades.Orm().Query().Exec(query, ru["role_id"], ru["user_id"], ru["created_at"], ru["updated_at"])
		if err != nil {
			return fmt.Errorf("关联角色用户失败: %w", err)
		}
	}

	// 8. 关联角色和权限
	rolePermissions := []map[string]interface{}{
{
"role_id": 1,
"permission_id": 1,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"role_id": 1,
"permission_id": 2,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"role_id": 1,
"permission_id": 3,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"role_id": 1,
"permission_id": 4,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"role_id": 1,
"permission_id": 5,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"role_id": 1,
"permission_id": 6,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"role_id": 1,
"permission_id": 7,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
}
	for _, rp := range rolePermissions {
		query := "INSERT INTO admin_role_permissions (role_id, permission_id, created_at, updated_at) VALUES (?, ?, ?, ?)"
		_, err := facades.Orm().Query().Exec(query, rp["role_id"], rp["permission_id"], rp["created_at"], rp["updated_at"])
		if err != nil {
			return fmt.Errorf("关联角色权限失败: %w", err)
		}
	}

	// 9. 关联权限和菜单
	permissionMenus := []map[string]interface{}{
{
"permission_id": 1,
"menu_id": 1,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 2,
"menu_id": 2,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 3,
"menu_id": 3,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 2,
"menu_id": 3,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 4,
"menu_id": 4,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 2,
"menu_id": 4,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 5,
"menu_id": 5,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 2,
"menu_id": 5,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 6,
"menu_id": 6,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 2,
"menu_id": 6,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 7,
"menu_id": 7,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
{
"permission_id": 2,
"menu_id": 7,
"created_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
"updated_at": time.Date(2024, time.August, 10, 7, 4, 53, 0, time.UTC),
},
}
	for _, pm := range permissionMenus {
		query := "INSERT INTO admin_permission_menu (permission_id, menu_id, created_at, updated_at) VALUES (?, ?, ?, ?)"
		_, err := facades.Orm().Query().Exec(query, pm["permission_id"], pm["menu_id"], pm["created_at"], pm["updated_at"])
		if err != nil {
			return fmt.Errorf("关联权限菜单失败: %w", err)
		}
	}

	return nil
}
