package services

import (
	"errors"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/goravel-admin/models/admin"
	"github.com/wcz0/goravel-admin/tools"
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
		return s.Fail(ctx, err.Error())
	}
	if !facades.Hash().Check(ctx.Request().Input("password"), adminUser.Password) {
		return s.Fail(ctx, "Password error.")
	}
	token, err := facades.Auth(ctx).Guard(tools.AuthGuard()).Login(&adminUser)
	if err != nil {
		return s.Fail(ctx, err.Error())
	}
	return s.SuccessData(ctx, map[string]any{
		"token":   token,
		"message": "登录成功",
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
		return a.Fail(ctx, err.Error())
	}

	items := make([]map[string]any, 0, len(users))
	loc := time.FixedZone("CST", 8*3600)
	for _, u := range users {
		var roles []admin.AdminRole
		_ = facades.Orm().Query().Model(&u).Association("AdminRoles").Find(&roles)
		roleItems := make([]map[string]any, 0, len(roles))
		for _, r := range roles {
			roleItems = append(roleItems, map[string]any{"id": r.ID, "name": r.Name})
		}
		items = append(items, map[string]any{
			"id":        u.ID,
			"username":  u.Username,
			"name":      u.Name,
			"avatar":    u.Avatar,
			"enabled":   u.Enabled,
			"createdAt": u.CreatedAt.In(loc),
			"roles":     roleItems,
		})
	}

	return a.SuccessData(ctx, map[string]any{
		"items": items,
	})
}

// Export 导出用户列表为 xlsx。
//
// 按当前筛选条件取全量数据（不分页），返回签名后的下载信息。
func (a *AdminUserService) Export(ctx http.Context) http.Response {
	query := facades.Orm().Query().Select("id", "username", "name", "enabled", "created_at").Order("id asc")

	if ids := ctx.Request().Input("_ids"); ids != "" {
		query.WhereIn("id", parseIDList(ids))
	}
	if keyword := ctx.Request().Input("keyword"); keyword != "" {
		query.Where("username", "like", "%"+keyword+"%").
			OrWhere("name", "like", "%"+keyword+"%")
	}

	var users []admin.AdminUser
	if err := query.Get(&users); err != nil {
		return a.Fail(ctx, err.Error())
	}

	headers := []string{"ID", "用户名", "姓名", "状态", "创建时间"}
	rows := make([][]any, 0, len(users))
	for _, u := range users {
		status := "禁用"
		if u.Enabled == admin.Enabled_ON {
			status = "启用"
		}
		rows = append(rows, []any{u.ID, u.Username, u.Name, status, u.CreatedAt.Format("2006-01-02 15:04:05")})
	}

	disk := facades.Config().GetString("admin.upload.disk")
	result, err := tools.BuildExport(disk, "users", headers, rows)
	if err != nil {
		return a.Fail(ctx, err.Error())
	}

	return a.SuccessData(ctx, map[string]any{
		"path":      result.Path,
		"expires":   result.Expires,
		"signature": result.Signature,
	})
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
	id := ctx.Request().InputInt("id")
	if id == 0 {
		return a.Fail(ctx, "ID不能为空")
	}
	var user admin.AdminUser
	if err := facades.Orm().Query().Find(&user, id); err != nil {
		return a.Fail(ctx, "用户不存在")
	}

	enabledVal := ctx.Request().Input("Enabled")
	if enabledVal == "" {
		enabledVal = ctx.Request().Input("enabled")
	}
	if enabledVal != "" {
		switch strings.ToLower(enabledVal) {
		case "true", "1":
			user.Enabled = 1
		case "false", "0":
			user.Enabled = 0
		default:
			return a.Fail(ctx, "enabled 参数不合法")
		}
		if err := facades.Orm().Query().Save(&user); err != nil {
			return a.Fail(ctx, err.Error())
		}
		return a.Success(ctx, "快速编辑成功")
	}

	return a.Fail(ctx, "未识别的快速编辑字段")
}

func (a *AdminUserService) QuickEditItem(ctx http.Context) http.Response {
	return a.QuickEdit(ctx)
}

// UpdateUserSetting 更新当前登录用户的个人资料。
//
// 与后台用户管理不同，这里只允许改头像、姓名和密码，且改密码必须校验旧密码，
// 避免拿到 token 后可直接改密。
func (a *AdminUserService) UpdateUserSetting(id uint, data map[string]any) error {
	var user admin.AdminUser
	if err := facades.Orm().Query().Find(&user, id); err != nil {
		return errors.New("用户不存在")
	}

	if name, ok := data["name"].(string); ok && strings.TrimSpace(name) != "" {
		user.Name = name
	}
	if avatar, ok := data["avatar"].(string); ok {
		user.Avatar = avatar
	}

	password, _ := data["password"].(string)
	if password != "" {
		confirm, _ := data["confirm_password"].(string)
		if password != confirm {
			return errors.New("两次输入的密码不一致")
		}
		oldPassword, _ := data["old_password"].(string)
		if oldPassword == "" {
			return errors.New("请输入原密码")
		}
		if !facades.Hash().Check(oldPassword, user.Password) {
			return errors.New("原密码错误")
		}
		hashed, err := facades.Hash().Make(password)
		if err != nil {
			return err
		}
		user.Password = hashed
	}

	return facades.Orm().Query().Save(&user)
}

func (a *AdminUserService) Update(ctx http.Context) http.Response {
	id := ctx.Request().InputInt("id")
	if id == 0 {
		if rid := ctx.Request().Route("id"); rid != "" {
			id = ctx.Request().RouteInt("id")
		}
	}
	if id == 0 {
		return a.Fail(ctx, "ID不能为空")
	}

	var user admin.AdminUser
	if err := facades.Orm().Query().Find(&user, id); err != nil {
		return a.Fail(ctx, "用户不存在")
	}

	username := ctx.Request().Input("username")
	name := ctx.Request().Input("name")
	avatar := ctx.Request().Input("avatar")
	enabled := ctx.Request().InputInt("enabled", int(user.Enabled))

	if username != "" {
		total, err := facades.Orm().Query().Model(&admin.AdminUser{}).Where("id != ? AND username = ?", id, username).Count()
		if err != nil {
			return a.Fail(ctx, "校验失败: "+err.Error())
		}
		if total > 0 {
			return a.Fail(ctx, "用户名已存在")
		}
		user.Username = username
	}

	if name != "" {
		user.Name = name
	}
	user.Avatar = avatar
	user.Enabled = uint8(enabled)

	if pwd := ctx.Request().Input("password"); pwd != "" {
		hashed, err := facades.Hash().Make(pwd)
		if err != nil {
			return a.Fail(ctx, err.Error())
		}
		user.Password = hashed
	}

	if err := facades.Orm().Query().Save(&user); err != nil {
		return a.Fail(ctx, err.Error())
	}

	roleIDs := []interface{}{}
	for _, v := range ctx.Request().InputArray("roles") {
		if v == "" {
			continue
		}
		roleIDs = append(roleIDs, v)
	}
	if len(roleIDs) == 0 {
		all := ctx.Request().All()
		if raw, ok := all["roles"]; ok {
			if arr, ok := raw.([]any); ok {
				for _, item := range arr {
					switch it := item.(type) {
					case map[string]any:
						if idVal, ok := it["id"]; ok {
							roleIDs = append(roleIDs, idVal)
						}
					}
				}
			}
		}
	}
	if len(roleIDs) > 0 {
		var roles []admin.AdminRole
		if err := facades.Orm().Query().WhereIn("id", roleIDs).Find(&roles); err == nil {
			rs := make([]*admin.AdminRole, 0, len(roles))
			for i := range roles {
				rs = append(rs, &roles[i])
			}
			_ = user.SyncRoles(rs)
		}
	}

	return a.Success(ctx, "更新成功")
}

func (a *AdminUserService) Store(ctx http.Context) http.Response {
	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")
	name := ctx.Request().Input("name")
	avatar := ctx.Request().Input("avatar")
	enabled := ctx.Request().InputInt("enabled", 1)
	rawRoles := ctx.Request().InputArray("roles")

	// 检查用户名是否已存在
	var existingUser admin.AdminUser
	if err := facades.Orm().Query().Where("username", username).First(&existingUser); err == nil {
		return a.Fail(ctx, "用户名已存在")
	}

	hashed, err := facades.Hash().Make(password)
	if err != nil {
		return a.Fail(ctx, err.Error())
	}

	user := &admin.AdminUser{
		Username: username,
		Password: hashed,
		Name:     name,
		Avatar:   avatar,
		Enabled:  uint8(enabled),
	}

	if err := facades.Orm().Query().Create(user); err != nil {
		return a.Fail(ctx, err.Error())
	}

	if len(rawRoles) > 0 {
		var roles []admin.AdminRole
		if err := facades.Orm().Query().Where("id IN ?", rawRoles).Get(&roles); err == nil {
			// 关联角色
			rs := make([]*admin.AdminRole, 0, len(roles))
			for i := range roles {
				rs = append(rs, &roles[i])
			}
			_ = user.SyncRoles(rs)
		}
	}

	return a.Success(ctx, "创建成功")
}

func (a *AdminUserService) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().InputInt("id")
	ridStr := ctx.Request().Route("id")
	if id == 0 && ridStr == "" {
		return a.Fail(ctx, "用户ID不能为空")
	}

	// 处理批量删除：路径 id 为逗号分隔
	if ridStr != "" && strings.Contains(ridStr, ",") {
		parts := strings.Split(ridStr, ",")
		ids := make([]interface{}, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if p == "1" {
				continue
			}
			ids = append(ids, p)
		}
		if len(ids) == 0 {
			return a.Fail(ctx, "没有可删除的用户")
		}
		if _, err := facades.Orm().Query().WhereIn("id", ids).Delete(&admin.AdminUser{}); err != nil {
			return a.Fail(ctx, err.Error())
		}
		return a.Success(ctx, "删除成功")
	}

	if id == 0 {
		id = ctx.Request().RouteInt("id")
	}
	if id == 0 {
		return a.Fail(ctx, "用户ID不能为空")
	}
	if id == 1 {
		return a.Fail(ctx, "超级管理员不可删除")
	}

	var user admin.AdminUser
	if err := facades.Orm().Query().Find(&user, id); err != nil {
		return a.Fail(ctx, "用户不存在")
	}
	if _, err := facades.Orm().Query().Delete(&user); err != nil {
		return a.Fail(ctx, err.Error())
	}
	return a.Success(ctx, "删除成功")
}
