package models

import (
	"github.com/goravel/framework/database/orm"
)

type AdminMenu struct {
	ParentId uint32
	Title string
	Icon string
	Uri string
	UrlType uint8
	Visible uint8
	IsHome uint8
	Component string
	IsFull uint8
	Extension string
	orm.Model
}