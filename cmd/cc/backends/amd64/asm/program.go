package asm

type Program struct {
	FuncDef Function
}

type Function struct {
	Name string
	Ins  []Instruction
}
