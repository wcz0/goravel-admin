package tools

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
)

func RequestIs(ctx http.Context, path string) bool {
	requestPath := strings.TrimSuffix(ctx.Request().Path(), "/")
	path = strings.TrimSuffix(path, "/")

	if requestPath == path {
		return true
	}

	if strings.HasSuffix(path, "*") {
		basePath := strings.TrimSuffix(path, "*")
		if strings.HasPrefix(requestPath, basePath) {
			return true
		}
	}

	return false
}
