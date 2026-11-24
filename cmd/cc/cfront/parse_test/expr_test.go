package parsetest

import (
	"testing"

	"cc260717/internal/unittest"
)

func TestTernaryExpr(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "tenary",
			Src:    "1 ? 2 : 3",
			Expect: "(ternary 1 2 3)",
		},
	}
	check_expr(t, tt)
}
