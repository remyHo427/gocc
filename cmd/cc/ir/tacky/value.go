package tacky

import "fmt"

type ValType int

const (
	CONST ValType = iota
	VAR
)

type Value interface {
	Type() ValType
	Node
}

type Constant struct {
	Value int
}

func (v *Constant) Type() ValType {
	return CONST
}
func (v *Constant) String() string {
	return fmt.Sprintf("%d", v.Value)
}

type Variable struct {
	Name string
}

func (v *Variable) Type() ValType {
	return VAR
}

func (v *Variable) String() string {
	return v.Name
}
