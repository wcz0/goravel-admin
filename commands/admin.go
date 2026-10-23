package commands

import (
	"fmt"
	"github.com/wcz0/goravel-admin/models/admin"

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

	// 4. 创建权限
	// http_path 必须与真实路由一致：匹配时只会自动补 /admin-api 前缀，
	// 写成 /auth/roles 这种旧路径会导致非超管用户即使有权限也全部 403。
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
			HttpPath:    []string{"/dashboard*"},
			CustomOrder: 2,
		},
		{
			Name:        "Admin users",
			Slug:        "admin_users",
			HttpPath:    []string{"/system/admin_users*", "/system/users*"},
			CustomOrder: 3,
		},
		{
			Name:        "Admin roles",
			Slug:        "admin_roles",
			HttpPath:    []string{"/system/admin_roles*"},
			CustomOrder: 4,
		},
		{
			Name:        "Admin permissions",
			Slug:        "admin_permissions",
			HttpPath:    []string{"/system/permissions*", "/system/admin_permissions*", "/system/_admin_permissions_auto_generate"},
			CustomOrder: 5,
		},
		{
			Name:        "Admin menus",
			Slug:        "admin_menus",
			HttpPath:    []string{"/system/menus*", "/system/admin_menus*"},
			CustomOrder: 6,
		},
		{
			Name:        "Settings",
			Slug:        "settings",
			HttpPath:    []string{"/_settings*", "/system/settings*"},
			CustomOrder: 7,
		},
		{
			// 登录后人人需要的公共接口：当前用户、个人设置、上传、页面结构、图标检索
			Name:        "Common",
			Slug:        "common",
			HttpPath:    []string{"/current-user", "/user_setting*", "/menus", "/upload_image", "/upload_file", "/upload_rich", "/upload_chunk*", "/page_schema*", "/_iconify_search*"},
			CustomOrder: 8,
		},
	}

	permissionsBySlug := make(map[string]*admin.AdminPermission, len(permissions))
	for _, permission := range permissions {
		if err := facades.Orm().Query().FirstOrCreate(permission, admin.AdminPermission{Slug: permission.Slug}); err != nil {
			return fmt.Errorf("创建权限失败: %w", err)
		}
		permissionsBySlug[permission.Slug] = permission
	}

	// 5. 关联角色和权限
	// 管理员角色拥有全部权限
	for _, permission := range permissions {
		rolePermission := &admin.AdminRolePermission{
			RoleID:       adminRole.ID,
			PermissionID: permission.ID,
		}
		if err := facades.Orm().Query().FirstOrCreate(rolePermission, admin.AdminRolePermission{
			RoleID:       adminRole.ID,
			PermissionID: permission.ID,
		}); err != nil {
			return fmt.Errorf("关联管理员角色和权限失败: %w", err)
		}
	}

	// 操作员角色只拿到日常够用的权限，不含菜单管理与系统设置
	operatorSlugs := []string{"dashboard", "admin_users", "admin_roles", "common"}
	for _, slug := range operatorSlugs {
		permission, ok := permissionsBySlug[slug]
		if !ok {
			continue
		}
		rolePermission := &admin.AdminRolePermission{
			RoleID:       operatorRole.ID,
			PermissionID: permission.ID,
		}
		if err := facades.Orm().Query().FirstOrCreate(rolePermission, admin.AdminRolePermission{
			RoleID:       operatorRole.ID,
			PermissionID: permission.ID,
		}); err != nil {
			return fmt.Errorf("关联操作员角色和权限失败: %w", err)
		}
	}

	// 6. 创建菜单
	// 先建父级再建子级，父级 ID 按 URL 回填，避免硬编码 2 这类脆弱假设。
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
			CustomOrder: 1,
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
			CustomOrder: 1,
			Title:       "admin_roles",
			Icon:        "carbon:user-role",
			Url:         "/system/admin_roles",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
		{
			CustomOrder: 2,
			Title:       "admin_permission",
			Icon:        "fluent-mdl2:permissions",
			Url:         "/system/admin_permission",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
		{
			CustomOrder: 3,
			Title:       "admin_menu",
			Icon:        "ant-design:menu-unfold-outlined",
			Url:         "/system/admin_menus",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
		{
			CustomOrder: 4,
			Title:       "admin_setting",
			Icon:        "akar-icons:settings-horizontal",
			Url:         "/system/settings",
			UrlType:     admin.TYPE_ROUTE,
			Visible:     admin.Enabled_ON,
			KeepAlive:   &[]uint8{0}[0],
			IFrameUrl:   &[]string{""}[0],
		},
	}

	menusByUrl := make(map[string]*admin.AdminMenu, len(menus))
	for _, menu := range menus {
		if err := facades.Orm().Query().FirstOrCreate(menu, admin.AdminMenu{Url: menu.Url}); err != nil {
			return fmt.Errorf("创建菜单失败: %w", err)
		}
		menusByUrl[menu.Url] = menu
	}

	// 子菜单统一挂到 /system 下
	if systemMenu, ok := menusByUrl["/system"]; ok {
		for _, menu := range menus {
			if menu.ParentId != 0 || menu.Url == "/system" || menu.Url == "/dashboard" {
				continue
			}
			menu.ParentId = uint32(systemMenu.ID)
			if _, err := facades.Orm().Query().Model(menu).Where("id", menu.ID).Update("parent_id", systemMenu.ID); err != nil {
				return fmt.Errorf("更新菜单父级失败: %w", err)
			}
		}
	}

	// 7. 关联权限和菜单
	// 非超管用户的菜单来自「角色 -> 权限 -> 菜单」这条链，缺了关联侧边栏就是空的。
	permMenus := map[string][]string{
		"*":                 {"/dashboard", "/system", "/system/admin_users", "/system/admin_roles", "/system/admin_permission", "/system/admin_menus", "/system/settings"},
		"dashboard":         {"/dashboard"},
		"admin_users":       {"/system", "/system/admin_users"},
		"admin_roles":       {"/system", "/system/admin_roles"},
		"admin_permissions": {"/system", "/system/admin_permission"},
		"admin_menus":       {"/system", "/system/admin_menus"},
		"settings":          {"/system", "/system/settings"},
		"common":            nil,
	}
	for slug, urls := range permMenus {
		permission, ok := permissionsBySlug[slug]
		if !ok {
			continue
		}
		for _, url := range urls {
			menu, ok := menusByUrl[url]
			if !ok {
				continue
			}
			pm := &admin.AdminPermissionMenu{PermissionID: permission.ID, MenuID: menu.ID}
			if err := facades.Orm().Query().FirstOrCreate(pm, admin.AdminPermissionMenu{
				PermissionID: permission.ID,
				MenuID:       menu.ID,
			}); err != nil {
				return fmt.Errorf("关联权限(%s)和菜单失败: %w", slug, err)
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
		if err := facades.Orm().Query().FirstOrCreate(setting, admin.AdminSetting{Key: setting.Key}); err != nil {
			return fmt.Errorf("创建系统设置失败: %w", err)
		}
	}

	// 10. 关联用户和角色
	// 管理员用户关联管理员角色 - 使用自定义 pivot 模型
	adminUserRole := &admin.AdminRoleUser{
		RoleID: adminRole.ID,
		UserID: adminUser.ID,
	}
	if err := facades.Orm().Query().FirstOrCreate(adminUserRole, admin.AdminRoleUser{
		RoleID: adminRole.ID,
		UserID: adminUser.ID,
	}); err != nil {
		return fmt.Errorf("关联管理员用户和角色失败: %w", err)
	}

	return nil
}
