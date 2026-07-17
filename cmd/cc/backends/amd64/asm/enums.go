package asm

//go:generate go run golang.org/x/tools/cmd/stringer -type=UnaryOpType
type UnaryOpType int

const (
	NEG UnaryOpType = iota
	NOT
	LOGICAL_NOT
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=BinaryOpType
type BinaryOpType int

const (
	ADD BinaryOpType = iota
	SUB
	MUL
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=Condcode
type Condcode int

const (
	E  Condcode = iota // equal
	NE                 // not equal
	G                  // greater than
	GE                 // greater than or equal
	L                  // less than
	LE                 // less than or equal
)
