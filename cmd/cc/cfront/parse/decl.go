package parse

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
)

func (p *Parser) parse_block_item() *ast.BlockItem {
	item := ast.BlockItem{}

	switch p.peek() {
	case lex.INT:
		if decl := p.ParseDecl(); decl == nil {
			return nil
		} else {
			item.Decl = decl
		}
	default:
		if stmt := p.ParseStmt(); stmt == nil {
			return nil
		} else {
			item.Stmt = stmt
		}
	}

	return &item
}

func (p *Parser) ParseDecl() ast.Decl {
	decl := ast.Declaration{}

	// vtype := p.curr.Type
	p.adv()

	p.expect(lex.IDENT)
	decl.Name = p.curr.Literal
	p.adv()

	if p.is(lex.SCOLON) {
		p.adv()
		decl.Init = nil
		return &decl
	}

	p.exadv(lex.ASSIGN)
	init := p.ParseExpr(LOWEST)
	if init == nil {
		return nil
	}
	p.adv()
	p.exadv(lex.SCOLON)

	decl.Init = init
	return &decl
}
