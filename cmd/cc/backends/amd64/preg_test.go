package amd64

import (
	"testing"

	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/backends/amd64/preg"
	"cc260717/internal/unittest"
)

func TestPseudoRegistersReplacement(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name: "return",
			Src:  "return 0;",
			Expect: []string{
				"move $0 ax",
				"return",
			},
		},
		{
			Name: "unary negate",
			Src:  "unary - 1 tmp.0;",
			Expect: []string{
				"move $1 s-4",
				"unary NEG s-4",
			},
		},
		{
			Name: "unary not",
			Src:  "unary ! 1 tmp.0;",
			Expect: []string{
				"comparison $0 $1",
				"move $0 s-4",
				"set_cc E s-4",
			},
		},
		{
			Name: "unary nested",
			Src:  "unary - 1 tmp.0; unary ~ tmp.0 tmp.1;",
			Expect: []string{
				"move $1 s-4",
				"unary NEG s-4",
				"move s-4 s-8",
				"unary NOT s-8",
			},
		},
		// binary
		{
			Name: "add",
			Src:  "binary + 1 2 tmp.0;",
			Expect: []string{
				"move $1 s-4",
				"binary ADD $2 s-4",
			},
		},
		{
			Name: "sub",
			Src:  "binary - 1 2 tmp.0;",
			Expect: []string{
				"move $1 s-4",
				"binary SUB $2 s-4",
			},
		},
		{
			Name: "mul",
			Src:  "binary * 1 2 tmp.0;",
			Expect: []string{
				"move $1 s-4",
				"binary MUL $2 s-4",
			},
		},
		{
			Name: "div",
			Src:  "binary / 1 2 tmp.0;",
			Expect: []string{
				"move $1 ax",
				"cdq",
				"intdiv $2",
				"move ax s-4",
			},
		},
		{
			Name: "remainder",
			Src:  "binary % 1 2 tmp.0;",
			Expect: []string{
				"move $1 ax",
				"cdq",
				"intdiv $2",
				"move dx s-4",
			},
		},
		{
			Name: "less than",
			Src:  "binary < 1 2 tmp.0;",
			Expect: []string{
				"comparison $2 $1",
				"move $0 s-4",
				"set_cc L s-4",
			},
		},
		{
			Name: "binary nested",
			Src:  "binary + 1 2 tmp.0; binary + tmp.0 3 tmp.1;",
			Expect: []string{
				"move $1 s-4",
				"binary ADD $2 s-4",
				"move s-4 s-8",
				"binary ADD $3 s-8",
			},
		},
	}
	check_pg(t, tt)
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

func check_pg(t *testing.T, tt []unittest.NamedTestObj) {
	unittest.TestObj(t, make_pg_program, tt)
}
func make_pg_program(src string) []string {
	pg := preg.New(4)
	asm := []asm.Instruction{}

	for _, tackyIns := range parseTacky(src) {
		asm = append(asm, ToAsmInstructions(tackyIns)...)
	}

	program := pg.Replace(asmProgram(asm))
	result := []string{}
	for _, ins := range program.FuncDef.Ins {
		result = append(result, ins.String())
	}

	return result
}
