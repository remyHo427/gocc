package amd64

import (
	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/backends/amd64/fix"
	"cc260717/cmd/cc/backends/amd64/preg"
	"testing"
)

func TestFixInstruction(t *testing.T) {
	checkFix(t, "move", "unary - 1 tmp.0; unary ~ tmp.0 tmp.1;",
		"move $1 s-4",
		"unary NEG s-4",
		"move s-4 r10", // cannot have two stack addresses in move
		"move r10 s-8",
		"unary NOT s-8",
	)
	checkFix(t, "binary add",
		"unary - 1 tmp.0; unary - 2 tmp.1; binary + tmp.0 tmp.1 tmp.2;",
		"move $1 s-4",
		"unary NEG s-4",
		"move $2 s-8",
		"unary NEG s-8",
		"move s-4 r10", // cannot have two stack addresses in move
		"move r10 s-12",
		"move s-8 r10", // cannot have two stack addresses in addl
		"binary ADD r10 s-12",
	)
	checkFix(t, "sub",
		"unary - 1 tmp.0; unary - 2 tmp.1; binary - tmp.0 tmp.1 tmp.2;",
		"move $1 s-4",
		"unary NEG s-4",
		"move $2 s-8",
		"unary NEG s-8",
		"move s-4 r10", // cannot have two stack addresses in move
		"move r10 s-12",
		"move s-8 r10", // cannot have two stack addresses in subl
		"binary SUB r10 s-12",
	)
	checkFix(t, "mul", "binary * 1 2 tmp.0;",
		"move $1 s-4",
		"move s-4 r11",
		"binary MUL $2 r11",
		"move r11 s-4",
	)
	checkFix(t, "div", "binary / 1 2 tmp.0;",
		"move $1 ax",
		"cdq",
		"move $2 r10", // cannot have immediate value in intdiv
		"intdiv r10",
		"move ax s-4",
	)
	checkFix(t, "remainder", "binary % 1 2 tmp.0;",
		"move $1 ax",
		"cdq",
		"move $2 r10", // cannot have immediate value in intdiv
		"intdiv r10",
		"move dx s-4",
	)
}

func checkFix(t *testing.T, name string, src string, expect ...string) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		pg := preg.New(4)
		asm := []asm.Instruction{}

		for _, tackyIns := range parseTacky(src) {
			asm = append(asm, ToAsmInstructions(tackyIns)...)
		}

		program := pg.Replace(asmProgram(asm))
		program = fix.FixInstructions(program)

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
