package tackyparse

import (
	"strconv"

	"cc260717/cmd/cc/ir/tacky"
	"cc260717/cmd/cc/util"
)

type Parser struct {
	l      *Lexer
	curr   Token
	greedy bool
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, greedy: false}
	p.adv()
	return p
}

func (p *Parser) ParseAll() []tacky.Instruction {
	result := []tacky.Instruction{}

	p.greedy = true
	for {
		if ins := p.Parse(); ins != nil {
			result = append(result, ins)
		} else {
			break
		}
	}
	p.greedy = false

	return result
}

func (p *Parser) Parse() tacky.Instruction {
	switch p.peek() {
	case RETURN:
		p.adv()
		return p.parse_return()
	case UNARY:
		p.adv()
		return p.parse_unary()
	case BINARY:
		p.adv()
		return p.parse_binary()
	case LABEL:
		p.adv()
		return p.parse_label()
	case COPY:
		p.adv()
		return p.parse_copy()
	case JUMP:
		p.adv()
		return p.parse_jump()
	case JUMP_IF_ZERO:
		p.adv()
		return p.parse_jump_if_zero()
	case JUMP_IF_NOT_ZERO:
		p.adv()
		return p.parse_jump_if_not_zero()
	default:
		if !p.greedy {
			util.Exit_with_println("unknown tacky instruction type\n")
		}
	}

	return nil
}

func (p *Parser) parse_return() *tacky.Return {
	v := p.parse_value()
	p.exadv(SCOLON)

	return &tacky.Return{Value: v}
}
func (p *Parser) parse_unary() *tacky.Unary {
	var operator tacky.UnaryOpType

	switch p.peek() {
	case SUB:
		operator = tacky.NEGATE
	case TILDE:
		operator = tacky.COMPLEMENT
	case NOT:
		operator = tacky.NOT
	default:
		util.Exit_with_printf("unknown unary operator %v\n", p.curr.Type)
	}
	p.adv()

	src1 := p.parse_value()
	src2 := p.parse_value()

	if src1 == nil || src2 == nil {
		util.Exit_with_println("unary missing value 1 or value 2")
	}
	p.exadv(SCOLON)

	return &tacky.Unary{Operator: operator, Src: src1, Dst: src2}
}
func (p *Parser) parse_binary() *tacky.Binary {
	var op_map = util.MakeOpmap(map[Toktype]tacky.BinaryOpType{
		ADD:    tacky.ADD,
		SUB:    tacky.SUB,
		MUL:    tacky.MUL,
		DIV:    tacky.DIV,
		MOD:    tacky.REM,
		BOR:    tacky.BOR,
		BAND:   tacky.BAND,
		BXOR:   tacky.BXOR,
		LSHIFT: tacky.LSHIFT,
		RSHIFT: tacky.RSHIFT,
		LT:     tacky.LT,
		GT:     tacky.GT,
		LEQ:    tacky.LEQ,
		GEQ:    tacky.GEQ,
		EQ:     tacky.EQ,
		NEQ:    tacky.NEQ,
		AND:    tacky.AND,
		OR:     tacky.OR,
	})

	operator := op_map(p.peek())
	p.adv()

	src1 := p.parse_value()
	src2 := p.parse_value()
	dst := p.parse_value()
	if src1 == nil || src2 == nil || dst == nil {
		util.Exit_with_println("binary operands count mismatch")
	}
	p.exadv(SCOLON)

	return &tacky.Binary{
		Operator: operator,
		Src1:     src1,
		Src2:     src2,
		Dst:      dst,
	}
}
func (p *Parser) parse_label() *tacky.Label {
	p.expect(IDENT)
	name := p.curr.Literal
	p.adv()
	p.exadv(SCOLON)

	return &tacky.Label{Name: name}
}
func (p *Parser) parse_jump() *tacky.Jump {
	p.expect(IDENT)
	name := p.curr.Literal
	p.adv()
	p.exadv(SCOLON)

	return &tacky.Jump{Target: name}
}
func (p *Parser) parse_jump_if_zero() *tacky.JumpIfZero {
	condition := p.parse_value()
	if condition == nil {
		util.Exit_with_println("expect value for jump condition\n")
		return nil
	}

	p.expect(IDENT)
	target := p.curr.Literal
	p.adv()

	p.exadv(SCOLON)
	return &tacky.JumpIfZero{
		Condition: condition,
		Target:    target,
	}
}
func (p *Parser) parse_jump_if_not_zero() *tacky.JumpIfNotZero {
	condition := p.parse_value()
	if condition == nil {
		util.Exit_with_println("expect value for jump condition\n")
		return nil
	}

	p.expect(IDENT)
	target := p.curr.Literal
	p.adv()

	p.exadv(SCOLON)
	return &tacky.JumpIfNotZero{
		Condition: condition,
		Target:    target,
	}
}
func (p *Parser) parse_copy() *tacky.Copy {
	src := p.parse_value()
	dst := p.parse_value()

	if src == nil || dst == nil {
		util.Exit_with_println("copy missing value 1 or value 2")
		return nil
	}

	p.exadv(SCOLON)
	return &tacky.Copy{
		Src: src,
		Dst: dst,
	}
}

func (p *Parser) parse_value() tacky.Value {
	switch tok := p.curr; tok.Type {
	case IDENT:
		p.adv()
		return &tacky.Variable{Name: tok.Literal}
	case INTEGER:
		p.adv()
		if n, err := strconv.ParseInt(tok.Literal, 10, 64); err != nil {
			util.Exit_with_error(err)
		} else {
			return &tacky.Constant{Value: int(n)}
		}
	}

	// should be allowed to return nil as sometimes values are optional
	return nil
}

func (p *Parser) peek() Toktype {
	return p.curr.Type
}
func (p *Parser) adv() {
	p.curr = p.l.Lex()
}
func (p *Parser) expect(ttype Toktype) {
	if tok := p.curr; tok.Type != ttype {
		util.Exit_with_printf("expect %d got %d\n", ttype, tok.Type)
	}
}
func (p *Parser) exadv(ttype Toktype) {
	p.expect(ttype)
	p.adv()
}
