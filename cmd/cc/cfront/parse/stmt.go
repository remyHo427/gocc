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
	case lex.IF:
		p.adv()
		return p.parse_if_stmt()
	case lex.LBRACE:
		return p.parse_compound_stmt()
	case lex.BREAK:
		p.adv()
		return p.parse_break_stmt()
	case lex.CONTINUE:
		p.adv()
		return p.parse_continue_stmt()
	case lex.WHILE:
		p.adv()
		return p.parse_while_stmt()
	case lex.DO:
		p.adv()
		return p.parse_do_while_stmt()
	case lex.FOR:
		p.adv()
		return p.parse_for_stmt()
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

	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()

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

	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()

	return stmt
}

func (p *Parser) parse_if_stmt() *ast.IfStmt {
	if !p.expect(lex.LPAREN) {
		return nil
	}
	p.adv()

	cond := p.ParseExpr(LOWEST)
	if cond == nil {
		return nil
	}
	p.adv()

	if !p.expect(lex.RPAREN) {
		return nil
	}
	p.adv()

	then := p.ParseStmt()
	if then == nil {
		return nil
	}

	var Else ast.Stmt
	if p.is(lex.ELSE) {
		p.adv()
		if Else = p.ParseStmt(); Else == nil {
			return nil
		}
	}

	return &ast.IfStmt{
		Condition: cond,
		Then:      then,
		Else:      Else,
	}
}

func (p *Parser) parse_compound_stmt() *ast.CompoundStmt {
	if block := p.ParseBlock(); block == nil {
		return nil
	} else {
		return &ast.CompoundStmt{
			Block: *block,
		}
	}
}

func (p *Parser) parse_break_stmt() *ast.BreakStmt {
	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()
	return &ast.BreakStmt{}
}

func (p *Parser) parse_continue_stmt() *ast.ContinueStmt {
	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()
	return &ast.ContinueStmt{}
}

func (p *Parser) parse_while_stmt() *ast.WhileStmt {
	if !p.expect(lex.LPAREN) {
		return nil
	}
	p.adv()

	condition := p.ParseExpr(LOWEST)
	if condition == nil {
		return nil
	}
	p.adv()

	if !p.expect(lex.RPAREN) {
		return nil
	}
	p.adv()

	body := p.ParseStmt()
	if body == nil {
		return nil
	}

	return &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parse_do_while_stmt() *ast.DoWhileStmt {
	body := p.ParseStmt()
	if body == nil {
		return nil
	}

	if !p.expect(lex.WHILE) {
		return nil
	}
	p.adv()

	if !p.expect(lex.LPAREN) {
		return nil
	}
	p.adv()

	condition := p.ParseExpr(LOWEST)
	if condition == nil {
		return nil
	}
	p.adv()

	if !p.expect(lex.RPAREN) {
		return nil
	}
	p.adv()

	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()

	return &ast.DoWhileStmt{
		Body:      body,
		Condition: condition,
	}
}

func (p *Parser) parse_for_stmt() *ast.ForStmt {
	if !p.expect(lex.LPAREN) {
		return nil
	}
	p.adv()

	init := p.parse_for_init()
	if init == nil {
		return nil
	}

	condition, ok := p.parse_optional_expr_stmt()
	if !ok {
		return nil
	}

	post, ok := p.parse_for_post_expr()
	if !ok {
		return nil
	}

	body := p.ParseStmt()
	if body == nil {
		return nil
	}

	return &ast.ForStmt{
		Init:      *init,
		Condition: condition,
		Post:      post,
		Body:      body,
	}
}

func (p *Parser) parse_for_post_expr() (*ast.Expr, bool) {
	if p.is(lex.RPAREN) { // post expr is omitted
		p.adv()
		return nil, true
	}

	expr := p.ParseExpr(LOWEST)
	if expr == nil {
		return nil, false
	}
	p.adv()

	if !p.expect(lex.RPAREN) {
		return nil, false
	}
	p.adv()

	return &expr, true
}

func (p *Parser) parse_for_init() *ast.ForInit {
	if p.is(lex.SCOLON) {
		p.adv()
		return &ast.ForInit{
			Type: ast.FOR_INIT_NULL,
		}
	}

	decl := p.ParseDecl()
	if decl != nil {
		return &ast.ForInit{
			Type:     ast.FOR_INIT_DECL,
			InitDecl: decl,
		}
	}

	expr := p.ParseExpr(LOWEST)
	if expr == nil {
		return nil
	}
	p.adv()

	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()

	return &ast.ForInit{
		Type:     ast.FOR_INIT_EXPR,
		InitExpr: expr,
	}
}

func (p *Parser) parse_optional_expr_stmt() (*ast.Expr, bool) {
	if p.is(lex.SCOLON) {
		p.adv()
		return nil, true // is optional, parsed successfully so return true
	}

	expr := p.ParseExpr(LOWEST)
	if expr == nil {
		return nil, false // not optional, failed to parse
	}
	p.adv()

	if !p.expect(lex.SCOLON) {
		return nil, false // missing scolon, failed to parse
	}
	p.adv()

	return &expr, true // parsed successfully so return true
}
