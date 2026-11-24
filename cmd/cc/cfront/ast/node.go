package ast

import "strings"

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
	Block Block
}

func (n *FunctionDefinition) String() string {
	return join(n.Name, n.Block)
}

type BlockItemType int

const (
	DECL BlockItemType = iota
	STMT
)

type BlockItem struct {
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
		return ""
	}
}

type Block struct {
	BlockItems []BlockItem
}

func (b *Block) String() string {
	items_strs := []string{}

	for _, item := range b.BlockItems {
		items_strs = append(items_strs, item.String())
	}

	return join("block", strings.Join(items_strs, " "))
}
