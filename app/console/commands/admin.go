package commands

import (
	"fmt"
	"goravel/app/models/admin"
	"time"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
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

func (receiver *Admin) Extend() command.Extend {
	return command.Extend{
		Category: "admin",
	}
}

// Handle 执行命令的逻辑
func (receiver *Admin) Handle(ctx console.Context) error {

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
	if err := facades.Orm().Query().FirstOrCreate(adminRole, admin.AdminRole{Slug: "administrator"}); err != nil {
		return fmt.Errorf("创建管理员角色失败: %w", err)
	}

	// 2. 创建操作员角色
	operatorRole := &admin.AdminRole{
		Name: "Operator",
		Slug: "operator",
	}
	if err := facades.Orm().Query().FirstOrCreate(operatorRole, admin.AdminRole{Slug: "operator"}); err != nil {
		return fmt.Errorf("创建操作员角色失败: %w", err)
	}

	// 3. 创建管理员用户
	adminUser := &admin.AdminUser{
		Username: "admin",
		Password: func() string {
			password, _ := facades.Hash().Make("admin")
			return password
		}(),
		Enabled: admin.Enabled_ON,
		Name:    "Administrator",
	}
	if err := facades.Orm().Query().FirstOrCreate(adminUser, admin.AdminUser{Username: "admin"}); err != nil {
		return fmt.Errorf("创建管理员用户失败: %w", err)
	}

	// 4. 关联管理员用户和角色
	// 使用 GORM 的原生方法处理 pivot table 时间戳
	now := time.Now()
	if _, err := facades.Orm().Query().Exec("INSERT INTO admin_role_users (role_id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?)",
		adminRole.ID, adminUser.ID, now, now); err != nil {
		return fmt.Errorf("关联管理员用户和角色失败: %w", err)
	}

    // 5. 创建权限（按用户指定的初始数据）
    permissions := []*admin.AdminPermission{
        {
            Name:        "首页",
            Slug:        "home",
            HttpMethod:  []string{},
            HttpPath:    []string{"/home*"},
            CustomOrder: 0,
            ParentId:    0,
        },
        {
            Name:        "系统",
            Slug:        "system",
            HttpMethod:  []string{},
            HttpPath:    []string{},
            CustomOrder: 0,
            ParentId:    0,
        },
        {
            Name:        "管理员",
            Slug:        "admin_users",
            HttpMethod:  []string{},
            HttpPath:    []string{"/admin_users*"},
            CustomOrder: 0,
            ParentId:    0, // 先设为0，创建后再更新
        },
        {
            Name:        "角色",
            Slug:        "roles",
            HttpMethod:  []string{},
            HttpPath:    []string{"/roles*"},
            CustomOrder: 0,
            ParentId:    0, // 先设为0，创建后再更新
        },
        {
            Name:        "权限",
            Slug:        "permissions",
            HttpMethod:  []string{},
            HttpPath:    []string{"/permissions*"},
            CustomOrder: 0,
            ParentId:    0, // 先设为0，创建后再更新
        },
        {
            Name:        "菜单",
            Slug:        "menus",
            HttpMethod:  []string{},
            HttpPath:    []string{"/menus*"},
            CustomOrder: 0,
            ParentId:    0, // 先设为0，创建后再更新
        },
        {
            Name:        "设置",
            Slug:        "settings",
            HttpMethod:  []string{},
            HttpPath:    []string{"/settings*"},
            CustomOrder: 0,
            ParentId:    0, // 先设为0，创建后再更新
        },
    }

    // 创建权限并保存ID用于后续关联
    var systemPermissionID uint
    for _, permission := range permissions {
        if err := facades.Orm().Query().Create(permission); err != nil {
            return fmt.Errorf("创建权限失败: %w", err)
        }
        // 保存"系统"权限的ID，用于设置其他权限的父级
        if permission.Slug == "system" {
            systemPermissionID = permission.ID
        }
    }

    // 更新子权限的ParentId
    childSlugs := []string{"admin_users", "roles", "permissions", "menus", "settings"}
    for _, slug := range childSlugs {
        if _, err := facades.Orm().Query().Model(&admin.AdminPermission{}).
            Where("slug = ?", slug).
            Update("parent_id", uint32(systemPermissionID)); err != nil {
            return fmt.Errorf("更新权限父级关系失败: %w", err)
        }
    }

	// 6. 关联角色和权限
	// 管理员角色拥有所有权限 - 使用自定义 pivot 模型处理时间戳
	for _, permission := range permissions {
		rolePermission := &admin.AdminRolePermission{
			RoleID:       adminRole.ID,
			PermissionID: permission.ID,
		}
		if err := facades.Orm().Query().Create(rolePermission); err != nil {
			return fmt.Errorf("关联管理员角色和权限失败: %w", err)
		}
	}

	// 操作员角色只有基本权限
	operatorPermissions := permissions[1:] // 除了All permission外的权限
	for _, permission := range operatorPermissions {
		rolePermission := &admin.AdminRolePermission{
			RoleID:       operatorRole.ID,
			PermissionID: permission.ID,
		}
		if err := facades.Orm().Query().Create(rolePermission); err != nil {
			return fmt.Errorf("关联操作员角色和权限失败: %w", err)
		}
	}

	// 7. 创建菜单
	menus := []*admin.AdminMenu{
		{
			ParentId:    0,
			CustomOrder: 0,
			Title:       "dashboard",
			Icon:        "mdi:chart-line",
			Url:         "/dashboard",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			IsHome:      admin.IS_HOME_ON,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
		{
			ParentId:    0,
			CustomOrder: 0,
			Title:       "admin_system",
			Icon:        "material-symbols:settings-outline",
			Url:         "/system",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			IsHome:      admin.IS_HOME_OFF,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
		{
			ParentId:    2,
			CustomOrder: 0,
			Title:       "admin_users",
			Icon:        "ph:user-gear",
			Url:         "/system/admin_users",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
		{
			ParentId:    2,
			CustomOrder: 0,
			Title:       "admin_roles",
			Icon:        "carbon:user-role",
			Url:         "/system/admin_roles",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
		{
			ParentId:    2,
			CustomOrder: 0,
			Title:       "admin_permission",
			Icon:        "fluent-mdl2:permissions",
			Url:         "/system/admin_permission",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
		{
			ParentId:    2,
			CustomOrder: 0,
			Title:       "admin_menu",
			Icon:        "ant-design:menu-unfold-outlined",
			Url:         "/system/admin_menus",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
		{
			ParentId:    2,
			CustomOrder: 0,
			Title:       "admin_setting",
			Icon:        "akar-icons:settings-horizontal",
			Url:         "/system/settings",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
	}

	for _, menu := range menus {
		if err := facades.Orm().Query().Create(menu); err != nil {
			return fmt.Errorf("创建菜单失败: %w", err)
		}
	}

    // 8. 关联权限和菜单（按 slug 与菜单对应关系）
    // home -> dashboard, system -> system, admin_users -> admin_users, roles -> admin_roles, permissions -> admin_permission, menus -> admin_menus, settings -> settings
    slugToMenuIndex := map[string]int{
        "home":         0,
        "system":       1,
        "admin_users":  2,
        "roles":        3,
        "permissions":  4,
        "menus":        5,
        "settings":     6,
    }
    for _, p := range permissions {
        if idx, ok := slugToMenuIndex[p.Slug]; ok && idx >= 0 && idx < len(menus) {
            pm := &admin.AdminPermissionMenu{PermissionID: p.ID, MenuID: menus[idx].ID}
            if err := facades.Orm().Query().Create(pm); err != nil {
                return fmt.Errorf("关联权限(%s)和菜单失败: %w", p.Slug, err)
            }
        }
    }

	// 9. 创建系统设置
	settings := []*admin.AdminSetting{
		{
			Key:    "site_title",
			Values: "Goravel Admin",
		},
		{
			Key:    "site_logo",
			Values: "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"feather feather-feather\"><path d=\"M20.24 12.24a6 6 0 0 0-8.49-8.49L5 10.5V19h8.5z\"></path><line x1=\"16\" y1=\"8\" x2=\"2\" y2=\"22\"></line><line x1=\"17.5\" y1=\"15\" x2=\"9\" y2=\"15\"></line></svg>",
		},
		{
			Key:    "site_logo_mini",
			Values: "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"feather feather-feather\"><path d=\"M20.24 12.24a6 6 0 0 0-8.49-8.49L5 10.5V19h8.5z\"></path><line x1=\"16\" y1=\"8\" x2=\"2\" y2=\"22\"></line><line x1=\"17.5\" y1=\"15\" x2=\"9\" y2=\"15\"></line></svg>",
		},
		{
			Key:    "site_footer",
			Values: "© Goravel Admin",
		},
	}

	for _, setting := range settings {
		if err := facades.Orm().Query().Create(setting); err != nil {
			return fmt.Errorf("创建系统设置失败: %w", err)
		}
	}

	// 10. 关联用户和角色
	// 管理员用户关联管理员角色 - 使用自定义 pivot 模型
	adminUserRole := &admin.AdminRoleUser{
		RoleID: adminRole.ID,
		UserID: adminUser.ID,
	}
	if err := facades.Orm().Query().Create(adminUserRole); err != nil {
		return fmt.Errorf("关联管理员用户和角色失败: %w", err)
	}

	return nil
}
