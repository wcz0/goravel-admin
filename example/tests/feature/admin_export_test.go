package feature

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"
	"github.com/xuri/excelize/v2"
	"github.com/wcz0/goravel-admin/tools"

	"goravel/tests"
)

// AdminExportTestSuite 覆盖导出链路的容器相关内容：
// 真实磁盘写入、xlsx 结构、签名校验。纯签名算法在 tools 包里已有单测。
type AdminExportTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestAdminExportTestSuite(t *testing.T) {
	suite.Run(t, new(AdminExportTestSuite))
}

// buildAndCleanup 生成导出文件，并在用例结束时删除。
func (s *AdminExportTestSuite) buildAndCleanup(sheet string, headers []string, rows [][]any) *tools.ExportFile {
	file, err := tools.BuildExport("public", sheet, headers, rows)
	s.Require().NoError(err)
	s.Require().NotNil(file)

	s.T().Cleanup(func() {
		_ = facades.Storage().Disk("public").Delete(file.Path)
	})

	return file
}

// cell 读取指定单元格，失败直接终止用例。
func (s *AdminExportTestSuite) cell(book *excelize.File, sheet, axis string) string {
	value, err := book.GetCellValue(sheet, axis)
	s.Require().NoError(err)

	return value
}

// TestBuildExportCreatesReadableXlsx 生成的文件能被 Excel 解析，内容与入参一致。
func (s *AdminExportTestSuite) TestBuildExportCreatesReadableXlsx() {
	file := s.buildAndCleanup("用户列表",
		[]string{"ID", "用户名", "角色"},
		[][]any{
			{1, "admin", "超级管理员"},
			{2, "smoke", "测试账号"},
		},
	)

	s.True(strings.HasPrefix(file.Path, "exports/"), "导出文件应落在 exports/ 下，实际 %q", file.Path)
	s.True(strings.HasSuffix(file.Path, ".xlsx"), "应以 .xlsx 结尾，实际 %q", file.Path)
	s.NotEmpty(file.Signature, "签名不能为空")
	s.NotEmpty(file.Expires, "过期时间不能为空")

	raw, err := facades.Storage().Disk("public").Get(file.Path)
	s.Require().NoError(err)
	s.True(strings.HasPrefix(raw, "PK"), "xlsx 是 zip 容器，应以 PK 开头")

	book, err := excelize.OpenReader(strings.NewReader(raw))
	s.Require().NoError(err)
	defer func() { _ = book.Close() }()

	sheets := book.GetSheetList()
	s.Require().Len(sheets, 1)
	s.Equal("用户列表", sheets[0], "中文工作表名应保留")

	for i, want := range []string{"ID", "用户名", "角色"} {
		axis, err := excelize.CoordinatesToCellName(i+1, 1)
		s.Require().NoError(err)
		s.Equal(want, s.cell(book, sheets[0], axis), "第 %d 列表头", i+1)
	}

	s.Equal("admin", s.cell(book, sheets[0], "B2"))
	s.Equal("超级管理员", s.cell(book, sheets[0], "C2"))
	s.Equal("smoke", s.cell(book, sheets[0], "B3"))
}

// TestBuildExportLongSheetName 超长中文工作表名被截到 31 字符，且不能出现乱码。
func (s *AdminExportTestSuite) TestBuildExportLongSheetName() {
	longName := strings.Repeat("统计报表", 10) // 40 个字符
	file := s.buildAndCleanup(longName, []string{"列"}, [][]any{{1}})

	raw, err := facades.Storage().Disk("public").Get(file.Path)
	s.Require().NoError(err)

	book, err := excelize.OpenReader(strings.NewReader(raw))
	s.Require().NoError(err)
	defer func() { _ = book.Close() }()

	sheets := book.GetSheetList()
	s.Require().Len(sheets, 1)
	// 40 个字符截断到 31 个（统计报表 ×7 + 统计报）
	s.Equal(string([]rune(longName)[:31]), sheets[0], "应截断到 31 个字符且不乱码")
	s.Equal(31, len([]rune(sheets[0])))
}

// TestBuildExportRejectsEmptyHeaders 没有表头直接报错，避免产出无意义的空表。
func (s *AdminExportTestSuite) TestBuildExportRejectsEmptyHeaders() {
	file, err := tools.BuildExport("public", "空表", nil, nil)
	s.Error(err)
	s.Nil(file)
}

// TestVerifyExportRoundTrip 正常生成的签名能通过校验。
func (s *AdminExportTestSuite) TestVerifyExportRoundTrip() {
	file := s.buildAndCleanup("签名校验", []string{"ID"}, [][]any{{1}})

	s.NoError(tools.VerifyExport(file.Path, file.Expires, file.Signature))
}

// TestVerifyExportRejectsTamperedSignature 篡改签名必须被拒。
func (s *AdminExportTestSuite) TestVerifyExportRejectsTamperedSignature() {
	file := s.buildAndCleanup("签名校验", []string{"ID"}, [][]any{{1}})

	tampered := "deadbeef" + file.Signature[8:]
	s.Require().NotEqual(file.Signature, tampered)
	s.Error(tools.VerifyExport(file.Path, file.Expires, tampered))
}

// TestVerifyExportRejectsExtendedExpiry 延长有效期但沿用旧签名必须被拒。
func (s *AdminExportTestSuite) TestVerifyExportRejectsExtendedExpiry() {
	file := s.buildAndCleanup("签名校验", []string{"ID"}, [][]any{{1}})

	future := strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10)
	s.Error(tools.VerifyExport(file.Path, future, file.Signature))
}

// TestVerifyExportRejectsExpired 过期链接被拒（过期判断发生在签名比较之前）。
func (s *AdminExportTestSuite) TestVerifyExportRejectsExpired() {
	past := strconv.FormatInt(time.Now().Add(-time.Minute).Unix(), 10)
	err := tools.VerifyExport("exports/anything.xlsx", past, strings.Repeat("0", 64))

	s.Require().Error(err)
	s.Contains(err.Error(), "过期")
}

// TestVerifyExportRejectsTraversal 任意路径读取必须被挡住。
func (s *AdminExportTestSuite) TestVerifyExportRejectsTraversal() {
	file := s.buildAndCleanup("签名校验", []string{"ID"}, [][]any{{1}})

	s.Error(tools.VerifyExport("../.env.xlsx", file.Expires, file.Signature))
	s.Error(tools.VerifyExport("exports/demo.txt", file.Expires, file.Signature))
}

// TestBuildExportSheetNameFallback 空表名回落到 Sheet1。
func (s *AdminExportTestSuite) TestBuildExportSheetNameFallback() {
	file := s.buildAndCleanup("", []string{"ID"}, [][]any{{1}})

	raw, err := facades.Storage().Disk("public").Get(file.Path)
	s.Require().NoError(err)

	book, err := excelize.OpenReader(strings.NewReader(raw))
	s.Require().NoError(err)
	defer func() { _ = book.Close() }()

	s.Equal([]string{"Sheet1"}, book.GetSheetList())
}
