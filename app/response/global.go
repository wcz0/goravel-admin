package response

var Unauthorized = map[string]any{
	"code":              401,
	"msg":               "Unauthenticated.",
	"data":              nil,
	"doNotDisplayToast": 0,
}

var FormError = map[string]any{
	"code":              400,
	"msg":               "Bad Request.",
	"data":              nil,
	"doNotDisplayToast": 0,
}

var LimitError = map[string]any{
	"code":              429,
	"msg":               "Forbidden.",
	"data":              nil,
	"doNotDisplayToast": 0,
}

var PermissionError = map[string]any{
	"code":              403,
	"msg":               "Permission denied.",
	"data":              nil,
	"doNotDisplayToast": 0,
}
