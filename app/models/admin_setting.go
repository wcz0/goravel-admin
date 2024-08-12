package models

import (
	"github.com/goravel/framework/support/carbon"
)

type AdminSetting struct {
	Key       string
	Values    string
	CreatedAt carbon.DateTime
	UpdateAt  carbon.DateTime
}

func NewAdminSetting() *AdminSetting {
	return &AdminSetting{}
}
