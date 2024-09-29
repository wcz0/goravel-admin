package admin

import (
	"goravel/app/models"
	"goravel/app/tools"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/facades"
)

type AdminPermission struct {
	Name string
	Slug string
	HttpMethod models.StringSlice `gorm:"type:json"`
	HttpPath models.StringSlice `gorm:"type:json"`
	CustomOrder int
	ParentId uint
	orm.Model
	AdminRoles []*AdminRole `gorm:"many2many:admin_role_permissions;joinForeignKey:permission_id;joinReferences:role_id"`
	AdminMenus []*AdminMenu `gorm:"many2many:admin_permission_menu;joinForeignKey:permission_id;joinReferences:menu_id"`
}

func NewPermission () *AdminPermission {
	return &AdminPermission{}
}

func (a *AdminPermission) ShouldPassThrough(ctx http.Context) bool {
	if len(a.HttpPath) == 0  &&  len(a.HttpMethod) == 0 {
		return true
	}
	routePrefix := facades.Config().GetString("admin.route.prefix")
	// 遍历path, 查看请求是否在其中
	for _, path := range a.HttpPath {
		path := strings.TrimSuffix(routePrefix, "/") + path
		if tools.RequestIs(ctx, path) {
			// 遍历method, 查看请求是否在其中
			for _, method := range a.HttpMethod {
				if tools.RequestMethodIs(ctx, method) {
					return true
				}
			}
		}
	}
	return false
}


func (a *AdminPermission) PrimaryKey() string {
	return "id"
}