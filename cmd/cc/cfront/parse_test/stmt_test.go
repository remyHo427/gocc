package parsetest

import (
	"testing"

	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"cc260717/internal/unittest"
)

func TestReturnStmt(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "void return",
			Src:    "return;",
			Expect: "(return )",
		},
		{
			Name:   "return with expr",
			Src:    "return 0;",
			Expect: "(return 0)",
		},
	}
	check_stmt(t, tt)
}
func TestExprStmt(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "expr stmt",
			Src:    "1;",
			Expect: "1",
		},
		{
			Name:   "expr stmt with no expression (null stmt)",
			Src:    ";",
			Expect: "(null)",
		},
	}
	check_stmt(t, tt)
}
func TestIfStmt(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "only if",
			Src:    "if (1) 1;",
			Expect: "(if 1 1 )",
		},
		{
			Name:   "if and else",
			Src:    "if (1) 1; else 0;",
			Expect: "(if 1 1 0)",
		},
		{
			Name:   "chained if",
			Src:    "if (0) 0; else if (1) 1; else 2;",
			Expect: "(if 0 0 (if 1 1 2))",
		},
	}
	check_stmt(t, tt)
}

func check_stmt(t *testing.T, tt []unittest.NamedTestPair) {
	unittest.TestPair(t, parse_stmt, tt)
}
func parse_stmt(src string) string {
	return parse.New(lex.New("", src)).ParseStmt().String()
}
