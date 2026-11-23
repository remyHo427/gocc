package parsetest

import (
	"testing"

	"cc260717/internal/unittest"
)

func TestPrimary(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "simple integer",
			Src:    "1",
			Expect: "1",
		},
		{
			Name:   "simple integer nested",
			Src:    "(1)",
			Expect: "1",
		},
	}
	check_expr(t, tt)
}
