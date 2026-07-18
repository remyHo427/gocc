package ast

type Decl interface {
	Node
	DeclNode()
}

type Declaration struct {
	Name string
	Init Expr
}

func (d *Declaration) DeclNode() {}
func (d *Declaration) String() string {
	return join("declaration", d.Name, d.Init)
}
