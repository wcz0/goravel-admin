package admin

import "time"

// AdminPage 对应「页面管理」里的自定义页面。
//
// 页面结构以 amis JSON schema 的形式存在库中，前端通过 page_schema?sign=xxx
// 拉取并直接渲染，因此不需要为每个页面单独写控制器。
type AdminPage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null;comment:页面名称" json:"title"`
	Sign      string    `gorm:"not null;comment:页面标识" json:"sign"`
	Schema    string    `gorm:"type:text;comment:页面结构" json:"schema"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func NewAdminPage() *AdminPage {
	return &AdminPage{}
}

// PrimaryKey 返回模型主键字段名
func (a *AdminPage) PrimaryKey() string {
	return "id"
}
