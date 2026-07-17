package lex

type Token struct {
	Type    Toktype
	Literal string
	IntVal  int64
}
