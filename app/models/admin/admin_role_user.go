package admin

import (
	"github.com/goravel/framework/database/orm"
)

// AdminRoleUser 角色用户关联表模型
type AdminRoleUser struct {
	RoleID uint `gorm:"primaryKey;column:role_id"`
	UserID uint `gorm:"primaryKey;column:user_id"`
	orm.Timestamps
}

// TableName 指定表名
func (AdminRoleUser) TableName() string {
	return "admin_role_users"
}