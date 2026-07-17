package tackyparse

var kw_map = map[string]Token{
	"return":           tok(RETURN),
	"unary":            tok(UNARY),
	"binary":           tok(BINARY),
	"copy":             tok(COPY),
	"jump":             tok(JUMP),
	"jump_if_zero":     tok(JUMP_IF_ZERO),
	"jump_if_not_zero": tok(JUMP_IF_NOT_ZERO),
	"label":            tok(LABEL),
}
