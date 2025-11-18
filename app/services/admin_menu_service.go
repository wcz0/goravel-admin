package services

import (
    "goravel/app/models/admin"

    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
)

type AdminMenuService struct {
    *Service
}

func NewAdminMenuService() *AdminMenuService {
    return &AdminMenuService{Service: NewService()}
}

func (s *AdminMenuService) List(ctx http.Context) http.Response {
    title := ctx.Request().Input("title")
    url := ctx.Request().Input("url")
    q := facades.Orm().Query().Model(&admin.AdminMenu{}).Order("custom_order asc").Order("id asc")
    if title != "" { q.Where("title", "like", "%"+title+"%") }
    if url != "" { q.Where("url", "like", "%"+url+"%") }
    var menus []admin.AdminMenu
    if err := q.Get(&menus); err != nil { return s.FailMsg(ctx, err.Error()) }
    parentMap := make(map[uint32][]admin.AdminMenu)
    for _, m := range menus { parentMap[m.ParentId] = append(parentMap[m.ParentId], m) }
    var build func(pid uint32) []map[string]any
    build = func(pid uint32) []map[string]any {
        children := parentMap[pid]
        list := make([]map[string]any, 0, len(children))
        for _, m := range children {
            node := map[string]any{
                "id": m.ID,
                "parent_id": m.ParentId,
                "title": m.Title,
                "icon": m.Icon,
                "url": m.Url,
                "visible": m.Visible,
                "is_home": m.IsHome,
                "component": m.Component,
                "custom_order": m.CustomOrder,
                "keep_alive": m.KeepAlive,
                "iframe_url": m.IFrameUrl,
                "is_full": m.IsFull,
                "created_at": m.CreatedAt,
                "updated_at": m.UpdatedAt,
            }
            if cs := build(uint32(m.ID)); len(cs) > 0 { node["children"] = cs }
            list = append(list, node)
        }
        return list
    }
    return s.SuccessData(ctx, map[string]any{"items": build(0)})
}

func (s *AdminMenuService) Show(ctx http.Context) http.Response {
    id := ctx.Request().InputInt("id")
    var m admin.AdminMenu
    if err := facades.Orm().Query().Find(&m, id); err != nil { return s.FailMsg(ctx, "菜单不存在") }
    return s.SuccessData(ctx, m)
}

func (s *AdminMenuService) Store(ctx http.Context) http.Response {
    m := admin.AdminMenu{
        ParentId:   uint32(ctx.Request().InputInt("parent_id", 0)),
        Title:      ctx.Request().Input("title"),
        Icon:       ctx.Request().Input("icon"),
        Url:        ctx.Request().Input("url"),
        Visible:    uint8(ctx.Request().InputInt("visible", 1)),
        IsHome:     uint8(ctx.Request().InputInt("is_home", 0)),
        Component:  ctx.Request().Input("component"),
        CustomOrder: ctx.Request().InputInt("custom_order", 0),
        IsFull:     uint8(ctx.Request().InputInt("is_full", 0)),
    }
    if keep := ctx.Request().InputInt("keep_alive", 0); keep >= 0 { v := uint8(keep); m.KeepAlive = &v }
    if iframe := ctx.Request().Input("iframe_url"); iframe != "" { m.IFrameUrl = &iframe }
    if err := facades.Orm().Query().Create(&m); err != nil { return s.FailMsg(ctx, err.Error()) }
    return s.SuccessMsg(ctx, "创建成功")
}

func (s *AdminMenuService) Update(ctx http.Context) http.Response {
    id := ctx.Request().InputInt("id")
    var m admin.AdminMenu
    if err := facades.Orm().Query().Find(&m, id); err != nil { return s.FailMsg(ctx, "菜单不存在") }
    m.ParentId = uint32(ctx.Request().InputInt("parent_id", int(m.ParentId)))
    if v := ctx.Request().Input("title"); v != "" { m.Title = v }
    if v := ctx.Request().Input("icon"); v != "" { m.Icon = v }
    if v := ctx.Request().Input("url"); v != "" { m.Url = v }
    if v := ctx.Request().InputInt("visible", int(m.Visible)); v >= 0 { m.Visible = uint8(v) }
    if v := ctx.Request().InputInt("is_home", int(m.IsHome)); v >= 0 { m.IsHome = uint8(v) }
    if v := ctx.Request().Input("component"); v != "" { m.Component = v }
    m.CustomOrder = ctx.Request().InputInt("custom_order", m.CustomOrder)
    if v := ctx.Request().InputInt("is_full", int(m.IsFull)); v >= 0 { m.IsFull = uint8(v) }
    if v := ctx.Request().InputInt("keep_alive", -1); v >= 0 { vv := uint8(v); m.KeepAlive = &vv }
    if v := ctx.Request().Input("iframe_url"); v != "" { m.IFrameUrl = &v }
    if err := facades.Orm().Query().Save(&m); err != nil { return s.FailMsg(ctx, err.Error()) }
    return s.SuccessMsg(ctx, "更新成功")
}

func (s *AdminMenuService) Destroy(ctx http.Context) http.Response {
    id := ctx.Request().InputInt("id")
    var m admin.AdminMenu
    if err := facades.Orm().Query().Find(&m, id); err != nil { return s.FailMsg(ctx, "菜单不存在") }
    if _, err := facades.Orm().Query().Delete(&m); err != nil { return s.FailMsg(ctx, err.Error()) }
    return s.SuccessMsg(ctx, "删除成功")
}

func (s *AdminMenuService) Reorder(ctx http.Context) http.Response {
    arr := ctx.Request().InputArray("ids")
    if len(arr) == 0 {
        idsStr := ctx.Request().Input("ids")
        if idsStr == "" { return s.SuccessMsg(ctx, "无排序更新") }
    }
    tx, err := facades.Orm().Query().Begin()
    if err != nil { return s.FailMsg(ctx, "开始事务失败: "+err.Error()) }
    for i, raw := range arr {
        var m admin.AdminMenu
        if err := tx.Find(&m, raw); err != nil { tx.Rollback(); return s.FailMsg(ctx, "菜单不存在") }
        m.CustomOrder = i
        if err := tx.Save(&m); err != nil { tx.Rollback(); return s.FailMsg(ctx, err.Error()) }
    }
    if err := tx.Commit(); err != nil { return s.FailMsg(ctx, "提交事务失败: "+err.Error()) }
    return s.SuccessMsg(ctx, "排序已更新")
}

func (s *AdminMenuService) QuickEdit(ctx http.Context) http.Response {
    return s.SuccessMsg(ctx, "快速编辑成功")
}

func (s *AdminMenuService) QuickEditItem(ctx http.Context) http.Response {
    id := ctx.Request().InputInt("id")
    name := ctx.Request().Input("name")
    value := ctx.Request().Input("value")
    var m admin.AdminMenu
    if err := facades.Orm().Query().Find(&m, id); err != nil { return s.FailMsg(ctx, "菜单不存在") }
    switch name {
    case "custom_order":
        m.CustomOrder = ctx.Request().InputInt("value", m.CustomOrder)
    case "visible":
        if value == "true" || value == "1" { m.Visible = 1 } else { m.Visible = 0 }
    case "is_home":
        if value == "true" || value == "1" { m.IsHome = 1 } else { m.IsHome = 0 }
    }
    if err := facades.Orm().Query().Save(&m); err != nil { return s.FailMsg(ctx, err.Error()) }
    return s.SuccessMsg(ctx, "已更新")
}
