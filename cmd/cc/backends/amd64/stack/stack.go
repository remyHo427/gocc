package stack

import "cc260717/cmd/cc/backends/amd64/asm"

func AllocateStack(program asm.Program, count int) asm.Program {
	ins := []asm.Instruction{}
	ins = append(ins, &asm.AllocateStack{Size: count})
	ins = append(ins, program.FuncDef.Ins[:]...)

	return asm.Program{
		FuncDef: asm.Function{
			Name: program.FuncDef.Name,
			Ins:  ins,
		},
	}
}
