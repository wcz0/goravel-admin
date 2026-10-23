package services

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/goravel-admin/models/admin"
)

// apiTemplate 描述一个内置的动态接口模板。
//
// owl-admin 的模板是 PHP 类，可以在运行时实例化；Go 不能动态加载代码，
// 所以模板必须是编译期注册的 handler，这里只保留模板元信息与参数 schema。
type apiTemplate struct {
	Title  string         `json:"title"`
	Name   string         `json:"name"`
	Method string         `json:"method"`
	Args   map[string]any `json:"args"`
}

// apiTemplates 是内置的动态接口模板表，对齐 owl-admin 的 Support/Apis 目录。
func apiTemplates() []apiTemplate {
	return []apiTemplate{
		{Title: "列表", Name: "list", Method: "get", Args: map[string]any{"model": "", "fields": ""}},
		{Title: "详情", Name: "detail", Method: "get", Args: map[string]any{"model": ""}},
		{Title: "新增", Name: "create", Method: "post", Args: map[string]any{"model": ""}},
		{Title: "更新", Name: "update", Method: "put", Args: map[string]any{"model": ""}},
		{Title: "删除", Name: "delete", Method: "delete", Args: map[string]any{"model": ""}},
		{Title: "状态切换", Name: "toggle", Method: "put", Args: map[string]any{"model": "", "field": "enabled"}},
		{Title: "下拉选项", Name: "options", Method: "get", Args: map[string]any{"model": "", "label": "name", "value": "id"}},
		{Title: "统计数据", Name: "aggregate", Method: "get", Args: map[string]any{"model": "", "field": "id", "type": "count"}},
		{Title: "读取配置", Name: "settings", Method: "get", Args: map[string]any{"keys": ""}},
		{Title: "保存配置", Name: "save_settings", Method: "post", Args: map[string]any{"keys": ""}},
	}
}

type AdminApiService struct {
	*Service
}

func NewAdminApiService() *AdminApiService {
	return &AdminApiService{Service: NewService()}
}

// Templates 返回内置模板列表，标题带方法前缀便于前端展示。
func (s *AdminApiService) Templates() []map[string]any {
	list := apiTemplates()
	options := make([]map[string]any, 0, len(list))
	for _, t := range list {
		options = append(options, map[string]any{
			"label":  t.Title,
			"value":  t.Name,
			"method": t.Method,
			"args":   t.Args,
		})
	}

	return options
}

// templateTitle 查模板标题，模板不存在时回落为原始 name。
func templateTitle(name string) string {
	for _, t := range apiTemplates() {
		if t.Name == name {
			return t.Title
		}
	}

	return name
}

// List 返回动态 API 列表。
func (s *AdminApiService) List(ctx http.Context) http.Response {
	query := facades.Orm().Query().Model(&admin.AdminApi{}).Order("id desc")
	if title := ctx.Request().Input("title"); title != "" {
		query = query.Where("title", "like", "%"+title+"%")
	}
	if path := ctx.Request().Input("path"); path != "" {
		query = query.Where("path", "like", "%"+path+"%")
	}

	var apis []admin.AdminApi
	if err := query.Get(&apis); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	items := make([]map[string]any, 0, len(apis))
	for _, api := range apis {
		items = append(items, map[string]any{
			"id":             api.ID,
			"title":          api.Title,
			"path":           api.Path,
			"template":       api.Template,
			"template_title": templateTitle(api.Template),
			"enabled":        api.Enabled,
			"args":           api.Args,
			"created_at":     api.CreatedAt,
			"updated_at":     api.UpdatedAt,
		})
	}

	return s.SuccessData(ctx, map[string]any{"items": items})
}

// Show 返回单条动态 API。
func (s *AdminApiService) Show(ctx http.Context) http.Response {
	var api admin.AdminApi
	if err := firstByID(ctx.Request().InputInt("id"), &api); err != nil {
		return s.FailMsg(ctx, "接口不存在")
	}

	return s.SuccessData(ctx, api)
}

// normalizePath 保证接口路径以 / 开头且不含空白。
func normalizePath(path string) (string, error) {
	if path == "" {
		return "", errors.New("接口路径不可为空")
	}
	if strings.ContainsAny(path, " \t\r\n") {
		return "", errors.New("接口路径不能包含空白字符")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return path, nil
}

func (s *AdminApiService) checkPath(path string, excludeID int) error {
	query := facades.Orm().Query().Model(&admin.AdminApi{}).Where("path", path)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	total, err := query.Count()
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("接口路径已存在")
	}

	return nil
}

// marshalArgs 把参数对象序列化成 JSON 文本。
func marshalArgs(value any) *string {
	if value == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	str := string(raw)

	return &str
}

// Store 新建动态 API。
func (s *AdminApiService) Store(ctx http.Context) http.Response {
	path, err := normalizePath(ctx.Request().Input("path"))
	if err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	if err := s.checkPath(path, 0); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	api := admin.AdminApi{
		Title:    ctx.Request().Input("title"),
		Path:     path,
		Template: ctx.Request().Input("template"),
		Enabled:  uint8(ctx.Request().InputInt("enabled", 1)),
		Args:     marshalArgs(ctx.Request().All()["args"]),
	}
	if err := facades.Orm().Query().Create(&api); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "创建成功")
}

// Update 更新动态 API。
func (s *AdminApiService) Update(ctx http.Context) http.Response {
	id := ctx.Request().InputInt("id")
	var api admin.AdminApi
	if err := firstByID(id, &api); err != nil {
		return s.FailMsg(ctx, "接口不存在")
	}

	path, err := normalizePath(ctx.Request().Input("path"))
	if err != nil {
		return s.FailMsg(ctx, err.Error())
	}
	if err := s.checkPath(path, id); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	api.Title = ctx.Request().Input("title")
	api.Path = path
	api.Template = ctx.Request().Input("template")
	api.Enabled = uint8(ctx.Request().InputInt("enabled", int(api.Enabled)))
	if args := ctx.Request().All()["args"]; args != nil {
		api.Args = marshalArgs(args)
	}
	if err := facades.Orm().Query().Save(&api); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "更新成功")
}

// Destroy 删除动态 API。
func (s *AdminApiService) Destroy(ctx http.Context) http.Response {
	if _, err := facades.Orm().Query().
		WhereIn("id", parseIDList(ctx.Request().Input("id"))).
		Delete(&admin.AdminApi{}); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "删除成功")
}

// Enabled 返回启用的动态 API，供启动阶段注册路由。
func (s *AdminApiService) Enabled() []admin.AdminApi {
	var apis []admin.AdminApi
	if err := facades.Orm().Query().
		Where("enabled", 1).
		Order("id asc").
		Find(&apis); err != nil {
		return nil
	}

	return apis
}
