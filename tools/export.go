package tools

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/xuri/excelize/v2"
)

// exportTTL 是导出文件链接的有效期。
const exportTTL = 10 * time.Minute

// exportDir 是导出文件的存储目录。
const exportDir = "exports"

// maxSheetNameLen 是 Excel 允许的工作表名最大字符数。
const maxSheetNameLen = 31

// ExportFile 描述一次导出的结果。
//
// 前端拿到这三个字段后会直接拼出下载地址，后端在下载时重新校验签名，
// 避免任意路径读取。
type ExportFile struct {
	Path      string `json:"path"`
	Expires   string `json:"expires"`
	Signature string `json:"signature"`
}

// signWithKey 用给定密钥对「路径 + 过期时间」做 HMAC 签名。
//
// 单独拆出来是为了能在不启动应用容器的情况下测试签名算法本身。
func signWithKey(key string, path string, expires string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(path + "|" + expires))

	return hex.EncodeToString(mac.Sum(nil))
}

// signExport 用 APP_KEY 对「路径 + 过期时间」做 HMAC 签名。
func signExport(path string, expires string) string {
	return signWithKey(facades.Config().GetString("app.key"), path, expires)
}

// BuildExport 把表头与数据行渲染成 xlsx 并返回签名后的下载信息。
func BuildExport(disk string, sheetName string, headers []string, rows [][]any) (*ExportFile, error) {
	if len(headers) == 0 {
		return nil, errors.New("导出的列不能为空")
	}
	if disk == "" {
		disk = "public"
	}

	file := excelize.NewFile()
	defer func() { _ = file.Close() }()

	// 直接重命名默认工作表，避免新建后再删除导致内部 sheet 文件名不连续
	sheet := sanitizeSheetName(sheetName)
	if sheet == "" || sheet == "export" {
		sheet = "Sheet1"
	}
	if err := file.SetSheetName("Sheet1", sheet); err != nil {
		return nil, err
	}

	// 表头
	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return nil, err
		}
		if err := file.SetCellValue(sheet, cell, header); err != nil {
			return nil, err
		}
	}

	// 数据行
	for rowIndex, row := range rows {
		for colIndex, value := range row {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			if err != nil {
				return nil, err
			}
			if err := file.SetCellValue(sheet, cell, value); err != nil {
				return nil, err
			}
		}
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	// 文件名带时间戳，避免同名覆盖
	path := fmt.Sprintf("%s/%d_%s.xlsx", exportDir, time.Now().UnixNano(), sanitizeSheetName(sheet))
	storage := facades.Storage().Disk(disk)
	if err := storage.Put(path, buffer.String()); err != nil {
		return nil, err
	}

	expires := strconv.FormatInt(time.Now().Add(exportTTL).Unix(), 10)

	return &ExportFile{
		Path:      path,
		Expires:   expires,
		Signature: signExport(path, expires),
	}, nil
}

// VerifyExport 校验下载请求的签名与有效期。
func VerifyExport(path string, expires string, signature string) error {
	return verifyExport(facades.Config().GetString("app.key"), path, expires, signature)
}

// verifyExport 是 VerifyExport 的实际实现，密钥由调用方传入。
//
// APP_KEY 为空时直接拒绝：此时任何人都能算出同样的签名，校验形同虚设。
func verifyExport(key string, path string, expires string, signature string) error {
	if path == "" || expires == "" || signature == "" {
		return errors.New("参数不完整")
	}
	if key == "" {
		return errors.New("APP_KEY 未配置，无法校验导出链接")
	}
	if strings.ContainsRune(path, 0) || strings.Contains(path, "..") {
		return errors.New("路径不合法")
	}
	if !strings.HasSuffix(strings.ToLower(path), ".xlsx") {
		return errors.New("只允许下载 xlsx 文件")
	}

	deadline, err := strconv.ParseInt(expires, 10, 64)
	if err != nil {
		return errors.New("过期时间不合法")
	}
	if time.Now().Unix() > deadline {
		return errors.New("下载链接已过期")
	}

	if !hmac.Equal([]byte(signWithKey(key, path, expires)), []byte(signature)) {
		return errors.New("签名校验失败")
	}

	return nil
}

// sanitizeSheetName 去掉工作表名里 Excel 不允许的字符。
//
// Excel 限制工作表名最多 31 个字符。这里按 rune 截断而不是按字节：
// 按字节切会把中文表名切成非法 UTF-8，open 时会报错。
func sanitizeSheetName(name string) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_", ":", "_", "*", "_", "?", "_", "[", "_", "]", "_")
	cleaned := replacer.Replace(name)
	if cleaned == "" {
		return "export"
	}

	runes := []rune(cleaned)
	if len(runes) > maxSheetNameLen {
		cleaned = string(runes[:maxSheetNameLen])
	}

	return cleaned
}
