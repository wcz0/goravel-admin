package admin

import (
	"github.com/goravel/framework/database/orm"
)

const (
	TYPE_ROUTE = iota + 1
	TYPE_LINK
	TYPE_IFRAME
	TYPE_PAGE
)

type AdminMenu struct {
	ParentId uint32
	Title string
	Icon string
	Url string
	UrlType uint8
	Visible uint8
	IsHome uint8
	KeepAlive uint8
	Component string
	CustomOrder int
	IFrameUrl string `json:"iframe_url"`
	IsFull uint8
	Extension string
	orm.Model
}