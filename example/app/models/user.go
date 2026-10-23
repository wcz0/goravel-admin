package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type User struct {
	Id uint64 `gorm:"primaryKey"`
	Username string
	Name   string
	Birthday carbon.Date
	CardId string
	Avatar string
	UpdatedAt carbon.DateTime
	CreatedAt carbon.DateTime
	orm.SoftDeletes
}
