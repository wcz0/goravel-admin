package admin

import (
	"github.com/goravel/framework/facades"
	"time"
)

const (
	Enabled_OFF uint8 = iota  // 0
	Enabled_ON                // 1
)

type AdminUser struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	Username      string      `gorm:"unique;not null;comment:用户名" json:"username"`
	Password      string      `gorm:"not null;comment:密码" json:"password"`
	Enabled       uint8       `gorm:"type:tinyint unsigned;default:1;comment:启用状态(1:启用,0:禁用)" json:"enabled"`
	Name          string      `gorm:"not null;comment:姓名" json:"name"`
	Avatar        string      `gorm:"comment:头像" json:"avatar"`
	RememberToken string      `gorm:"comment:记住登录" json:"remember_token"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
	AdminRoles []*AdminRole `gorm:"many2many:admin_role_users;joinForeignKey:user_id;joinReferences:role_id"`
}

func NewAdminUser() *AdminUser {
	return &AdminUser{}
}

// IsAdministrator 检查用户是否为超级管理员
func (a *AdminUser) IsAdministrator() bool {
	hasRole, err := a.IsRole("administrator")
	if err != nil {
		return false
	}
	return hasRole
}

// IsRole 检查用户是否拥有指定角色
func (a *AdminUser) IsRole(role string) (bool, error) {
	var adminRoles []AdminRole
	if err := facades.Orm().Query().Model(a).Association("AdminRoles").Find(&adminRoles); err != nil {
		return false, err
	}
	for _, v := range adminRoles {
		if v.Slug == role {
			return true, nil
		}
	}
	return false, nil
}

// HasRole 通过 Slug 检查用户是否拥有指定角色
func (a *AdminUser) HasRole(slug string) (bool, error) {
	return a.IsRole(slug)
}

// AttachRole 为用户分配角色
func (a *AdminUser) AttachRole(role *AdminRole) error {
	return facades.Orm().Query().Model(a).Association("AdminRoles").Append(role)
}

// DetachRole 移除用户角色
func (a *AdminUser) DetachRole(role *AdminRole) error {
	return facades.Orm().Query().Model(a).Association("AdminRoles").Delete(role)
}

// SyncRoles 同步用户角色
func (a *AdminUser) SyncRoles(roles []*AdminRole) error {
	// 先删除所有关联
	if err := facades.Orm().Query().Model(a).Association("AdminRoles").Clear(); err != nil {
		return err
	}

	// 重新添加所有角色
	for _, role := range roles {
		if err := facades.Orm().Query().Model(a).Association("AdminRoles").Append(role); err != nil {
			return err
		}
	}
	return nil
}

// AllPermissions 获取用户所有权限
// AllPermissions 汇总用户所有角色下的权限（按权限 ID 去重）。
//
// 之前这里写成 Model(a).With("AdminPermissions").Association("AdminRoles").Find(&permissions)：
// 关联名是 AdminRoles，接收的却是 []AdminPermission，类型对不上必然报错，
// 于是非超管用户永远拿到空权限、所有接口 403。
func (a *AdminUser) AllPermissions() []AdminPermission {
	var roles []AdminRole
	if err := facades.Orm().Query().Model(a).With("AdminPermissions").Association("AdminRoles").Find(&roles); err != nil {
		return nil
	}

	seen := make(map[uint]bool)
	result := make([]AdminPermission, 0, len(roles))
	for i := range roles {
		for _, permission := range roles[i].AdminPermissions {
			if permission == nil || seen[permission.ID] {
				continue
			}
			seen[permission.ID] = true
			result = append(result, *permission)
		}
	}

	return result
}

// HasPermission 检查用户是否拥有指定权限
func (a *AdminUser) HasPermission(permission *AdminPermission) (bool, error) {
	permissions := a.AllPermissions()
	for _, p := range permissions {
		if p.ID == permission.ID {
			return true, nil
		}
	}
	return false, nil
}

// PrimaryKey 返回模型主键字段名
func (a *AdminUser) PrimaryKey() string {
	return "id"
}