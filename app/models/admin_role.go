package models

import (
	"github.com/goravel/framework/database/orm"
)

type AdminRole struct {
	Name string
	Value string
	orm.Model
}
