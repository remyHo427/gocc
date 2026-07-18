package ast

import (
	"cc260717/cmd/cc/cfront/lex"
	"fmt"
)

type Expr interface {
	Node
	ExprNode()
}

type TestExpr struct{ Value int64 }
type ConstantExpr struct {
	Value int64
}

func (e *ConstantExpr) ExprNode() {}
func (e *ConstantExpr) String() string {
	return fmt.Sprintf("%d", e.Value)
}

type UnaryExpr struct {
	Operator lex.Toktype
	Expr     Expr
}

func (e *UnaryExpr) ExprNode() {}
func (e *UnaryExpr) String() string {
	return join("unary", e.Operator, e.Expr)
}

type BinaryExpr struct {
	Operator lex.Toktype
	Left     Expr
	Right    Expr
}

func (e *BinaryExpr) ExprNode() {}
func (e *BinaryExpr) String() string {
	return join("binary", e.Operator, e.Left, e.Right)
}

type VarExpr struct {
	Name string
}

func (e *VarExpr) ExprNode() {}
func (e *VarExpr) String() string {
	return join("var", e.Name)
}

type AssignmentExpr struct {
	Left  Expr
	Right Expr
}

func (e *AssignmentExpr) ExprNode() {}
func (e *AssignmentExpr) String() string {
	return join("assign", e.Left.String(), e.Right.String())
}
