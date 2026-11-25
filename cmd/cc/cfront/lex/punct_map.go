package lex

var punct_map = map[rune]Toktype{
	'?': QMARK,
	';': SCOLON,
	':': COLON,
	'{': LBRACE,
	'}': RBRACE,
	'(': LPAREN,
	')': RPAREN,
	'[': LBRACKET,
	']': RBRACKET,
	',': COMMA,
	'~': BCOMP,
}
