package parsetest

import (
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"cc260717/internal/unittest"
	"testing"
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
	MakeTester(t, tt)()
}

func MakeTester(t *testing.T, tt []unittest.NamedTestPair) func() {
	return unittest.MakePairTester(t, make_ast, tt)
}
func make_ast(src string) string {
	return parse.New(lex.New("", src)).ParseDecl().String()
}
