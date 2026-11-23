package tackyparse

import (
	"cc260717/internal/unittest"
	"testing"
)

func TestReturn(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "null",
			Src:    "return ;",
			Expect: "return ;",
		},
		{
			Name:   "integer",
			Src:    "unary ~ 1 tmp.0;",
			Expect: "unary ~ 1 tmp.0;",
		},
		{
			Name:   "identifier",
			Src:    "return tmp.0;",
			Expect: "return tmp.0;",
		},
	}
	check(t, tt)
}
func TestUnary(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "negate",
			Src:    "unary - 1 tmp.0;",
			Expect: "unary - 1 tmp.0;",
		},
		{
			Name:   "complement",
			Src:    "unary ~ 1 tmp.0;",
			Expect: "unary ~ 1 tmp.0;",
		},
		{
			Name:   "logical not",
			Src:    "unary ! 1 tmp.0;",
			Expect: "unary ! 1 tmp.0;",
		},
	}
	check(t, tt)
}
func TestBinary(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "add",
			Src:    "binary + 1 2 tmp.0;",
			Expect: "binary + 1 2 tmp.0;",
		},
		{
			Name:   "sub",
			Src:    "binary - 1 2 tmp.0;",
			Expect: "binary - 1 2 tmp.0;",
		},
		{
			Name:   "mul",
			Src:    "binary * 1 2 tmp.0;",
			Expect: "binary * 1 2 tmp.0;",
		},
		{
			Name:   "div",
			Src:    "binary / 1 2 tmp.0;",
			Expect: "binary / 1 2 tmp.0;",
		},
		{
			Name:   "rem",
			Src:    "binary % 1 2 tmp.0;",
			Expect: "binary % 1 2 tmp.0;",
		},
		{
			Name:   "bitwise OR",
			Src:    "binary | 1 2 tmp.0;",
			Expect: "binary | 1 2 tmp.0;",
		},
		{
			Name:   "bitwise AND",
			Src:    "binary & 1 2 tmp.0;",
			Expect: "binary & 1 2 tmp.0;",
		},
		{
			Name:   "bitwise XOR",
			Src:    "binary ^ 1 2 tmp.0;",
			Expect: "binary ^ 1 2 tmp.0;",
		},
		{
			Name:   "left shift",
			Src:    "binary << 1 2 tmp.0;",
			Expect: "binary << 1 2 tmp.0;",
		},
		{
			Name:   "right shift",
			Src:    "binary >> 1 2 tmp.0;",
			Expect: "binary >> 1 2 tmp.0;",
		},
		{
			Name:   "less than",
			Src:    "binary < 1 2 tmp.0;",
			Expect: "binary < 1 2 tmp.0;",
		},
		{
			Name:   "greater than",
			Src:    "binary > 1 2 tmp.0;",
			Expect: "binary > 1 2 tmp.0;",
		},
		{
			Name:   "less than or equal",
			Src:    "binary <= 1 2 tmp.0;",
			Expect: "binary <= 1 2 tmp.0;",
		},
		{
			Name:   "greater than or equal",
			Src:    "binary >= 1 2 tmp.0;",
			Expect: "binary >= 1 2 tmp.0;",
		},
		{
			Name:   "not equal",
			Src:    "binary != 1 2 tmp.0;",
			Expect: "binary != 1 2 tmp.0;",
		},
		{
			Name:   "equal",
			Src:    "binary == 1 2 tmp.0;",
			Expect: "binary == 1 2 tmp.0;",
		},
		{
			Name:   "logical and",
			Src:    "binary && 1 2 tmp.0;",
			Expect: "binary && 1 2 tmp.0;",
		},
		{
			Name:   "logical or",
			Src:    "binary || 1 2 tmp.0;",
			Expect: "binary || 1 2 tmp.0;",
		},
	}
	check(t, tt)
}
func TestMisc(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "label",
			Src:    "label lb.0;",
			Expect: "label lb.0;",
		},
		{
			Name:   "jump",
			Src:    "jump lb.0;",
			Expect: "jump lb.0;",
		},
		{
			Name:   "copy",
			Src:    "copy 0 tmp.0;",
			Expect: "copy 0 tmp.0;",
		},
		{
			Name:   "jump_if_zero",
			Src:    "jump_if_zero 0 lb.0;",
			Expect: "jump_if_zero 0 lb.0;",
		},
		{
			Name:   "jump_if_not_zero",
			Src:    "jump_if_not_zero 0 lb.0;",
			Expect: "jump_if_not_zero 0 lb.0;",
		},
	}

	check(t, tt)
}

func TestParseAll(t *testing.T) {
	t.Helper()

	src := "unary - 2 tmp.0; return tmp.0;"
	p := NewParser(New(src))
	tacky := p.ParseAll()

	if len(tacky) != 2 {
		t.Errorf("expect 2 instructions, got %d", len(tacky))
	}
	if actual := tacky[0].String(); actual != "unary - 2 tmp.0;" {
		t.Errorf("expect \"unary - 2 tmp.0;\", got \"%s\"", actual)
	}
	if actual := tacky[1].String(); actual != "return tmp.0;" {
		t.Errorf("expect \"return tmp.0\", got \"%s\"", actual)
	}
}

func check(t *testing.T, tt []unittest.NamedTestPair) {
	unittest.TestPair(t, parse_tacky, tt)
}
func parse_tacky(src string) string {
	return NewParser(New(src)).Parse().String()
}
