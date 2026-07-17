package ast

type Node interface {
	String() string
}

type Program struct {
	FuncDef FunctionDefinition
}

func (n *Program) String() string {
	return n.FuncDef.String()
}

type FunctionDefinition struct {
	Name string
	Body Stmt
}

func (n *FunctionDefinition) String() string {
	return join(n.Name, n.Body)
}
