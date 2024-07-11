package models

import "github.com/goravel/framework/database/orm"

type AdminPermission struct {
	ParentId uint
	Name string
	Value string
	Method string
	orm.Model
}