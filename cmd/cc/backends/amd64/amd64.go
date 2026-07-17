package amd64

import (
	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/backends/amd64/fix"
	"cc260717/cmd/cc/backends/amd64/preg"
	"cc260717/cmd/cc/backends/amd64/stack"
	"cc260717/cmd/cc/ir/tacky"
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

func ToAsm(src tacky.Program) asm.Program {
	pg := preg.New(4)

	program := asm.Program{
		FuncDef: ToAsmFunction(src.FuncDef),
	}
	program = pg.Replace(program)
	program = fix.FixInstructions(program)
	program = stack.AllocateStack(program, -pg.GetCount())

	return program
}

func ToAsmFunction(src tacky.Function) asm.Function {
	ins := []asm.Instruction{}

	for _, tackyInstruction := range src.Ins {
		ins = append(ins, ToAsmInstructions(tackyInstruction)...)
	}

	return asm.Function{
		Name: src.Name,
		Ins:  ins,
	}
}
func ToAsmInstructions(src tacky.Instruction) []asm.Instruction {
	switch t := src.(type) {
	case *tacky.Return:
		return []asm.Instruction{
			&asm.Move{Src: ToAsmOperand(t.Value), Dst: &REG_AX},
			&INS_RET,
		}
	case *tacky.Copy:
		return []asm.Instruction{
			&asm.Move{Src: ToAsmOperand(t.Src), Dst: ToAsmOperand(t.Dst)},
		}
	case *tacky.Label:
		return []asm.Instruction{
			&asm.Label{Name: t.Name},
		}
	case *tacky.Unary:
		return ToAsmUnaryInstruction(*t)
	case *tacky.Binary:
		return ToAsmBinaryInstruction(*t)
	case *tacky.Jump:
		return ToAsmJumpInstruction(*t)
	case *tacky.JumpIfZero:
		return ToAsmJumpIfZeroInstruction(*t)
	case *tacky.JumpIfNotZero:
		return ToAsmJumpIfNotZeroInstruction(*t)
	default:
		util.Exit_with_printf("unknown tacky instruction %v\n", t)
		return nil
	}
}
func ToAsmUnaryInstruction(src tacky.Unary) []asm.Instruction {
	var op_map = util.MakeOpmap(map[tacky.UnaryOpType]asm.UnaryOpType{
		tacky.NEGATE:     asm.NEG,
		tacky.COMPLEMENT: asm.NOT,
		tacky.NOT:        asm.LOGICAL_NOT,
	})

	switch op := op_map(src.Operator); op {
	case asm.NEG, asm.NOT:
		Dst := ToAsmOperand(src.Dst)
		return []asm.Instruction{
			&asm.Move{
				Src: ToAsmOperand(src.Src),
				Dst: Dst,
			},
			&asm.Unary{
				Operator: op,
				Value:    Dst,
			},
		}
	case asm.LOGICAL_NOT:
		return []asm.Instruction{
			&asm.Comparison{
				Value1: &asm.Immediate{Value: 0},
				Value2: ToAsmOperand(src.Src),
			},
			&asm.Move{
				Src: &asm.Immediate{Value: 0},
				Dst: ToAsmOperand(src.Dst),
			},
			&asm.SetCC{
				Code:  asm.E,
				Value: ToAsmOperand(src.Dst),
			},
		}
	default:
		util.Exit_with_printf("unknown unary operator\n")
		return nil
	}
}
func ToAsmBinaryInstruction(src tacky.Binary) []asm.Instruction {
	arithmetic_opmap := map[tacky.BinaryOpType]asm.BinaryOpType{
		tacky.ADD: asm.ADD,
		tacky.SUB: asm.SUB,
		tacky.MUL: asm.MUL,
	}
	if a_op, ok := arithmetic_opmap[src.Operator]; ok {
		dst := ToAsmOperand(src.Dst)
		return []asm.Instruction{
			&asm.Move{
				Src: ToAsmOperand(src.Src1),
				Dst: dst,
			},
			&asm.Binary{
				Operator: a_op,
				Src:      ToAsmOperand(src.Src2),
				Dst:      dst,
			},
		}
	}

	relational_opmap := map[tacky.BinaryOpType]asm.Condcode{
		tacky.EQ:  asm.E,
		tacky.NEQ: asm.NE,
		tacky.LT:  asm.L,
		tacky.LEQ: asm.LE,
		tacky.GT:  asm.G,
		tacky.GEQ: asm.GE,
	}
	if r_op, ok := relational_opmap[src.Operator]; ok {
		dst := ToAsmOperand(src.Dst)
		return []asm.Instruction{
			&asm.Comparison{
				Value1: ToAsmOperand(src.Src2),
				Value2: ToAsmOperand(src.Src1),
			},
			&asm.Move{
				Src: &asm.Immediate{Value: 0},
				Dst: dst,
			},
			&asm.SetCC{
				Code:  r_op,
				Value: dst,
			},
		}
	}

	switch src.Operator {
	case tacky.DIV:
		return []asm.Instruction{
			&asm.Move{
				Src: ToAsmOperand(src.Src1),
				Dst: &REG_AX,
			},
			&INS_CDQ,
			&asm.IntDiv{
				Src: ToAsmOperand(src.Src2),
			},
			&asm.Move{
				Src: &REG_AX,
				Dst: ToAsmOperand(src.Dst),
			},
		}
	case tacky.REM:
		return []asm.Instruction{
			&asm.Move{
				Src: ToAsmOperand(src.Src1),
				Dst: &REG_AX,
			},
			&INS_CDQ,
			&asm.IntDiv{
				Src: ToAsmOperand(src.Src2),
			},
			&asm.Move{
				Src: &REG_DX,
				Dst: ToAsmOperand(src.Dst),
			},
		}
	default:
		util.Exit_with_printf("unknown binary operator %v\n", src.Operator)
		return nil
	}
}
func ToAsmJumpInstruction(src tacky.Jump) []asm.Instruction {
	return []asm.Instruction{
		&asm.Jump{Name: src.Target},
	}
}
func ToAsmJumpIfZeroInstruction(src tacky.JumpIfZero) []asm.Instruction {
	return []asm.Instruction{
		&asm.Comparison{
			Value1: &asm.Immediate{Value: 0},
			Value2: ToAsmOperand(src.Condition),
		},
		&asm.JumpCC{
			Code: asm.E,
			Name: src.Target,
		},
	}
}
func ToAsmJumpIfNotZeroInstruction(src tacky.JumpIfNotZero) []asm.Instruction {
	return []asm.Instruction{
		&asm.Comparison{
			Value1: &asm.Immediate{Value: 0},
			Value2: ToAsmOperand(src.Condition),
		},
		&asm.JumpCC{
			Code: asm.NE,
			Name: src.Target,
		},
	}
}

func ToAsmOperand(src tacky.Value) asm.Operand {
	switch t := src.(type) {
	case *tacky.Constant:
		return &asm.Immediate{Value: int64(t.Value)}
	case *tacky.Variable:
		return &asm.PseudoRegister{Name: t.Name}
	default:
		util.Exit_with_printf("unknown tacky operand %v\n", t)
		return nil
	}
}
