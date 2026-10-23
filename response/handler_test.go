package response

import (
	"testing"

	"github.com/wcz0/goravel-admin/enums"
)

// assertEnvelope 断言前端依赖的响应外壳字段。
func assertEnvelope(t *testing.T, body map[string]any, wantStatus enums.StatusEnum) {
	t.Helper()

	gotStatus, ok := body["status"]
	if !ok {
		t.Fatal("响应缺少 status 字段")
	}
	if gotStatus != wantStatus {
		t.Errorf("status = %v, want %v", gotStatus, wantStatus)
	}

	for _, field := range []string{"msg", "data", "doNotDisplayToast"} {
		if _, ok := body[field]; !ok {
			t.Errorf("响应缺少 %s 字段: %v", field, body)
		}
	}

	if body["doNotDisplayToast"] != 0 {
		t.Errorf("doNotDisplayToast 默认应为 0，实际 %v", body["doNotDisplayToast"])
	}
}

// TestSuccess 无参成功响应：空消息、data 为空数组。
func TestSuccess(t *testing.T) {
	ctx := newCaptureContext()
	NewErrorHandler().Success(ctx)

	body := payload(t, ctx)
	assertEnvelope(t, body, enums.StatusSuccess)

	if body["msg"] != "" {
		t.Errorf("msg 应为空，实际 %v", body["msg"])
	}
	if data, ok := body["data"].([]any); !ok || len(data) != 0 {
		t.Errorf("data 应为空数组，实际 %#v", body["data"])
	}
	if body["code"] != enums.Success {
		t.Errorf("code 应为 %v，实际 %v", enums.Success, body["code"])
	}
}

// TestSuccessMsg 成功但带提示语。
func TestSuccessMsg(t *testing.T) {
	ctx := newCaptureContext()
	NewErrorHandler().SuccessMsg(ctx, "保存成功")

	body := payload(t, ctx)
	assertEnvelope(t, body, enums.StatusSuccess)

	if body["msg"] != "保存成功" {
		t.Errorf("msg = %v", body["msg"])
	}
}

// TestSuccessData data 原样透传，不能包一层。
func TestSuccessData(t *testing.T) {
	ctx := newCaptureContext()
	NewErrorHandler().SuccessData(ctx, map[string]any{"items": []any{1, 2}})

	body := payload(t, ctx)
	assertEnvelope(t, body, enums.StatusSuccess)

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("data 应为对象，实际 %#v", body["data"])
	}
	items, ok := data["items"].([]any)
	if !ok || len(items) != 2 {
		t.Errorf("data.items 应透传，实际 %#v", data["items"])
	}
}

// TestSuccessMsgData 同时带消息与数据。
func TestSuccessMsgData(t *testing.T) {
	ctx := newCaptureContext()
	NewErrorHandler().SuccessMsgData(ctx, "ok", []any{"a"})

	body := payload(t, ctx)
	assertEnvelope(t, body, enums.StatusSuccess)

	if body["msg"] != "ok" {
		t.Errorf("msg = %v", body["msg"])
	}
	if data, ok := body["data"].([]any); !ok || len(data) != 1 {
		t.Errorf("data = %#v", body["data"])
	}
}

// TestSuccessCodeMsg 自定义业务码。
func TestSuccessCodeMsg(t *testing.T) {
	ctx := newCaptureContext()
	NewErrorHandler().SuccessCodeMsg(ctx, 2001, "created")

	body := payload(t, ctx)
	assertEnvelope(t, body, enums.StatusSuccess)

	if body["code"] != 2001 {
		t.Errorf("code = %v", body["code"])
	}
}

// TestFail 默认失败响应带 Failed 文案。
func TestFail(t *testing.T) {
	ctx := newCaptureContext()
	NewErrorHandler().Fail(ctx)

	body := payload(t, ctx)
	assertEnvelope(t, body, enums.StatusFailed)

	if body["msg"] != enums.Failed.Message() {
		t.Errorf("msg = %v，期望 %q", body["msg"], enums.Failed.Message())
	}
	if body["code"] != enums.Failed {
		t.Errorf("code = %v", body["code"])
	}
}

// TestFailMsg 失败提示语原样透传。
func TestFailMsg(t *testing.T) {
	ctx := newCaptureContext()
	NewErrorHandler().FailMsg(ctx, "用户名已存在")

	body := payload(t, ctx)
	assertEnvelope(t, body, enums.StatusFailed)

	if body["msg"] != "用户名已存在" {
		t.Errorf("msg = %v", body["msg"])
	}
}

// TestFailMsgData 失败时也允许带数据（例如字段级错误）。
func TestFailMsgData(t *testing.T) {
	ctx := newCaptureContext()
	NewErrorHandler().FailMsgData(ctx, "校验失败", map[string]any{"field": "name"})

	body := payload(t, ctx)
	assertEnvelope(t, body, enums.StatusFailed)

	data, ok := body["data"].(map[string]any)
	if !ok || data["field"] != "name" {
		t.Errorf("data = %#v", body["data"])
	}
}

// TestFailData 只有数据没有消息时用默认文案。
func TestFailData(t *testing.T) {
	ctx := newCaptureContext()
	NewErrorHandler().FailData(ctx, []any{1})

	body := payload(t, ctx)
	assertEnvelope(t, body, enums.StatusFailed)

	if body["msg"] != "Error." {
		t.Errorf("msg = %v", body["msg"])
	}
}

// TestFailCodeMsg 同步设置 HTTP 状态码与业务码。
func TestFailCodeMsg(t *testing.T) {
	ctx := newCaptureContext()
	NewErrorHandler().FailCodeMsg(ctx, 422, "参数错误")

	body := payload(t, ctx)
	assertEnvelope(t, body, enums.StatusFailed)

	if body["code"] != 422 {
		t.Errorf("业务码 = %v", body["code"])
	}
	if ctx.httpStatus != 422 {
		t.Errorf("HTTP 状态码 = %d，期望 422", ctx.httpStatus)
	}
}

// TestErrorHelpersUseDefaultMessage 未传消息时回落到模板里的默认文案。
func TestErrorHelpersUseDefaultMessage(t *testing.T) {
	handler := NewErrorHandler()

	cases := []struct {
		name   string
		call   func(ctx *captureContext)
		want   string
		status enums.StatusEnum
	}{
		{"FormError", func(ctx *captureContext) { handler.FormError(ctx, "") }, "Bad Request.", enums.StatusFailed},
		{"LimitError", func(ctx *captureContext) { handler.LimitError(ctx, "") }, "Forbidden.", enums.StatusFailed},
		{"UnauthorizedError", func(ctx *captureContext) { handler.UnauthorizedError(ctx, "") }, "Unauthenticated.", enums.StatusFailed},
		{"PermissionError", func(ctx *captureContext) { handler.PermissionError(ctx, "") }, "Permission denied.", enums.StatusFailed},
	}

	for _, c := range cases {
		ctx := newCaptureContext()
		c.call(ctx)
		body := payload(t, ctx)

		assertEnvelope(t, body, c.status)
		if body["msg"] != c.want {
			t.Errorf("%s: msg = %v，期望默认文案 %q", c.name, body["msg"], c.want)
		}
	}
}

// TestErrorHelpersAcceptCustomMessage 传了消息就应该用传入的文案，
// 而不是被固定文案覆盖。
func TestErrorHelpersAcceptCustomMessage(t *testing.T) {
	handler := NewErrorHandler()

	cases := []struct {
		name string
		call func(ctx *captureContext, msg string)
	}{
		{"FormError", func(ctx *captureContext, msg string) { handler.FormError(ctx, msg) }},
		{"LimitError", func(ctx *captureContext, msg string) { handler.LimitError(ctx, msg) }},
		{"UnauthorizedError", func(ctx *captureContext, msg string) { handler.UnauthorizedError(ctx, msg) }},
		{"PermissionError", func(ctx *captureContext, msg string) { handler.PermissionError(ctx, msg) }},
	}

	for _, c := range cases {
		ctx := newCaptureContext()
		c.call(ctx, "自定义提示")
		body := payload(t, ctx)

		if body["msg"] != "自定义提示" {
			t.Errorf("%s: msg = %v，期望自定义提示", c.name, body["msg"])
		}

		// 全局模板不能被就地改写，否则后续请求会串味
		if FormError["msg"] == "自定义提示" || Unauthorized["msg"] == "自定义提示" ||
			PermissionError["msg"] == "自定义提示" || LimitError["msg"] == "自定义提示" {
			t.Fatalf("%s: 修改了全局响应模板", c.name)
		}
	}
}

// TestGlobalTemplates 全局模板的字段值锁定，中间件与宿主可能直接引用。
func TestGlobalTemplates(t *testing.T) {
	cases := []struct {
		name string
		temp map[string]any
		code int
	}{
		{"Unauthorized", Unauthorized, 401},
		{"TokenExpired", TokenExpired, 401},
		{"FormError", FormError, 400},
		{"LimitError", LimitError, 429},
		{"PermissionError", PermissionError, 403},
	}

	for _, c := range cases {
		// 前端靠 status 是否为数字判断 AMis 包装，模板里必须有
		if c.temp["status"] != enums.StatusFailed {
			t.Errorf("%s: status = %v，期望 %v（缺少会被前端当成原始响应）", c.name, c.temp["status"], enums.StatusFailed)
		}
		if c.temp["code"] != c.code {
			t.Errorf("%s: code = %v，期望 %d", c.name, c.temp["code"], c.code)
		}
		if _, ok := c.temp["data"].([]any); !ok {
			t.Errorf("%s: data 应为空数组，实际 %#v", c.name, c.temp["data"])
		}
		if c.temp["doNotDisplayToast"] != 0 {
			t.Errorf("%s: doNotDisplayToast 应为 0", c.name)
		}
	}
}
