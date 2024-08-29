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
}

func NewPermission () *AdminPermission {
	return &AdminPermission{}
}

func (a *AdminPermission) ShouldPassThrough(ctx http.Context) bool {
	if len(a.HttpPath) == 0  &&  len(a.HttpMethod) == 0 {
		return true
	}
	routePrefix := facades.Config().GetString("admin.route.prefix")
	for _, path := range a.HttpPath {
		path := strings.Trim(routePrefix, "/") + path
		if strings.Contains(path, ":") {

		}
	}
	// TODO: 未写完
	return true
}

func (a *AdminPermission) matchRequest(ctx http.Context, match map[string]any) bool {
	path := match["path"].(string)
	if path == "/" {
		path = "/"
	}else {
		path = strings.Trim(path, "/")
	}
	if !tools.RequestIs(ctx, path) {
		return false
	}
	methods, ok := match["methods"].([]string)
	if !ok || len(methods) == 0 {
		return true
	}
	for _, m := range methods {
		if strings.ToUpper(m) == ctx.Request().Method() {
			return true
		}
	}
	return false
}
