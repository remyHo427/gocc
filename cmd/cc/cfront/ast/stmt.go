package ast

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
