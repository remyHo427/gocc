package parse

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/util"
)

type Prec uint

const (
	LOWEST  Prec = iota
	COMMA        // , 							(left-to-right)
	ASSIGN       // = += -= etc.				(right-to-left)
	COND         // a ? 1 : 0					(left-to-right)
	OR           // ||							(left-to-right)
	AND          // &&							(left-to-right)
	BOR          // |							(left-to-right)
	BXOR         // ^							(left-to-right)
	BAND         // &							(left-to-right)
	EQ           // == !=						(left-to-right)
	ORDER        // < <= > >=					(left-to-right)
	SHIFT        // << >>						(left-to-right)
	SUM          // + -							(left-to-right)
	PRODUCT      // * / %						(left-to-right)
	PREFIX       // ++i --i +i -i !a ~a			(right-to-left)
	POSTFIX      // i++ i-- printf() a[] . ->	(left-to-right)
)

var prec = map[lex.Toktype]Prec{
	lex.INC:        POSTFIX,
	lex.DEC:        POSTFIX,
	lex.LPAREN:     POSTFIX,
	lex.LBRACKET:   POSTFIX,
	lex.DOT:        POSTFIX,
	lex.ARROW:      POSTFIX,
	lex.MUL:        PRODUCT,
	lex.DIV:        PRODUCT,
	lex.MOD:        PRODUCT,
	lex.ADD:        SUM,
	lex.SUB:        SUM,
	lex.LSHIFT:     SHIFT,
	lex.RSHIFT:     SHIFT,
	lex.GT:         ORDER,
	lex.LT:         ORDER,
	lex.GEQ:        ORDER,
	lex.LEQ:        ORDER,
	lex.EQ:         EQ,
	lex.NEQ:        EQ,
	lex.BAND:       BAND,
	lex.BXOR:       BXOR,
	lex.BOR:        BOR,
	lex.AND:        AND,
	lex.OR:         OR,
	lex.QMARK:      COND,
	lex.ASSIGN:     ASSIGN,
	lex.ADD_ASSIGN: ASSIGN,
	lex.SUB_ASSIGN: ASSIGN,
	lex.MUL_ASSIGN: ASSIGN,
	lex.DIV_ASSIGN: ASSIGN,
	lex.MOD_ASSIGN: ASSIGN,
	lex.RS_ASSIGN:  ASSIGN,
	lex.LS_ASSIGN:  ASSIGN,
	lex.BA_ASSIGN:  ASSIGN,
	lex.XO_ASSIGN:  ASSIGN,
	lex.BO_ASSIGN:  ASSIGN,
	lex.COMMA:      COMMA,
}

func (p *Parser) ParseExpr(currPrec Prec) ast.Expr {
	left := p.parsePrefix()

	if left == nil {
		return nil
	}

	for !p.is(lex.SCOLON) && currPrec < p.precn() {
		p.adv()

		switch p.peek() {
		case lex.ADD, lex.SUB, lex.MUL, lex.DIV, lex.MOD, lex.RSHIFT,
			lex.LSHIFT, lex.LT, lex.GT, lex.LEQ, lex.GEQ, lex.EQ,
			lex.NEQ, lex.BAND, lex.BXOR, lex.BOR, lex.AND, lex.OR,
			lex.DOT:
			if expr := p.parseInfixOperator(left); expr == nil {
				return nil
			} else {
				left = expr
			}
		case lex.ASSIGN:
			p.adv()
			if right := p.ParseExpr(p.prec()); right == nil {
				return nil
			} else {
				left = &ast.AssignmentExpr{
					Left:  left,
					Right: right,
				}
			}
		case lex.QMARK:
			p.adv()
			if expr := p.parseTernaryOperator(left); expr == nil {
				return nil
			} else {
				left = expr
			}
		default:
			util.Exit_with_printf("illegal state\n")
		}
	}

	return left
}

func (p *Parser) parseInfixOperator(left ast.Expr) *ast.BinaryExpr {
	expr := &ast.BinaryExpr{
		Operator: p.peek(),
		Left:     left,
	}

	currPrec := p.prec()
	p.adv()

	if right := p.ParseExpr(currPrec); right == nil {
		return nil
	} else {
		expr.Right = right
	}

	return expr
}

func (p *Parser) parseTernaryOperator(left ast.Expr) *ast.TernaryExpr {
	expr := &ast.TernaryExpr{
		Condition: left,
	}

	if Then := p.ParseExpr(LOWEST); Then == nil {
		return nil
	} else {
		expr.Then = Then
	}
	p.adv()

	p.exadv(lex.COLON)

	if Else := p.ParseExpr(p.prec()); Else == nil {
		return nil
	} else {
		expr.Else = Else
	}

	return expr
}

func (p *Parser) parsePrefix() ast.Expr {
	switch ttype := p.peek(); ttype {
	case lex.INT_CONST:
		return &ast.ConstantExpr{
			Value: p.curr.IntVal,
		}
	case lex.IDENT:
		return &ast.VarExpr{
			Name: p.curr.Literal,
		}
	case lex.DEC, lex.SUB, lex.BCOMP, lex.NOT:
		return p.parse_prefix_operator()
	case lex.LPAREN:
		p.adv()
		if inner := p.ParseExpr(LOWEST); inner == nil {
			return nil
		} else {
			p.adv()
			p.expect(lex.RPAREN)
			return inner
		}
	default:
		util.Exit_with_printf("unknown token type %s\n", ttype.String())
	}

	return nil
}

func (p *Parser) parse_prefix_operator() ast.Expr {
	expr := &ast.UnaryExpr{
		Operator: p.peek(),
	}
	p.adv()

	if right := p.ParseExpr(PREFIX); right == nil {
		return nil
	} else {
		expr.Expr = right
	}

	return expr
}

func (p *Parser) prec() Prec {
	return prec[p.curr.Type]
}
func (p *Parser) precn() Prec {
	return prec[p.next.Type]
}
