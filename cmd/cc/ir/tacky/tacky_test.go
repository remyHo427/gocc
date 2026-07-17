package tacky

import (
	"testing"

	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
)

type TestPair struct {
	name     string
	src      string
	expected []string
}

func TestReturn(t *testing.T) {
	check(t, "return with value", "return 3;", "return 3;")
	check(t, "void return", "return;", "return ;")
}
func TestUnary(t *testing.T) {
	check(t, "negate", "-2;", "unary - 2 tmp.0;")
	check(t, "complement", "~2;", "unary ~ 2 tmp.0;")
	check(t, "logical not", "!2;", "unary ! 2 tmp.0;")

	check(t, "nested", "-(-2);",
		"unary - 2 tmp.0;",
		"unary - tmp.0 tmp.1;",
	)
}

func TestBinary(t *testing.T) {
	check(t, "plus", "1 + 2;", "binary + 1 2 tmp.0;")
	check(t, "sub", "1 - 2;", "binary - 1 2 tmp.0;")
	check(t, "mul", "1 * 2;", "binary * 1 2 tmp.0;")
	check(t, "div", "1 / 2;", "binary / 1 2 tmp.0;")
	check(t, "remainder", "1 % 2;", "binary % 1 2 tmp.0;")
	check(t, "bitwise AND", "1 & 2;", "binary & 1 2 tmp.0;")
	check(t, "bitwise OR", "1 | 2;", "binary | 1 2 tmp.0;")
	check(t, "bitwise XOR", "1 ^ 2;", "binary ^ 1 2 tmp.0;")
	check(t, "bitwise left shift", "1 << 2;", "binary << 1 2 tmp.0;")
	check(t, "bitwise right shift", "1 >> 2;", "binary >> 1 2 tmp.0;")
	check(t, "equal", "1 == 2;", "binary == 1 2 tmp.0;")
	check(t, "not equal", "1 != 2;", "binary != 1 2 tmp.0;")
	check(t, "less than", "1 < 2;", "binary < 1 2 tmp.0;")
	check(t, "greater than", "1 > 2;", "binary > 1 2 tmp.0;")
	check(t, "less than or equal", "1 <= 2;", "binary <= 1 2 tmp.0;")
	check(t, "greater than or equal", "1 >= 2;", "binary >= 1 2 tmp.0;")

	check(t, "nested", "1 + 2 + 3;",
		"binary + 1 2 tmp.0;",
		"binary + tmp.0 3 tmp.1;")
}
func TestLogicalBinary(t *testing.T) {
	check(t, "and", "1 && 0;",
		"jump_if_zero 1 lb.0;",
		"jump_if_zero 0 lb.0;",
		"copy 1 tmp.0;",
		"jump lb.1;",
		"label lb.0;",
		"copy 0 tmp.0;",
		"label lb.1;",
	)
	check(t, "or", "1 || 0;",
		"jump_if_not_zero 1 lb.0;",
		"jump_if_not_zero 0 lb.0;",
		"copy 0 tmp.0;",
		"jump lb.1;",
		"label lb.0;",
		"copy 1 tmp.0;",
		"label lb.1;",
	)
}

func check(t *testing.T, name string, src string, expected ...string) {
	t.Run(name, func(t *testing.T) {
		t.Helper()

		tacky := GenerateTacky(src)
		if len(tacky.FuncDef.Ins) == 0 {
			t.Errorf("no tacky instructions emitted")
		} else if len(tacky.FuncDef.Ins) != len(expected) {
			t.Errorf("actual instructions count does not match expected")
		}

		for i, expect := range expected {
			if actual := tacky.FuncDef.Ins[i].String(); actual != expect {
				t.Errorf("expect \"%s\" got \"%s\"", expect, actual)
			}
		}
	})
}

func GenerateTacky(src string) Program {
	generator := New()
	p := parse.New(lex.New(src))

	return generator.Generate(ast.Program{
		FuncDef: ast.FunctionDefinition{
			Name: "main",
			Body: p.ParseStmt(),
		},
	})
}
