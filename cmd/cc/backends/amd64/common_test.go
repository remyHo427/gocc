package amd64

import (
	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/ir/tacky"

	tackyparse "cc260717/cmd/cc/ir/tacky_parse"
)

func parseTacky(src string) []tacky.Instruction {
	return tackyparse.NewParser(tackyparse.New(src)).ParseAll()
}
func asmProgram(ins []asm.Instruction) asm.Program {
	return asm.Program{
		FuncDef: asm.Function{
			Name: "main",
			Ins:  ins,
		},
	}
}
