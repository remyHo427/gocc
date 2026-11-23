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
	return parse.New(lex.New("", src)).ParseExpr(parse.LOWEST).String()
}
