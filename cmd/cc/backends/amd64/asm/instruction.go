package asm

import (
	"fmt"
)

type InsType int

const (
	MOV InsType = iota
	UNARY
	BINARY
	CMP
	IDIV
	CDQ
	JMP
	JMP_CC
	SET_CC
	LABEL
	ALLOCATE_STACK
	RET
)

type Node interface {
	String() string
}

type Instruction interface {
	Type() InsType
	Node
}

type Move struct {
	Src Operand
	Dst Operand
}

func (i *Move) Type() InsType {
	return MOV
}
func (i *Move) String() string {
	return join("move", i.Src, i.Dst)
}

type Unary struct {
	Operator UnaryOpType
	Value    Operand
}

func (i *Unary) Type() InsType {
	return UNARY
}
func (i *Unary) String() string {
	return join("unary", i.Operator, i.Value)
}

type Binary struct {
	Operator BinaryOpType
	Src      Operand
	Dst      Operand
}

func (i *Binary) Type() InsType {
	return BINARY
}
func (i *Binary) String() string {
	return join("binary", i.Operator, i.Src, i.Dst)
}

type Comparison struct {
	Value1 Operand
	Value2 Operand
}

func (i *Comparison) Type() InsType {
	return BINARY
}
func (i *Comparison) String() string {
	return join("comparison", i.Value1, i.Value2)
}

type IntDiv struct {
	Src Operand
}

func (i *IntDiv) Type() InsType {
	return IDIV
}
func (i *IntDiv) String() string {
	return join("intdiv", i.Src)
}

type Cdq struct{}

func (i *Cdq) Type() InsType {
	return CDQ
}
func (i *Cdq) String() string {
	return "cdq"
}

type Jump struct {
	Name string
}

func (i *Jump) Type() InsType {
	return JMP
}
func (i *Jump) String() string {
	return join("jump", i.Name)
}

type JumpCC struct {
	Code Condcode
	Name string
}

func (i *JumpCC) Type() InsType {
	return JMP_CC
}
func (i *JumpCC) String() string {
	return join("jump_cc", i.Code, i.Name)
}

type SetCC struct {
	Code  Condcode
	Value Operand
}

func (i *SetCC) Type() InsType {
	return SET_CC
}
func (i *SetCC) String() string {
	return join("set_cc", i.Code, i.Value)
}

type Label struct {
	Name string
}

func (i *Label) Type() InsType {
	return LABEL
}
func (i *Label) String() string {
	return join("label", i.Name)
}

type AllocateStack struct {
	Size int
}

func (i *AllocateStack) Type() InsType {
	return ALLOCATE_STACK
}
func (i *AllocateStack) String() string {
	return join("allocate_stack", fmt.Sprintf("%d", i.Size))
}

type Return struct{}

func (i *Return) Type() InsType {
	return RET
}
func (i *Return) String() string {
	return "return"
}
