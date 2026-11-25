package lex

import "errors"

var ErrInvalidCharacter = errors.New("Invalid character")
var ErrInvalidIntLiteral = errors.New("Invalid integer literal")

type LexerError struct {
	File  string
	Col   int
	Line  int
	Error error // underlying error
}
