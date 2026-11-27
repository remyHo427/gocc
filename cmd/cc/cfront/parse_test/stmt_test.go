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

func TestCompoundStmt(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "compound stmt",
			Src:    "{ int a; 1; }",
			Expect: "(compound (block (declaration a ) 1))",
		},
	}
	check_stmt(t, tt)
}

func TestWhileStmt(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "with simple condition",
			Src:    "while (1) 2;",
			Expect: "(while 1 2)",
		},
		{
			Name:   "with complex condition",
			Src:    "while (a != 10) a = a + 1;",
			Expect: "(while (binary NEQ a 10) (assign a (binary ADD a 1)))",
		},
		{
			Name:   "with null stmt",
			Src:    "while (1) ;",
			Expect: "(while 1 (null))",
		},
		{
			Name:   "with block body",
			Src:    "while (1) { return 1; }",
			Expect: "(while 1 (compound (block (return 1))))",
		},
	}
	check_stmt(t, tt)
}
func TestDoWhileStmt(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "with simple condition",
			Src:    "do 2; while (1);",
			Expect: "(do-while 2 1)",
		},
		{
			Name:   "with complex condition",
			Src:    "do a = a + 1; while (a != 10);",
			Expect: "(do-while (assign a (binary ADD a 1)) (binary NEQ a 10))",
		},
		{
			Name:   "with null stmt",
			Src:    "do ; while (1);",
			Expect: "(do-while (null) 1)",
		},
		{
			Name:   "with block body",
			Src:    "do { return 2; } while (1);",
			Expect: "(do-while (compound (block (return 2))) 1)",
		},
	}
	check_stmt(t, tt)
}
func TestForStmt(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "no init, no cond, no post",
			Src:    "for (;;) ;",
			Expect: "(for (null) (null) (null) (null))",
		},
		{
			Name:   "decl init, no cond, no post",
			Src:    "for (int a;;) ;",
			Expect: "(for (declaration a ) (null) (null) (null))",
		},
		{
			Name:   "expr init, no cond, no post",
			Src:    "for (1;;) ;",
			Expect: "(for 1 (null) (null) (null))",
		},
		{
			Name:   "no init, expr cond, no post",
			Src:    "for (;2;) ;",
			Expect: "(for (null) 2 (null) (null))",
		},
		{
			Name:   "no init, no cond, expr post",
			Src:    "for (;;3) ;",
			Expect: "(for (null) (null) 3 (null))",
		},
		{
			Name:   "init, cond, post",
			Src:    "for (1;2;3) ;",
			Expect: "(for 1 2 3 (null))",
		},
	}
	check_stmt(t, tt)
}

func TestMiscStmts(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "break",
			Src:    "break;",
			Expect: "(break)",
		},
		{
			Name:   "continue",
			Src:    "continue;",
			Expect: "(continue)",
		},
	}
	check_stmt(t, tt)
}

func check_stmt(t *testing.T, tt []unittest.NamedTestPair) {
	unittest.TestPair(t, parse_stmt, tt)
}
func parse_stmt(src string) string {
	l := lex.New("", src)
	toks, _ := l.LexAll()

	return parse.New(toks).ParseStmt().String()
}
