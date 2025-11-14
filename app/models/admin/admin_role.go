package admin

import (
	"github.com/goravel/framework/facades"
	"time"
)

type AdminRole struct {
	ID   uint
	Name string `gorm:"unique"`
	Slug string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	AdminUsers []*AdminUser `gorm:"many2many:admin_role_users;joinForeignKey:role_id;joinReferences:user_id"`
	AdminPermissions []*AdminPermission `gorm:"many2many:admin_role_permissions;joinForeignKey:role_id;joinReferences:permission_id"`
}

func NewAdminRole() *AdminRole {
	return &AdminRole{}
}

// IsAdministrator 检查是否为超级管理员
func (r *AdminRole) IsAdministrator() bool {
	return r.Slug == "administrator"
}

// AllPermissions 获取角色的所有权限
func (r *AdminRole) AllPermissions() []AdminPermission {
	var permissions []AdminPermission
	if err := facades.Orm().Query().Model(r).Association("AdminPermissions").Find(&permissions); err != nil {
		return nil
	}
	return permissions
}

// AttachPermission 为角色分配权限
func (r *AdminRole) AttachPermission(permission *AdminPermission) error {
	return facades.Orm().Query().Model(r).Association("AdminPermissions").Append(permission)
}

// DetachPermission 移除角色权限
func (r *AdminRole) DetachPermission(permission *AdminPermission) error {
	return facades.Orm().Query().Model(r).Association("AdminPermissions").Delete(permission)
}

// SyncPermissions 同步角色权限
func (r *AdminRole) SyncPermissions(permissions []*AdminPermission) error {
	// 先删除所有关联
	if err := facades.Orm().Query().Model(r).Association("AdminPermissions").Clear(); err != nil {
		return err
	}
	
	// 重新添加所有权限
	for _, permission := range permissions {
		if err := facades.Orm().Query().Model(r).Association("AdminPermissions").Append(permission); err != nil {
			return err
		}
	}
	return nil
}

// HasUser 检查角色是否拥有指定用户
func (r *AdminRole) HasUser(user *AdminUser) (bool, error) {
	var users []AdminUser
	if err := facades.Orm().Query().Model(r).Association("AdminUsers").Find(&users); err != nil {
		return false, err
	}
	for _, u := range users {
		if u.ID == user.ID {
			return true, nil
		}
	}
	return false, nil
}

// PrimaryKey 返回模型主键字段名
func (r *AdminRole) PrimaryKey() string {
	return "id"
}
