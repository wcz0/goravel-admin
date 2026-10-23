package response

import "github.com/wcz0/goravel-admin/enums"

// 全局响应模板。
//
// 必须带 status 字段：前端拦截器用 `typeof backend.status === 'number'` 判断
// 这是不是 AMis 风格的包装，缺了会把错误响应当成原始 axios 响应返回给业务层。
var Unauthorized = map[string]any{
	"status":            enums.StatusFailed,
	"code":              401,
	"msg":               "Unauthenticated.",
	"data":              []any{},
	"doNotDisplayToast": 0,
}

var TokenExpired = map[string]any{
	"status":            enums.StatusFailed,
	"code":              401,
	"msg":               "Authentication failed: Token has expired.",
	"data":              []any{},
	"doNotDisplayToast": 0,
}

var FormError = map[string]any{
	"status":            enums.StatusFailed,
	"code":              400,
	"msg":               "Bad Request.",
	"data":              []any{},
	"doNotDisplayToast": 0,
}

var LimitError = map[string]any{
	"status":            enums.StatusFailed,
	"code":              429,
	"msg":               "Forbidden.",
	"data":              []any{},
	"doNotDisplayToast": 0,
}

var PermissionError = map[string]any{
	"status":            enums.StatusFailed,
	"code":              403,
	"msg":               "Permission denied.",
	"data":              []any{},
	"doNotDisplayToast": 0,
}
