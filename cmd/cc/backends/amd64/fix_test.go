package amd64

import (
	"testing"

	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/backends/amd64/fix"
	"cc260717/cmd/cc/backends/amd64/preg"
	"cc260717/internal/unittest"
)

func TestFixInstruction(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name: "move",
			Src:  "unary - 1 tmp.0; unary ~ tmp.0 tmp.1;",
			Expect: []string{
				"move $1 s-4",
				"unary NEG s-4",
				"move s-4 r10",
				"move r10 s-8", // cannot have two stack addresses in move
				"unary NOT s-8",
			},
		},
		{
			Name: "binary add",
			Src:  "unary - 1 tmp.0; unary - 2 tmp.1; binary + tmp.0 tmp.1 tmp.2;",
			Expect: []string{
				"move $1 s-4",
				"unary NEG s-4",
				"move $2 s-8",
				"unary NEG s-8",
				"move s-4 r10", // cannot have two stack addresses in move
				"move r10 s-12",
				"move s-8 r10", // cannot have two stack addresses in addl
				"binary ADD r10 s-12",
			},
		},
		{
			Name: "sub",
			Src:  "unary - 1 tmp.0; unary - 2 tmp.1; binary - tmp.0 tmp.1 tmp.2;",
			Expect: []string{
				"move $1 s-4",
				"unary NEG s-4",
				"move $2 s-8",
				"unary NEG s-8",
				"move s-4 r10", // cannot have two stack addresses in move
				"move r10 s-12",
				"move s-8 r10", // cannot have two stack addresses in subl
				"binary SUB r10 s-12",
			},
		},
		{
			Name: "mul",
			Src:  "binary * 1 2 tmp.0;",
			Expect: []string{
				"move $1 s-4",
				"move s-4 r11",
				"binary MUL $2 r11",
				"move r11 s-4",
			},
		},
		{
			Name: "div",
			Src:  "binary / 1 2 tmp.0;",
			Expect: []string{
				"move $1 ax",
				"cdq",
				"move $2 r10", // cannot have immediate value in intdiv
				"intdiv r10",
				"move ax s-4",
			},
		},
		{
			Name: "remainder",
			Src:  "binary % 1 2 tmp.0;",
			Expect: []string{
				"move $1 ax",
				"cdq",
				"move $2 r10", // cannot have immediate value in intdiv
				"intdiv r10",
				"move dx s-4",
			},
		},
	}
	check_fix(t, tt)
}

func check_fix(t *testing.T, tt []unittest.NamedTestObj) {
	unittest.TestObj(t, make_fixed_program, tt)
}
func make_fixed_program(src string) []string {
	pg := preg.New(4)
	asm := []asm.Instruction{}

	for _, tackyIns := range parseTacky(src) {
		asm = append(asm, ToAsmInstructions(tackyIns)...)
	}

	program := pg.Replace(asmProgram(asm))
	program = fix.FixInstructions(program)

	result := []string{}
	for _, ins := range program.FuncDef.Ins {
		result = append(result, ins.String())
	}

	return result
}
