package parsetest

import "testing"

func TestPrimary(t *testing.T) {
	tt := []TestPair{
		{"simple integer", "1", "1"},
		{"simple integer nested", "(1)", "1"},
	}
	check_expr(t, tt)
}
