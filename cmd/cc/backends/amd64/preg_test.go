package amd64

import (
	"testing"

	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/backends/amd64/preg"
)

func TestPseudoRegistersReplacement(t *testing.T) {
	checkPG(t, "return", "return 0;",
		"move $0 ax",
		"return",
	)
	checkPG(t, "unary", "unary - 1 tmp.0;",
		"move $1 s-4",
		"unary NEG s-4",
	)
	checkPG(t, "unary not", "unary ! 1 tmp.0;",
		"comparison $0 $1",
		"move $0 s-4",
		"set_cc E s-4",
	)
	checkPG(t, "unary nested", "unary - 1 tmp.0; unary ~ tmp.0 tmp.1;",
		"move $1 s-4",
		"unary NEG s-4",
		"move s-4 s-8",
		"unary NOT s-8",
	)

	// binary
	checkPG(t, "binary", "binary + 1 2 tmp.0;",
		"move $1 s-4",
		"binary ADD $2 s-4",
	)
	checkPG(t, "binary division", "binary / 1 2 tmp.0;",
		"move $1 ax",
		"cdq",
		"intdiv $2",
		"move ax s-4",
	)
	checkPG(t, "binary remainder", "binary % 1 2 tmp.0;",
		"move $1 ax",
		"cdq",
		"intdiv $2",
		"move dx s-4",
	)
	checkPG(t, "less than", "binary < 1 2 tmp.0;",
		"comparison $2 $1",
		"move $0 s-4",
		"set_cc L s-4",
	)
	checkPG(t, "binary nested", "binary + 1 2 tmp.0; binary + tmp.0 3 tmp.1;",
		"move $1 s-4",
		"binary ADD $2 s-4",
		"move s-4 s-8",
		"binary ADD $3 s-8",
	)
}

func TestMiscInstructionsPG(t *testing.T) {
	checkPG(t, "copy", "copy 1 tmp.0;", "move $1 s-4")
}

func checkPG(t *testing.T, name string, src string, expect ...string) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		pg := preg.New(4)
		asm := []asm.Instruction{}

		for _, tackyIns := range parseTacky(src) {
			asm = append(asm, ToAsmInstructions(tackyIns)...)
		}

		program := pg.Replace(asmProgram(asm))

		if len(program.FuncDef.Ins) != len(expect) {
			t.Errorf("expect %d instructions, got %d instead",
				len(program.FuncDef.Ins), len(expect))
		}

		for i, ins := range program.FuncDef.Ins {
			if actual := ins.String(); expect[i] != actual {
				t.Errorf("expect \"%s\", got \"%s\"",
					expect[i], actual)
			}
		}
	})
}
