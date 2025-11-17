package admin

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type MenuController struct {
}

func (m *MenuController) GetCurrentMenus() {

}

func NewMenuController() *MenuController {
	return &MenuController{}
}

func (m *MenuController) Index(ctx http.Context) http.Response {
	return nil
}

func (m *MenuController) Show(ctx http.Context) http.Response {
	// 验证参数
	validator, err := facades.Validation().Make(ctx.Request().All(), map[string]string{
		"id": "required|number",
	})
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "验证器创建失败",
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"message": "参数验证失败",
			"errors":  validator.Errors().All(),
		})
	}
	
	return nil
}

func (m *MenuController) Store(ctx http.Context) http.Response {
	// 验证参数
	validator, err := facades.Validation().Make(ctx.Request().All(), map[string]string{
		"title":      "required|string|max:50",
		"icon":       "string|max:50",
		"uri":        "string|max:255",
		"parent_id":  "number",
		"order":      "number",
		"permission": "string|max:255",
	})
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "验证器创建失败",
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"message": "参数验证失败",
			"errors":  validator.Errors().All(),
		})
	}

	return nil
}

func (m *MenuController) Update(ctx http.Context) http.Response {
	// 验证参数
	validator, err := facades.Validation().Make(ctx.Request().All(), map[string]string{
		"id":         "required|number",
		"title":      "required|string|max:50",
		"icon":       "string|max:50",
		"uri":        "string|max:255",
		"parent_id":  "number",
		"order":      "number",
		"permission": "string|max:255",
	})
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "验证器创建失败",
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"message": "参数验证失败",
			"errors":  validator.Errors().All(),
		})
	}

	return nil
}

func (m *MenuController) Destroy(ctx http.Context) http.Response {
	// 验证参数
	validator, err := facades.Validation().Make(ctx.Request().All(), map[string]string{
		"id": "required|number",
	})
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "验证器创建失败",
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"message": "参数验证失败",
			"errors":  validator.Errors().All(),
		})
	}

	return nil
}