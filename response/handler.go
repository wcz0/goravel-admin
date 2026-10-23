package response

import (
	"github.com/wcz0/goravel-admin/enums"

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
		"msg":               "",
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
		"msg":               "",
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

// failWith 复制一份全局模板并可选覆盖 msg。
//
// 必须复制：全局模板是包级变量，直接改会污染后续请求。
func (e *ErrorHandler) failWith(ctx http.Context, template map[string]any, msg any) http.Response {
	payload := make(map[string]any, len(template)+1)
	for key, value := range template {
		payload[key] = value
	}

	// 传了自定义消息就用传入的，否则保留模板里的默认文案
	if custom, ok := msg.(string); ok && custom != "" {
		payload["msg"] = custom
	}

	return ctx.Response().Success().Json(payload)
}

func (e *ErrorHandler) FormError(ctx http.Context, msg any) http.Response {
	return e.failWith(ctx, FormError, msg)
}

func (e *ErrorHandler) LimitError(ctx http.Context, msg any) http.Response {
	return e.failWith(ctx, LimitError, msg)
}

func (e *ErrorHandler) UnauthorizedError(ctx http.Context, msg any) http.Response {
	return e.failWith(ctx, Unauthorized, msg)
}

func (e *ErrorHandler) PermissionError(ctx http.Context, msg any) http.Response {
	return e.failWith(ctx, PermissionError, msg)
}
