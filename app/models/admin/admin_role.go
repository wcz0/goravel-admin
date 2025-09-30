package admin

import (
	"github.com/goravel/framework/database/orm"
)

// AdminRole 角色模型
type AdminRole struct {
	Name string `gorm:"size:50;not null;unique"`
	Slug string `gorm:"size:50;not null;unique"`
	orm.Model

	// 关联关系
	AdminUsers       []*AdminUser       `gorm:"many2many:admin_role_users;joinForeignKey:role_id;joinReferences:user_id"`
	AdminPermissions []*AdminPermission `gorm:"many2many:admin_role_permissions;joinForeignKey:role_id;joinReferences:permission_id"`
}
