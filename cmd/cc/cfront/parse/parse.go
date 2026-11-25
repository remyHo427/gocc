package parse

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
)

type ParserError struct {
	File  string
	Line  int
	Error error
}
type Parser struct {
	curr lex.Token
	next lex.Token
	tp   int
	toks []lex.Token
	Err  []ParserError
}

func New(toks []lex.Token) *Parser {
	p := &Parser{toks: toks}

	p.adv()
	p.adv()

	return p
}

func (p *Parser) Parse() ast.Program {
	ast := ast.Program{
		FuncDef: *p.parse_function_definition(),
	}
	p.expect(lex.EOF)

	return ast
}

func (p *Parser) parse_function_definition() *ast.FunctionDefinition {
	p.exadv(lex.INT)

	p.expect(lex.IDENT)
	id := p.curr.Literal
	p.adv()

	p.exadv(lex.LPAREN)
	p.exadv(lex.VOID)
	p.exadv(lex.RPAREN)

	block := p.ParseBlock()
	if block == nil {
		return nil
	}

	return &ast.FunctionDefinition{
		Name:  id,
		Block: *block,
	}
}
func (p *Parser) ParseBlock() *ast.Block {
	p.exadv(lex.LBRACE)
	block_items := p.ParseBlockItems()
	p.exadv(lex.RBRACE)

	return &ast.Block{
		BlockItems: block_items,
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
	p.next = p.pop()
}
func (p *Parser) pop() lex.Token {
	if p.tp < len(p.toks) {
		tok := p.toks[p.tp]
		p.tp++
		return tok
	} else {
		return lex.EOFTOK
	}
}
func (p *Parser) expect(ttype lex.Toktype) {
	if tok := p.curr; tok.Type != ttype {
		e := ParserError{
			File:  tok.File,
			Line:  tok.Line,
			Error: ErrUnexpectedToken,
		}
		p.Err = append(p.Err, e)
		p.panic()
	}
}
func (p *Parser) panic() {
	for {
		if tok := p.curr; tok.Type == lex.SCOLON || tok.Type == lex.EOF {
			break
		}
		p.adv()
	}
}
