package tools

import (
	"os"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// KebabToCamel 中划线转小驼峰
func KebabToCamel(str string) string {
	arr := strings.Split(str, "-")

	for i, v := range arr {
		if i > 0 {
			arr[i] = cases.Title(language.English, cases.NoLower).String(v)
		} else {
			arr[i] = strings.ToLower(v)
		}
	}

	return strings.Join(arr, "")
}

// ReadDir 读取目录下的所有文件, 返回文件名列表
func ReadDir(dir string) []string {
	var files []string

	dirList, _ := os.ReadDir(dir)

	for _, v := range dirList {
		if !v.IsDir() {
			files = append(files, v.Name())
		}
	}

	return files
}


// PascalToCamel 大驼峰转小驼峰
func PascalToCamel(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// PascalToSnake 大驼峰转蛇形
func PascalToSnake(s string) string {
	var result string
	for i, v := range s {
		if unicode.IsUpper(v) {
			if i != 0 {
				result += "_"
			}
			result += string(unicode.ToLower(v))
		} else {
			result += string(v)
		}
	}
	return result
}

// SnakeToPascal 蛇形转大驼峰
func SnakeToPascal(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_'
	})

	for i := 0; i < len(words); i++ {
		runes := []rune(words[i])
		runes[0] = unicode.ToUpper(runes[0])
		words[i] = string(runes)
	}

	return strings.Join(words, "")
}