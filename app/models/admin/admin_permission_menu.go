package admin

import "github.com/goravel/framework/database/orm"

type AdminPermissionMenu struct {
	PermissionId uint32
	MenuId uint32
	orm.Model
}