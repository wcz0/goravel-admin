package response

import (
	"context"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	contractshttp "github.com/goravel/framework/contracts/http"
)

// captureContext 是最小可用的 http.Context 实现，只负责接住控制器写出的 JSON，
// 便于在不起应用容器的情况下断言响应契约。
type captureContext struct {
	ctx        context.Context
	httpStatus int
	body       any
}

func newCaptureContext() *captureContext {
	return &captureContext{ctx: context.Background()}
}

func (c *captureContext) Deadline() (time.Time, bool) { return c.ctx.Deadline() }
func (c *captureContext) Done() <-chan struct{}       { return c.ctx.Done() }
func (c *captureContext) Err() error                  { return c.ctx.Err() }
func (c *captureContext) Value(key any) any           { return c.ctx.Value(key) }

func (c *captureContext) Context() context.Context        { return c.ctx }
func (c *captureContext) WithContext(ctx context.Context) { c.ctx = ctx }
func (c *captureContext) WithValue(key any, value any)    {}
func (c *captureContext) Request() contractshttp.ContextRequest {
	return nil
}
func (c *captureContext) Response() contractshttp.ContextResponse {
	return &captureResponse{ctx: c}
}

type captureResponse struct {
	ctx *captureContext
}

func (r *captureResponse) Cookie(cookie contractshttp.Cookie) contractshttp.ContextResponse {
	return r
}
func (r *captureResponse) Data(code int, contentType string, data []byte) contractshttp.AbortableResponse {
	return &noopAbortable{}
}
func (r *captureResponse) Download(filepath, filename string) contractshttp.Response {
	return &noopAbortable{}
}
func (r *captureResponse) File(filepath string) contractshttp.Response { return &noopAbortable{} }
func (r *captureResponse) Header(key, value string) contractshttp.ContextResponse {
	return r
}
func (r *captureResponse) Json(code int, obj any) contractshttp.AbortableResponse {
	r.ctx.httpStatus = code
	r.ctx.body = obj

	return &noopAbortable{}
}
func (r *captureResponse) NoContent(code ...int) contractshttp.AbortableResponse {
	return &noopAbortable{}
}
func (r *captureResponse) Origin() contractshttp.ResponseOrigin { return nil }
func (r *captureResponse) Redirect(code int, location string) contractshttp.AbortableResponse {
	return &noopAbortable{}
}
func (r *captureResponse) String(code int, format string, values ...any) contractshttp.AbortableResponse {
	return &noopAbortable{}
}
func (r *captureResponse) Success() contractshttp.ResponseStatus {
	return &captureStatus{ctx: r.ctx, httpStatus: 200}
}
func (r *captureResponse) Status(code int) contractshttp.ResponseStatus {
	return &captureStatus{ctx: r.ctx, httpStatus: code}
}
func (r *captureResponse) Stream(code int, step func(w contractshttp.StreamWriter) error) contractshttp.Response {
	return &noopAbortable{}
}
func (r *captureResponse) View() contractshttp.ResponseView { return nil }
func (r *captureResponse) Writer() nethttp.ResponseWriter {
	return httptest.NewRecorder()
}
func (r *captureResponse) WithoutCookie(name string) contractshttp.ContextResponse { return r }
func (r *captureResponse) Flush()                                                  {}

type captureStatus struct {
	ctx        *captureContext
	httpStatus int
}

func (s *captureStatus) Data(contentType string, data []byte) contractshttp.AbortableResponse {
	return &noopAbortable{}
}
func (s *captureStatus) Json(obj any) contractshttp.AbortableResponse {
	s.ctx.httpStatus = s.httpStatus
	s.ctx.body = obj

	return &noopAbortable{}
}
func (s *captureStatus) String(format string, values ...any) contractshttp.AbortableResponse {
	return &noopAbortable{}
}
func (s *captureStatus) Stream(step func(w contractshttp.StreamWriter) error) contractshttp.Response {
	return &noopAbortable{}
}

type noopAbortable struct{}

func (n *noopAbortable) Render() error { return nil }
func (n *noopAbortable) Abort() error  { return nil }

// payload 取出捕获到的 JSON 并断言是对象。
//
// 控制器写出的是 contractshttp.Json（map[string]any 的命名类型），
// 这里统一转成 map 方便断言。
func payload(t *testing.T, ctx *captureContext) map[string]any {
	t.Helper()

	switch body := ctx.body.(type) {
	case contractshttp.Json:
		return map[string]any(body)
	case map[string]any:
		return body
	default:
		t.Fatalf("响应体不是 JSON 对象: %#v", ctx.body)

		return nil
	}
}
