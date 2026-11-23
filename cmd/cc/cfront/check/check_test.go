package check

import (
	"strings"
	"testing"

	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"cc260717/internal/unittest"
)

func TestDeclaration(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "declaration with no init",
			Src:    "int a;}", // the trailing } is required
			Expect: "(declaration a.0 )",
		},
		{
			Name:   "declaration with integer init",
			Src:    "int a = 0;}",
			Expect: "(declaration a.0 0)",
		},
		{
			Name:   "declaration with expr init",
			Src:    "int a = 1 + 2;}",
			Expect: "(declaration a.0 (binary ADD 1 2))",
		},
		{
			Name:   "declaration with later init",
			Src:    "int a; a = 0;}",
			Expect: "(declaration a.0 ) (assign a.0 0)",
		},
		{
			Name:   "multiple declarations",
			Src:    "int a; int b;}",
			Expect: "(declaration a.0 ) (declaration b.1 )",
		},
	}

	check(t, tt)
}

func check(t *testing.T, tt []unittest.NamedTestPair) {
	unittest.TestPair(t, make_program, tt)
}
func make_program(src string) string {
	l := lex.New("", src)
	p := parse.New(l)
	c := New()

	//
	program := ast.Program{
		FuncDef: ast.FunctionDefinition{
			Name: "main",
			Block: ast.Block{
				Blocks: p.ParseBlockItems(),
			},
		},
	}
	program = c.Check(program)

	//
	items := []string{}
	for _, item := range program.FuncDef.Block.Blocks {
		items = append(items, item.String())
	}

	return strings.Join(items, " ")
}
