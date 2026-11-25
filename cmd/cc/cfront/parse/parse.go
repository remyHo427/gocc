package parse

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/util"
	"fmt"
	"os"
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
	var program ast.Program

	func_def := p.parse_function_definition()
	if func_def == nil {
		return program
	}
	program.FuncDef = *func_def

	if !p.expect(lex.EOF) {
		util.Exit_with_println("failed to parse program")
		os.Exit(1)
	}

	return program
}

func (p *Parser) parse_function_definition() *ast.FunctionDefinition {
	if !p.expect(lex.INT) {
		return nil
	}
	p.adv()

	if !p.expect(lex.IDENT) {
		return nil
	}
	id := p.curr.Literal
	p.adv()

	if !p.expect(lex.LPAREN) {
		return nil
	}
	p.adv()
	if !p.expect(lex.VOID) {
		return nil
	}
	p.adv()
	if !p.expect(lex.RPAREN) {
		return nil
	}
	p.adv()

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
	if !p.expect(lex.LBRACE) {
		return nil
	}
	p.adv()

	block_items := p.ParseBlockItems()
	if !p.expect(lex.RBRACE) {
		return nil
	}
	p.adv()

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
func (p *Parser) expect(ttype lex.Toktype) bool {
	if tok := p.curr; tok.Type != ttype {
		e := ParserError{
			File: tok.File,
			Line: tok.Line,
			Error: fmt.Errorf("%s \"%s\"",
				ERR_UNEXPECTED_TOKEN, tok.Type.String()),
		}
		p.Err = append(p.Err, e)
		return false
	}
	return true
}
