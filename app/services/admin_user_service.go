package services

import (
	"goravel/app/models/admin"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type AdminUserService struct {
	*AdminService[*admin.AdminUser]
}

func NewAdminUserService() *AdminUserService {
	return &AdminUserService{
		AdminService: NewAdminService[*admin.AdminUser](admin.NewAdminUser()),
	}
}

func (s *AdminUserService) Login(ctx http.Context) http.Response {
	var adminUser admin.AdminUser
	if err := facades.Orm().Query().Where("username", ctx.Request().Input("username")).
		Where("enabled", admin.Enabled_ON).
		First(&adminUser); err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	if !facades.Hash().Check(ctx.Request().Input("password"), adminUser.Password) {
		return s.FailMsg(ctx, "Password error.")
	}
	token, err := facades.Auth(ctx).Login(&adminUser)
	if err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	return s.SuccessMsgData(ctx, "登录成功", map[string]string{
		"token": token,
	})
}

func (a *AdminUserService) List(ctx http.Context) http.Response {
	query := facades.Orm().Query().Select("id", "username", "name", "avatar", "enabled", "created_at")

	if keyword := ctx.Request().Input("keyword"); keyword != "" {
		query.Where("username", "like", "%"+keyword+"%").
			OrWhere("name", "like", "%"+keyword+"%")
	}

	var users []admin.AdminUser
	if err := query.Get(&users); err != nil {
		return a.FailMsg(ctx, err.Error())
	}

	return a.SuccessData(ctx, users)
}

func (a *AdminUserService) Export(ctx http.Context) http.Response {
	return a.SuccessData(ctx, []any{})
}

func (a *AdminUserService) RoleOptions(ctx http.Context) []map[string]any {
	var roles []admin.AdminRole
	query := facades.Orm().Query()

	// 获取当前用户
	currentUser, ok := ctx.Value("admin_user").(*admin.AdminUser)
	if !ok || currentUser == nil {
		return []map[string]any{}
	}

	// 如果不是超级管理员，则不能分配超级管理员角色
	if !currentUser.IsAdministrator() {
		query.Where("slug", "!=", "administrator")
	}

	if err := query.Get(&roles); err != nil {
		return []map[string]any{}
	}

	options := make([]map[string]any, 0)
	for _, role := range roles {
		options = append(options, map[string]any{
			"id":   role.ID,
			"name": role.Name,
		})
	}

	return options
}

func (a *AdminUserService) QuickEdit(ctx http.Context) http.Response {
	// 添加调试日志查看请求数据
	allData := ctx.Request().All()
	facades.Log().Info("QuickEdit方法接收到的所有数据", map[string]any{
		"all_data": allData,
		"method": ctx.Request().Method(),
		"content_type": ctx.Request().Header("Content-Type"),
	})

	// 获取请求参数
	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")
	name := ctx.Request().Input("name")
	avatar := ctx.Request().Input("avatar")
	enabled := ctx.Request().InputInt("enabled", 1)
	roleIds := ctx.Request().Input("role_ids")

	// 添加调试日志查看获取到的参数
	facades.Log().Info("QuickEdit获取到的参数", map[string]any{
		"username": username,
		"password": password,
		"name":     name,
		"avatar":   avatar,
		"enabled":  enabled,
		"role_ids": roleIds,
	})

	return a.Success(ctx)
}

func (a *AdminUserService) QuickEditItem(ctx http.Context) http.Response {
	// 添加调试日志查看请求数据
	allData := ctx.Request().All()
	facades.Log().Info("QuickEditItem方法接收到的所有数据", map[string]any{
		"all_data": allData,
		"method": ctx.Request().Method(),
		"content_type": ctx.Request().Header("Content-Type"),
	})

	// 获取请求参数
	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")
	name := ctx.Request().Input("name")
	avatar := ctx.Request().Input("avatar")
	enabled := ctx.Request().InputInt("enabled", 1)
	roleIds := ctx.Request().Input("role_ids")

	// 添加调试日志查看获取到的参数
	facades.Log().Info("QuickEditItem获取到的参数", map[string]any{
		"username": username,
		"password": password,
		"name":     name,
		"avatar":   avatar,
		"enabled":  enabled,
		"role_ids": roleIds,
	})

	return a.Success(ctx)
}

func (a *AdminUserService) Update(ctx http.Context) http.Response {
	// 获取用户ID
	id := ctx.Request().Input("id")
	if id == "" {
		return a.FailMsg(ctx, "用户ID不能为空")
	}

	// 获取请求参数
	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")
	name := ctx.Request().Input("name")
	avatar := ctx.Request().Input("avatar")
	enabled := ctx.Request().InputInt("enabled", 1)
	roleIds := ctx.Request().Input("roles") // 角色ID数组

	// 查找用户
	var user admin.AdminUser
	if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
		return a.FailMsg(ctx, "用户不存在")
	}

	// 检查用户名是否被其他用户使用
	if username != "" && username != user.Username {
		var existingUser admin.AdminUser
		if err := facades.Orm().Query().Where("username", username).Where("id", "!=", id).First(&existingUser); err == nil {
			return a.FailMsg(ctx, "用户名已存在")
		}
	}

	// 开始事务
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return a.FailMsg(ctx, "开始事务失败: "+err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户基本信息
	if username != "" {
		user.Username = username
	}
	if name != "" {
		user.Name = name
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	user.Enabled = uint8(enabled)

	// 如果有新密码，则加密后更新
	if password != "" {
		hashedPassword, err := facades.Hash().Make(password)
		if err != nil {
			tx.Rollback()
			return a.FailMsg(ctx, "密码加密失败: "+err.Error())
		}
		user.Password = hashedPassword
	}

	// 保存用户信息
	if err := tx.Save(&user); err != nil {
		tx.Rollback()
		return a.FailMsg(ctx, "更新用户失败: "+err.Error())
	}

	// 更新角色关联
	if roleIds != "" {
		// 先删除现有的角色关联
		if _, err := tx.Where("user_id", user.ID).Delete(&admin.AdminRoleUser{}); err != nil {
			tx.Rollback()
			return a.FailMsg(ctx, "删除原有角色关联失败: "+err.Error())
		}

		// 处理新的角色ID
		var roleIdList []interface{}
		parts := strings.Split(roleIds, ",")
		for _, part := range parts {
			if strings.TrimSpace(part) != "" {
				roleIdList = append(roleIdList, strings.TrimSpace(part))
			}
		}

		// 获取角色并重新关联
		if len(roleIdList) > 0 {
			var roles []admin.AdminRole
			if err := tx.WhereIn("id", roleIdList).Find(&roles); err != nil {
				tx.Rollback()
				return a.FailMsg(ctx, "获取角色失败: "+err.Error())
			}

			// 使用中间表模型直接插入关联关系
			for _, role := range roles {
				roleUser := &admin.AdminRoleUser{
					RoleID: role.ID,
					UserID: user.ID,
				}
				if err := tx.Create(roleUser); err != nil {
					tx.Rollback()
					return a.FailMsg(ctx, "分配角色失败: "+err.Error())
				}
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return a.FailMsg(ctx, "提交事务失败: "+err.Error())
	}

	return a.SuccessData(ctx, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"enabled":  user.Enabled,
	})
}

func (a *AdminUserService) Store(ctx http.Context) http.Response {
	// 获取请求参数
	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")
	name := ctx.Request().Input("name")
	avatar := ctx.Request().Input("avatar")
	enabled := ctx.Request().InputInt("enabled", 1) // 默认启用
	roles := ctx.Request().InputArray("roles") // 角色ID数组

	// 检查用户名是否已存在
	var existingUser admin.AdminUser
	if err := facades.Orm().Query().Where("username", username).First(&existingUser); err != nil {
		return a.FailMsg(ctx, err.Error())
	}

	// 密码加密
	hashedPassword, err := facades.Hash().Make(password)
	if err != nil {
		return a.FailMsg(ctx, "密码加密失败: "+err.Error())
	}

	// 创建用户
	user := &admin.AdminUser{
		Username: username,
		Password: hashedPassword,
		Name:     name,
		Avatar:   avatar,
		Enabled:  uint8(enabled),
	}

	// 开始事务
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return a.FailMsg(ctx, "开始事务失败: "+err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存用户
	if err := tx.Create(&user); err != nil {
		tx.Rollback()
		return a.FailMsg(ctx, "创建用户失败: "+err.Error())
	}

	// 分配角色
	if len(roles) > 0 {
		// 处理角色ID，支持字符串数组或逗号分隔的字符串
		var roleIdList []interface{}

		// 将角色ID数组转换为interface{}数组
		for _, roleId := range roles {
			if roleId != "" {
				roleIdList = append(roleIdList, roleId)
			}
		}

		// 获取角色并关联
		if len(roleIdList) > 0 {
			var validRoles []admin.AdminRole
			if err := tx.WhereIn("id", roleIdList).Find(&validRoles); err != nil {
				tx.Rollback()
				return a.FailMsg(ctx, "获取角色失败: "+err.Error())
			}

			// 使用模型关联方法直接关联角色
			if err := tx.Model(&user).Association("AdminRoles").Append(&validRoles); err != nil {
				tx.Rollback()
				return a.FailMsg(ctx, "分配角色失败: "+err.Error())
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return a.FailMsg(ctx, "提交事务失败: "+err.Error())
	}

	return a.SuccessData(ctx, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"name":     user.Name,
		"enabled":  user.Enabled,
	})
}

func (a *AdminUserService) Destroy(ctx http.Context) http.Response {
	// 获取用户ID
	id := ctx.Request().InputInt("id")
	if id == 0 {
		return a.FailMsg(ctx, "用户ID不能为空")
	}

	// 查找用户
	var user admin.AdminUser
	if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
		return a.FailMsg(ctx, err.Error())
	}

	if user.IsAdministrator() {
		return a.FailMsg(ctx, "超级管理员不能被删除")
	}

	// 开始事务
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return a.FailMsg(ctx, "开始事务失败: "+err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除用户角色关联
	if _, err := tx.Where("user_id", id).Delete(&admin.AdminRoleUser{}); err != nil {
		tx.Rollback()
		return a.FailMsg(ctx, "删除用户角色关联失败: "+err.Error())
	}

	// 删除用户
	if _, err := tx.Delete(&user); err != nil {
		tx.Rollback()
		return a.FailMsg(ctx, "删除用户失败: "+err.Error())
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return a.FailMsg(ctx, "提交事务失败: "+err.Error())
	}

	return a.SuccessMsg(ctx, "用户删除成功")
}
