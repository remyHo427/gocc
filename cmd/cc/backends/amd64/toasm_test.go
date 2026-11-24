package amd64

import (
	"testing"

	"cc260717/internal/unittest"
)

func TestToAsmInstructionReturn(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name: "return integer",
			Src:  "return 0;",
			Expect: []string{
				"move $0 ax",
				"return",
			},
		},
		{
			Name: "return unary",
			Src:  "unary - 1 tmp.0; return tmp.0;",
			Expect: []string{
				"move $1 tmp.0",
				"unary NEG tmp.0",
				"move tmp.0 ax",
				"return",
			},
		},
		{
			Name: "return binary",
			Src:  "binary + 1 2 tmp.0; return tmp.0;",
			Expect: []string{
				"move $1 tmp.0",
				"binary ADD $2 tmp.0",
				"move tmp.0 ax",
				"return",
			},
		},
	}
	check(t, tt)
}
func TestToAsmInstructionUnary(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name: "complement",
			Src:  "unary -1 tmp.0;",
			Expect: []string{
				"move $1 tmp.0",
				"unary NEG tmp.0",
			},
		},
		{
			Name: "negate",
			Src:  "unary ~ 1 tmp.0;",
			Expect: []string{
				"move $1 tmp.0",
				"unary NOT tmp.0",
			},
		},
		{
			Name: "logical not",
			Src:  "unary ! 1 tmp.0;",
			Expect: []string{
				"comparison $0 $1",
				"move $0 tmp.0",
				"set_cc E tmp.0",
			},
		},
		{
			Name: "nested",
			Src:  "unary - 1 tmp.0; unary - tmp.0 tmp.1;",
			Expect: []string{
				"move $1 tmp.0",
				"unary NEG tmp.0",
				"move tmp.0 tmp.1",
				"unary NEG tmp.1",
			},
		},
	}
	check(t, tt)
}
func TestToAsmInstructionBinary(t *testing.T) {
	tt := []unittest.NamedTestObj{
		// arithmetic
		{
			Name: "plus",
			Src:  "binary + 1 2 tmp.0;",
			Expect: []string{
				"move $1 tmp.0",
				"binary ADD $2 tmp.0",
			},
		},
		{
			Name: "sub",
			Src:  "binary - 1 2 tmp.0;",
			Expect: []string{
				"move $1 tmp.0",
				"binary SUB $2 tmp.0",
			},
		},
		{
			Name: "mul",
			Src:  "binary * 1 2 tmp.0;",
			Expect: []string{
				"move $1 tmp.0",
				"binary MUL $2 tmp.0",
			},
		},
		{
			Name: "div",
			Src:  "binary / 1 2 tmp.0;",
			Expect: []string{
				"move $1 ax",
				"cdq",
				"intdiv $2",
				"move ax tmp.0",
			},
		},
		{
			Name: "remainder",
			Src:  "binary % 1 2 tmp.0;",
			Expect: []string{
				"move $1 ax",
				"cdq",
				"intdiv $2",
				"move dx tmp.0",
			},
		},
		// relational
		{
			Name: "less than",
			Src:  "binary < 1 2 tmp.0;",
			Expect: []string{
				"comparison $2 $1",
				"move $0 tmp.0",
				"set_cc L tmp.0",
			},
		},
		{
			Name: "less than or equal",
			Src:  "binary <= 1 2 tmp.0;",
			Expect: []string{
				"comparison $2 $1",
				"move $0 tmp.0",
				"set_cc LE tmp.0",
			},
		},
		{
			Name: "greater than",
			Src:  "binary > 1 2 tmp.0;",
			Expect: []string{
				"comparison $2 $1",
				"move $0 tmp.0",
				"set_cc G tmp.0",
			},
		},
		{
			Name: "greater than or equal",
			Src:  "binary >= 1 2 tmp.0;",
			Expect: []string{
				"comparison $2 $1",
				"move $0 tmp.0",
				"set_cc GE tmp.0",
			},
		},
		{
			Name: "equal",
			Src:  "binary == 1 2 tmp.0;",
			Expect: []string{
				"comparison $2 $1",
				"move $0 tmp.0",
				"set_cc E tmp.0",
			},
		},
		{
			Name: "not equal",
			Src:  "binary != 1 2 tmp.0;",
			Expect: []string{
				"comparison $2 $1",
				"move $0 tmp.0",
				"set_cc NE tmp.0",
			},
		},
		{
			Name: "nested",
			Src:  "binary + 1 2 tmp.0; binary + tmp.0 3 tmp.1;",
			Expect: []string{
				"move $1 tmp.0",
				"binary ADD $2 tmp.0",
				"move tmp.0 tmp.1",
				"binary ADD $3 tmp.1",
			},
		},
	}
	check(t, tt)
}

func TestMiscInstructions(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name: "label",
			Src:  "label lb.0;",
			Expect: []string{
				"label lb.0",
			},
		},
		{
			Name: "copy",
			Src:  "copy 1 tmp.0;",
			Expect: []string{
				"move $1 tmp.0",
			},
		},
		{
			Name: "jump",
			Src:  "jump lb.0;",
			Expect: []string{
				"jump lb.0",
			},
		},
		{
			Name: "jump_if_zero",
			Src:  "jump_if_zero 1 lb.0;",
			Expect: []string{
				"comparison $0 $1",
				"jump_cc E lb.0",
			},
		},
		{
			Name: "jump_if_not_zero",
			Src:  "jump_if_not_zero 0 lb.0;",
			Expect: []string{
				"comparison $0 $0",
				"jump_cc NE lb.0",
			},
		},
	}

	check(t, tt)
}

func check(t *testing.T, tt []unittest.NamedTestObj) {
	unittest.TestObj(t, make_asm, tt)
}
func make_asm(src string) []string {
	result := []string{}

	for _, tackyIns := range parseTacky(src) {
		for _, asmIns := range ToAsmInstructions(tackyIns) {
			result = append(result, asmIns.String())
		}
	}

	return result
}
