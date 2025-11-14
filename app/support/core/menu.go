package core

import (
	"goravel/app/models/admin"
	"goravel/app/tools"
	"strconv"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type Meta struct {
	Title        string `json:"title"`
	Icon         string `json:"icon"`
	Hide         bool  `json:"hide"`
	// SingleLayout string `json:"singleLayout"`
	Order int `json:"order"`
}

type Route struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Component string `json:"component"`
	IsHome uint8 `json:"is_home"`
	IFrameUrl string `json:"iframe_url"`
	UrlType uint8 `json:"url_type"`
	KeepAlive uint8 `json:"keep_alive"`
	IsFull uint8 `json:"is_full"`
	IsLink bool `json:"is_link"`
	PageSign string `json:"page_sign"`
	Meta      Meta  `json:"meta"`
	Children  []Route `json:"children"`
}

type Menu struct {
	menus []admin.AdminMenu
}

func NewMenu() *Menu {
	return &Menu{
		menus: []admin.AdminMenu{},
	}
}

func (m *Menu) All(ctx http.Context) []Route {
	menus := slice.Unique(m.userMenus(ctx))
	if err := slice.SortByField(menus, "CustomOrder", "asc"); err != nil {
		return []Route{}
	}
	// m.menus = append(m.menus, menus...)
	routes := m.list2Menu(ctx, menus, 0, "")
	return routes
}

func (m *Menu) userMenus(ctx http.Context) []admin.AdminMenu {
	config := facades.Config()
	if !config.GetBool("admin.auth.enable") {
		facades.Log().Debug("menus: auth disabled")
		return []admin.AdminMenu{}
	}
	user := ctx.Value("admin_user").(*admin.AdminUser)
	if user == nil {
		facades.Log().Debug("menus: user is nil")
		return []admin.AdminMenu{}
	}
	
	// 调试信息
	facades.Log().Debug("menus: user", "username", user.Username, "isAdmin", user.IsAdministrator())
	facades.Log().Debug("menus: permission enabled", config.GetBool("admin.auth.permission"))
	
	var list []admin.AdminMenu
	if user.IsAdministrator() || !config.GetBool("admin.auth.permission") {
		facades.Log().Debug("menus: querying all menus")
		if err := facades.Orm().Query().OrderBy("custom_order").Get(&list); err != nil {
			facades.Log().Error("menus: query error", "error", err.Error())
			return []admin.AdminMenu{}
		}
		facades.Log().Debug("menus: found menus", "count", len(list))
	} else {
		// 获取用户角色权限菜单
		facades.Log().Debug("menus: getting user role menus")
		if err := facades.Orm().Query().With("AdminRoles.AdminPermissions.AdminMenus").Find(&user); err != nil {
			facades.Log().Error("menus: user query error", "error", err.Error())
			return []admin.AdminMenu{}
		}
		for _, role := range user.AdminRoles {
			for _, permission := range role.AdminPermissions {
				for _, menu := range permission.AdminMenus {
					list = append(list, *menu)
				}
			}
		}
		facades.Log().Debug("menus: permission menus found", "count", len(list))
	}
	return list
}

func (m *Menu) list2Menu(ctx http.Context, list []admin.AdminMenu, parentId int, parentName string) []Route {
	routes := []Route{}
	for _, v := range list {
		if v.ParentId == uint32(parentId) {
			var _component string
			switch v.UrlType {
				case admin.TYPE_IFRAME:
					_component = "iframe"
				case admin.TYPE_PAGE:
					_component = "amis"
				default:
					if v.Component == "" {
						_component = "amis"
					} else {
						_component = v.Component
					}
				}
			if parentName != "" {
				parentName = strings.TrimRight(parentName, "-") + "-"
			}
			route := Route{
				Name: parentName + "["+ strconv.Itoa(int(v.ID)) +"]" ,
				Path: v.Url,
				Component: _component,
				IsHome: v.IsHome,
				IFrameUrl: v.IFrameUrl,
				UrlType: v.UrlType,
				KeepAlive: v.KeepAlive,
				IsFull: v.IsFull,
				IsLink: v.UrlType == admin.TYPE_LINK,
				PageSign: func () string {
					if v.UrlType == admin.TYPE_PAGE {
						return v.Component
					}
					return ""
				}(),
				Meta: Meta{
					Title: facades.Lang(ctx).Get("menu."+v.Title),
					Icon: v.Icon,
					Hide: v.Visible == 0,
					Order: v.CustomOrder,
				},
			}
			children := m.list2Menu(ctx, list, int(v.ID), route.Name)
			if len(children) > 0 {
				route.Component = "amis"
				route.Children = children
			}
			routes = append(routes, route)
			if slice.IndexOf(facades.Config().Get("admin.route.without_extra_routes").([]string), route.Path) == -1 && v.UrlType != admin.TYPE_PAGE {
				routes = append(routes, m.generateRoute(ctx, route)...)
			}
			// routes = append(routes[:i], routes[i+1:]...)
		}
	}
	return routes
}

func (m *Menu) generateRoute(ctx http.Context, route Route) []Route {
	url := route.Path
	url = strings.Split(url, "?")[0]
	if url == "" || len(route.Children) > 0 {
		return []Route{}
	}
	menu := func (action string, path string) Route  {
		return Route{
			Name: route.Name+"-"+action,
			Path: url+path,
			Component: func () string {
				if route.UrlType == admin.TYPE_IFRAME {
					return "iframe"
				}else {
					if route.Component == "" {
						return "amis"
					}else{
						return route.Component
					}
			}}(),
			Meta: Meta{
				Hide: true,
				Icon:	route.Meta.Icon,
				Title: route.Meta.Title + " - " + tools.AdminLang(ctx, action),
			},
		}
	}

	routes := []Route{
		menu("create", "/create"),
		menu("edit", "/edit"),
		menu("show", "/show"),
	}
	return routes
}


func (m *Menu) Extra(ctx http.Context) []Route {
	route := []Route{}
	config := facades.Config()
	if config.GetBool("admin.auth.enable") {
		route = append(route, Route{
			Name:      "user_setting",
			Path:      "/user_setting",
			Component: "user_setting",
			Meta: Meta{
				Hide:         true,
				Title:        tools.AdminLang(ctx, "user_setting"),
				Icon:         "material-symbols:manage-accounts",
				// SingleLayout: "basic",
			},
		})
	}
	if config.GetBool("admin.show_development_tools") {
		route = append(route, m.devToolMenus(ctx)...)
	}
	return route
}

func (m *Menu) devToolMenus(ctx http.Context) []Route {
	return []Route{
		{
			Name:      "dev_tools",
			Path:      "/dev_tools",
			Component: "amis",
			Meta: Meta{
				Title: tools.AdminLang(ctx, "developer"),
				Icon:  "fluent:window-dev-tools-20-regular",
			},
			Children: []Route{
				{
					Name:      "dev_tools_extensions",
					Path:      "/dev_tools/extensions",
					Component: "amis",
					Meta: Meta{
						Title: tools.AdminLang(ctx, "extensions.menu"),
						Icon:  "ion:extension-puzzle-outline",
					},
				},
				{
					Name:      "dev_tools_pages",
					Path:      "/dev_tools/pages",
					Component: "amis",
					Meta: Meta{
						Title: tools.AdminLang(ctx, "pages.menu"),
						Icon:  "iconoir:multiple-pages",
					},
				},
				{
					Name:      "dev_tools_relationships",
					Path:      "/dev_tools/relationships",
					Component: "amis",
					Meta: Meta{
						Title: tools.AdminLang(ctx, "relationships.menu"),
						Icon:  "ant-design:node-index-outlined",
					},
				},
				{
					Name:      "dev_tools_apis",
					Path:      "/dev_tools/apis",
					Component: "amis",
					Meta: Meta{
						Title: tools.AdminLang(ctx, "apis.menu"),
						Icon:  "ant-design:api-outlined",
					},
				},
				{
					Name:      "dev_tools_code_generator",
					Path:      "/dev_tools/code_generator",
					Component: "amis",
					Meta: Meta{
						Title: tools.AdminLang(ctx, "code_generator"),
						Icon:  "ic:baseline-code",
					},
				},
				{
					Name:      "dev_tools_editor",
					Path:      "/dev_tools/editor",
					Component: "editor",
					Meta: Meta{
						Title: tools.AdminLang(ctx,"visual_editor"),
						Icon:  "mdi:monitor-edit",
					},
				},
			},
		},
	}
}
