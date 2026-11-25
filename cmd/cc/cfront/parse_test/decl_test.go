package parsetest

import (
	"testing"

	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"cc260717/internal/unittest"
)

func TestDecl(t *testing.T) {
	tt := []unittest.NamedTestPair{
		{
			Name:   "declare int",
			Src:    "int a = 0;",
			Expect: "(declaration a 0)",
		},
		{
			Name:   "declare int no init",
			Src:    "int a;",
			Expect: "(declaration a )",
		},
		{
			Name:   "declare int expr",
			Src:    "int a = 1 + 1;",
			Expect: "(declaration a (binary ADD 1 1))",
		},
	}

	check_decl(t, tt)
}

func check_decl(t *testing.T, tt []unittest.NamedTestPair) {
	unittest.TestPair(t, make_ast, tt)
}
func make_ast(src string) string {
	l := lex.New("", src)
	toks, _ := l.LexAll()

	return parse.New(toks).ParseDecl().String()
}
