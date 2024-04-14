package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Success(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"msg": "Success.",
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) MsgSuccess(ctx http.Context, msg string) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"msg": msg,
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) DataSuccess(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"msg": "Success.",
		"data": data,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) MsgDataSuccess(ctx http.Context, msg string, data any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"msg": msg,
		"data": data,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) Error(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusInternalServerError, http.Json{
		"msg": "Error.",
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) MsgError(ctx http.Context, msg string) http.Response {
	return ctx.Response().Json(http.StatusInternalServerError, http.Json{
		"msg": msg,
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) DataError(ctx http.Context, data any) http.Response {
	return ctx.Response().Json(http.StatusInternalServerError, http.Json{
		"msg": "Error.",
		"data": data,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) MsgDataError(ctx http.Context, msg string, data any) http.Response {
	return ctx.Response().Json(http.StatusInternalServerError, http.Json{
		"msg": msg,
		"data": data,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) CodeMsgError(ctx http.Context, code int, msg string) http.Response {
	return ctx.Response().Json(code, http.Json{
		"msg": msg,
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) FormError(ctx http.Context, msg string) http.Response {
	if msg != "" {
		msg = "Form validation error."
	}
	return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
		"msg": msg,
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) LimitError(ctx http.Context, msg string) http.Response {
	if msg != "" {
		msg = "Too many requests."
	}
	return ctx.Response().Json(http.StatusTooManyRequests, http.Json{
		"msg": msg,
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) UnauthorizedError(ctx http.Context, msg string) http.Response {
	if msg != "" {
		msg = "Unauthorized."
	}
	return ctx.Response().Json(http.StatusUnauthorized, http.Json{
		"msg": msg,
		"data": nil,
		"doNotDisplayToast": 0,
	})
}

func (c *Controller) AuthError(ctx http.Context, msg string) http.Response {
	if msg != "" {
		msg = "Unauthenticated."
	}
	return ctx.Response().Json(http.StatusUnauthorized, http.Json{
		"msg": msg,
		"data": nil,
		"doNotDisplayToast": 0,
	})
}