package parsetest

import (
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"testing"
)

func TestDecl(t *testing.T) {
	tt := []TestPair{
		{"declare int", "int a = 0;", "(declaration a 0)"},
		{"declare int no init", "int a;", "(declaration a )"},
		{"declare int expr", "int a = 1 + 1;", "(declaration a (binary ADD 1 1))"},
	}
	check_decl(t, tt)
}

func check_decl(t *testing.T, tt []TestPair) {
	t.Helper()

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			l := lex.New(test.src)
			p := parse.New(l)

			ast := p.ParseDecl()
			actual := ast.String()
			if actual != test.expect {
				t.Errorf("expected \"%s\", got \"%s\"", test.expect, actual)
			}
		})
	}
}
