package admin

import "time"

const (
	Extension_Enabled_OFF uint8 = iota
	Extension_Enabled_ON
)

// AdminExtension 记录扩展的启用状态。
//
// Go 是静态编译，无法像 PHP 那样把扩展目录丢进去就动态加载，因此这里的
// 表只承担「已安装扩展名单 + 启用状态」的语义：扩展本身需要在宿主编译期
// 注册，启停只影响菜单、静态资源与 _settings 下发的 enabled_extensions。
type AdminExtension struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null;comment:扩展名称" json:"name"`
	IsEnabled uint8     `gorm:"default:0;comment:是否启用" json:"is_enabled"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func NewAdminExtension() *AdminExtension {
	return &AdminExtension{}
}

// PrimaryKey 返回模型主键字段名
func (a *AdminExtension) PrimaryKey() string {
	return "id"
}
