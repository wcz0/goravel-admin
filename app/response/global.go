package response

var Unauthorized = map[string]any{
	"code":              401,
	"msg":               "Unauthenticated.",
	"data":              []any{},
	"doNotDisplayToast": 0,
}

var TokenExpired = map[string]any{
	"code":              401,
	"msg":               "Authentication failed: Token has expired.",
	"data":              []any{},
	"doNotDisplayToast": 0,
}

var FormError = map[string]any{
	"code":              400,
	"msg":               "Bad Request.",
	"data":              []any{},
	"doNotDisplayToast": 0,
}

var LimitError = map[string]any{
	"code":              429,
	"msg":               "Forbidden.",
	"data":              []any{},
	"doNotDisplayToast": 0,
}

var PermissionError = map[string]any{
	"code":              403,
	"msg":               "Permission denied.",
	"data":              []any{},
	"doNotDisplayToast": 0,
}
