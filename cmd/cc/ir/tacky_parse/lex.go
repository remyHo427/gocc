package tackyparse

import (
	"unicode"

	"cc260717/cmd/cc/util"
	"cc260717/internal/lexer"
)

type Lexer struct {
	lexer.BaseLexer
}

func New(src string) *Lexer {
	return &Lexer{
		BaseLexer: lexer.New(src),
	}
}

func (l *Lexer) Lex() Token {
	for !l.IsEnd() {
		c := l.Peek()

		if unicode.IsSpace(c) {
			l.Adv()
			continue
		}
		if unicode.IsLetter(c) || c == '_' {
			start := l.Sp
			for {
				l.Adv()
				if l.IsEnd() {
					break
				} else if c := l.Peek(); !unicode.IsLetter(c) &&
					!unicode.IsDigit(c) &&
					c != '.' && c != '_' {
					break
				}
			}
			end := l.Sp

			s := string(l.Src[start:end])
			if kwtok, ok := kw_map[s]; ok {
				return kwtok
			}

			return Token{
				Type:    IDENT,
				Literal: s,
			}
		}
		if unicode.IsDigit(c) {
			start := l.Sp
			for {
				l.Adv()
				if l.IsEnd() {
					break
				} else if !unicode.IsDigit(l.Peek()) {
					break
				}
			}
			end := l.Sp

			return Token{
				Type:    INTEGER,
				Literal: string(l.Src[start:end]),
			}
		}

		var ttype Toktype
		switch c {
		case '-':
			l.Adv()
			return tok(SUB)
		case '~':
			l.Adv()
			return tok(TILDE)
		case '+':
			l.Adv()
			return tok(ADD)
		case '*':
			l.Adv()
			return tok(MUL)
		case '/':
			l.Adv()
			return tok(DIV)
		case '%':
			l.Adv()
			return tok(MOD)
		case '^':
			l.Adv()
			return tok(BXOR)
		case ';':
			l.Adv()
			return tok(SCOLON)
		case '!':
			l.Adv()
			l.match("=", NEQ, &ttype)
			l.match("", NOT, &ttype)
			return tok(ttype)
		case '|':
			l.Adv()
			l.match("|", OR, &ttype)
			l.match("", BOR, &ttype)
			return tok(ttype)
		case '&':
			l.Adv()
			l.match("&", AND, &ttype)
			l.match("", BAND, &ttype)
			return tok(ttype)
		case '=':
			l.Adv()
			l.match("=", EQ, &ttype)
			l.match("", ASSIGN, &ttype)
			return tok(ttype)
		case '<':
			l.Adv()
			l.match("<", LSHIFT, &ttype)
			l.match("=", LEQ, &ttype)
			l.match("", LT, &ttype)
			return tok(ttype)
		case '>':
			l.Adv()
			l.match(">", RSHIFT, &ttype)
			l.match("=", GEQ, &ttype)
			l.match("", GT, &ttype)
			return tok(ttype)
		default:
			util.Exit_with_printf("unknown character '%c'\n", c)
		}
	}

	return tok(EOF)
}

func (l *Lexer) match(s string, ttype Toktype, curr *Toktype) {
	if *curr != EOF {
		return
	}

	start := l.Sp
	for _, c := range s {
		if l.IsEnd() || l.Peek() != rune(c) {
			l.Sp = start
			return
		} else {
			l.Adv()
		}
	}

	*curr = ttype
}
func tok(ttype Toktype) Token {
	return Token{Type: ttype}
}
