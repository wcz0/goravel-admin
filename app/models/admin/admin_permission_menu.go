package admin

import "github.com/goravel/framework/database/orm"

// AdminPermissionMenu 权限菜单关联表模型
type AdminPermissionMenu struct {
	PermissionID uint `gorm:"primaryKey;column:permission_id"`
	MenuID       uint `gorm:"primaryKey;column:menu_id"`
	orm.Timestamps
}

// TableName 指定表名
func (AdminPermissionMenu) TableName() string {
	return "admin_permission_menu"
}