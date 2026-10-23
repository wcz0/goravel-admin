package services

import "testing"

// TestClampPage 页码下限为 1，避免出现负数 offset。
func TestClampPage(t *testing.T) {
	cases := map[int]int{-5: 1, 0: 1, 1: 1, 7: 7}
	for in, want := range cases {
		if got := clampPage(in); got != want {
			t.Errorf("clampPage(%d) = %d, want %d", in, got, want)
		}
	}
}

// TestClampPerPage 每页条数来自请求，必须收敛，否则一次请求就能拖全表。
func TestClampPerPage(t *testing.T) {
	if got := clampPerPage(0); got != defaultPageSize {
		t.Errorf("clampPerPage(0) = %d, want %d", got, defaultPageSize)
	}
	if got := clampPerPage(-3); got != defaultPageSize {
		t.Errorf("clampPerPage(-3) = %d, want %d", got, defaultPageSize)
	}
	if got := clampPerPage(10); got != 10 {
		t.Errorf("clampPerPage(10) = %d, want 10", got)
	}
	if got := clampPerPage(maxPageSize + 1000); got != maxPageSize {
		t.Errorf("clampPerPage(超限) = %d, want %d", got, maxPageSize)
	}
}
