package ast

import "cc260717/cmd/cc/util"

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

func (bi *BlockItem) String() string {
	switch bi.Type {
	case STMT:
		return bi.Stmt.String()
	case DECL:
		return bi.Decl.String()
	default:
		util.Exit_with_println("cannot string block item")
		return ""
	}
}
