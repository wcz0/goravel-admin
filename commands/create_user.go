package commands

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/goravel-admin/models/admin"
)

// CreateUser 以非交互方式创建后台管理员，便于部署脚本使用。
type CreateUser struct{}

func (receiver *CreateUser) Signature() string {
	return "admin:create-user"
}

func (receiver *CreateUser) Description() string {
	return "创建后台管理员账号"
}

func (receiver *CreateUser) Extend() command.Extend {
	return command.Extend{
		Category: "admin",
		Flags: []command.Flag{
			&command.StringFlag{Name: "username", Aliases: []string{"u"}, Usage: "登录用户名"},
			&command.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "登录密码"},
			&command.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "显示名称，默认与用户名相同"},
			&command.StringSliceFlag{Name: "roles", Usage: "角色 slug 列表，逗号分隔"},
		},
	}
}

func (receiver *CreateUser) Handle(ctx console.Context) error {
	username := strings.TrimSpace(ctx.Option("username"))
	if username == "" {
		value, err := ctx.Ask("请输入用户名")
		if err != nil {
			return err
		}
		username = strings.TrimSpace(value)
	}
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	password := ctx.Option("password")
	if password == "" {
		value, err := ctx.Ask("请输入密码")
		if err != nil {
			return err
		}
		password = value
	}
	if len(password) < 6 {
		return fmt.Errorf("密码长度不能少于 6 位")
	}

	name := strings.TrimSpace(ctx.Option("name"))
	if name == "" {
		name = username
	}

	total, err := facades.Orm().Query().Model(&admin.AdminUser{}).Where("username", username).Count()
	if err != nil {
		return err
	}
	if total > 0 {
		return fmt.Errorf("用户名 %s 已存在", username)
	}

	hashed, err := facades.Hash().Make(password)
	if err != nil {
		return err
	}

	user := admin.AdminUser{
		Username: username,
		Password: hashed,
		Name:     name,
		Enabled:  admin.Enabled_ON,
	}
	if err := facades.Orm().Query().Create(&user); err != nil {
		return err
	}

	// 未显式指定角色时默认给超级管理员，避免新建的账号登录后什么都看不到
	roleSlugs := ctx.OptionSlice("roles")
	if len(roleSlugs) == 0 {
		roleSlugs = []string{"administrator"}
	}

	slugs := make([]any, 0, len(roleSlugs))
	for _, slug := range roleSlugs {
		slugs = append(slugs, slug)
	}
	var roles []admin.AdminRole
	if err := facades.Orm().Query().WhereIn("slug", slugs).Find(&roles); err != nil {
		return err
	}
	if len(roles) == 0 {
		ctx.Warning(fmt.Sprintf("没有找到角色 %s，账号已创建但未分配角色", strings.Join(roleSlugs, ",")))
		return nil
	}
	for i := range roles {
		if err := user.AttachRole(&roles[i]); err != nil {
			return err
		}
	}

	ctx.Info(fmt.Sprintf("管理员 %s 创建成功，已分配角色：%s", username, strings.Join(roleSlugs, ",")))

	return nil
}
