package lex

import (
	"cc260717/internal/unittest"
	"testing"
)

func TestEmptyString(t *testing.T) {
	MakeTest(t, *New("", ""), []Toktype{})()
}
func TestSpaceOnly(t *testing.T) {
	l := New("", "\f\n\r\t\v ")
	seq := []Toktype{}
	MakeTest(t, *l, seq)()
}

func TestKeywords(t *testing.T) {
	l := New("", `auto break case char const continue default do
	double else enum extern float for goto if int long register
	return short signed sizeof static struct switch typedef
	union unsigned void volatile while`)
	seq := []Toktype{
		AUTO, BREAK, CASE, CHAR, CONST, CONTINUE, DEFAULT, DO,
		DOUBLE, ELSE, ENUM, EXTERN, FLOAT, FOR, GOTO, IF, INT,
		LONG, REGISTER, RETURN, SHORT, SIGNED, SIZEOF, STATIC,
		STRUCT, SWITCH, TYPEDEF, UNION, UNSIGNED, VOID, VOLATILE,
		WHILE,
	}
	MakeTest(t, *l, seq)()
}

func TestIdentifiers(t *testing.T) {
	l := New("", `a Toktype a0 a00 __test__`)
	seq := []Toktype{
		IDENT, IDENT, IDENT, IDENT, IDENT,
	}

	MakeTest(t, *l, seq)()
}
func TestInteger(t *testing.T) {
	l := New("", `0 1 10 010`)
	seq := []Toktype{
		INT_CONST, INT_CONST, INT_CONST, INT_CONST,
	}
	MakeTest(t, *l, seq)()
}

func TestOperators(t *testing.T) {
	l := New("", `
		[ ] ( ) . -> ++ -- & * + - ~ ! / % << >> < > <=
		>= == != ^ && || ? : = *= /= %= -= <<= >>= &= ^=
		|= , { } ; ...
	`)
	seq := []Toktype{
		LBRACKET, RBRACKET, LPAREN, RPAREN, DOT, ARROW, INC, DEC,
		BAND, MUL, ADD, SUB, BCOMP, NOT, DIV, MOD, LSHIFT, RSHIFT,
		LT, GT, LEQ, GEQ, EQ, NEQ, BXOR, AND, OR, QMARK, COLON,
		ASSIGN, MUL_ASSIGN, DIV_ASSIGN, MOD_ASSIGN, SUB_ASSIGN,
		LS_ASSIGN, RS_ASSIGN, BA_ASSIGN, XO_ASSIGN, BO_ASSIGN,
		COMMA, LBRACE, RBRACE, SCOLON, ELLIP,
	}
	MakeTest(t, *l, seq)()
}

func MakeTest(t *testing.T, l Lexer, seq []Toktype) func() {
	return unittest.MakeTester(t, next(l), seq)
}
func next(l Lexer) func() Toktype {
	return func() Toktype { return l.Lex().Type }
}
