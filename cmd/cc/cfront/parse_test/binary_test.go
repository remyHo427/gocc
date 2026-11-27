package parsetest

import (
	"testing"

	"cc260717/internal/unittest"
)

func TestAllBinaryOperators(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "arithmetic addition",
			Src:    "1 + 1",
			Expect: "(binary ADD 1 1)",
		},
		{
			Name:   "arithmetic subtraction",
			Src:    "1 - 1",
			Expect: "(binary SUB 1 1)",
		},
		{
			Name:   "arithmetic multiplication",
			Src:    "1 * 1",
			Expect: "(binary MUL 1 1)",
		},
		{
			Name:   "arithmetic division",
			Src:    "1 / 1",
			Expect: "(binary DIV 1 1)",
		},
		{
			Name:   "arithmetic remainder",
			Src:    "1 % 1",
			Expect: "(binary MOD 1 1)",
		},
		{
			Name:   "bitwise AND",
			Src:    "1 & 1",
			Expect: "(binary BAND 1 1)",
		},
		{
			Name:   "bitwise OR",
			Src:    "1 | 1",
			Expect: "(binary BOR 1 1)",
		},
		{
			Name:   "bitwise XOR",
			Src:    "1 ^ 1",
			Expect: "(binary BXOR 1 1)",
		},
		{
			Name:   "left shift",
			Src:    "1 << 1",
			Expect: "(binary LSHIFT 1 1)",
		},
		{
			Name:   "right shift",
			Src:    "1 >> 1",
			Expect: "(binary RSHIFT 1 1)",
		},
		{
			Name:   "logical or",
			Src:    "1 || 1",
			Expect: "(binary OR 1 1)",
		},
		{
			Name:   "logical and",
			Src:    "1 && 1",
			Expect: "(binary AND 1 1)",
		},
		{
			Name:   "equal",
			Src:    "1 == 1",
			Expect: "(binary EQ 1 1)",
		},
		{
			Name:   "not equal",
			Src:    "1 != 1",
			Expect: "(binary NEQ 1 1)",
		},
		{
			Name:   "less than",
			Src:    "1 < 1",
			Expect: "(binary LT 1 1)",
		},
		{
			Name:   "greater than",
			Src:    "1 > 1",
			Expect: "(binary GT 1 1)",
		},
		{
			Name:   "less than or equal",
			Src:    "1 <= 1",
			Expect: "(binary LEQ 1 1)",
		},
		{
			Name:   "greater than or equal",
			Src:    "1 >= 1",
			Expect: "(binary GEQ 1 1)",
		},
		{
			Name:   "ternary",
			Src:    "1 ? 1 : 0",
			Expect: "(ternary 1 1 0)",
		},
	}
	check_expr(t, tt)
}

func TestBinaryOperatorsAssociativity(t *testing.T) {
	tt := []unittest.NamedTestPair{
		// left-to-right
		{
			Name:   "arithmetic add",
			Src:    "1 + 2 + 3",
			Expect: "(binary ADD (binary ADD 1 2) 3)",
		},
		{
			Name:   "arithmetic sub",
			Src:    "1 - 2 - 3",
			Expect: "(binary SUB (binary SUB 1 2) 3)",
		},
		{
			Name:   "arithmetic mul",
			Src:    "1 * 2 * 3",
			Expect: "(binary MUL (binary MUL 1 2) 3)",
		},
		{
			Name:   "arithmetic div",
			Src:    "1 / 2 / 3",
			Expect: "(binary DIV (binary DIV 1 2) 3)",
		},
		{
			Name:   "arithmetic remainder",
			Src:    "1 % 2 % 3",
			Expect: "(binary MOD (binary MOD 1 2) 3)",
		},
		{
			Name:   "bitwise AND",
			Src:    "1 & 2 & 3",
			Expect: "(binary BAND (binary BAND 1 2) 3)",
		},
		{
			Name:   "bitwise OR",
			Src:    "1 | 2 | 3",
			Expect: "(binary BOR (binary BOR 1 2) 3)",
		},
		{
			Name:   "bitwise XOR",
			Src:    "1 ^ 2 ^ 3",
			Expect: "(binary BXOR (binary BXOR 1 2) 3)",
		},
		{
			Name:   "left shift",
			Src:    "1 << 2 << 3",
			Expect: "(binary LSHIFT (binary LSHIFT 1 2) 3)",
		},
		{
			Name:   "right shift",
			Src:    "1 >> 2 >> 3",
			Expect: "(binary RSHIFT (binary RSHIFT 1 2) 3)",
		},
		{
			Name:   "logical or",
			Src:    "1 || 2 || 3",
			Expect: "(binary OR (binary OR 1 2) 3)",
		},
		{
			Name:   "logical and",
			Src:    "1 && 2 && 3",
			Expect: "(binary AND (binary AND 1 2) 3)",
		},
		{
			Name:   "equal",
			Src:    "1 == 2 == 3",
			Expect: "(binary EQ (binary EQ 1 2) 3)",
		},
		{
			Name:   "not equal",
			Src:    "1 != 2 != 3",
			Expect: "(binary NEQ (binary NEQ 1 2) 3)",
		},
		{
			Name:   "less than",
			Src:    "1 < 2 < 3",
			Expect: "(binary LT (binary LT 1 2) 3)",
		},
		{
			Name:   "greater than",
			Src:    "1 > 2 > 3",
			Expect: "(binary GT (binary GT 1 2) 3)",
		},
		{
			Name:   "less than or equal",
			Src:    "1 <= 2 <= 3",
			Expect: "(binary LEQ (binary LEQ 1 2) 3)",
		},
		{
			Name:   "greater than or equal",
			Src:    "1 >= 2 >= 3",
			Expect: "(binary GEQ (binary GEQ 1 2) 3)",
		},

		// right-to-left
		{
			Name:   "assign",
			Src:    "1 = 2 = 3",
			Expect: "(assign 1 (assign 2 3))",
		},

		// 1 ? (2 ? 3 : 4) : (5 ? 6 : 7)
		{
			Name:   "ternary",
			Src:    "1 ? 2 ? 3 : 4 : 5 ? 6 : 7",
			Expect: "(ternary 1 (ternary 2 3 4) (ternary 5 6 7))",
		},
	}
	check_expr(t, tt)
}

// test mixing different operators with the same precedence level
func TestEqualPrecedenceAssociativity(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "arithmetic mul, div, reminder",
			Src:    "1 * 2 / 3 % 4",
			Expect: "(binary MOD (binary DIV (binary MUL 1 2) 3) 4)",
		},
		{
			Name:   "arithmetic add and sub",
			Src:    "1 + 2 - 3",
			Expect: "(binary SUB (binary ADD 1 2) 3)",
		},
		{
			Name:   "left shift and right shift",
			Src:    "1 << 2 >> 3",
			Expect: "(binary RSHIFT (binary LSHIFT 1 2) 3)",
		},
		{
			Name:   "less than, greater than, less than or equal, greater than or equal",
			Src:    "1 < 2 > 3 <= 4 >= 5",
			Expect: "(binary GEQ (binary LEQ (binary GT (binary LT 1 2) 3) 4) 5)",
		},
		{
			Name:   "equal and not equal",
			Src:    "1 == 2 != 3",
			Expect: "(binary NEQ (binary EQ 1 2) 3)",
		},
		// bitwise AND, bitwise XOR, bitwise OR, logical AND and logical OR
		// each have their own precedence level, therefore not tested here
	}
	check_expr(t, tt)
}

// test mixing operators with different precedence levels
func TestPrecedenceLevels(t *testing.T) {
	tt := []unittest.NamedTestPair{
		// + with * / %
		{
			Name:   "arithmetic add with arithmetic mul",
			Src:    "1 + 2 * 3",
			Expect: "(binary ADD 1 (binary MUL 2 3))",
		},
		{
			Name:   "arithmetic add with arithmetic div",
			Src:    "1 + 2 / 3",
			Expect: "(binary ADD 1 (binary DIV 2 3))",
		},
		{
			Name:   "arithmetic add with arithmetic remainder",
			Src:    "1 + 2 % 3",
			Expect: "(binary ADD 1 (binary MOD 2 3))",
		},
		// - with * / %
		{
			Name:   "arithmetic sub with arithmetic mul",
			Src:    "1 - 2 * 3",
			Expect: "(binary SUB 1 (binary MUL 2 3))",
		},
		{
			Name:   "arithmetic sub with arithmetic div",
			Src:    "1 - 2 / 3",
			Expect: "(binary SUB 1 (binary DIV 2 3))",
		},
		{
			Name:   "arithmetic sub with arithmetic remainder",
			Src:    "1 - 2 % 3",
			Expect: "(binary SUB 1 (binary MOD 2 3))",
		},
		// << with + -
		{
			Name:   "left shift with arithmetic add",
			Src:    "1 << 2 + 3",
			Expect: "(binary LSHIFT 1 (binary ADD 2 3))",
		},
		{
			Name:   "left shift with arithmetic sub",
			Src:    "1 << 2 - 3",
			Expect: "(binary LSHIFT 1 (binary SUB 2 3))",
		},
		// >> with + -
		{
			Name:   "right shift with arithmetic add",
			Src:    "1 >> 2 + 3",
			Expect: "(binary RSHIFT 1 (binary ADD 2 3))",
		},
		{
			Name:   "right shift with arithmetic sub",
			Src:    "1 >> 2 + 3",
			Expect: "(binary RSHIFT 1 (binary ADD 2 3))",
		},
		// < with << >>
		{
			Name:   "less than with left shift",
			Src:    "1 < 2 << 3",
			Expect: "(binary LT 1 (binary LSHIFT 2 3))",
		},
		{
			Name:   "less than with right shift",
			Src:    "1 < 2 >> 3",
			Expect: "(binary LT 1 (binary RSHIFT 2 3))",
		},
		// > with << >>
		{
			Name:   "greater than with left shift",
			Src:    "1 > 2 << 3",
			Expect: "(binary GT 1 (binary LSHIFT 2 3))",
		},
		{
			Name:   "greater than with right shift",
			Src:    "1 > 2 >> 3",
			Expect: "(binary GT 1 (binary RSHIFT 2 3))",
		},
		// <= with << >>
		{
			Name:   "less than or equal with left shift",
			Src:    "1 <= 2 << 3",
			Expect: "(binary LEQ 1 (binary LSHIFT 2 3))",
		},
		{
			Name:   "less than or equal with right shift",
			Src:    "1 <= 2 >> 3",
			Expect: "(binary LEQ 1 (binary RSHIFT 2 3))",
		},
		// >= with << >>
		{
			Name:   "greater than or equal with left shift",
			Src:    "1 >= 2 << 3",
			Expect: "(binary GEQ 1 (binary LSHIFT 2 3))",
		},
		{
			Name:   "greater than or euqal with right shift",
			Src:    "1 >= 2 >> 3",
			Expect: "(binary GEQ 1 (binary RSHIFT 2 3))",
		},
		// == with < > <= >=
		{
			Name:   "equal with less than",
			Src:    "1 == 2 < 3",
			Expect: "(binary EQ 1 (binary LT 2 3))",
		},
		{
			Name:   "equal with greater than",
			Src:    "1 == 2 > 3",
			Expect: "(binary EQ 1 (binary GT 2 3))",
		},
		{
			Name:   "equal with less than or equal",
			Src:    "1 == 2 <= 3",
			Expect: "(binary EQ 1 (binary LEQ 2 3))",
		},
		{
			Name:   "equal with greater than or equal",
			Src:    "1 == 2 >= 3",
			Expect: "(binary EQ 1 (binary GEQ 2 3))",
		},
		// != with < > <= >=
		{
			Name:   "not equal with less than",
			Src:    "1 != 2 < 3",
			Expect: "(binary NEQ 1 (binary LT 2 3))",
		},
		{
			Name:   "not equal with greater than",
			Src:    "1 != 2 > 3",
			Expect: "(binary NEQ 1 (binary GT 2 3))",
		},
		{
			Name:   "not equal with less than or equal",
			Src:    "1 != 2 <= 3",
			Expect: "(binary NEQ 1 (binary LEQ 2 3))",
		},
		{
			Name:   "not equal with greater than or equal",
			Src:    "1 != 2 >= 3",
			Expect: "(binary NEQ 1 (binary GEQ 2 3))",
		},
		// & with == !=
		{
			Name:   "bitwise AND with equal",
			Src:    "1 & 2 == 3",
			Expect: "(binary BAND 1 (binary EQ 2 3))",
		},
		{
			Name:   "bitwise AND with not equal",
			Src:    "1 & 2 != 3",
			Expect: "(binary BAND 1 (binary NEQ 2 3))",
		},
		// ^ with &
		{
			Name:   "bitwise XOR with bitwise AND",
			Src:    "1 ^ 2 & 3",
			Expect: "(binary BXOR 1 (binary BAND 2 3))",
		},
		// | with ^
		{
			Name:   "bitwise OR with bitwise XOR",
			Src:    "1 | 2 ^ 3",
			Expect: "(binary BOR 1 (binary BXOR 2 3))",
		},
		// && with |
		{
			Name:   "bitwise AND with bitwise OR ",
			Src:    "1 && 2 | 3",
			Expect: "(binary AND 1 (binary BOR 2 3))",
		},
		// || with &&
		{
			Name:   "logical OR with logical AND",
			Src:    "1 || 2 && 3",
			Expect: "(binary OR 1 (binary AND 2 3))",
		},
	}
	check_expr(t, tt)
}
