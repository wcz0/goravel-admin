package response

import	"github.com/goravel/framework/contracts/http"


type Error interface {
	Response(ctx http.Context) http.Response
}

type UnauthorizedError struct {
}

func NewUnauthorizedError() *UnauthorizedError {
	return &UnauthorizedError{}
}

func (a *UnauthorizedError) Response(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusUnauthorized, http.Json{
		"code": 401,
		"msg": "Unauthenticated.",
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

type FormError struct {
}

func NewFormError() *FormError {
	return &FormError{}
}

func (f *FormError) Response(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
		"code": 400,
		"msg": "Bad Request.",
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

type LimitError struct {
}

func NewLimitError() *LimitError {
	return &LimitError{}
}

func (l *LimitError) Response(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusTooManyRequests, http.Json{
		"code": 429,
		"msg": "Forbidden.",
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

type PermissionError struct {
}

func NewPermissionError() *PermissionError {
	return &PermissionError{}
}

func (p *PermissionError) Response(ctx http.Context) http.Response{
	return ctx.Response().Json(http.StatusForbidden, http.Json{
		"code": 403,
		"msg": "Permission denied.",
		"data": nil,
		"doNotDisplayToast": 0,
	})
}