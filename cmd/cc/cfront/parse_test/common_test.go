package parsetest

import (
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"testing"
)

type TestPair struct {
	name   string
	src    string
	expect string
}

func check_expr(t *testing.T, tt []TestPair) {
	t.Helper()

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			l := lex.New(test.src)
			p := parse.New(l)

			ast := p.ParseExpr(parse.LOWEST)
			actual := ast.String()

			if actual != test.expect {
				t.Errorf("expected \"%s\", got \"%s\"", test.expect, actual)
			}
		})
	}
}
