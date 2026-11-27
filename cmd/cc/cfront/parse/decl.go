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

	for !p.is(lex.EOF) && !p.is(lex.RBRACE) {
		if b := p.ParseBlockItem(); b == nil {
			for !p.is(lex.SCOLON) && !p.is(lex.EOF) {
				p.adv()
			}
			p.adv()
		} else {
			items = append(items, *b)
		}
	}

	return items
}

func (p *Parser) ParseDecl() ast.Decl {
	switch p.curr.Type {
	case lex.INT:
		return p.parse_simple_decl()
	default:
		return nil
	}
}
func (p *Parser) parse_simple_decl() *ast.Declaration {
	decl := ast.Declaration{}

	// vtype := p.curr.Type
	p.adv()

	if !p.expect(lex.IDENT) {
		return nil
	}
	decl.Name = p.curr.Literal
	p.adv()

	if p.is(lex.SCOLON) {
		p.adv()
		decl.Init = nil
		return &decl
	}

	if !p.expect(lex.ASSIGN) {
		return nil
	}
	p.adv()

	init := p.ParseExpr(LOWEST)
	if init == nil {
		return nil
	}
	p.adv()

	if !p.expect(lex.SCOLON) {
		return nil
	}
	p.adv()

	decl.Init = init
	return &decl
}
