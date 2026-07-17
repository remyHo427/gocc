package tackyparse

import "testing"

func TestIdent(t *testing.T) {
	l := New(`a b tmp.0`)
	seq := []Toktype{IDENT, IDENT, IDENT}
	tokseq(*l, seq, t)
}
func TestValues(t *testing.T) {
	l := New(`0 10 100`)
	seq := []Toktype{INTEGER, INTEGER, INTEGER}
	tokseq(*l, seq, t)
}
func TestOperators(t *testing.T) {
	l := New(`
		- ~ + * / % ^ ; ! != || | && & == = 
		<< <= < >> >= >
	`)
	seq := []Toktype{
		SUB, TILDE, ADD, MUL, DIV, MOD, BXOR,
		SCOLON, NOT, NEQ, OR, BOR, AND, BAND, EQ,
		ASSIGN, LSHIFT, LEQ, LT, RSHIFT, GEQ, GT,
	}
	tokseq(*l, seq, t)
}
func TestKeywords(t *testing.T) {
	l := New(`
		unary binary return copy jump jump_if_zero 
		jump_if_not_zero label
	`)
	seq := []Toktype{
		UNARY, BINARY, RETURN, COPY, JUMP,
		JUMP_IF_ZERO, JUMP_IF_NOT_ZERO,
		LABEL,
	}
	tokseq(*l, seq, t)
}

func tokseq(l Lexer, seq []Toktype, t *testing.T) {
	t.Helper()

	i := 0
	for _, ttype := range seq {
		if tok := l.Lex(); tok.Type != ttype {
			t.Errorf("expected %d, got %d at seq[%d]",
				ttype, tok.Type, i)
		}
		i++
	}

	if i != len(seq) {
		t.Errorf("input does not match expected amount of tokens")
	}
	if l.Lex().Type != EOF {
		t.Errorf("input not exhausted")
	}
}
