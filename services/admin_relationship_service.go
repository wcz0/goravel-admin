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

// previewCode 生成关联查询的预览代码。
//
// Goravel 的 ORM 没有 Eloquent 那种 resolveRelationUsing 动态关联，关联只能
// 手写成查询方法，所以这里给的是「可复制进模型、再补两处 TODO 就能跑」的
// 显式查询骨架，而不是伪造一个 orm.Relation 方法。标识符统一转成合法的 Go
// 命名，避免标题里带空格或中文时拿到无法编译的片段。
func previewCode(item admin.AdminRelationship) string {
	var params map[string]any
	if item.Args != nil {
		_ = json.Unmarshal([]byte(*item.Args), &params)
	}

	model := goIdent(lastSegment(item.Model), true)
	method := goIdent(item.Title, true)
	related := goIdent(firstArg(params, "related", "related_model", "relatedModel", "model"), true)
	if related == "" {
		related = "Related"
	}

	var b strings.Builder
	b.WriteString("// " + method + " 是 " + model + " 的 " + item.Type + " 关联，关联模型 " + related + "。\n")
	b.WriteString("// Goravel ORM 未提供 Eloquent 式关联方法，这里用显式查询实现。\n")
	b.WriteString("func (m *" + model + ") " + method + "() (items []" + related + ", err error) {\n")

	pivot := firstArg(params, "table", "pivot_table", "pivotTable")
	foreignKey := firstArg(params, "foreignKey", "foreign_key", "foreignPivotKey", "foreign_pivot_key")
	relatedKey := firstArg(params, "relatedKey", "related_key", "relatedPivotKey", "related_pivot_key", "ownerKey", "owner_key")

	switch {
	case pivot != "":
		// belongsToMany / morphToMany：先走中间表再用关联表的键打平
		b.WriteString("    err = facades.Orm().Query().\n")
		b.WriteString("        Join(\"inner join " + pivot + " on " + pivot + "." + orDefault(relatedKey, "related_id") + " = <关联表>.id\"). // TODO 补关联表名\n")
		b.WriteString("        Where(\"" + pivot + "." + orDefault(foreignKey, "model_id") + "\", m.ID).\n")
		b.WriteString("        Get(&items)\n")
	case foreignKey != "":
		b.WriteString("    err = facades.Orm().Query().\n")
		b.WriteString("        Where(\"" + foreignKey + "\", m.ID).\n")
		b.WriteString("        Get(&items)\n")
	default:
		b.WriteString("    err = facades.Orm().Query().Get(&items) // TODO 补上关联条件\n")
	}

	b.WriteString("    return\n")

	return b.String() + "}"
}

// goIdent 把任意文本转成合法的 Go 标识符。
//
// exported 为真时按 PascalCase 拼接（admin_users -> AdminUsers，
// Smoke Roles -> SmokeRoles），这样表名和带空格的标题都能得到像样的名字。
func goIdent(raw string, exported bool) string {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9')
	})

	if !exported {
		ident := strings.ToLower(strings.Join(parts, "_"))

		return prefixDigit(ident)
	}

	var b strings.Builder
	for _, part := range parts {
		b.WriteString(strings.ToUpper(part[:1]) + part[1:])
	}

	return prefixDigit(b.String())
}

// prefixDigit 保证标识符不以数字开头。
func prefixDigit(ident string) string {
	if ident == "" {
		return ""
	}
	if ident[0] >= '0' && ident[0] <= '9' {
		return "R" + ident
	}

	return ident
}

// lastSegment 取 FQCN 风格字符串的最后一段（App\Models\User -> User）。
func lastSegment(raw string) string {
	raw = strings.ReplaceAll(raw, "\\", "/")
	if idx := strings.LastIndex(raw, "/"); idx >= 0 {
		raw = raw[idx+1:]
	}

	return strings.TrimSpace(raw)
}

// firstArg 按顺序取第一个非空的字符串参数。
func firstArg(params map[string]any, keys ...string) string {
	for _, key := range keys {
		if v, ok := params[key].(string); ok && v != "" {
			return v
		}
	}

	return ""
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
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
		query = query.Where("model", "like", "%"+model+"%")
	}
	if title := ctx.Request().Input("title"); title != "" {
		query = query.Where("title", "like", "%"+title+"%")
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
	if err := firstByID(ctx.Request().InputInt("id"), &item); err != nil {
		return s.FailMsg(ctx, "关联配置不存在")
	}

	return s.SuccessData(ctx, item)
}

// checkUnique 校验同一模型下关联名唯一。
func (s *AdminRelationshipService) checkUnique(model, title string, excludeID int) error {
	query := facades.Orm().Query().Model(&admin.AdminRelationship{}).
		Where("model", model).Where("title", title)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
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
	if err := firstByID(id, &item); err != nil {
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
