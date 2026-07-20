package parse

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/util"
	"fmt"
)

type Parser struct {
	l      *lex.Lexer
	curr   lex.Token
	next   lex.Token
	Errors []error
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

	return ast.FunctionDefinition{
		Name:  id,
		Items: p.ParseBlockItems(),
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
	tok := p.l.Lex()
	if tok.Type != lex.ERR {
		p.curr = p.next
		p.next = tok
		return
	} else {
		p.tok_err(tok)
		p.adv()
	}
}
func (p *Parser) expect(ttype lex.Toktype) {
	if tok := p.curr; tok.Type != ttype {
		util.Exit_with_printf("expect %s got %s\n",
			ttype.String(), tok.Type.String())
	}
}
func (p *Parser) tok_err(tok lex.Token) {
	p.Errors = append(p.Errors, fmt.Errorf("%s:%d %s: %s\n",
		tok.File, tok.Line, tok.Error.Error(), tok.ErrorStr))
}
