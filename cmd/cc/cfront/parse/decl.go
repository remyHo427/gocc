package parse

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
)

func (p *Parser) ParseBlockItem() *ast.BlockItem {
	item := ast.BlockItem{}

	switch p.peek() {
	case lex.INT:
		if decl := p.ParseDecl(); decl != nil {
			item.Type = ast.DECL
			item.Decl = decl
			return &item
		} else {
			return nil
		}
	default:
		if stmt := p.ParseStmt(); stmt != nil {
			item.Type = ast.STMT
			item.Stmt = stmt
			return &item
		} else {
			return nil
		}
	}
}
func (p *Parser) ParseBlockItems() []ast.BlockItem {
	items := []ast.BlockItem{}

	for !p.is(lex.RBRACE) {
		if item := p.ParseBlockItem(); item != nil {
			items = append(items, *item)
		} else {
			break
		}
	}

	return items
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
