// Package icons 提供后台图标选择器所需的 iconify 图标名检索。
//
// 数据来自 iconify 全量集合快照（22 万个图标名），编译期通过 embed 打进二进制，
// 首次调用时解析一次并常驻内存，避免像 PHP 版那样每次请求都读 6MB 文件。
package icons

import (
	_ "embed"
	"encoding/json"
	"strings"
	"sync"
)

//go:embed resources/iconify.json
var rawData []byte

// maxResults 是单次搜索返回的条数上限，与 owl-admin 保持一致。
const maxResults = 1000

var (
	once     sync.Once
	iconList []string
)

// load 懒加载图标名列表。
func load() []string {
	once.Do(func() {
		_ = json.Unmarshal(rawData, &iconList)
	})

	return iconList
}

// Search 按关键字做子串匹配，返回匹配到的图标名（最多 1000 条）。
//
// 空关键字回落到 "home"，与前端默认行为一致。
func Search(keyword string) []string {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		keyword = "home"
	}
	keyword = strings.ToLower(keyword)

	all := load()
	result := make([]string, 0, 32)
	for _, name := range all {
		if strings.Contains(strings.ToLower(name), keyword) {
			result = append(result, name)
			if len(result) >= maxResults {
				break
			}
		}
	}

	return result
}

// Count 返回图标总数，供诊断使用。
func Count() int {
	return len(load())
}
