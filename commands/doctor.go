package commands

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

// randomPassword 生成指定长度的随机密码（大小写字母 + 数字）。
func randomPassword(length int) string {
	const charset = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "Admin123456"
	}
	for i := range buf {
		buf[i] = charset[int(buf[i])%len(charset)]
	}

	return string(buf)
}

// coreTables 是后台运行所依赖的表。
var coreTables = []string{
	"admin_users",
	"admin_roles",
	"admin_permissions",
	"admin_menus",
	"admin_role_users",
	"admin_role_permissions",
	"admin_permission_menu",
	"admin_settings",
	"admin_extensions",
	"admin_pages",
	"admin_apis",
	"admin_relationships",
}

// Doctor 巡检后台运行环境，输出问题与修复建议。
type Doctor struct{}

func (receiver *Doctor) Signature() string {
	return "admin:doctor"
}

func (receiver *Doctor) Description() string {
	return "检查后台运行环境与数据完整性"
}

func (receiver *Doctor) Extend() command.Extend {
	return command.Extend{
		Category: "admin",
	}
}

func (receiver *Doctor) Handle(ctx console.Context) error {
	issues := 0

	// 1. APP_KEY
	if facades.Config().GetString("app.key") == "" {
		ctx.Error("[APP_KEY] 未配置，导出下载签名与 JWT 会失效，请执行 artisan key:generate")
		issues++
	} else {
		ctx.Info("[APP_KEY] 已配置")
	}

	// 2. 数据库连通性
	if db, err := facades.Orm().Query().Exec("select 1"); err != nil || db == nil {
		ctx.Error(fmt.Sprintf("[DATABASE] 连接失败: %v", err))
		ctx.Warning("后续检查已跳过")

		return nil
	}
	ctx.Info("[DATABASE] 连接正常")

	// 3. 核心表
	missing := make([]string, 0)
	for _, table := range coreTables {
		if !facades.Schema().HasTable(table) {
			missing = append(missing, table)
		}
	}
	if len(missing) > 0 {
		ctx.Error(fmt.Sprintf("[MIGRATION] 缺少 %d 张表: %s", len(missing), strings.Join(missing, ", ")))
		ctx.Warning("请执行 artisan migrate")
		issues++
	} else {
		ctx.Info(fmt.Sprintf("[MIGRATION] %d 张核心表齐全", len(coreTables)))
	}

	// 4. 管理员账号
	userCount, err := facades.Orm().Query().Model(&struct{}{}).Table("admin_users").Count()
	if err != nil {
		ctx.Error(fmt.Sprintf("[USERS] 查询失败: %v", err))
		issues++
	} else if userCount == 0 {
		ctx.Error("[USERS] 没有任何后台账号，请执行 artisan admin:create-user")
		issues++
	} else {
		ctx.Info(fmt.Sprintf("[USERS] 共 %d 个账号", userCount))
	}

	// 5. 菜单首页唯一性
	homeCount, err := facades.Orm().Query().Model(&struct{}{}).Table("admin_menus").Where("is_home", 1).Count()
	if err != nil {
		ctx.Error(fmt.Sprintf("[MENUS] 查询失败: %v", err))
		issues++
	} else if homeCount == 0 {
		ctx.Warning("[MENUS] 没有设置首页菜单，登录后会落到空白页")
	} else if homeCount > 1 {
		ctx.Error(fmt.Sprintf("[MENUS] 有 %d 个首页菜单，应当只保留一个", homeCount))
		issues++
	} else {
		ctx.Info("[MENUS] 首页菜单正常")
	}

	// 6. 菜单父级断链
	if missingParents := findOrphanMenus(); len(missingParents) > 0 {
		ctx.Error(fmt.Sprintf("[MENUS] %d 个菜单的父级不存在: %v", len(missingParents), missingParents))
		issues++
	} else {
		ctx.Info("[MENUS] 父级引用正常")
	}

	ctx.NewLine()
	if issues == 0 {
		ctx.Info("检查完成，未发现问题")
	} else {
		ctx.Warning(fmt.Sprintf("检查完成，共发现 %d 个问题", issues))
	}

	return nil
}

// findOrphanMenus 找出父级 id 不存在的菜单（父级为 0 表示顶级，不算断链）。
func findOrphanMenus() []uint {
	var menus []struct {
		ID       uint
		ParentId uint32
	}
	if err := facades.Orm().Query().Table("admin_menus").Select("id", "parent_id").Find(&menus); err != nil {
		return nil
	}

	ids := make(map[uint]bool, len(menus))
	for _, m := range menus {
		ids[m.ID] = true
	}

	orphans := make([]uint, 0)
	for _, m := range menus {
		if m.ParentId != 0 && !ids[uint(m.ParentId)] {
			orphans = append(orphans, m.ID)
		}
	}

	return orphans
}
