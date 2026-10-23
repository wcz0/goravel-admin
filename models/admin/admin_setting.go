package admin

import (
	"github.com/goravel/framework/database/orm"

)

type AdminSetting struct {
	Key       string 
	Values    string
	orm.Timestamps
}

func NewAdminSetting() *AdminSetting {
	return &AdminSetting{}
}
