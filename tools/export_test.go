package tools

import (
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

const testExportKey = "base64:test-app-key-for-export-signature"

// futureExpires 返回一个未来时间戳字符串，避免测试随时间腐化。
func futureExpires(d time.Duration) string {
	return strconv.FormatInt(time.Now().Add(d).Unix(), 10)
}

// TestSignWithKeyIsDeterministic 签名必须稳定：同 key 同输入得到同结果，换任一要素就不同。
func TestSignWithKeyIsDeterministic(t *testing.T) {
	path := "exports/1_demo.xlsx"
	expires := "1760000000"

	base := signWithKey(testExportKey, path, expires)
	if len(base) != 64 {
		t.Fatalf("HMAC-SHA256 十六进制应是 64 字符，实际 %d 字符", len(base))
	}
	if again := signWithKey(testExportKey, path, expires); again != base {
		t.Error("相同输入应得到相同签名")
	}
	if got := signWithKey(testExportKey+"x", path, expires); got == base {
		t.Error("换密钥后签名应改变")
	}
	if got := signWithKey(testExportKey, "exports/2_demo.xlsx", expires); got == base {
		t.Error("换路径后签名应改变")
	}
	if got := signWithKey(testExportKey, path, "1760000001"); got == base {
		t.Error("换过期时间后签名应改变")
	}
}

// TestVerifyExportAcceptsFreshSignature 正常链接应通过校验。
func TestVerifyExportAcceptsFreshSignature(t *testing.T) {
	path := "exports/1729000000000000000_用户列表.xlsx"
	expires := futureExpires(5 * time.Minute)
	signature := signWithKey(testExportKey, path, expires)

	if err := verifyExport(testExportKey, path, expires, signature); err != nil {
		t.Fatalf("有效链接应通过校验，实际报错: %v", err)
	}
}

// TestVerifyExportAcceptsUppercaseExtension 后缀判断不区分大小写。
func TestVerifyExportAcceptsUppercaseExtension(t *testing.T) {
	path := "exports/demo.XLSX"
	expires := futureExpires(time.Minute)
	signature := signWithKey(testExportKey, path, expires)

	if err := verifyExport(testExportKey, path, expires, signature); err != nil {
		t.Fatalf("大写后缀应通过，实际报错: %v", err)
	}
}

// TestVerifyExportRejectsTamperedSignature 篡改签名必须失败。
func TestVerifyExportRejectsTamperedSignature(t *testing.T) {
	path := "exports/demo.xlsx"
	expires := futureExpires(time.Minute)
	signature := signWithKey(testExportKey, path, expires)

	tampered := "0" + signature[1:]
	if tampered == signature {
		tampered = "1" + signature[1:]
	}

	err := verifyExport(testExportKey, path, expires, tampered)
	if err == nil {
		t.Fatal("篡改签名不应通过")
	}
	if !strings.Contains(err.Error(), "签名") {
		t.Errorf("应提示签名错误，实际: %v", err)
	}
}

// TestVerifyExportRejectsExtendedExpiry 改长有效期但沿用旧签名必须失败，
// 否则拿到的链接可以被无限续期。
func TestVerifyExportRejectsExtendedExpiry(t *testing.T) {
	path := "exports/demo.xlsx"
	expires := futureExpires(time.Minute)
	signature := signWithKey(testExportKey, path, expires)

	longer := futureExpires(24 * time.Hour)
	if err := verifyExport(testExportKey, path, longer, signature); err == nil {
		t.Fatal("续期后的链接不应通过")
	}
}

// TestVerifyExportRejectsExpired 过期链接必须失败。
func TestVerifyExportRejectsExpired(t *testing.T) {
	path := "exports/demo.xlsx"
	expires := futureExpires(-time.Minute)
	signature := signWithKey(testExportKey, path, expires)

	err := verifyExport(testExportKey, path, expires, signature)
	if err == nil {
		t.Fatal("过期链接不应通过")
	}
	if !strings.Contains(err.Error(), "过期") {
		t.Errorf("应提示已过期，实际: %v", err)
	}
}

// TestVerifyExportRejectsBadPath 路径穿越、NUL、非 xlsx 后缀都必须挡住。
func TestVerifyExportRejectsBadPath(t *testing.T) {
	expires := futureExpires(time.Minute)

	cases := []struct {
		name string
		path string
	}{
		{"上级目录", "exports/../../.env.xlsx"},
		{"Windows 上级目录", "exports\\..\\secret.xlsx"},
		{"NUL 截断", "exports/demo.xlsx\x00.png"},
		{"非 xlsx", "exports/config.json"},
		{"无后缀", "exports/demo"},
	}

	for _, c := range cases {
		signature := signWithKey(testExportKey, c.path, expires)
		if err := verifyExport(testExportKey, c.path, expires, signature); err == nil {
			t.Errorf("%s: 路径 %q 不应通过校验", c.name, c.path)
		}
	}
}

// TestVerifyExportRejectsMissingParams 缺参直接拒绝，不进入签名比较。
func TestVerifyExportRejectsMissingParams(t *testing.T) {
	valid := futureExpires(time.Minute)
	path := "exports/demo.xlsx"
	signature := signWithKey(testExportKey, path, valid)

	cases := []struct{ path, expires, signature string }{
		{"", valid, signature},
		{path, "", signature},
		{path, valid, ""},
	}

	for i, c := range cases {
		if err := verifyExport(testExportKey, c.path, c.expires, c.signature); err == nil {
			t.Errorf("第 %d 个缺参用例不应通过", i)
		}
	}
}

// TestVerifyExportRejectsEmptyKey APP_KEY 为空等价于无签名，必须 fail closed。
func TestVerifyExportRejectsEmptyKey(t *testing.T) {
	path := "exports/demo.xlsx"
	expires := futureExpires(time.Minute)
	// 空密钥下任何人算出的签名都"正确"，所以必须在比较前就拒绝
	signature := signWithKey("", path, expires)

	err := verifyExport("", path, expires, signature)
	if err == nil {
		t.Fatal("空 APP_KEY 不应通过校验")
	}
	if !strings.Contains(err.Error(), "APP_KEY") {
		t.Errorf("应提示 APP_KEY 未配置，实际: %v", err)
	}
}

// TestVerifyExportRejectsBadExpires 过期时间必须是十进制时间戳。
func TestVerifyExportRejectsBadExpires(t *testing.T) {
	path := "exports/demo.xlsx"
	signature := signWithKey(testExportKey, path, "not-a-number")

	err := verifyExport(testExportKey, path, "not-a-number", signature)
	if err == nil {
		t.Fatal("非法过期时间不应通过")
	}
	if !strings.Contains(err.Error(), "过期时间") {
		t.Errorf("应提示过期时间不合法，实际: %v", err)
	}
}

// TestSanitizeSheetName 工作表名要清掉 Excel 不允许的字符。
func TestSanitizeSheetName(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"用户列表", "用户列表"},
		{"a/b", "a_b"},
		{`a\b`, "a_b"},
		{"a:b", "a_b"},
		{"a*b", "a_b"},
		{"a?b", "a_b"},
		{"a[b]", "a_b_"},
		{"a[b]c", "a_b_c"},
		{"", "export"},
		// Excel 工作表名上限 31 字符；按 rune 截断，不能把中文切成半个字
		{strings.Repeat("长", 40), strings.Repeat("长", 31)},
		{"a" + strings.Repeat("b", 40), "a" + strings.Repeat("b", 30)},
	}

	for _, c := range cases {
		got := sanitizeSheetName(c.in)
		if got != c.want {
			t.Errorf("sanitizeSheetName(%q) = %q, want %q", c.in, got, c.want)
		}
		if len([]rune(got)) > maxSheetNameLen {
			t.Errorf("sanitizeSheetName(%q) 长度 %d 超过上限 %d", c.in, len([]rune(got)), maxSheetNameLen)
		}
		if !utf8.ValidString(got) {
			t.Errorf("sanitizeSheetName(%q) 产生了非法 UTF-8: %q", c.in, got)
		}
	}
}
