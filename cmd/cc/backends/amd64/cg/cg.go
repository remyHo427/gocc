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

func (g *CodeGenerator) generate_instructions(src []asm.Instruction) {
	for _, ins := range src {
		switch t := ins.(type) {
		case *asm.Cdq:
			g.emit_instruction("cdq")
		case *asm.Move:
			g.emit_instruction("movl",
				g.map_operand(t.Src), g.map_operand(t.Dst))
		case *asm.AllocateStack:
			g.emit_instruction("subq", fmt.Sprintf("$%d", t.Size), "%rsp")
		case *asm.Unary:
			s, ok := asm_unary_op_map[t.Operator]
			if !ok {
				util.Exit_with_printf(
					"assembly operator %v does not have a corresponding string",
					t.Operator)
			}
			g.emit_instruction(s, g.map_operand(t.Value))
		case *asm.Binary:
			s, ok := asm_binary_op_map[t.Operator]
			if !ok {
				util.Exit_with_printf(
					"assembly operator %v does not have a corresponding string",
					t.Operator,
				)
			}
			g.emit_instruction(s, g.map_operand(t.Src), g.map_operand(t.Dst))
		case *asm.IntDiv:
			g.emit_instruction("idivl", g.map_operand(t.Src))
		case *asm.Return:
			// epilogue
			g.emit_instruction("movq", "%rbp", "%rsp")
			g.emit_instruction("popq", "%rbp")
			g.emit_instruction("ret")
		}
	}
}

var register_map = map[asm.Reg]string{
	asm.AX:  "%eax",
	asm.DX:  "%edx",
	asm.R10: "%r10d",
	asm.R11: "%r11d",
}

func (g *CodeGenerator) map_operand(src asm.Operand) string {
	switch t := src.(type) {
	case *asm.Register:
		if s, ok := register_map[t.Reg]; ok {
			return s
		}
		util.Exit_with_printf(
			"register %v does not have a corresponding string\n", t.Reg)
	case *asm.Stack:
		return fmt.Sprintf("%d(%%rbp)", t.Size)
	case *asm.Immediate:
		return fmt.Sprintf("$%d", t.Value)
	default:
		util.Exit_with_printf("unknown operand %v\n", t)
	}

	return ""
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
