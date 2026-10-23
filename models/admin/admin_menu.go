package admin

import (
	"time"
)

const (
	TYPE_ROUTE = iota + 1
	TYPE_LINK
	TYPE_IFRAME
	TYPE_PAGE
)

const (
	IS_HOME_ON = 1
	IS_HOME_OFF = 0
)

type AdminMenu struct {
	ID          uint
	ParentId    uint32
	Title       string
	Icon        string
	Url         string
	UrlType     uint8
	Visible     uint8
	IsHome      uint8
	Component   string
	CustomOrder int `gorm:"column:custom_order"`
	KeepAlive   *uint8 `gorm:"column:keep_alive"`
	IFrameUrl   *string `gorm:"column:iframe_url"`
	IsFull      uint8
	Extension   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
