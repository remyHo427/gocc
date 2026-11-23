package parsetest

import (
	"testing"

	"cc260717/internal/unittest"
)

func TestAllUnaryOperaotrs(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "decrement",
			Src:    "--1",
			Expect: "(unary DEC 1)",
		},
		{
			Name:   "negate",
			Src:    "-1",
			Expect: "(unary SUB 1)",
		},
		{
			Name:   "bitwise not",
			Src:    "~1",
			Expect: "(unary BCOMP 1)",
		},
		{
			Name:   "logical not",
			Src:    "!1",
			Expect: "(unary NOT 1)",
		},
	}
	check_expr(t, tt)
}

func TestUnaryOperatorsAssociativity(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "decrement",
			Src:    "----1",
			Expect: "(unary DEC (unary DEC 1))",
		},
		{
			Name:   "negate",
			Src:    "-(-1)",
			Expect: "(unary SUB (unary SUB 1))",
		},
		{
			Name:   "bitwise not",
			Src:    "~~1",
			Expect: "(unary BCOMP (unary BCOMP 1))",
		},
		{
			Name:   "logical not",
			Src:    "!!1",
			Expect: "(unary NOT (unary NOT 1))",
		},
	}
	check_expr(t, tt)
}

func TestUnaryEqualPrecedenceAssociativity(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "decrement and negate",
			Src:    "--(-1)",
			Expect: "(unary DEC (unary SUB 1))",
		},
		{
			Name:   "negate and bitwise not",
			Src:    "-~1",
			Expect: "(unary SUB (unary BCOMP 1))",
		},
		{
			Name:   "bitwise not and logical not",
			Src:    "~!1",
			Expect: "(unary BCOMP (unary NOT 1))",
		},
		{
			Name:   "all",
			Src:    "--(-~!1))",
			Expect: "(unary DEC (unary SUB (unary BCOMP (unary NOT 1))))",
		},
	}
	check_expr(t, tt)
}
