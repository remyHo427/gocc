package parsetest

import "testing"

func TestAllUnaryOperaotrs(t *testing.T) {
	tt := []TestPair{
		{"decrement", "--1", "(unary DEC 1)"},
		{"negate", "-1", "(unary SUB 1)"},
		{"bitwise not", "~1", "(unary BCOMP 1)"},
		{"logical not", "!1", "(unary NOT 1)"},
	}
	check_expr(t, tt)
}

func TestUnaryOperatorsAssociativity(t *testing.T) {
	tt := []TestPair{
		{"decrement", "----1", "(unary DEC (unary DEC 1))"},
		{"negate", "-(-1)", "(unary SUB (unary SUB 1))"},
		{"bitwise not", "~~1", "(unary BCOMP (unary BCOMP 1))"},
		{"logical not", "!!1", "(unary NOT (unary NOT 1))"},
	}
	check_expr(t, tt)
}

func TestUnaryEqualPrecedenceAssociativity(t *testing.T) {
	tt := []TestPair{
		{"decrement and negate", "--(-1)", "(unary DEC (unary SUB 1))"},
		{"negate and bitwise not", "-~1", "(unary SUB (unary BCOMP 1))"},
		{"bitwise not and logical not", "~!1", "(unary BCOMP (unary NOT 1))"},
		{"all", "--(-~!1))",
			"(unary DEC (unary SUB (unary BCOMP (unary NOT 1))))"},
	}
	check_expr(t, tt)
}
