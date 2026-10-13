package tacky

import (
	"testing"

	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/check"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
)

type TestPair struct {
	name     string
	src      string
	expected []string
}

func TestReturn(t *testing.T) {
	tt := []TestPair{
		{
			"return with value",
			"return 3;",
			[]string{
				"return 3;",
			},
		},
		{
			"void return",
			"return;",
			[]string{
				"return ;",
			},
		},
	}

	check_tt(t, tt)
}
func TestUnary(t *testing.T) {
	tt := []TestPair{
		{
			"negate",
			"-2;",
			[]string{"unary - 2 tmp.0;"},
		},
		{
			"complement",
			"~2;",
			[]string{"unary ~ 2 tmp.0;"},
		},
		{
			"logical NOT",
			"!2;",
			[]string{"unary ! 2 tmp.0;"},
		},
		{
			"nested",
			"-(-2)",
			[]string{
				"unary - 2 tmp.0;",
				"unary - tmp.0 tmp.1;",
			},
		},
	}

	check_tt(t, tt)
}
func TestDecl(t *testing.T) {
	tt := []TestPair{
		{
			"declaration",
			"int a;",
			[]string{},
		},
		{
			"declaration with init",
			"int a = 0;",
			[]string{"copy 0 a.0;"},
		},
	}

	check_tt(t, tt)
}

func TestBinary(t *testing.T) {
	tt := []TestPair{
		{
			"plus",
			"1 + 2;",
			[]string{"binary + 1 2 tmp.0;"},
		},
		{
			"sub",
			"1 - 2;",
			[]string{"binary - 1 2 tmp.0;"},
		},
		{
			"mul",
			"1 * 2;",
			[]string{"binary * 1 2 tmp.0;"},
		},
		{
			"div",
			"1 / 2;",
			[]string{"binary / 1 2 tmp.0;"},
		},
		{
			"remainder",
			"1 % 2;",
			[]string{"binary % 1 2 tmp.0;"},
		},
		{
			"bitwise AND",
			"1 & 2;",
			[]string{"binary & 1 2 tmp.0;"},
		},
		{
			"bitwise OR",
			"1 | 2;",
			[]string{"binary | 1 2 tmp.0;"},
		},
		{
			"left shift",
			"1 << 2;",
			[]string{"binary << 1 2 tmp.0;"},
		},
		{
			"right shift",
			"1 >> 2;",
			[]string{"binary >> 1 2 tmp.0;"},
		},
		{
			"equal",
			"1 == 2;",
			[]string{"binary == 1 2 tmp.0;"},
		},
		{
			"not equal",
			"1 != 2;",
			[]string{"binary != 1 2 tmp.0;"},
		},
		{
			"less than",
			"1 < 2;",
			[]string{"binary < 1 2 tmp.0;"},
		},
		{
			"greater than",
			"1 > 2;",
			[]string{"binary > 1 2 tmp.0;"},
		},
		{
			"less than or equal",
			"1 <= 2;",
			[]string{"binary <= 1 2 tmp.0;"},
		},
		{
			"greater than or equal",
			"1 >= 2;",
			[]string{"binary >= 1 2 tmp.0;"},
		},
		{
			"nested",
			"1 + 2 + 3",
			[]string{
				"binary + 1 2 tmp.0;",
				"binary + tmp.0 3 tmp.1;",
			},
		},
	}

	check_tt(t, tt)
}

func TestLogicalBinary(t *testing.T) {
	tt := []TestPair{
		{
			"logical and",
			"1 && 0;",
			[]string{
				"jump_if_zero 1 lb.0;",
				"jump_if_zero 0 lb.0;",
				"copy 1 tmp.0;",
				"jump lb.1;",
				"label lb.0;",
				"copy 0 tmp.0;",
				"label lb.1;",
			},
		},
		{
			"logical or",
			"1 || 0",
			[]string{
				"jump_if_not_zero 1 lb.0;",
				"jump_if_not_zero 0 lb.0;",
				"copy 0 tmp.0;",
				"jump lb.1;",
				"label lb.0;",
				"copy 1 tmp.0;",
				"label lb.1;",
			},
		},
	}

	check_tt(t, tt)
}

func check_tt(t *testing.T, pairs []TestPair) {
	for _, pair := range pairs {
		t.Run(pair.name, func(t *testing.T) {
			t.Helper()

			tacky := GenerateTacky(pair.src)
			if len(pair.expected) == 0 && len(tacky.FuncDef.Ins) == 0 {
				// expected the tacky output to be empty
				return
			}

			if len(tacky.FuncDef.Ins) == 0 {
				t.Errorf("no tacky instructions emitted")
			} else if len(tacky.FuncDef.Ins) != len(pair.expected) {
				t.Errorf("actual instructions count does not match expected")
			}

			for i, expect := range pair.expected {
				if actual := tacky.FuncDef.Ins[i].String(); actual != expect {
					t.Errorf("expect \"%s\" got \"%s\"", expect, actual)
				}
			}
		})
	}
}

func GenerateTacky(src string) Program {
	generator := New()
	p := parse.New(lex.New("", src))
	checker := check.New()

	program := ast.Program{
		FuncDef: ast.FunctionDefinition{
			Name: "main",
			Block: ast.Block{
				Blocks: []ast.BlockItem{*p.ParseBlockItem()},
			},
		},
	}
	program = checker.Check(program)

	return generator.Generate(program)
}
