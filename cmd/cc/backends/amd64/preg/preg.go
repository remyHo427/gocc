package preg

import (
	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/util"
)

type PseudoRegisterReplacer struct {
	store  map[string]asm.Stack
	step   int
	offset int
}

func New(step int) PseudoRegisterReplacer {
	return PseudoRegisterReplacer{
		store:  map[string]asm.Stack{},
		step:   step,
		offset: 0,
	}
}

func (pg *PseudoRegisterReplacer) Replace(program asm.Program) asm.Program {
	ins := []asm.Instruction{}

	for _, instruction := range program.FuncDef.Ins {
		switch t := instruction.(type) {
		case *asm.Cdq, *asm.Return, *asm.Jump, *asm.JumpCC, *asm.Label:
			ins = append(ins, t)
		case *asm.Move:
			ins = append(ins, &asm.Move{
				Src: pg.replace_operand(t.Src),
				Dst: pg.replace_operand(t.Dst),
			})
		case *asm.Unary:
			ins = append(ins, &asm.Unary{
				Operator: t.Operator,
				Value:    pg.replace_operand(t.Value),
			})
		case *asm.Binary:
			ins = append(ins, &asm.Binary{
				Operator: t.Operator,
				Src:      pg.replace_operand(t.Src),
				Dst:      pg.replace_operand(t.Dst),
			})
		case *asm.IntDiv:
			ins = append(ins, &asm.IntDiv{
				Src: pg.replace_operand(t.Src),
			})
		case *asm.Comparison:
			ins = append(ins, &asm.Comparison{
				Value1: pg.replace_operand(t.Value1),
				Value2: pg.replace_operand(t.Value2),
			})
		case *asm.SetCC:
			ins = append(ins, &asm.SetCC{
				Code:  t.Code,
				Value: pg.replace_operand(t.Value),
			})
		default:
			util.Exit_with_printf("unknown instruction type %v\n", t)
		}
	}

	return asm.Program{
		FuncDef: asm.Function{
			Name: program.FuncDef.Name,
			Ins:  ins,
		},
	}
}

func (p *PseudoRegisterReplacer) replace_operand(operand asm.Operand) asm.Operand {
	if operand.Type() != asm.PSEUDO_REG {
		return operand
	}

	pg := operand.(*asm.PseudoRegister)
	if stack, ok := p.store[pg.Name]; ok {
		return &stack
	} else {
		p.offset -= p.step
		new_stack := asm.Stack{Size: p.offset}
		p.store[pg.Name] = new_stack
		return &new_stack
	}
}

func (pg PseudoRegisterReplacer) GetCount() int {
	return pg.offset
}
