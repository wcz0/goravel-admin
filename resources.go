package admin

import (
	"embed"
	"io/fs"
)

//go:embed all:assets/dist
var distEmbed embed.FS

//go:embed all:lang
var langEmbed embed.FS

// DistFS 返回内嵌的后台前端构建产物。
func DistFS() fs.FS {
	sub, err := fs.Sub(distEmbed, "assets/dist")
	if err != nil {
		panic(err)
	}

	return sub
}

// LangFS 返回内嵌的多语言文件（en、zh）。
func LangFS() fs.FS {
	sub, err := fs.Sub(langEmbed, "lang")
	if err != nil {
		panic(err)
	}

	return sub
}
