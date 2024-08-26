package models

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/database/orm"
)

type AdminPermission struct {
	ParentId uint
	Name string
	Value string
	Method string
	orm.Model
	AdminRoles []*AdminRole `gorm:"many2many:admin_role_permissions;joinForeignKey:permission_id;joinReferences:role_id"`
}

func NewPermission () *AdminPermission {
	return &AdminPermission{}
}

func (a *AdminPermission) ShouldPassThrough(ctx http.Context) bool {
	
}