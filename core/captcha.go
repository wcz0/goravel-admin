package core

import (
	"strings"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/mojocn/base64Captcha"
)

// captchaTTL 是验证码在缓存中的存活时间，与 owl-admin 的 600 秒保持一致。
const captchaTTL = 10 * time.Minute

// captchaCacheKey 为验证码缓存键前缀。
const captchaCacheKey = "admin_captcha_"

// cacheStore 把验证码答案存进 Goravel 的 Cache，而不是库自带的内存 store。
//
// 内置的内存 store 在进程重启或多实例部署时会丢答案；走 Cache 之后可以复用
// 宿主配置的 file / redis 驱动，行为与 owl-admin 的 cache()->put() 对齐。
type cacheStore struct{}

func (s *cacheStore) Set(id string, value string) error {
	return facades.Cache().Put(captchaCacheKey+id, value, captchaTTL)
}

func (s *cacheStore) Get(id string, clear bool) string {
	key := captchaCacheKey + id
	value := facades.Cache().Get(key, "")
	if clear {
		facades.Cache().Forget(key)
	}

	str, ok := value.(string)
	if !ok {
		return ""
	}

	return str
}

func (s *cacheStore) Verify(id, answer string, clear bool) bool {
	return strings.EqualFold(strings.TrimSpace(s.Get(id, clear)), strings.TrimSpace(answer))
}

// captchaStore 是验证码答案的存储实现。
//
// 默认走 Cache，方便复用宿主配置的驱动；单测里会替换成内存实现，
// 这样无需启动应用也能验证「生成 → 校验 → 一次性失效」的完整链路。
var captchaStore base64Captcha.Store = &cacheStore{}

// newCaptcha 构造一个 4 位数字验证码实例。
//
// 宽高按后台登录页的布局取值（100x40），与 owl-admin 的 Captcha 默认尺寸一致。
func newCaptcha() *base64Captcha.Captcha {
	driver := base64Captcha.NewDriverDigit(40, 100, 4, 0.7, 80)

	return base64Captcha.NewCaptcha(driver, captchaStore)
}

// GenerateCaptcha 生成一张验证码。
//
// 返回的 image 是可直接放进 <img src> 的 data URI，id 需要随登录请求一并提交，
// 用于从缓存里取回正确答案。
func GenerateCaptcha() (id string, image string, err error) {
	id, image, _, err = newCaptcha().Generate()

	return id, image, err
}

// VerifyCaptcha 校验验证码，校验后立即失效（一次性使用）。
func VerifyCaptcha(id string, answer string) bool {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(answer) == "" {
		return false
	}

	return newCaptcha().Verify(id, answer, true)
}
