package services

import (
	"errors"
	"fmt"
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

// List 返回用户列表（分页）。
//
// amis 的 CRUD 依赖 total 才能渲染分页条，缺了它前端会一直停在第 1 页并显示
// 全量数据；角色关联以前是逐行查一次（N+1），这里改为 With 预加载。
func (a *AdminUserService) List(ctx http.Context) http.Response {
	page := clampPage(ctx.Request().InputInt("page", 1))
	perPage := clampPerPage(ctx.Request().InputInt("perPage", defaultPageSize))

	query := facades.Orm().Query().Model(&admin.AdminUser{})

	// 关键字同时匹配用户名与姓名，必须用括号包住 OR，否则会逃出其它条件。
	if keyword := strings.TrimSpace(ctx.Request().Input("keyword")); keyword != "" {
		query = query.WhereAny([]string{"username", "name"}, "like", "%"+keyword+"%")
	}

	total, err := query.Count()
	if err != nil {
		return a.Fail(ctx, err.Error())
	}

	var users []admin.AdminUser
	if err := query.
		Select("id", "username", "name", "avatar", "enabled", "created_at").
		With("AdminRoles").
		Order("id asc").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Get(&users); err != nil {
		return a.Fail(ctx, err.Error())
	}

	items := make([]map[string]any, 0, len(users))
	loc := time.FixedZone("CST", 8*3600)
	for _, u := range users {
		roleItems := make([]map[string]any, 0, len(u.AdminRoles))
		// 关联表历史上没有唯一约束，重复行会让前端出现重复角色标签，这里按角色 ID 去重。
		seen := make(map[uint]bool, len(u.AdminRoles))
		for _, r := range u.AdminRoles {
			if r == nil || seen[r.ID] {
				continue
			}
			seen[r.ID] = true
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
		"items":   items,
		"total":   total,
		"page":    page,
		"perPage": perPage,
	})
}

// Export 导出用户列表为 xlsx。
//
// 按当前筛选条件取全量数据（不分页），返回签名后的下载信息。
func (a *AdminUserService) Export(ctx http.Context) http.Response {
	query := facades.Orm().Query().Select("id", "username", "name", "enabled", "created_at").Order("id asc")

	if ids := ctx.Request().Input("_ids"); ids != "" {
		query = query.WhereIn("id", parseIDList(ids))
	}
	// 关键字必须用 WhereAny 包成 (username like ? or name like ?)：
	// 直接 Where(...).OrWhere(...) 会让 OR 逃出上面的 _ids 限制，
	// 勾选两行导出却把全表都导出来。
	if keyword := strings.TrimSpace(ctx.Request().Input("keyword")); keyword != "" {
		query = query.WhereAny([]string{"username", "name"}, "like", "%"+keyword+"%")
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
		query = query.Where("slug", "!=", "administrator")
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
	if err := firstByID(id, &user); err != nil {
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
	if err := firstByID(id, &user); err != nil {
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
	if err := firstByID(id, &user); err != nil {
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

	if roleIDs := roleIDsFrom(ctx); len(roleIDs) > 0 {
		if err := user.SyncRoles(loadRoles(roleIDs)); err != nil {
			return a.Fail(ctx, err.Error())
		}
	}

	return a.Success(ctx, "更新成功")
}

// roleIDsFrom 从请求里取出角色 ID，兼容 ["1","2"] 与 [{"id":1}] 两种前端写法，并去重。
//
// 去重是必须的：关联表没有唯一约束，同一个角色 ID 出现两次会插两行，
// 用户列表就会渲染出重复的角色标签。
func roleIDsFrom(ctx http.Context) []any {
	ids := make([]any, 0, 4)
	var appendID func(v any)
	appendID = func(v any) {
		switch it := v.(type) {
		case nil:
			return
		case string:
			if strings.TrimSpace(it) == "" {
				return
			}
			ids = append(ids, it)
		case float64, int, int64, uint:
			ids = append(ids, it)
		case map[string]any:
			if idVal, ok := it["id"]; ok {
				appendID(idVal)
			}
		}
	}

	for _, v := range ctx.Request().InputArray("roles") {
		appendID(v)
	}
	if len(ids) == 0 {
		if raw, ok := ctx.Request().All()["roles"]; ok {
			if arr, ok := raw.([]any); ok {
				for _, item := range arr {
					appendID(item)
				}
			}
		}
	}

	seen := make(map[string]bool, len(ids))
	unique := ids[:0]
	for _, id := range ids {
		key := fmt.Sprintf("%v", id)
		if seen[key] {
			continue
		}
		seen[key] = true
		unique = append(unique, id)
	}

	return unique
}

// loadRoles 按 ID 取角色实体，供 SyncRoles 使用。
func loadRoles(ids []any) []*admin.AdminRole {
	var roles []admin.AdminRole
	if err := facades.Orm().Query().WhereIn("id", ids).Find(&roles); err != nil {
		return nil
	}

	result := make([]*admin.AdminRole, 0, len(roles))
	for i := range roles {
		result = append(result, &roles[i])
	}

	return result
}

func (a *AdminUserService) Store(ctx http.Context) http.Response {
	username := ctx.Request().Input("username")
	password := ctx.Request().Input("password")
	name := ctx.Request().Input("name")
	avatar := ctx.Request().Input("avatar")
	enabled := ctx.Request().InputInt("enabled", 1)


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

	if roleIDs := roleIDsFrom(ctx); len(roleIDs) > 0 {
		if err := user.SyncRoles(loadRoles(roleIDs)); err != nil {
			return a.Fail(ctx, err.Error())
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
	if err := firstByID(id, &user); err != nil {
		return a.Fail(ctx, "用户不存在")
	}
	if _, err := facades.Orm().Query().Delete(&user); err != nil {
		return a.Fail(ctx, err.Error())
	}
	return a.Success(ctx, "删除成功")
}
