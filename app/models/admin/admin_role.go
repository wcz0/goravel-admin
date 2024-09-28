package admin

import (
	"github.com/goravel/framework/database/orm"
)

type AdminRole struct {
	Name string
	Slug string
	orm.Model
	AdminUsers []*AdminUser `gorm:"many2many:admin_role_users;joinForeignKey:role_id;joinReferences:user_id"`
	AdminPermissions []*AdminPermission `gorm:"many2many:admin_role_permissions;joinForeignKey:role_id;joinReferences:permission_id"`
}
