package core

import (
	"strings"
	"testing"
	"time"

	"github.com/mojocn/base64Captcha"
)

// useMemoryStore 把验证码存储临时切换到内存实现，避免依赖应用容器。
func useMemoryStore(t *testing.T) base64Captcha.Store {
	t.Helper()

	memory := base64Captcha.NewMemoryStore(100, 10*time.Minute)
	previous := captchaStore
	captchaStore = memory

	t.Cleanup(func() { captchaStore = previous })

	return memory
}

// TestCaptchaRoundTrip 验证「生成 → 取答案 → 校验通过 → 二次校验失败」。
func TestCaptchaRoundTrip(t *testing.T) {
	memory := useMemoryStore(t)

	id, image, err := GenerateCaptcha()
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}
	if id == "" {
		t.Fatal("验证码 id 为空")
	}
	if !strings.HasPrefix(image, "data:image/png;base64,") {
		t.Fatalf("验证码图片不是 data URI: %.40s", image)
	}

	answer := memory.Get(id, false)
	if answer == "" {
		t.Fatal("答案没有写进存储")
	}

	if !VerifyCaptcha(id, answer) {
		t.Fatal("正确答案应当校验通过")
	}
	// 校验后应当立即失效，防止同一张验证码被复用
	if VerifyCaptcha(id, answer) {
		t.Fatal("验证码应当只能使用一次")
	}
}

// TestVerifyCaptchaRejectsEmpty 验证空参数直接拒绝。
func TestVerifyCaptchaRejectsEmpty(t *testing.T) {
	useMemoryStore(t)

	if VerifyCaptcha("", "1234") {
		t.Fatal("空 id 不应通过")
	}
	if VerifyCaptcha("someid", "  ") {
		t.Fatal("空答案不应通过")
	}
}

// TestVerifyCaptchaRejectsWrongAnswer 验证错误答案被拒。
func TestVerifyCaptchaRejectsWrongAnswer(t *testing.T) {
	memory := useMemoryStore(t)

	id, _, err := GenerateCaptcha()
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	answer := memory.Get(id, false)
	if answer == "" {
		t.Fatal("答案没有写进存储")
	}

	wrong := "!!!!"
	if strings.EqualFold(answer, wrong) {
		wrong = "####"
	}
	if VerifyCaptcha(id, wrong) {
		t.Fatal("错误答案不应通过")
	}
}
