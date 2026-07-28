package check

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"strings"
	"testing"
)

type TestPair struct {
	name   string
	src    string
	expect string
}

func TestDeclaration(t *testing.T) {
	tt := []TestPair{
		{
			"declaration with no init",
			"int a;}", // the trailing } is required
			"(declaration a.0 )",
		},
		{
			"declaration with integer init",
			"int a = 0;}",
			"(declaration a.0 0)",
		},
		{
			"declaration with expr init",
			"int a = 1 + 2;}",
			"(declaration a.0 (binary ADD 1 2))",
		},
		{
			"declaration with later init",
			"int a; a = 0;}",
			"(declaration a.0 ) (assign a.0 0)",
		},
		{
			"multiple declarations",
			"int a; int b;}",
			"(declaration a.0 ) (declaration b.1 )",
		},
	}

	checkProgram(t, tt)
}

func checkProgram(t *testing.T, tt []TestPair) {
	t.Helper()

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			l := lex.New("", test.src)
			p := parse.New(l)
			c := New()

			program := ast.Program{
				FuncDef: ast.FunctionDefinition{
					Name: "main",
					Block: ast.Block{
						Blocks: p.ParseBlockItems(),
					},
				},
			}
			program = c.Check(program)
			items := []string{}
			for _, item := range program.FuncDef.Block.Blocks {
				items = append(items, item.String())
			}

			actual := strings.Join(items, " ")
			if actual != test.expect {
				t.Errorf("expected \"%s\", got \"%s\"", test.expect, actual)
			}
		})
	}
}
