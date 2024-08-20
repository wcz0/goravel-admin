package models

import "github.com/goravel/framework/database/orm"

type AdminUser struct {
	Username string
	Password string
	Enabled  int8
	Name     string
	Avatar   string
	RememberToken string
	orm.Model
	AdminRole []*AdminRole
}