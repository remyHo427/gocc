package lex

import (
	"cc260717/internal/lexer"
	"strconv"
	"unicode"
)

var ERRTOK = Token{Type: ERR}
var EOFTOK = Token{Type: EOF}

type Lexer struct {
	file string
	col  int
	line int
	err  []LexerError
	lexer.BaseLexer
}

func New(filename string, src string) *Lexer {
	return &Lexer{
		file:      filename,
		col:       0,
		line:      1,
		BaseLexer: lexer.New(src),
	}
}

func (l *Lexer) LexAll() ([]Token, []LexerError) {
	toks := []Token{}

	for {
		if tok := l.Lex(); tok.Type == ERR {
			continue
		} else if tok.Type == EOF {
			break
		} else {
			toks = append(toks, tok)
			continue
		}
	}

	return toks, l.err
}
func (l *Lexer) Lex() Token {
	for !l.IsEnd() {
		c := l.Peek()

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
			return l.tok(punct)
		}

		var ttype Toktype
		switch c {
		case '.':
			l.adv()
			l.match("..", ELLIP, &ttype)
			l.match("", DOT, &ttype)
			return l.tok(ttype)
		case '+':
			l.adv()
			l.match("+", INC, &ttype)
			l.match("=", ADD_ASSIGN, &ttype)
			l.match("", ADD, &ttype)
			return l.tok(ttype)
		case '-':
			l.adv()
			l.match("-", DEC, &ttype)
			l.match("=", SUB_ASSIGN, &ttype)
			l.match(">", ARROW, &ttype)
			l.match("", SUB, &ttype)
			return l.tok(ttype)
		case '*':
			l.adv()
			l.match("=", MUL_ASSIGN, &ttype)
			l.match("", MUL, &ttype)
			return l.tok(ttype)
		case '%':
			l.adv()
			l.match("=", MOD_ASSIGN, &ttype)
			l.match("", MOD, &ttype)
			return l.tok(ttype)
		case '/':
			l.adv()
			l.match("=", DIV_ASSIGN, &ttype)
			l.match("", DIV, &ttype)
			return l.tok(ttype)
		case '<':
			l.adv()
			l.match("<=", LS_ASSIGN, &ttype)
			l.match("<", LSHIFT, &ttype)
			l.match("=", LEQ, &ttype)
			l.match("", LT, &ttype)
			return l.tok(ttype)
		case '>':
			l.adv()
			l.match(">=", RS_ASSIGN, &ttype)
			l.match(">", RSHIFT, &ttype)
			l.match("=", GEQ, &ttype)
			l.match("", GT, &ttype)
			return l.tok(ttype)
		case '&':
			l.adv()
			l.match("&", AND, &ttype)
			l.match("=", BA_ASSIGN, &ttype)
			l.match("", BAND, &ttype)
			return l.tok(ttype)
		case '|':
			l.adv()
			l.match("|", OR, &ttype)
			l.match("=", BO_ASSIGN, &ttype)
			l.match("", BOR, &ttype)
			return l.tok(ttype)
		case '^':
			l.adv()
			l.match("=", XO_ASSIGN, &ttype)
			l.match("", BXOR, &ttype)
			return l.tok(ttype)
		case '!':
			l.adv()
			l.match("=", NEQ, &ttype)
			l.match("", NOT, &ttype)
			return l.tok(ttype)
		case '=':
			l.adv()
			l.match("=", EQ, &ttype)
			l.match("", ASSIGN, &ttype)
			return l.tok(ttype)
		default:
			l.error(ErrInvalidCharacter)
			l.adv()
			return ERRTOK
		}
	}

	return EOFTOK
}

func (l *Lexer) word() Token {
	start := l.Sp
	for {
		l.adv()
		if l.IsEnd() {
			break
		} else if c := l.Peek(); !unicode.IsLetter(c) &&
			!unicode.IsDigit(c) &&
			c != '_' {
			break
		}
	}
	end := l.Sp

	s := string(l.Src[start:end])
	if kwtok, ok := kw_map[s]; ok {
		return l.tok(kwtok)
	}

	return Token{
		Type:    IDENT,
		Literal: s,
	}
}
func (l *Lexer) num() Token {
	start := l.Sp
	for {
		l.adv()
		if l.IsEnd() {
			break
		} else if !unicode.IsDigit(l.Peek()) {
			break
		}
	}
	end := l.Sp

	s := string(l.Src[start:end])
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		l.error(ErrInvalidIntLiteral)
		return ERRTOK
	}

	return Token{Type: INT_CONST, IntVal: n, Literal: s}
}
func (l *Lexer) error(e error) {
	lex_err := LexerError{
		File:  l.file,
		Line:  l.line,
		Col:   l.col,
		Error: e,
	}
	l.err = append(l.err, lex_err)
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
			l.adv()
		}
	}

	*curr = ttype
}
func (l *Lexer) adv() {
	l.col++
	l.Sp++
}
func (l *Lexer) tok(ttype Toktype) Token {
	return Token{
		File: l.file,
		Line: l.line,
		Type: ttype,
	}
}
