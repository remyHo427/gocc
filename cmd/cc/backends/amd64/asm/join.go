package asm

import (
	"bytes"

	"cc260717/cmd/cc/util"
)

func join(args ...any) string {
	var out bytes.Buffer

	for i, arg := range args {
		switch t := arg.(type) {
		case Node:
			out.WriteString(t.String())
		case string:
			out.WriteString(t)
		default:
			util.Exit_with_printf("unknown value %v\n", t)
		}

		if i != len(args)-1 {
			out.WriteString(" ")
		}
	}

	return out.String()
}
