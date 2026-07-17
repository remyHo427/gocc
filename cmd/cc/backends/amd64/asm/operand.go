package asm

import (
	"fmt"

	"cc260717/cmd/cc/util"
)

type OperandType int

const (
	IMM OperandType = iota
	REG
	PSEUDO_REG // pseudo register
	STACK
)

type Operand interface {
	Type() OperandType
	Node
}

type Immediate struct {
	Value int64
}

func (o *Immediate) Type() OperandType {
	return IMM
}
func (o *Immediate) String() string {
	return fmt.Sprintf("$%d", o.Value)
}

type Reg int

const (
	AX Reg = iota
	DX
	R10
	R11
)

var register_strmap = map[Reg]string{
	AX:  "ax",
	DX:  "dx",
	R10: "r10",
	R11: "r11",
}

type Register struct {
	Reg Reg
}

func (o *Register) Type() OperandType {
	return REG
}
func (o *Register) String() string {
	s, ok := register_strmap[o.Reg]
	if !ok {
		util.Exit_with_printf("unknown register %v", o.Reg)
	}
	return s
}

type PseudoRegister struct {
	Name string
}

func (o *PseudoRegister) Type() OperandType {
	return PSEUDO_REG
}
func (o *PseudoRegister) String() string {
	return o.Name
}

type Stack struct {
	Size int
}

func (o *Stack) Type() OperandType {
	return STACK
}
func (o *Stack) String() string {
	return fmt.Sprintf("s%d", o.Size)
}
