package tools

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func AdminLang(ctx http.Context, key string) string {
	return facades.Lang(ctx).Get("admin." + key)
}