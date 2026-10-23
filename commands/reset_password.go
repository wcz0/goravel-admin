package commands

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/goravel-admin/models/admin"
)

// ResetPassword 重置后台用户密码，用于登录不上时的应急处理。
type ResetPassword struct{}

func (receiver *ResetPassword) Signature() string {
	return "admin:reset-password"
}

func (receiver *ResetPassword) Description() string {
	return "重置后台用户密码"
}

func (receiver *ResetPassword) Extend() command.Extend {
	return command.Extend{
		Category: "admin",
		Flags: []command.Flag{
			&command.StringFlag{Name: "username", Aliases: []string{"u"}, Usage: "登录用户名"},
			&command.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "新密码，留空则随机生成"},
		},
	}
}

func (receiver *ResetPassword) Handle(ctx console.Context) error {
	username := strings.TrimSpace(ctx.Option("username"))
	if username == "" {
		value, err := ctx.Ask("请输入要重置的用户名")
		if err != nil {
			return err
		}
		username = strings.TrimSpace(value)
	}
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	var user admin.AdminUser
	if err := facades.Orm().Query().Where("username", username).First(&user); err != nil {
		return fmt.Errorf("用户 %s 不存在", username)
	}

	password := ctx.Option("password")
	generated := false
	if password == "" {
		password = randomPassword(12)
		generated = true
	} else if len(password) < 6 {
		return fmt.Errorf("密码长度不能少于 6 位")
	}

	hashed, err := facades.Hash().Make(password)
	if err != nil {
		return err
	}

	user.Password = hashed
	if err := facades.Orm().Query().Save(&user); err != nil {
		return err
	}

	if generated {
		ctx.Info(fmt.Sprintf("用户 %s 的密码已重置为：%s", username, password))
		ctx.Warning("请登录后立即修改密码")
	} else {
		ctx.Info(fmt.Sprintf("用户 %s 的密码已重置", username))
	}

	return nil
}
