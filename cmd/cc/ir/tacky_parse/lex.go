package tackyparse

import (
	"unicode"

	"cc260717/cmd/cc/util"
)

type Lexer struct {
	src []rune
	sp  int
}

func New(src string) *Lexer {
	return &Lexer{
		src: []rune(src),
	}
}

func (l *Lexer) Lex() Token {
	for !l.isend() {
		c := l.peek()

		if unicode.IsSpace(c) {
			l.adv()
			continue
		}
		if unicode.IsLetter(c) || c == '_' {
			start := l.sp
			for {
				l.adv()
				if l.isend() {
					break
				} else if c := l.peek(); !unicode.IsLetter(c) &&
					!unicode.IsDigit(c) &&
					c != '.' && c != '_' {
					break
				}
			}
			end := l.sp

			s := string(l.src[start:end])
			if kwtok, ok := kw_map[s]; ok {
				return kwtok
			}

			return Token{
				Type:    IDENT,
				Literal: s,
			}
		}
		if unicode.IsDigit(c) {
			start := l.sp
			for {
				l.adv()
				if l.isend() {
					break
				} else if !unicode.IsDigit(l.peek()) {
					break
				}
			}
			end := l.sp

			return Token{
				Type:    INTEGER,
				Literal: string(l.src[start:end]),
			}
		}

		var ttype Toktype
		switch c {
		case '-':
			l.adv()
			return tok(SUB)
		case '~':
			l.adv()
			return tok(TILDE)
		case '+':
			l.adv()
			return tok(ADD)
		case '*':
			l.adv()
			return tok(MUL)
		case '/':
			l.adv()
			return tok(DIV)
		case '%':
			l.adv()
			return tok(MOD)
		case '^':
			l.adv()
			return tok(BXOR)
		case ';':
			l.adv()
			return tok(SCOLON)
		case '!':
			l.adv()
			l.match("=", NEQ, &ttype)
			l.match("", NOT, &ttype)
			return tok(ttype)
		case '|':
			l.adv()
			l.match("|", OR, &ttype)
			l.match("", BOR, &ttype)
			return tok(ttype)
		case '&':
			l.adv()
			l.match("&", AND, &ttype)
			l.match("", BAND, &ttype)
			return tok(ttype)
		case '=':
			l.adv()
			l.match("=", EQ, &ttype)
			l.match("", ASSIGN, &ttype)
			return tok(ttype)
		case '<':
			l.adv()
			l.match("<", LSHIFT, &ttype)
			l.match("=", LEQ, &ttype)
			l.match("", LT, &ttype)
			return tok(ttype)
		case '>':
			l.adv()
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

	start := l.sp
	for _, c := range s {
		if l.isend() || l.peek() != rune(c) {
			l.sp = start
			return
		} else {
			l.adv()
		}
	}

	*curr = ttype
}
func (l *Lexer) peek() rune {
	return l.src[l.sp]
}
func (l *Lexer) adv() {
	l.sp++
}
func (l *Lexer) isend() bool {
	return l.sp >= len(l.src)
}
func tok(ttype Toktype) Token {
	return Token{Type: ttype}
}
