package services

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/goravel-admin/models/admin"
)

// relationshipTypes 是支持的关联类型，对齐 owl-admin 的 9 种。
func relationshipTypes() []map[string]string {
	return []map[string]string{
		{"label": "hasOne", "value": "hasOne"},
		{"label": "hasMany", "value": "hasMany"},
		{"label": "belongsTo", "value": "belongsTo"},
		{"label": "belongsToMany", "value": "belongsToMany"},
		{"label": "hasOneThrough", "value": "hasOneThrough"},
		{"label": "hasManyThrough", "value": "hasManyThrough"},
		{"label": "morphOne", "value": "morphOne"},
		{"label": "morphMany", "value": "morphMany"},
		{"label": "morphToMany", "value": "morphToMany"},
	}
}

type AdminRelationshipService struct {
	*Service
}

func NewAdminRelationshipService() *AdminRelationshipService {
	return &AdminRelationshipService{Service: NewService()}
}

// Types 返回关联类型下拉选项。
func (s *AdminRelationshipService) Types() []map[string]string {
	return relationshipTypes()
}

// previewCode 生成 Goravel ORM 的关联方法预览。
//
// 这是一个「可复制到模型里」的代码片段：Goravel 的关联必须写在 struct 上，
// 没法像 Eloquent 那样运行时挂载，所以只能给开发者提供参考实现。
func previewCode(item admin.AdminRelationship) string {
	var params map[string]any
	if item.Args != nil {
		_ = json.Unmarshal([]byte(*item.Args), &params)
	}

	related := "RelatedModel"
	if v, ok := params["related_model"].(string); ok && v != "" {
		related = v
	}

	return "func (m *" + item.Model + ") " + item.Title + "() orm.Relation {\n" +
		"    return m." + item.Type + "(" + related + ")\n" +
		"}"
}

// castToString 安全取指针字符串。
func castToString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

// List 返回关联配置列表。
func (s *AdminRelationshipService) List(ctx http.Context) http.Response {
	query := facades.Orm().Query().Model(&admin.AdminRelationship{}).Order("id desc")
	if model := ctx.Request().Input("model"); model != "" {
		query.Where("model", "like", "%"+model+"%")
	}
	if title := ctx.Request().Input("title"); title != "" {
		query.Where("title", "like", "%"+title+"%")
	}

	var list []admin.AdminRelationship
	if err := query.Get(&list); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	items := make([]map[string]any, 0, len(list))
	for _, item := range list {
		items = append(items, map[string]any{
			"id":           item.ID,
			"model":        item.Model,
			"title":        item.Title,
			"type":         item.Type,
			"remark":       item.Remark,
			"args":         item.Args,
			"extra":        item.Extra,
			"preview_code": previewCode(item),
			"created_at":   item.CreatedAt,
			"updated_at":   item.UpdatedAt,
		})
	}

	return s.SuccessData(ctx, map[string]any{"items": items})
}

// Show 返回单条关联配置。
func (s *AdminRelationshipService) Show(ctx http.Context) http.Response {
	var item admin.AdminRelationship
	if err := facades.Orm().Query().Find(&item, ctx.Request().InputInt("id")); err != nil {
		return s.FailMsg(ctx, "关联配置不存在")
	}

	return s.SuccessData(ctx, item)
}

// checkUnique 校验同一模型下关联名唯一。
func (s *AdminRelationshipService) checkUnique(model, title string, excludeID int) error {
	query := facades.Orm().Query().Model(&admin.AdminRelationship{}).
		Where("model", model).Where("title", title)
	if excludeID > 0 {
		query.Where("id != ?", excludeID)
	}

	total, err := query.Count()
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("该模型下存在同名关联")
	}

	return nil
}

// Store 新建关联配置。
func (s *AdminRelationshipService) Store(ctx http.Context) http.Response {
	model := ctx.Request().Input("model")
	title := ctx.Request().Input("title")
	if model == "" || title == "" {
		return s.FailMsg(ctx, "模型与关联名不可为空")
	}
	if err := s.checkUnique(model, title, 0); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	item := admin.AdminRelationship{
		Model:  model,
		Title:  title,
		Type:   ctx.Request().Input("type"),
		Args:   marshalArgs(ctx.Request().All()["args"]),
		Extra:  marshalArgs(ctx.Request().All()["extra"]),
		Remark: optionalString(ctx.Request().Input("remark")),
	}
	if err := facades.Orm().Query().Create(&item); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "创建成功")
}

// Update 更新关联配置。
func (s *AdminRelationshipService) Update(ctx http.Context) http.Response {
	id := ctx.Request().InputInt("id")
	var item admin.AdminRelationship
	if err := facades.Orm().Query().Find(&item, id); err != nil {
		return s.FailMsg(ctx, "关联配置不存在")
	}

	model := ctx.Request().Input("model")
	title := ctx.Request().Input("title")
	if model == "" || title == "" {
		return s.FailMsg(ctx, "模型与关联名不可为空")
	}
	if err := s.checkUnique(model, title, id); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	item.Model = model
	item.Title = title
	item.Type = ctx.Request().Input("type")
	item.Remark = optionalString(ctx.Request().Input("remark"))
	if args := ctx.Request().All()["args"]; args != nil {
		item.Args = marshalArgs(args)
	}
	if extra := ctx.Request().All()["extra"]; extra != nil {
		item.Extra = marshalArgs(extra)
	}
	if err := facades.Orm().Query().Save(&item); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "更新成功")
}

// Destroy 删除关联配置。
func (s *AdminRelationshipService) Destroy(ctx http.Context) http.Response {
	if _, err := facades.Orm().Query().
		WhereIn("id", parseIDList(ctx.Request().Input("id"))).
		Delete(&admin.AdminRelationship{}); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "删除成功")
}

// optionalString 把空字符串转成 nil，避免写入无意义的空值。
func optionalString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	return &value
}
