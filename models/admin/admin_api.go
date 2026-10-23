package admin

import "time"

// AdminApi 对应「动态 API」里的一条接口配置。
//
// template 指向内置的动态接口模板（列表 / 详情 / 新增 / 更新 / 删除 / 排序等），
// args 以 JSON 形式保存模板所需参数。与 owl-admin 的差异在于模板必须是编译期
// 注册的 handler，不能像 PHP 那样按类名动态实例化。
type AdminApi struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null;comment:接口名称" json:"title"`
	Path      string    `gorm:"not null;comment:接口路径" json:"path"`
	Template  string    `gorm:"not null;comment:接口模板" json:"template"`
	Enabled   uint8     `gorm:"default:1;comment:是否启用" json:"enabled"`
	Args      *string   `gorm:"type:text;comment:接口参数" json:"args"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func NewAdminApi() *AdminApi {
	return &AdminApi{}
}

// PrimaryKey 返回模型主键字段名
func (a *AdminApi) PrimaryKey() string {
	return "id"
}
