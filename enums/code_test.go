package enums

import "testing"

// TestCodeEnumValues 锁定对外协议里的数字，改这些值会破坏前端约定。
func TestCodeEnumValues(t *testing.T) {
	cases := []struct {
		name string
		got  int
		want int
	}{
		{"Success", int(Success), 200},
		{"Failed", int(Failed), 400},
		{"Unauthorized", int(Unauthorized), 401},
		{"PermissionDenied", int(PermissionDenied), 403},
		{"NotFound", int(NotFound), 404},
		{"ServerError", int(ServerError), 500},
	}

	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %d, want %d", c.name, c.got, c.want)
		}
	}
}

// TestCodeEnumMessage 每个状态码都要有可读文案，不能落到 Unknown error。
func TestCodeEnumMessage(t *testing.T) {
	cases := []struct {
		code CodeEnum
		want string
	}{
		{Success, "Success"},
		{Failed, "Failed"},
		{Unauthorized, "Unauthorized"},
		{PermissionDenied, "PermissionDenied"},
		{NotFound, "NotFound"},
		{ServerError, "ServerError"},
		// 未定义值统一回落
		{CodeEnum(9999), "Unknown error"},
		{CodeEnum(0), "Unknown error"},
	}

	for _, c := range cases {
		if got := c.code.Message(); got != c.want {
			t.Errorf("CodeEnum(%d).Message() = %q, want %q", int(c.code), got, c.want)
		}
	}
}

// TestStatusEnumValues 前端靠 status 判断成功与否，0 成功 / 1 失败是硬约定。
func TestStatusEnumValues(t *testing.T) {
	if int(StatusSuccess) != 0 {
		t.Errorf("StatusSuccess 应为 0，实际 %d", int(StatusSuccess))
	}
	if int(StatusFailed) != 1 {
		t.Errorf("StatusFailed 应为 1，实际 %d", int(StatusFailed))
	}
}
