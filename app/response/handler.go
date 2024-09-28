package response

import (
	"goravel/app/enums"

	"github.com/goravel/framework/contracts/http"
)

type ErrorHandler struct {
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (e *ErrorHandler) Json(ctx http.Context, json any) http.Response {
	return ctx.Response().Success().Json(json)
}

func (e *ErrorHandler) Success(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status":            enums.StatusSuccess,
		"code":              enums.Success,
		"msg":               enums.CodeEnum(enums.Success).Message(),
		"data":              []any{},
		"doNotDisplayToast": 0,
	})
}

func (e *ErrorHandler) SuccessMsg(ctx http.Context, msg any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status":            enums.StatusSuccess,
		"code":              enums.Success,
		"msg":               msg,
		"data":              []any{},
		"doNotDisplayToast": 0,
	})
}

func (e *ErrorHandler) SuccessData(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status":            enums.StatusSuccess,
		"code":              enums.Success,
		"msg":               enums.CodeEnum(enums.Success).Message(),
		"data":              data,
		"doNotDisplayToast": 0,
	})
}

func (e *ErrorHandler) SuccessMsgData(ctx http.Context, msg any, data any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status":            enums.StatusSuccess,
		"code":              enums.Success,
		"msg":               msg,
		"data":              data,
		"doNotDisplayToast": 0,
	})
}

func (e *ErrorHandler) SuccessCodeMsg(ctx http.Context, code int, msg any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status":            enums.StatusSuccess,
		"code":              code,
		"msg":               msg,
		"data":              []any{},
		"doNotDisplayToast": 0,
	})
}

func (e *ErrorHandler) Fail(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status":            enums.StatusFailed,
		"code":              enums.Failed,
		"msg":               enums.CodeEnum(enums.Failed).Message(),
		"data":              []any{},
		"doNotDisplayToast": 0,
	})
}

func (e *ErrorHandler) FailMsg(ctx http.Context, msg any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status":            enums.StatusFailed,
		"code":              enums.Failed,
		"msg":               msg,
		"data":              []any{},
		"doNotDisplayToast": 0,
	})
}

func (e *ErrorHandler) FailData(ctx http.Context, data any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status":            enums.StatusFailed,
		"code":              enums.Failed,
		"msg":               "Error.",
		"data":              data,
		"doNotDisplayToast": 0,
	})
}

func (e *ErrorHandler) FailMsgData(ctx http.Context, msg any, data any) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"status":            enums.StatusFailed,
		"code":              enums.Failed,
		"msg":               msg,
		"data":              data,
		"doNotDisplayToast": 0,
	})
}

func (e *ErrorHandler) FailCodeMsg(ctx http.Context, code int, msg string) http.Response {
	return ctx.Response().Json(code, http.Json{
		"status":            enums.StatusFailed,
		"code":              code,
		"msg":               msg,
		"data":              []any{},
		"doNotDisplayToast": 0,
	})
}

func (e *ErrorHandler) FormError(ctx http.Context, msg any) http.Response {
	if msg != "" {
		msg = "Form validation error."
	}
	return ctx.Response().Success().Json(FormError)
}

func (e *ErrorHandler) LimitError(ctx http.Context, msg any) http.Response {
	if msg != "" {
		msg = "Too many requests."
	}
	return ctx.Response().Success().Json(LimitError)
}

func (e *ErrorHandler) UnauthorizedError(ctx http.Context, msg any) http.Response {
	if msg != "" {
		msg = "Unauthorized."
	}
	return ctx.Response().Success().Json(Unauthorized)
}

func (e *ErrorHandler) PermissionError(ctx http.Context, msg any) http.Response {
	if msg != "" {
		msg = "Permission denied."
	}
	return ctx.Response().Success().Json(PermissionError)
}
