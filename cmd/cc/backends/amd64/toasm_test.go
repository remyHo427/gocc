package amd64

import (
	"cc260717/cmd/cc/backends/amd64/asm"
	"testing"
)

func TestToAsmInstructionReturn(t *testing.T) {
	checkToAsm(t, "return integer", "return 0;",
		"move $0 ax",
		"return",
	)
	checkToAsm(t, "return unary", "unary - 1 tmp.0; return tmp.0;",
		"move $1 tmp.0",
		"unary NEG tmp.0",
		"move tmp.0 ax",
		"return",
	)
	checkToAsm(t, "return binary", "binary + 1 2 tmp.0; return tmp.0;",
		"move $1 tmp.0",
		"binary ADD $2 tmp.0",
		"move tmp.0 ax",
		"return",
	)
}
func TestToAsmInstructionUnary(t *testing.T) {
	checkToAsm(t, "complement", "unary - 1 tmp.0;",
		"move $1 tmp.0",
		"unary NEG tmp.0",
	)
	checkToAsm(t, "negate", "unary ~ 1 tmp.0;",
		"move $1 tmp.0",
		"unary NOT tmp.0",
	)
	checkToAsm(t, "logical not", "unary ! 1 tmp.0;",
		"comparison $0 $1",
		"move $0 tmp.0",
		"set_cc E tmp.0",
	)
	checkToAsm(t, "nested", "unary - 1 tmp.0; unary - tmp.0 tmp.1;",
		"move $1 tmp.0",
		"unary NEG tmp.0",
		"move tmp.0 tmp.1",
		"unary NEG tmp.1",
	)
}
func TestToAsmInstructionBinary(t *testing.T) {
	// arithmetic
	checkToAsm(t, "plus", "binary + 1 2 tmp.0;",
		"move $1 tmp.0",
		"binary ADD $2 tmp.0",
	)
	checkToAsm(t, "sub", "binary - 1 2 tmp.0;",
		"move $1 tmp.0",
		"binary SUB $2 tmp.0",
	)
	checkToAsm(t, "mul", "binary * 1 2 tmp.0;",
		"move $1 tmp.0",
		"binary MUL $2 tmp.0",
	)
	checkToAsm(t, "div", "binary / 1 2 tmp.0;",
		"move $1 ax",
		"cdq",
		"intdiv $2",
		"move ax tmp.0",
	)
	checkToAsm(t, "remainder", "binary % 1 2 tmp.0;",
		"move $1 ax",
		"cdq",
		"intdiv $2",
		"move dx tmp.0",
	)

	// relational
	checkToAsm(t, "less than", "binary < 1 2 tmp.0;",
		"comparison $2 $1",
		"move $0 tmp.0",
		"set_cc L tmp.0",
	)
	checkToAsm(t, "less than", "binary < 1 2 tmp.0;",
		"comparison $2 $1",
		"move $0 tmp.0",
		"set_cc L tmp.0",
	)
	checkToAsm(t, "less than or equal", "binary <= 1 2 tmp.0;",
		"comparison $2 $1",
		"move $0 tmp.0",
		"set_cc LE tmp.0",
	)
	checkToAsm(t, "greater than", "binary > 1 2 tmp.0;",
		"comparison $2 $1",
		"move $0 tmp.0",
		"set_cc G tmp.0",
	)
	checkToAsm(t, "greater than or equal", "binary >= 1 2 tmp.0;",
		"comparison $2 $1",
		"move $0 tmp.0",
		"set_cc GE tmp.0",
	)
	checkToAsm(t, "equal", "binary == 1 2 tmp.0;",
		"comparison $2 $1",
		"move $0 tmp.0",
		"set_cc E tmp.0",
	)
	checkToAsm(t, "not equal", "binary != 1 2 tmp.0;",
		"comparison $2 $1",
		"move $0 tmp.0",
		"set_cc NE tmp.0",
	)

	checkToAsm(t, "nested", "binary + 1 2 tmp.0; binary + tmp.0 3 tmp.1;",
		"move $1 tmp.0",
		"binary ADD $2 tmp.0",
		"move tmp.0 tmp.1",
		"binary ADD $3 tmp.1",
	)
}

func TestMiscInstructions(t *testing.T) {
	checkToAsm(t, "label", "label lb.0;", "label lb.0")
	checkToAsm(t, "copy", "copy 1 tmp.0;", "move $1 tmp.0")
	checkToAsm(t, "jump", "jump lb.0;", "jump lb.0")
	checkToAsm(t, "jump_if_zero", "jump_if_zero 1 lb.0;",
		"comparison $0 $1",
		"jump_cc E lb.0",
	)
	checkToAsm(t, "jump_if_not_zero", "jump_if_not_zero 0 lb.0;",
		"comparison $0 $0",
		"jump_cc NE lb.0",
	)
}

func checkToAsm(t *testing.T, name string, src string, expect ...string) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		asm := []asm.Instruction{}

		for _, tackyIns := range parseTacky(src) {
			asm = append(asm, ToAsmInstructions(tackyIns)...)
		}

		if len(asm) != len(expect) {
			t.Errorf("expect %d instructions, got %d instead",
				len(expect), len(asm))
		}

		for i, ins := range asm {
			if actual := ins.String(); expect[i] != actual {
				t.Errorf("expect \"%s\", got \"%s\"",
					expect[i], actual)
			}
		}
	})
}
