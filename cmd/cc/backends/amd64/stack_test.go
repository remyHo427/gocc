package amd64

import (
	"testing"

	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/backends/amd64/fix"
	"cc260717/cmd/cc/backends/amd64/preg"
	"cc260717/cmd/cc/backends/amd64/stack"
	"cc260717/internal/unittest"
)

func TestReturn(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name: "simple integer",
			Src:  "return 0;",
			Expect: []string{
				"allocate_stack 0",
				"move $0 ax",
				"return",
			},
		},
	}
	check_stack(t, tt)
}

func TestUnary(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name: "negate",
			Src:  "unary - 1 tmp.0;",
			Expect: []string{
				"allocate_stack 4",
				"move $1 s-4",
				"unary NEG s-4",
			},
		},
		{
			Name: "negate nested",
			Src:  "unary - 1 tmp.0; unary - tmp.0 tmp.1;",
			Expect: []string{
				"allocate_stack 8",
				"move $1 s-4",
				"unary NEG s-4",
				"move s-4 r10",
				"move r10 s-8",
				"unary NEG s-8",
			},
		},
		{
			Name: "complement",
			Src:  "unary ~ 1 tmp.0;",
			Expect: []string{
				"allocate_stack 4",
				"move $1 s-4",
				"unary NOT s-4",
			},
		},
		{
			Name: "complement nested",
			Src:  "unary ~ 1 tmp.0; unary ~ tmp.0 tmp.1;",
			Expect: []string{
				"allocate_stack 8",
				"move $1 s-4",
				"unary NOT s-4",
				"move s-4 r10",
				"move r10 s-8",
				"unary NOT s-8",
			},
		},
	}
	check_stack(t, tt)
}

func TestBinary(t *testing.T) {
	tt := []unittest.NamedTestObj{
		{
			Name: "add",
			Src:  "binary + 1 2 tmp.0;",
			Expect: []string{
				"allocate_stack 4",
				"move $1 s-4",
				"binary ADD $2 s-4",
			},
		},
		{
			Name: "add nested",
			Src:  "binary + 1 2 tmp.0; binary + tmp.0 3 tmp.1;",
			Expect: []string{
				"allocate_stack 8",
				"move $1 s-4",
				"binary ADD $2 s-4",
				"move s-4 r10",
				"move r10 s-8",
				"binary ADD $3 s-8",
			},
		},
		{
			Name: "sub",
			Src:  "binary - 1 2 tmp.0;",
			Expect: []string{
				"allocate_stack 4",
				"move $1 s-4",
				"binary SUB $2 s-4",
			},
		},

		{
			Name: "sub nested",
			Src:  "binary - 1 2 tmp.0; binary - tmp.0 3 tmp.1;",
			Expect: []string{
				"allocate_stack 8",
				"move $1 s-4",
				"binary SUB $2 s-4",
				"move s-4 r10",
				"move r10 s-8",
				"binary SUB $3 s-8",
			},
		},
	}
	check_stack(t, tt)
}

func check_stack(t *testing.T, tt []unittest.NamedTestObj) {
	unittest.TestObj(t, make_stack_program, tt)
}

func make_stack_program(src string) []string {
	pg := preg.New(4)
	asm := []asm.Instruction{}

	for _, tackyIns := range parseTacky(src) {
		asm = append(asm, ToAsmInstructions(tackyIns)...)
	}

	program := pg.Replace(asmProgram(asm))
	program = fix.FixInstructions(program)
	program = stack.AllocateStack(program, -pg.GetCount())

	result := []string{}
	for _, ins := range program.FuncDef.Ins {
		result = append(result, ins.String())
	}

	return result
}
