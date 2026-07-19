package parse

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/util"
)

type Parser struct {
	l    *lex.Lexer
	curr lex.Token
	next lex.Token
}

func New(l *lex.Lexer) *Parser {
	p := &Parser{l: l}

	p.adv()
	p.adv()

	return p
}

func (p *Parser) Parse() ast.Program {
	ast := ast.Program{
		FuncDef: p.parse_function_definition(),
	}
	p.expect(lex.EOF)

	return ast
}

func (p *Parser) parse_function_definition() ast.FunctionDefinition {
	p.exadv(lex.INT)

	p.expect(lex.IDENT)
	id := p.curr.Literal
	p.adv()

	p.exadv(lex.LPAREN)
	p.exadv(lex.VOID)
	p.exadv(lex.RPAREN)

	p.exadv(lex.LBRACE)

	items := []ast.BlockItem{}
	for !p.is(lex.RBRACE) {
		if item := p.ParseBlockItem(); item == nil {
			util.Exit_with_println("failed to parse function body\n")
		} else {
			items = append(items, *item)
		}
	}

	p.exadv(lex.RBRACE)

	return ast.FunctionDefinition{
		Name:  id,
		Items: items,
	}
}

func (p *Parser) peek() lex.Toktype {
	return p.curr.Type
}
func (p *Parser) is(ttype lex.Toktype) bool {
	return p.curr.Type == ttype
}
func (p *Parser) exadv(ttype lex.Toktype) {
	p.expect(ttype)
	p.adv()
}
func (p *Parser) adv() {
	p.curr = p.next
	p.next = p.l.Lex()
}
func (p *Parser) expect(ttype lex.Toktype) {
	if tok := p.curr; tok.Type != ttype {
		util.Exit_with_printf("expect %s got %s\n",
			ttype.String(), tok.Type.String())
	}
}
