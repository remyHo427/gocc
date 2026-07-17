package tacky

import (
	"bytes"
	"cc260717/cmd/cc/util"
)

func (op UnaryOpType) String() string {
	var strmap = map[UnaryOpType]string{
		NEGATE:     "-",
		COMPLEMENT: "~",
		NOT:        "!",
	}

	s, ok := strmap[op]
	if !ok {
		util.Exit_with_printf(
			"UnaryOpType %d doesn't have a corresponding string\n", op)
	}

	return s
}
func (op BinaryOpType) String() string {
	var strmap = map[BinaryOpType]string{
		ADD:    "+",
		SUB:    "-",
		MUL:    "*",
		DIV:    "/",
		REM:    "%",
		BAND:   "&",
		BOR:    "|",
		BXOR:   "^",
		LSHIFT: "<<",
		RSHIFT: ">>",
		EQ:     "==",
		NEQ:    "!=",
		LT:     "<",
		GT:     ">",
		LEQ:    "<=",
		GEQ:    ">=",
		AND:    "&&",
		OR:     "||",
	}

	s, ok := strmap[op]
	if !ok {
		util.Exit_with_printf(
			"BinaryOpType %d doesn't have a corresponding string\n", op)
	}

	return s
}

func join(args ...any) string {
	var out bytes.Buffer

	for i, a := range args {
		switch t := a.(type) {
		case string:
			out.WriteString(t)
		case Node:
			out.WriteString(t.String())
		case nil:
			// do nothing
		default:
			util.Exit_with_printf("cannot stringify unknown value %v", t)
		}

		if i < len(args)-1 {
			out.WriteString(" ")
		}
	}
	out.WriteString(";")

	return out.String()
}
