package lex

type Token struct {
	Type    Toktype
	Literal string

	//
	IntVal int64

	//
	Error    error
	ErrorStr string
	Col      int
	Line     int
	File     string
}
