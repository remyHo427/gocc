package parsetest

import "testing"

func TestAllBinaryOperators(t *testing.T) {
	tt := []TestPair{
		{"arithmetic addition", "1 + 1", "(binary ADD 1 1)"},
		{"arithmetic subtraction", "1 - 1", "(binary SUB 1 1)"},
		{"arithmetic multiplication", "1 * 1", "(binary MUL 1 1)"},
		{"arithmetic division", "1 / 1", "(binary DIV 1 1)"},
		{"arithmetic remainder", "1 % 1", "(binary MOD 1 1)"},
		{"bitwise AND", "1 & 1", "(binary BAND 1 1)"},
		{"bitwise OR", "1 | 1", "(binary BOR 1 1)"},
		{"bitwise XOR", "1 ^ 1", "(binary BXOR 1 1)"},
		{"left shift", "1 << 1", "(binary LSHIFT 1 1)"},
		{"right shift", "1 >> 1", "(binary RSHIFT 1 1)"},
		{"logical or", "1 || 1", "(binary OR 1 1)"},
		{"logical and", "1 && 1", "(binary AND 1 1)"},
		{"equal", "1 == 1", "(binary EQ 1 1)"},
		{"not equal", "1 != 1", "(binary NEQ 1 1)"},
		{"less than", "1 < 1", "(binary LT 1 1)"},
		{"greater than", "1 > 1", "(binary GT 1 1)"},
		{"less than or equal", "1 <= 1", "(binary LEQ 1 1)"},
		{"greater than or equal", "1 >= 1", "(binary GEQ 1 1)"},
		{"ternary", "1 ? 1 : 0", "(ternary 1 1 0)"},
	}
	check_expr(t, tt)
}

func TestBinaryOperatorsAssociativity(t *testing.T) {
	tt := []TestPair{
		// left-to-right
		{"arithmetic add", "1 + 2 + 3", "(binary ADD (binary ADD 1 2) 3)"},
		{"arithmetic sub", "1 - 2 - 3", "(binary SUB (binary SUB 1 2) 3)"},
		{"arithmetic mul", "1 * 2 * 3", "(binary MUL (binary MUL 1 2) 3)"},
		{"arithmetic div", "1 / 2 / 3", "(binary DIV (binary DIV 1 2) 3)"},
		{"arithmetic remainder", "1 % 2 % 3", "(binary MOD (binary MOD 1 2) 3)"},
		{"bitwise AND", "1 & 2 & 3", "(binary BAND (binary BAND 1 2) 3)"},
		{"bitwise OR", "1 | 2 | 3", "(binary BOR (binary BOR 1 2) 3)"},
		{"bitwise XOR", "1 ^ 2 ^ 3", "(binary BXOR (binary BXOR 1 2) 3)"},
		{"left shift", "1 << 2 << 3", "(binary LSHIFT (binary LSHIFT 1 2) 3)"},
		{"right shift", "1 >> 2 >> 3", "(binary RSHIFT (binary RSHIFT 1 2) 3)"},
		{"logical or", "1 || 2 || 3", "(binary OR (binary OR 1 2) 3)"},
		{"logical and", "1 && 2 && 3", "(binary AND (binary AND 1 2) 3)"},
		{"equal", "1 == 2 == 3", "(binary EQ (binary EQ 1 2) 3)"},
		{"not equal", "1 != 2 != 3", "(binary NEQ (binary NEQ 1 2) 3)"},
		{"less than", "1 < 2 < 3", "(binary LT (binary LT 1 2) 3)"},
		{"greater than", "1 > 2 > 3", "(binary GT (binary GT 1 2) 3)"},
		{"less than or equal", "1 <= 2 <= 3", "(binary LEQ (binary LEQ 1 2) 3)"},
		{"greater than or equal", "1 >= 2 >= 3", "(binary GEQ (binary GEQ 1 2) 3)"},

		// right-to-left
		{"assign", "1 = 2 = 3", "(assign 1 (assign 2 3))"},

		// 1 ? (2 ? 3 : 4) : (5 ? 6 : 7)
		{"ternary", "1 ? 2 ? 3 : 4 : 5 ? 6 : 7",
			"(ternary 1 (ternary 2 3 4) (ternary 5 6 7))"},
	}
	check_expr(t, tt)
}

// test mixing different operators with the same precedence level
func TestEqualPrecedenceAssociativity(t *testing.T) {
	tt := []TestPair{
		{"arithmetic mul, div, reminder", "1 * 2 / 3 % 4",
			"(binary MOD (binary DIV (binary MUL 1 2) 3) 4)"},
		{"arithmetic add and sub", "1 + 2 - 3",
			"(binary SUB (binary ADD 1 2) 3)"},
		{"left shift and right shift", "1 << 2 >> 3",
			"(binary RSHIFT (binary LSHIFT 1 2) 3)"},
		{"less than, greater than, less than or equal, greater than or equal",
			"1 < 2 > 3 <= 4 >= 5",
			"(binary GEQ (binary LEQ (binary GT (binary LT 1 2) 3) 4) 5)"},
		{"equal and not equal", "1 == 2 != 3",
			"(binary NEQ (binary EQ 1 2) 3)"},
		// bitwise AND, bitwise XOR, bitwise OR, logical AND and logical OR
		// each have their own precedence level, therefore not tested here
	}
	check_expr(t, tt)
}

// test mixing operators with different precedence levels
func TestPrecedenceLevels(t *testing.T) {
	tt := []TestPair{
		// + with * / %
		{"arithmetic add with arithmetic mul",
			"1 + 2 * 3", "(binary ADD 1 (binary MUL 2 3))"},
		{"arithmetic add with arithmetic div",
			"1 + 2 / 3", "(binary ADD 1 (binary DIV 2 3))"},
		{"arithmetic add with arithmetic remainder",
			"1 + 2 % 3", "(binary ADD 1 (binary MOD 2 3))"},
		// - with * / %
		{"arithmetic sub with arithmetic mul",
			"1 - 2 * 3", "(binary SUB 1 (binary MUL 2 3))"},
		{"arithmetic sub with arithmetic div",
			"1 - 2 / 3", "(binary SUB 1 (binary DIV 2 3))"},
		{"arithmetic sub with arithmetic remainder",
			"1 - 2 % 3", "(binary SUB 1 (binary MOD 2 3))"},
		// << with + -
		{"left shift with arithmetic add",
			"1 << 2 + 3", "(binary LSHIFT 1 (binary ADD 2 3))"},
		{"left shift with arithmetic sub",
			"1 << 2 - 3", "(binary LSHIFT 1 (binary SUB 2 3))"},
		// >> with + -
		{"right shift with arithmetic add",
			"1 >> 2 + 3", "(binary RSHIFT 1 (binary ADD 2 3))"},
		{"right shift with arithmetic sub",
			"1 >> 2 + 3", "(binary RSHIFT 1 (binary ADD 2 3))"},
		// < with << >>
		{"less than with left shift",
			"1 < 2 << 3", "(binary LT 1 (binary LSHIFT 2 3))"},
		{"less than with right shift",
			"1 < 2 >> 3", "(binary LT 1 (binary RSHIFT 2 3))"},
		// > with << >>
		{"greater than with left shift",
			"1 > 2 << 3", "(binary GT 1 (binary LSHIFT 2 3))"},
		{"greater than with right shift",
			"1 > 2 >> 3", "(binary GT 1 (binary RSHIFT 2 3))"},
		// <= with << >>
		{"less than or equal with left shift",
			"1 <= 2 << 3", "(binary LEQ 1 (binary LSHIFT 2 3))"},
		{"less than or equal with right shift",
			"1 <= 2 >> 3", "(binary LEQ 1 (binary RSHIFT 2 3))"},
		// >= with << >>
		{"greater than with left shift",
			"1 >= 2 << 3", "(binary GEQ 1 (binary LSHIFT 2 3))"},
		{"greater than with right shift",
			"1 >= 2 >> 3", "(binary GEQ 1 (binary RSHIFT 2 3))"},
		// == with < > <= >=
		{"equal with less than", "1 == 2 < 3",
			"(binary EQ 1 (binary LT 2 3))"},
		{"equal with greater than", "1 == 2 > 3",
			"(binary EQ 1 (binary GT 2 3))"},
		{"equal with less than or equal", "1 == 2 <= 3",
			"(binary EQ 1 (binary LEQ 2 3))"},
		{"equal with greater than or equal", "1 == 2 >= 3",
			"(binary EQ 1 (binary GEQ 2 3))"},
		// != with < > <= >=
		{"not equal with less than", "1 != 2 < 3",
			"(binary NEQ 1 (binary LT 2 3))"},
		{"not equal with greater than", "1 != 2 > 3",
			"(binary NEQ 1 (binary GT 2 3))"},
		{"not equal with less than or equal", "1 != 2 <= 3",
			"(binary NEQ 1 (binary LEQ 2 3))"},
		{"not equal with greater than or equal", "1 != 2 >= 3",
			"(binary NEQ 1 (binary GEQ 2 3))"},
		// & with == !=
		{"bitwise AND with equal", "1 & 2 == 3",
			"(binary BAND 1 (binary EQ 2 3))"},
		{"bitwise AND with not equal", "1 & 2 != 3",
			"(binary BAND 1 (binary NEQ 2 3))"},
		// ^ with &
		{"bitwise XOR with bitwise AND", "1 ^ 2 & 3",
			"(binary BXOR 1 (binary BAND 2 3))"},
		// | with ^
		{"bitwise OR with bitwise XOR", "1 | 2 ^ 3",
			"(binary BOR 1 (binary BXOR 2 3))"},
		// && with |
		{"bitwise AND with bitwise OR ", "1 && 2 | 3",
			"(binary AND 1 (binary BOR 2 3))"},
		// || with &&
		{"logical OR with logical AND", "1 || 2 && 3",
			"(binary OR 1 (binary AND 2 3))"},
	}
	check_expr(t, tt)
}
