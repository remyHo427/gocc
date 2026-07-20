package lex

import "testing"

func TestEmptyString(t *testing.T) {
	l := New("", "")
	seq := []Toktype{}
	tokseq(*l, seq, t)
}
func TestSpaceOnly(t *testing.T) {
	l := New("", "\f\n\r\t\v ")
	seq := []Toktype{}
	tokseq(*l, seq, t)
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
	tokseq(*l, seq, t)
}

func TestIdentifiers(t *testing.T) {
	l := New("", `a Toktype a0 a00 __test__`)
	seq := []Toktype{
		IDENT, IDENT, IDENT, IDENT, IDENT,
	}
	tokseq(*l, seq, t)
}
func TestInteger(t *testing.T) {
	l := New("", `0 1 10 010`)
	seq := []Toktype{
		INT_CONST, INT_CONST, INT_CONST, INT_CONST,
	}
	tokseq(*l, seq, t)
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
	tokseq(*l, seq, t)
}

func tokseq(l Lexer, seq []Toktype, t *testing.T) {
	t.Helper()

	i := 0
	for _, ttype := range seq {
		if tok := l.Lex(); tok.Type != ttype {
			t.Errorf("expected %s, got %s at seq[%d]",
				ttype.String(), ttype.String(), i)
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
