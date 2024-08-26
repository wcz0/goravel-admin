package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
)

type AdminUser struct {
	Username      string
	Password      string
	Enabled       int8
	Name          string
	Avatar        string
	RememberToken string
	orm.Model
	AdminRoles []*AdminRole `gorm:"many2many:admin_role_users;joinForeignKey:user_id;joinReferences:role_id"`
}

func NewAdminUser() *AdminUser {
	return &AdminUser{}
}

func (a *AdminUser) IsAdministrator() bool {
	bool, _ := a.IsRole("administrator")
	return bool
}

func (a *AdminUser) IsRole(role string) (bool, error) {
	var adminRole []AdminRole
	if err := facades.Orm().Query().Model(a).Association("AdminRoles").Find(&adminRole); err != nil {
		return false, err
	}
	for _, v := range adminRole {
		if v.Name == role {
			return true, nil
		}
	}
	return false, nil
}

func (a *AdminUser) AllPermissions() []AdminPermission {
	var adminPermissions []AdminPermission
	if err := facades.Orm().Query().Model(a).With("AdminPermissions").Association("AdminRoles").Find(&adminPermissions); err != nil {
		return nil
	}
	return adminPermissions
}