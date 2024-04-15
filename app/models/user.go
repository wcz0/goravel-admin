package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type User struct {
	orm.Model
	Id uint64
	Username string
	Name   string
	Birthday carbon.Date
	CardId string
	Avatar string
	orm.SoftDeletes
}
