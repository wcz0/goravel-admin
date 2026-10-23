package admin

import (
	"github.com/goravel/framework/database/orm"
)

// AdminRolePermission 角色权限关联表模型
type AdminRolePermission struct {
	RoleID       uint `gorm:"primaryKey;column:role_id"`
	PermissionID uint `gorm:"primaryKey;column:permission_id"`
	orm.Timestamps
}

// TableName 指定表名
func (AdminRolePermission) TableName() string {
	return "admin_role_permissions"
}