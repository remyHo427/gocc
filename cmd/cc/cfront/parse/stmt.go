package parse

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
)

func (p *Parser) ParseStmt() ast.Stmt {
	switch p.curr.Type {
	case lex.RETURN:
		p.adv()
		return p.parse_return_stmt()
	default:
		return p.parse_expr_stmt()
	}
}

func (p *Parser) parse_return_stmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{}

	if p.is(lex.SCOLON) {
		p.adv()
		return stmt
	}

	if exp := p.ParseExpr(LOWEST); exp == nil {
		return nil
	} else {
		stmt.Expr = exp
	}
	p.adv()

	p.exadv(lex.SCOLON)
	return stmt
}

func (p *Parser) parse_expr_stmt() ast.Stmt {
	stmt := &ast.ExprStmt{}

	// can be a null stmt
	if p.is(lex.SCOLON) {
		p.adv()
		return &ast.NullStmt{}
	}

	if expr := p.ParseExpr(LOWEST); expr == nil {
		return nil
	} else {
		stmt.Expr = expr
	}
	p.adv()

	p.exadv(lex.SCOLON)
	return stmt
}
