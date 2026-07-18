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
	Name  string
	Items []BlockItem
}

func (n *FunctionDefinition) String() string {
	return join(n.Name, n.Items)
}

type BlockItemType int

const (
	DECL BlockItemType = iota
	STMT
)

type BlockItem struct {
	Node
	Type BlockItemType
	Stmt Stmt
	Decl Decl
}
