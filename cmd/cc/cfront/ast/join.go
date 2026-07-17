package ast

import (
	"bytes"
	"cc260717/cmd/cc/util"
)

func join(args ...any) string {
	var out bytes.Buffer

	out.WriteString("(")
	for i, a := range args {
		switch t := a.(type) {
		case Node:
			out.WriteString(t.String())
		case string:
			out.WriteString(t)
		case []any:
			out.WriteString(join(t))
		case nil:
			// do nothing
		default:
			util.Exit_with_printf("cannot stringify unknown value %v", t)
		}

		if i < len(args)-1 {
			out.WriteString(" ")
		}
	}
	out.WriteString(")")

	return out.String()
}
