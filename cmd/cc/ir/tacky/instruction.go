package tacky

type InsType int

const (
	RETURN InsType = iota
	UNARY
	BINARY
	COPY
	JUMP
	JUMP_IF_ZERO
	JUMP_IF_NOT_ZERO
	LABEL
)

type Instruction interface {
	Type() InsType
	Node
}

type Return struct {
	Value Value
}

func (i *Return) Type() InsType {
	return RETURN
}
func (i *Return) String() string {
	return join("return", i.Value)
}

type UnaryOpType int
type Unary struct {
	Operator UnaryOpType
	Src      Value
	Dst      Value
}

const (
	NEGATE UnaryOpType = iota
	COMPLEMENT
	NOT
)

func (i *Unary) Type() InsType {
	return UNARY
}
func (i *Unary) String() string {
	return join("unary", i.Operator, i.Src, i.Dst)
}

type BinaryOpType int
type Binary struct {
	Operator BinaryOpType
	Src1     Value
	Src2     Value
	Dst      Value
}

const (
	ADD BinaryOpType = iota
	SUB
	MUL
	DIV
	REM
	BAND
	BOR
	BXOR
	LSHIFT
	RSHIFT
	EQ
	NEQ
	LT
	LEQ
	GT
	GEQ
	AND
	OR
)

func (i *Binary) Type() InsType {
	return BINARY
}
func (i *Binary) String() string {
	return join("binary", i.Operator, i.Src1, i.Src2, i.Dst)
}

type Copy struct {
	Src Value
	Dst Value
}

func (i *Copy) Type() InsType {
	return COPY
}
func (i *Copy) String() string {
	return join("copy", i.Src, i.Dst)
}

type Jump struct {
	Target string
}

func (i *Jump) Type() InsType {
	return JUMP
}
func (i *Jump) String() string {
	return join("jump", i.Target)
}

type JumpIfZero struct {
	Condition Value
	Target    string
}

func (i *JumpIfZero) Type() InsType {
	return JUMP_IF_ZERO
}
func (i *JumpIfZero) String() string {
	return join("jump_if_zero", i.Condition, i.Target)
}

type JumpIfNotZero struct {
	Condition Value
	Target    string
}

func (i *JumpIfNotZero) Type() InsType {
	return JUMP_IF_NOT_ZERO
}
func (i *JumpIfNotZero) String() string {
	return join("jump_if_not_zero", i.Condition, i.Target)
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
