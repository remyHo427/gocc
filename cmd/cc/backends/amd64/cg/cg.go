package cg

import (
	"bytes"
	"fmt"
	"strings"

	"cc260717/cmd/cc/backends/amd64/asm"
	"cc260717/cmd/cc/util"
)

type CodeGenerator struct {
	buf bytes.Buffer
}

func New() CodeGenerator {
	return CodeGenerator{}
}

func (g *CodeGenerator) Generate(src asm.Program) string {
	g.emit_instruction(".globl", src.FuncDef.Name)
	g.emit_str(fmt.Sprintf("%s:", src.FuncDef.Name))

	// prologue
	g.emit_instruction("pushq", "%rbp")
	g.emit_instruction("movq", "%rsp", "%rbp")

	g.generate_instructions(src.FuncDef.Ins)

	// for linux only
	g.emit_str("\t.section\t.note.GNU-stack, \"\", @progbits")

	return g.buf.String()
}

var asm_unary_op_map = map[asm.UnaryOpType]string{
	asm.NEG: "negl",
	asm.NOT: "notl",
}
var asm_binary_op_map = map[asm.BinaryOpType]string{
	asm.ADD: "addl",
	asm.SUB: "subl",
	asm.MUL: "imull",
}
var asm_condcode_map = map[asm.Condcode]string{
	asm.E:  "e",
	asm.NE: "ne",
	asm.L:  "l",
	asm.LE: "le",
	asm.G:  "g",
	asm.GE: "ge",
}

func (g *CodeGenerator) generate_instructions(src []asm.Instruction) {
	for _, ins := range src {
		switch t := ins.(type) {
		case *asm.Cdq:
			g.emit_instruction("cdq")
		case *asm.Move:
			g.emit_instruction("movl",
				g.map_operand(t.Src, 4), g.map_operand(t.Dst, 4))
		case *asm.AllocateStack:
			g.emit_instruction("subq", fmt.Sprintf("$%d", t.Size), "%rsp")
		case *asm.Unary:
			s, ok := asm_unary_op_map[t.Operator]
			if !ok {
				util.Exit_with_printf(
					"assembly operator %v does not have a corresponding string",
					t.Operator)
			}
			g.emit_instruction(s, g.map_operand(t.Value, 4))
		case *asm.Binary:
			s, ok := asm_binary_op_map[t.Operator]
			if !ok {
				util.Exit_with_printf(
					"assembly operator %v does not have a corresponding string",
					t.Operator,
				)
			}
			g.emit_instruction(s, g.map_operand(t.Src, 4), g.map_operand(t.Dst, 4))
		case *asm.IntDiv:
			g.emit_instruction("idivl", g.map_operand(t.Src, 4))
		case *asm.Comparison:
			g.emit_instruction("cmpl",
				g.map_operand(t.Value1, 4), g.map_operand(t.Value2, 4))
		case *asm.Jump:
			g.emit_instruction("jmp", ".L"+t.Name)
		case *asm.JumpCC:
			cc_str, ok := asm_condcode_map[t.Code]
			if !ok {
				util.Exit_with_printf("unknown CondCode %d", t.Code)
			}
			g.emit_instruction("j"+cc_str, ".L"+t.Name)
		case *asm.SetCC:
			cc_str, ok := asm_condcode_map[t.Code]
			if !ok {
				util.Exit_with_printf("unknown CondCode %d", t.Code)
			}
			g.emit_instruction("set"+cc_str, g.map_operand(t.Value, 4))
		case *asm.Label:
			g.emit_str(".L" + t.Name + ":")
		case *asm.Return:
			// epilogue
			g.emit_instruction("movq", "%rbp", "%rsp")
			g.emit_instruction("popq", "%rbp")
			g.emit_instruction("ret")
		}
	}
}

func (g *CodeGenerator) map_operand(src asm.Operand, byte_count int) string {
	switch t := src.(type) {
	case *asm.Register:
		return g.map_register(t, byte_count)
	case *asm.Stack:
		return fmt.Sprintf("%d(%%rbp)", t.Size)
	case *asm.Immediate:
		return fmt.Sprintf("$%d", t.Value)
	default:
		util.Exit_with_printf("unknown operand %v\n", t)
	}
	return ""
}

var one_byte_rmap = map[asm.Reg]string{
	asm.AX:  "%al",
	asm.DX:  "%dl",
	asm.R10: "%r10b",
	asm.R11: "%r11b",
}
var four_bytes_rmap = map[asm.Reg]string{
	asm.AX:  "%eax",
	asm.DX:  "%edx",
	asm.R10: "%r10d",
	asm.R11: "%r11d",
}
var rmaps = []map[asm.Reg]string{nil, one_byte_rmap, nil, nil, four_bytes_rmap}

func (g *CodeGenerator) map_register(reg *asm.Register, byte_count int) string {
	rmap := rmaps[byte_count]
	if rmap == nil {
		util.Exit_with_printf("no register strmap for %d bytes\n", byte_count)
	}

	rstr, ok := rmap[reg.Reg]
	if !ok {
		util.Exit_with_printf("unknown register %d\n", reg.Reg)
	}

	return rstr
}

func (g *CodeGenerator) emit_instruction(ins string, operands ...string) {
	line := fmt.Sprintf("\t%s\t\t%s\n",
		ins, strings.Join(operands, ", "))

	g.buf.WriteString(line)
}
func (g *CodeGenerator) emit_str(s string) {
	g.buf.WriteString(s)
	g.buf.WriteString("\n")
}
