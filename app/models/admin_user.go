package models

import "github.com/goravel/framework/database/orm"

type AdminUser struct {
	Username string
	Password string
	Name     string
	Avatar   string
	RememberToken string
	orm.Model
}