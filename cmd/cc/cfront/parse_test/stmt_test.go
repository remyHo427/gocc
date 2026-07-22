package parsetest

import (
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"testing"
)

func TestReturnStmt(t *testing.T) {
	tt := []TestPair{
		{"void return", "return;", "(return )"},
		{"return with expr", "return 0;", "(return 0)"},
	}
	check_stmt(t, tt)
}
func TestExprStmt(t *testing.T) {
	tt := []TestPair{
		{"expr stmt", "1;", "1"},
		{"expr stmt with no expression (null stmt)", ";", "(null)"},
	}
	check_stmt(t, tt)
}
func TestIfStmt(t *testing.T) {
	tt := []TestPair{
		{"only if", "if (1) 1;", "(if 1 1 )"},
		{"if and else", "if (1) 1; else 0;", "(if 1 1 0)"},
		{"chained if", "if (0) 0; else if (1) 1; else 2;",
			"(if 0 0 (if 1 1 2))"},
	}
	check_stmt(t, tt)
}

func check_stmt(t *testing.T, tt []TestPair) {
	t.Helper()

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			l := lex.New("", test.src)
			p := parse.New(l)

			ast := p.ParseStmt()
			actual := ast.String()
			if actual != test.expect {
				t.Errorf("expected \"%s\", got \"%s\"", test.expect, actual)
			}
		})
	}
}
