package tacky

import (
	"testing"

	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/check"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"cc260717/internal/unittest"
)

func TestReturn(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name: "return with value",
			Src:  "return 3;",
			Expect: []string{
				"return 3;",
			},
		},
		{
			Name: "void return",
			Src:  "return;",
			Expect: []string{
				"return ;",
			},
		},
	}

	Check(t, tt)
}
func TestUnary(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name:   "negate",
			Src:    "-2;",
			Expect: []string{"unary - 2 tmp.0;"},
		},
		{
			Name:   "complement",
			Src:    "~2;",
			Expect: []string{"unary ~ 2 tmp.0;"},
		},
		{
			Name:   "logical NOT",
			Src:    "!2;",
			Expect: []string{"unary ! 2 tmp.0;"},
		},
		{
			Name: "nested",
			Src:  "-(-2);",
			Expect: []string{
				"unary - 2 tmp.0;",
				"unary - tmp.0 tmp.1;",
			},
		},
	}

	Check(t, tt)
}
func TestDecl(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name:   "declaration",
			Src:    "int a;",
			Expect: []string{},
		},
		{
			Name:   "declaration with init",
			Src:    "int a = 0;",
			Expect: []string{"copy 0 a.0;"},
		},
	}

	Check(t, tt)
}

func TestBinary(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name:   "plus",
			Src:    "1 + 2;",
			Expect: []string{"binary + 1 2 tmp.0;"},
		},
		{
			Name:   "sub",
			Src:    "1 - 2;",
			Expect: []string{"binary - 1 2 tmp.0;"},
		},
		{
			Name:   "mul",
			Src:    "1 * 2;",
			Expect: []string{"binary * 1 2 tmp.0;"},
		},
		{
			Name:   "div",
			Src:    "1 / 2;",
			Expect: []string{"binary / 1 2 tmp.0;"},
		},
		{
			Name:   "remainder",
			Src:    "1 % 2;",
			Expect: []string{"binary % 1 2 tmp.0;"},
		},
		{
			Name:   "bitwise AND",
			Src:    "1 & 2;",
			Expect: []string{"binary & 1 2 tmp.0;"},
		},
		{
			Name:   "bitwise OR",
			Src:    "1 | 2;",
			Expect: []string{"binary | 1 2 tmp.0;"},
		},
		{
			Name:   "left shift",
			Src:    "1 << 2;",
			Expect: []string{"binary << 1 2 tmp.0;"},
		},
		{
			Name:   "right shift",
			Src:    "1 >> 2;",
			Expect: []string{"binary >> 1 2 tmp.0;"},
		},
		{
			Name:   "equal",
			Src:    "1 == 2;",
			Expect: []string{"binary == 1 2 tmp.0;"},
		},
		{
			Name:   "not equal",
			Src:    "1 != 2;",
			Expect: []string{"binary != 1 2 tmp.0;"},
		},
		{
			Name:   "less than",
			Src:    "1 < 2;",
			Expect: []string{"binary < 1 2 tmp.0;"},
		},
		{
			Name:   "greater than",
			Src:    "1 > 2;",
			Expect: []string{"binary > 1 2 tmp.0;"},
		},
		{
			Name:   "less than or equal",
			Src:    "1 <= 2;",
			Expect: []string{"binary <= 1 2 tmp.0;"},
		},
		{
			Name:   "greater than or equal",
			Src:    "1 >= 2;",
			Expect: []string{"binary >= 1 2 tmp.0;"},
		},
		{
			Name: "nested",
			Src:  "1 + 2 + 3;",
			Expect: []string{
				"binary + 1 2 tmp.0;",
				"binary + tmp.0 3 tmp.1;",
			},
		},
	}

	Check(t, tt)
}

func TestLogicalBinary(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name: "logical and",
			Src:  "1 && 0;",
			Expect: []string{
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
			Name: "logical or",
			Src:  "1 || 0;",
			Expect: []string{
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

	Check(t, tt)
}

func Check(t *testing.T, tt []unittest.NamedTestObj) {
	unittest.TestObj(t, make_tacky, tt)
}
func make_tacky(src string) []string {
	generator := New()
	p := parse.New(lex.New("", src))
	checker := check.New()

	//
	program := ast.Program{
		FuncDef: ast.FunctionDefinition{
			Name: "main",
			Block: ast.Block{
				Blocks: []ast.BlockItem{*p.ParseBlockItem()},
			},
		},
	}
	program = checker.Check(program)

	//
	result := []string{}
	for _, ins := range generator.Generate(program).FuncDef.Ins {
		result = append(result, ins.String())
	}

	return result
}
