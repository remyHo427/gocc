package lex

import (
	"fmt"
	"strconv"
	"unicode"
)

type Lexer struct {
	sp  int
	len int

	file string
	col  int
	line int

	src []rune
}

func New(filename string, src string) *Lexer {
	return &Lexer{
		sp:   0,
		file: filename,
		col:  0,
		line: 1,
		len:  len(src),
		src:  []rune(src),
	}
}
func (l *Lexer) Lex() Token {
	for !l.isend() {
		c := l.peek()

		if unicode.IsSpace(c) {
			if c == '\n' {
				l.line++
				l.col = 0
			}
			l.adv()
			continue
		}
		if unicode.IsLetter(c) || c == '_' {
			return l.word()
		}
		if unicode.IsDigit(c) {
			return l.num()
		}
		if punct, ok := punct_map[c]; ok {
			l.adv()
			return punct
		}

		var ttype Toktype
		switch c {
		case '.':
			l.adv()
			l.match("..", ELLIP, &ttype)
			l.match("", DOT, &ttype)
			return tok(ttype)
		case '+':
			l.adv()
			l.match("+", INC, &ttype)
			l.match("=", ADD_ASSIGN, &ttype)
			l.match("", ADD, &ttype)
			return tok(ttype)
		case '-':
			l.adv()
			l.match("-", DEC, &ttype)
			l.match("=", SUB_ASSIGN, &ttype)
			l.match(">", ARROW, &ttype)
			l.match("", SUB, &ttype)
			return tok(ttype)
		case '*':
			l.adv()
			l.match("=", MUL_ASSIGN, &ttype)
			l.match("", MUL, &ttype)
			return tok(ttype)
		case '%':
			l.adv()
			l.match("=", MOD_ASSIGN, &ttype)
			l.match("", MOD, &ttype)
			return tok(ttype)
		case '/':
			l.adv()
			l.match("=", DIV_ASSIGN, &ttype)
			l.match("", DIV, &ttype)
			return tok(ttype)
		case '<':
			l.adv()
			l.match("<=", LS_ASSIGN, &ttype)
			l.match("<", LSHIFT, &ttype)
			l.match("=", LEQ, &ttype)
			l.match("", LT, &ttype)
			return tok(ttype)
		case '>':
			l.adv()
			l.match(">=", RS_ASSIGN, &ttype)
			l.match(">", RSHIFT, &ttype)
			l.match("=", GEQ, &ttype)
			l.match("", GT, &ttype)
			return tok(ttype)
		case '&':
			l.adv()
			l.match("&", AND, &ttype)
			l.match("=", BA_ASSIGN, &ttype)
			l.match("", BAND, &ttype)
			return tok(ttype)
		case '|':
			l.adv()
			l.match("|", OR, &ttype)
			l.match("=", BO_ASSIGN, &ttype)
			l.match("", BOR, &ttype)
			return tok(ttype)
		case '^':
			l.adv()
			l.match("=", XO_ASSIGN, &ttype)
			l.match("", BXOR, &ttype)
			return tok(ttype)
		case '!':
			l.adv()
			l.match("=", NEQ, &ttype)
			l.match("", NOT, &ttype)
			return tok(ttype)
		case '=':
			l.adv()
			l.match("=", EQ, &ttype)
			l.match("", ASSIGN, &ttype)
			return tok(ttype)
		default:
			l.adv()
			return Token{
				Type:     ERR,
				Error:    ErrInvalidCharacter,
				ErrorStr: fmt.Sprintf("'%c'", c),
				File:     l.file,
				Col:      l.col,
				Line:     l.line,
			}
		}
	}

	return tok(EOF)
}

func (l *Lexer) word() Token {
	start := l.sp
	for {
		l.adv()
		if l.isend() {
			break
		} else if c := l.peek(); !unicode.IsLetter(c) &&
			!unicode.IsDigit(c) &&
			c != '_' {
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
func (l *Lexer) num() Token {
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

	s := string(l.src[start:end])
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return Token{
			Type:     ERR,
			Error:    ErrInvalidIntLiteral,
			ErrorStr: err.Error(),
			File:     l.file,
			Col:      l.col,
			Line:     l.line,
		}
	}

	return Token{Type: INT_CONST, IntVal: n, Literal: s}
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
func (l Lexer) peek() rune {
	return l.src[l.sp]
}
func (l *Lexer) adv() {
	l.col++
	l.sp++
}
func (l *Lexer) isend() bool {
	return l.sp >= l.len
}
func tok(ttype Toktype) Token {
	return Token{Type: ttype}
}
