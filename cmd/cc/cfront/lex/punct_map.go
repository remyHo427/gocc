package lex

var punct_map = map[rune]Token{
	'?': tok(QMARK),
	';': tok(SCOLON),
	':': tok(COLON),
	'{': tok(LBRACE),
	'}': tok(RBRACE),
	'(': tok(LPAREN),
	')': tok(RPAREN),
	'[': tok(LBRACKET),
	']': tok(RBRACKET),
	',': tok(COMMA),
	'~': tok(BCOMP),
}
