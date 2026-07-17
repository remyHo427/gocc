package amd64

import (
	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/backends/amd64/fix"
	"cc260717/cmd/cc/backends/amd64/preg"
	"cc260717/cmd/cc/backends/amd64/stack"
	"testing"
)

func TestReturn(t *testing.T) {
	checkAllocateStack(t, "simple integer", "return 0;",
		"allocate_stack 0",
		"move $0 ax",
		"return",
	)
}

func TestUnary(t *testing.T) {
	checkAllocateStack(t, "negate", "unary - 1 tmp.0;",
		"allocate_stack 4",
		"move $1 s-4",
		"unary NEG s-4",
	)
	checkAllocateStack(t, "negate nested",
		"unary - 1 tmp.0; unary - tmp.0 tmp.1;",

		"allocate_stack 8",
		"move $1 s-4",
		"unary NEG s-4",
		"move s-4 r10",
		"move r10 s-8",
		"unary NEG s-8",
	)
	checkAllocateStack(t, "complement", "unary ~ 1 tmp.0;",
		"allocate_stack 4",
		"move $1 s-4",
		"unary NOT s-4",
	)
	checkAllocateStack(t, "complement nested",
		"unary ~ 1 tmp.0; unary ~ tmp.0 tmp.1;",

		"allocate_stack 8",
		"move $1 s-4",
		"unary NOT s-4",
		"move s-4 r10",
		"move r10 s-8",
		"unary NOT s-8",
	)
}

func TestBinary(t *testing.T) {
	checkAllocateStack(t, "add", "binary + 1 2 tmp.0;",
		"allocate_stack 4",
		"move $1 s-4",
		"binary ADD $2 s-4",
	)
	checkAllocateStack(t, "add nested",
		"binary + 1 2 tmp.0; binary + tmp.0 3 tmp.1;",

		"allocate_stack 8",
		"move $1 s-4",
		"binary ADD $2 s-4",
		"move s-4 r10",
		"move r10 s-8",
		"binary ADD $3 s-8",
	)

	checkAllocateStack(t, "sub", "binary - 1 2 tmp.0;",
		"allocate_stack 4",
		"move $1 s-4",
		"binary SUB $2 s-4",
	)
	checkAllocateStack(t, "sub nested",
		"binary - 1 2 tmp.0; binary - tmp.0 3 tmp.1;",

		"allocate_stack 8",
		"move $1 s-4",
		"binary SUB $2 s-4",
		"move s-4 r10",
		"move r10 s-8",
		"binary SUB $3 s-8",
	)
}

func checkAllocateStack(t *testing.T, name string, src string, expect ...string) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		pg := preg.New(4)
		asm := []asm.Instruction{}

		for _, tackyIns := range parseTacky(src) {
			asm = append(asm, ToAsmInstructions(tackyIns)...)
		}

		program := pg.Replace(asmProgram(asm))
		program = fix.FixInstructions(program)
		program = stack.AllocateStack(program, -pg.GetCount())

		if len(program.FuncDef.Ins) != len(expect) {
			t.Errorf("expect %d instructions, got %d instead",
				len(program.FuncDef.Ins), len(expect))
		}

		for i, ins := range program.FuncDef.Ins {
			if actual := ins.String(); expect[i] != actual {
				t.Errorf("expect \"%s\", got \"%s\"\n",
					expect[i], actual)
			}
		}
	})
}
