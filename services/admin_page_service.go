package services

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"github.com/wcz0/goravel-admin/models/admin"
)

type AdminPageService struct {
	*Service
}

func NewAdminPageService() *AdminPageService {
	return &AdminPageService{Service: NewService()}
}

// pageCacheKey 是自定义页面 schema 的缓存键前缀。
const pageCacheKey = "admin_page_"

// Get 按 sign 取页面结构。
//
// 返回的是解析后的 JSON 对象而不是原始字符串：前端拿到 data 之后会直接
// 交给 amis 渲染，字符串会导致 schema 解析失败。
func (s *AdminPageService) Get(sign string) any {
	if sign == "" {
		return nil
	}

	var page admin.AdminPage
	// First 找不到记录时不会报错，返回的是零值结构体，得自己判断 ID
	if err := facades.Orm().Query().Where("sign", sign).First(&page); err != nil || page.ID == 0 {
		return nil
	}

	var schema any
	if err := json.Unmarshal([]byte(page.Schema), &schema); err != nil {
		return nil
	}

	return schema
}

// List 返回页面列表，支持按名称和标识筛选。
func (s *AdminPageService) List(ctx http.Context) http.Response {
	query := facades.Orm().Query().Model(&admin.AdminPage{}).Order("id desc")
	if title := ctx.Request().Input("title"); title != "" {
		query.Where("title", "like", "%"+title+"%")
	}
	if sign := ctx.Request().Input("sign"); sign != "" {
		query.Where("sign", "like", "%"+sign+"%")
	}

	var pages []admin.AdminPage
	if err := query.Get(&pages); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	items := make([]map[string]any, 0, len(pages))
	for _, p := range pages {
		items = append(items, map[string]any{
			"id":         p.ID,
			"title":      p.Title,
			"sign":       p.Sign,
			"created_at": p.CreatedAt,
			"updated_at": p.UpdatedAt,
		})
	}

	return s.SuccessData(ctx, map[string]any{"items": items})
}

// Show 返回单条页面记录。
//
// 与 List 不同，编辑时需要把 schema 包进 page.schema，配合前端的
// SubFormControl("page") 表单结构。
func (s *AdminPageService) Show(ctx http.Context) http.Response {
	var page admin.AdminPage
	if err := firstByID(ctx.Request().InputInt("id"), &page); err != nil {
		return s.FailMsg(ctx, "页面不存在")
	}

	var schema any
	if page.Schema != "" {
		_ = json.Unmarshal([]byte(page.Schema), &schema)
	}

	return s.SuccessData(ctx, map[string]any{
		"id":         page.ID,
		"title":      page.Title,
		"sign":       page.Sign,
		"created_at": page.CreatedAt,
		"updated_at": page.UpdatedAt,
		"page":       map[string]any{"schema": schema},
	})
}

// extractSchema 从请求里取出页面结构并序列化成字符串。
//
// 前端提交的是 page.schema 嵌套对象；同时兼容直接提交 schema 字段的情况，
// 该字段既能是 amis 解析后的对象，也能是 JSON 文本（curl 直接调用时常见）。
func extractSchema(all map[string]any) (string, error) {
	if page, ok := all["page"].(map[string]any); ok {
		if schema, ok := page["schema"]; ok && schema != nil {
			raw, err := json.Marshal(schema)
			if err != nil {
				return "", err
			}

			return string(raw), nil
		}
	}

	schema, ok := all["schema"]
	if !ok || schema == nil {
		return "", errors.New("页面结构不可为空")
	}

	// JSON 文本原样落库，避免二次序列化后出现转义字符串
	if raw, isText := schema.(string); isText {
		if raw == "" {
			return "", errors.New("页面结构不可为空")
		}
		if !json.Valid([]byte(raw)) {
			return "", errors.New("页面结构不是合法的 JSON")
		}

		return raw, nil
	}

	raw, err := json.Marshal(schema)
	if err != nil {
		return "", err
	}
	if len(raw) <= 2 {
		return "", errors.New("页面结构不可为空")
	}

	return string(raw), nil
}

// checkSign 校验 sign 唯一性（排除自身）。
func (s *AdminPageService) checkSign(sign string, excludeID int) error {
	query := facades.Orm().Query().Model(&admin.AdminPage{}).Where("sign", sign)
	if excludeID > 0 {
		query.Where("id != ?", excludeID)
	}

	total, err := query.Count()
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("页面标识已存在")
	}

	return nil
}

// Store 新建页面。
func (s *AdminPageService) Store(ctx http.Context) http.Response {
	all := ctx.Request().All()
	sign := ctx.Request().Input("sign")
	if sign == "" {
		return s.FailMsg(ctx, "页面标识不可为空")
	}
	if err := s.checkSign(sign, 0); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	schema, err := extractSchema(all)
	if err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	page := admin.AdminPage{
		Title:  ctx.Request().Input("title"),
		Sign:   sign,
		Schema: schema,
	}
	if err := facades.Orm().Query().Create(&page); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "创建成功")
}

// Update 更新页面。
func (s *AdminPageService) Update(ctx http.Context) http.Response {
	id := ctx.Request().InputInt("id")
	var page admin.AdminPage
	if err := firstByID(id, &page); err != nil {
		return s.FailMsg(ctx, "页面不存在")
	}

	sign := ctx.Request().Input("sign")
	if sign == "" {
		return s.FailMsg(ctx, "页面标识不可为空")
	}
	if err := s.checkSign(sign, id); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	schema, err := extractSchema(ctx.Request().All())
	if err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	oldSign := page.Sign
	page.Title = ctx.Request().Input("title")
	page.Sign = sign
	page.Schema = schema
	if err := facades.Orm().Query().Save(&page); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	// sign 变了要同时清掉新旧两个缓存键
	s.clearCache(oldSign)
	s.clearCache(sign)

	return s.SuccessMsg(ctx, "更新成功")
}

// Destroy 删除页面（支持逗号分隔的批量 id）。
func (s *AdminPageService) Destroy(ctx http.Context) http.Response {
	ids := parseIDList(ctx.Request().Input("id"))
	if len(ids) == 0 {
		return s.FailMsg(ctx, "ID不能为空")
	}

	var pages []admin.AdminPage
	if err := facades.Orm().Query().WhereIn("id", ids).Find(&pages); err == nil {
		for _, p := range pages {
			s.clearCache(p.Sign)
		}
	}

	if _, err := facades.Orm().Query().WhereIn("id", ids).Delete(&admin.AdminPage{}); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	return s.SuccessMsg(ctx, "删除成功")
}

// Options 供菜单表单里的「页面」下拉使用。
func (s *AdminPageService) Options(ctx http.Context) http.Response {
	var pages []admin.AdminPage
	if err := facades.Orm().Query().Order("id desc").Find(&pages); err != nil {
		return s.FailMsg(ctx, err.Error())
	}

	options := make([]map[string]any, 0, len(pages))
	for _, p := range pages {
		options = append(options, map[string]any{"label": p.Title, "value": p.Sign})
	}

	return s.SuccessData(ctx, options)
}

func (s *AdminPageService) clearCache(sign string) {
	if sign == "" {
		return
	}
	func() {
		defer func() { _ = recover() }()
		facades.Cache().Forget(pageCacheKey + sign)
	}()
}

// parseIDList 把 "1,2,3" 或 "1" 解析成 []any，忽略空值与非法值。
func parseIDList(raw string) []any {
	parts := strings.Split(raw, ",")
	ids := make([]any, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.Atoi(part)
		if err != nil || id <= 0 {
			continue
		}
		ids = append(ids, id)
	}

	return ids
}
