package tacky

type Node interface {
	String() string
}

type Program struct {
	FuncDef Function
}

func (n *Program) String() string {
	return join("program", n.FuncDef)
}

type Function struct {
	Name string
	Ins  []Instruction
}

func (n *Function) String() string {
	return join("function", n.Name, n.Ins)
}
