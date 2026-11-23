package tackyparse

import (
	"testing"

	"cc260717/internal/unittest"
)

func TestIdent(t *testing.T) {
	l := New(`a b tmp.0`)
	seq := []Toktype{IDENT, IDENT, IDENT}
	tokseq(t, *l, seq)
}
func TestValues(t *testing.T) {
	l := New(`0 10 100`)
	seq := []Toktype{INTEGER, INTEGER, INTEGER}
	tokseq(t, *l, seq)
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
	tokseq(t, *l, seq)
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
	tokseq(t, *l, seq)
}

func tokseq(t *testing.T, l Lexer, seq []Toktype) {
	unittest.Test(t, next(l), seq)
}
func next(l Lexer) func() Toktype {
	return func() Toktype { return l.Lex().Type }
}
