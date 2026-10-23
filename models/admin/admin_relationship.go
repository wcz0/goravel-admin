package admin

import "time"

// AdminRelationship 记录「动态关联」配置。
//
// Goravel 的 ORM 关联必须在 struct 上静态声明，没有 Eloquent 那样的运行时
// resolveRelationUsing，因此这张表在 Go 版的定位是「关联元数据登记 +
// 代码预览」：后台可以把关联关系登记清楚、导出成 struct 方法，但不会像
// PHP 版那样登记后立刻生效。
type AdminRelationship struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Model     string    `gorm:"not null;comment:模型" json:"model"`
	Title     string    `gorm:"not null;comment:关联名称" json:"title"`
	Type      string    `gorm:"not null;comment:关联类型" json:"type"`
	Remark    *string   `gorm:"comment:备注" json:"remark"`
	Args      *string   `gorm:"type:text;comment:关联参数" json:"args"`
	Extra     *string   `gorm:"type:text;comment:额外参数" json:"extra"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func NewAdminRelationship() *AdminRelationship {
	return &AdminRelationship{}
}

// PrimaryKey 返回模型主键字段名
func (a *AdminRelationship) PrimaryKey() string {
	return "id"
}
