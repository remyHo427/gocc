package ast

import "cc260717/cmd/cc/util"

type Stmt interface {
	Node
	StmtNode()
}

type ReturnStmt struct {
	Expr Expr
}

func (s *ReturnStmt) StmtNode() {}
func (s *ReturnStmt) String() string {
	return join("return", s.Expr)
}

type ExprStmt struct {
	Expr Expr
}

func (s *ExprStmt) StmtNode() {}
func (s *ExprStmt) String() string {
	return s.Expr.String()
}

type NullStmt struct{}

func (s *NullStmt) StmtNode() {}
func (s *NullStmt) String() string {
	return join("null")
}

type IfStmt struct {
	Condition Expr
	Then      Stmt
	Else      Stmt
}

func (s *IfStmt) StmtNode() {}
func (s *IfStmt) String() string {
	return join("if", s.Condition, s.Then, s.Else)
}

type CompoundStmt struct {
	Block Block
}

func (s *CompoundStmt) StmtNode() {}
func (s *CompoundStmt) String() string {
	return join("compound", s.Block.String())
}

type BreakStmt struct {
	Label string
}

func (s *BreakStmt) StmtNode() {}
func (s *BreakStmt) String() string {
	if s.Label != "" {
		return join("break", s.Label)
	} else {
		return join("break")
	}
}

type ContinueStmt struct {
	Label string
}

func (s *ContinueStmt) StmtNode() {}
func (s *ContinueStmt) String() string {
	if s.Label != "" {
		return join("continue", s.Label)
	} else {
		return join("continue")
	}
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (s *WhileStmt) StmtNode() {}
func (s *WhileStmt) String() string {
	return join("while", s.Condition, s.Body)
}

type DoWhileStmt struct {
	Body      Stmt
	Condition Expr
}

func (s *DoWhileStmt) StmtNode() {}
func (s *DoWhileStmt) String() string {
	return join("do-while", s.Body, s.Condition)
}

type ForStmt struct {
	Init      ForInit
	Condition *Expr
	Post      *Expr
	Body      Stmt
}

func (s *ForStmt) StmtNode() {}
func (s *ForStmt) String() string {
	var cond, post string

	if s.Condition == nil {
		cond = "(null)"
	} else {
		cond = (*s.Condition).String()
	}
	if s.Post == nil {
		post = "(null)"
	} else {
		post = (*s.Post).String()
	}

	return join("for", s.Init.String(), cond, post, s.Body)
}

type ForInitType int

const (
	FOR_INIT_DECL ForInitType = iota
	FOR_INIT_EXPR
	FOR_INIT_NULL
)

type ForInit struct {
	Type     ForInitType
	InitDecl Decl
	InitExpr Expr
}

func (s *ForInit) String() string {
	switch s.Type {
	case FOR_INIT_DECL:
		return s.InitDecl.String()
	case FOR_INIT_EXPR:
		return s.InitExpr.String()
	case FOR_INIT_NULL:
		return "(null)"
	default:
		util.Exit_with_println("unknown ForInit type")
		return ""
	}
}
