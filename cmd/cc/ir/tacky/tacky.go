package tacky

import (
	"fmt"

	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/util"
)

type TackyGenerator struct {
	ins       []Instruction
	new_name  func() string
	new_label func() string
}

func New() TackyGenerator {
	return TackyGenerator{
		ins:       []Instruction{},
		new_name:  make_name("tmp", 0),
		new_label: make_name("lb", 0),
	}
}

func (g *TackyGenerator) Generate(tree ast.Program) Program {
	return Program{
		FuncDef: g.generate_function(tree.FuncDef),
	}
}

func (g *TackyGenerator) generate_function(tree ast.FunctionDefinition) Function {
	g.clear()

	for _, item := range tree.Items {
		switch item.Type {
		case ast.DECL:
			g.generate_declaration(item.Decl)
		case ast.STMT:
			g.generate_instructions(item.Stmt)
		default:
			util.Exit_with_printf("unknown block item %v (%T)\n", item, item)
		}
	}

	return Function{Name: tree.Name, Ins: g.ins}
}

func (g *TackyGenerator) generate_declaration(decl ast.Decl) {
	switch t := decl.(type) {
	case *ast.Declaration:
		if t.Init != nil {
			v := Variable{Name: t.Name}
			rhs := g.generate_value(t.Init)
			g.push(&Copy{Src: rhs, Dst: &v})
		}
	default:
		util.Exit_with_printf("unknown decl %v (%T)\n", t, t)
	}
}

func (g *TackyGenerator) generate_instructions(body ast.Stmt) {
	switch t := body.(type) {
	case *ast.ReturnStmt:
		g.push(&Return{Value: g.generate_value(t.Expr)})
	case *ast.ExprStmt:
		g.generate_value(t.Expr)
	case *ast.NullStmt:
		// do nothing
	default:
		util.Exit_with_printf("unknown stmt %v (%T)\n", t, t)
	}
}

func (g *TackyGenerator) generate_value(expr ast.Expr) Value {
	switch t := expr.(type) {
	case *ast.ConstantExpr:
		return &Constant{Value: int(t.Value)}
	case *ast.UnaryExpr:
		return g.generate_unary_expr_value(*t)
	case *ast.BinaryExpr:
		return g.generate_binary_expr_value(*t)
	case *ast.AssignmentExpr:
		rhs := g.generate_value(t.Right)
		lhs := g.generate_value(t.Left)
		g.push(&Copy{Src: rhs, Dst: lhs})
		return lhs
	case *ast.VarExpr:
		return &Variable{Name: t.Name}
	case nil:
		// do nothing
		return nil
	default:
		util.Exit_with_printf("illegal state\n")
		return nil
	}
}

func (g *TackyGenerator) generate_unary_expr_value(expr ast.UnaryExpr) Value {
	op_map := util.MakeOpmap(map[lex.Toktype]UnaryOpType{
		lex.BCOMP: COMPLEMENT,
		lex.SUB:   NEGATE,
		lex.NOT:   NOT,
	})

	src := g.generate_value(expr.Expr)
	dst_name := g.new_name()
	dst := &Variable{dst_name}
	op := op_map(expr.Operator)
	g.push(&Unary{op, src, dst})

	return dst
}
func (g *TackyGenerator) generate_binary_expr_value(expr ast.BinaryExpr) Value {
	op_map := util.MakeOpmap(map[lex.Toktype]BinaryOpType{
		lex.ADD:    ADD,
		lex.SUB:    SUB,
		lex.MUL:    MUL,
		lex.DIV:    DIV,
		lex.MOD:    REM,
		lex.BAND:   BAND,
		lex.BOR:    BOR,
		lex.BXOR:   BXOR,
		lex.LSHIFT: LSHIFT,
		lex.RSHIFT: RSHIFT,
		lex.EQ:     EQ,
		lex.NEQ:    NEQ,
		lex.LT:     LT,
		lex.GT:     GT,
		lex.LEQ:    LEQ,
		lex.GEQ:    GEQ,
		lex.AND:    AND,
		lex.OR:     OR,
	})

	switch op := op_map(expr.Operator); op {
	case ADD, SUB, MUL, DIV, REM, BAND, BOR,
		BXOR, LSHIFT, RSHIFT, EQ, NEQ, LT, GT,
		LEQ, GEQ:
		return g.generate_binary(op, expr)
	case AND:
		return g.generate_logical_and(expr)
	case OR:
		return g.generate_logical_or(expr)
	default:
		util.Exit_with_printf("unknown tacky operator %v\n", op)
		return nil
	}
}

func (g *TackyGenerator) generate_binary(op BinaryOpType, expr ast.BinaryExpr) Value {
	v1 := g.generate_value(expr.Left)
	v2 := g.generate_value(expr.Right)
	dst := &Variable{g.new_name()}
	g.push(&Binary{op, v1, v2, dst})
	return dst
}
func (g *TackyGenerator) generate_logical_and(expr ast.BinaryExpr) Value {
	if_false_label := g.new_label()
	end_label := g.new_label()

	v1 := g.generate_value(expr.Left)
	g.push(&JumpIfZero{
		Condition: v1,
		Target:    if_false_label,
	})

	v2 := g.generate_value(expr.Right)
	g.push(&JumpIfZero{
		Condition: v2,
		Target:    if_false_label,
	})

	result := &Variable{Name: g.new_name()}
	g.push(&Copy{&Constant{1}, result})
	g.push(&Jump{Target: end_label})

	g.push(&Label{Name: if_false_label})
	g.push(&Copy{&Constant{0}, result})

	g.push(&Label{Name: end_label})

	return result
}
func (g *TackyGenerator) generate_logical_or(expr ast.BinaryExpr) Value {
	if_true_label := g.new_label()
	end_label := g.new_label()

	v1 := g.generate_value(expr.Left)
	g.push(&JumpIfNotZero{
		Condition: v1,
		Target:    if_true_label,
	})

	v2 := g.generate_value(expr.Right)
	g.push(&JumpIfNotZero{
		Condition: v2,
		Target:    if_true_label,
	})

	result := &Variable{Name: g.new_name()}
	g.push(&Copy{&Constant{0}, result})
	g.push(&Jump{Target: end_label})

	g.push(&Label{Name: if_true_label})
	g.push(&Copy{&Constant{1}, result})

	g.push(&Label{Name: end_label})

	return result
}

func (g *TackyGenerator) clear() {
	g.ins = []Instruction{}
}
func (g *TackyGenerator) push(ins Instruction) {
	g.ins = append(g.ins, ins)
}

func make_name(prefix string, startAt int) func() string {
	counter := startAt

	return func() string {
		s := fmt.Sprintf("%s.%d", prefix, counter)
		counter++
		return s
	}
}
