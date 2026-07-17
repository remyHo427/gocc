package fix

import (
	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/util"
)

var (
	REG_AX   = asm.Register{Reg: asm.AX}
	REG_DX   = asm.Register{Reg: asm.DX}
	REG_R10D = asm.Register{Reg: asm.R10}
	REG_R11D = asm.Register{Reg: asm.R11}

	// singleton instructions
	INS_RET = asm.Return{}
	INS_CDQ = asm.Cdq{}
)

func FixInstructions(program asm.Program) asm.Program {
	ins := []asm.Instruction{}

	for _, instruction := range program.FuncDef.Ins {
		switch t := instruction.(type) {
		case *asm.Return, *asm.Unary, *asm.AllocateStack, *asm.Cdq:
			ins = append(ins, t)
		case *asm.Comparison:
			if t.Value2.Type() == asm.IMM {
				ins = append(ins,
					&asm.Move{
						Src: t.Value1,
						Dst: &REG_R11D,
					},
					&asm.Comparison{
						Value1: t.Value2,
						Value2: &REG_R11D,
					},
				)
				continue
			}
			if t.Value1.Type() == asm.STACK && t.Value2.Type() == asm.STACK {
				ins = append(ins,
					&asm.Move{
						Src: t.Value1,
						Dst: &REG_R10D,
					},
					&asm.Move{
						Src: &REG_R10D,
						Dst: t.Value2,
					},
				)
				continue
			}
			ins = append(ins, t)
		case *asm.Move:
			if t.Src.Type() != asm.STACK || t.Dst.Type() != asm.STACK {
				ins = append(ins, t)
			} else {
				ins = append(ins,
					&asm.Move{
						Src: t.Src,
						Dst: &REG_R10D,
					},
					&asm.Move{
						Src: &REG_R10D,
						Dst: t.Dst,
					},
				)
			}
		case *asm.Binary:
			ins = append(ins, fix_arithmetic_instruction(t)...)
		case *asm.IntDiv:
			if t.Src.Type() == asm.IMM {
				ins = append(ins,
					&asm.Move{
						Src: t.Src,
						Dst: &REG_R10D,
					},
					&asm.IntDiv{
						Src: &REG_R10D,
					},
				)
			}
		default:
			util.Exit_with_printf("unknown asm instruction %v", t)
		}
	}

	return asm.Program{
		FuncDef: asm.Function{
			Name: program.FuncDef.Name,
			Ins:  ins,
		},
	}
}

func fix_arithmetic_instruction(src *asm.Binary) []asm.Instruction {
	ins := []asm.Instruction{}

	switch src.Operator {
	case asm.ADD, asm.SUB:
		if src.Src.Type() != asm.STACK || src.Dst.Type() != asm.STACK {
			ins = append(ins, src)
		} else {
			ins = append(ins,
				&asm.Move{
					Src: src.Src,
					Dst: &REG_R10D,
				}, &asm.Binary{
					Operator: src.Operator,
					Src:      &REG_R10D,
					Dst:      src.Dst,
				},
			)
		}
	case asm.MUL:
		if src.Dst.Type() != asm.STACK {
			ins = append(ins, src)
		} else {
			ins = append(ins,
				&asm.Move{
					Src: src.Dst,
					Dst: &REG_R11D,
				}, &asm.Binary{
					Operator: asm.MUL,
					Src:      src.Src,
					Dst:      &REG_R11D,
				}, &asm.Move{
					Src: &REG_R11D,
					Dst: src.Dst,
				},
			)
		}
	default:
		util.Exit_with_printf("unknown binary operator type %d", src.Operator)
	}

	return ins
}
