package parsetest

import (
	"testing"

	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"cc260717/internal/unittest"
)

func check_expr(t *testing.T, tt []unittest.NamedTestPair) {
	unittest.TestPair(t, make_expr, tt)
}
func make_expr(src string) string {
	l := lex.New("", src)
	toks, _ := l.LexAll()

	return parse.New(toks).ParseExpr(parse.LOWEST).String()
}
