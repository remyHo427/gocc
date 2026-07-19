package check

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/util"
	"fmt"
)

type Checker struct {
	store     map[string]string
	make_name func(string) string
}

func New() Checker {
	return Checker{
		store:     map[string]string{},
		make_name: make_new_namer("var_"),
	}
}

func (c *Checker) check(tree ast.Program) ast.Program {
	result := []ast.BlockItem{}

	for _, item := range tree.FuncDef.Items {
		result = append(result, *c.resolve_block_item(item))
	}

	return tree
}

func (c *Checker) resolve_block_item(item ast.BlockItem) *ast.BlockItem {
	result := ast.BlockItem{
		Type: item.Type,
	}

	switch item.Type {
	case ast.DECL:
		result.Decl = c.resolve_decl(item.Decl)
	case ast.STMT:
		result.Stmt = c.resolve_stmt(item.Stmt)
	default:
		util.Exit_with_printf("unknown block item type\n")
		return nil
	}

	return &result
}

func (c *Checker) resolve_decl(decl ast.Decl) ast.Decl {
	switch t := decl.(type) {
	case *ast.Declaration:
		return c.resolve_declaration(*t)
	default:
		util.Exit_with_printf("unknown decl %v (%T)\n", t, t)
		return nil
	}
}
func (c *Checker) resolve_stmt(stmt ast.Stmt) ast.Stmt {
	switch t := stmt.(type) {
	case *ast.ExprStmt:
		return &ast.ExprStmt{Expr: c.resolve_expr(t.Expr)}
	case *ast.ReturnStmt:
		return &ast.ReturnStmt{Expr: c.resolve_expr(t.Expr)}
	case *ast.NullStmt:
		return &ast.NullStmt{}
	default:
		util.Exit_with_printf("unknown stmt %v (%T)\n", t, t)
		return nil
	}
}

func (c *Checker) resolve_declaration(decl ast.Declaration) *ast.Declaration {
	if _, ok := c.store[decl.Name]; ok {
		util.Exit_with_printf("variable %s is already declared!\n", decl.Name)
		return nil
	}

	result := ast.Declaration{}
	unique_name := c.make_name(decl.Name)
	c.store[decl.Name] = unique_name
	result.Name = unique_name

	if decl.Init == nil {
		return &result
	}
	result.Init = c.resolve_expr(decl.Init)

	return &result
}

func (c *Checker) resolve_expr(expr ast.Expr) ast.Expr {
	switch t := expr.(type) {
	case *ast.AssignmentExpr:
		if _, ok := t.Left.(*ast.VarExpr); !ok {
			util.Exit_with_printf("Invalid lvalue!, got %v (%T)",
				t.Left, t.Left)
			return nil
		} else {
			return &ast.AssignmentExpr{
				Left:  c.resolve_expr(t.Left),
				Right: c.resolve_expr(t.Right),
			}
		}
	case *ast.VarExpr:
		if v, ok := c.store[t.Name]; !ok {
			util.Exit_with_printf("Undeclared variable!, got %v (%T)",
				t.Name, t.Name)
			return nil
		} else {
			return &ast.VarExpr{Name: v}
		}
	case *ast.BinaryExpr:
		return &ast.BinaryExpr{
			Operator: t.Operator,
			Left:     c.resolve_expr(t.Left),
			Right:    c.resolve_expr(t.Right),
		}
	case *ast.UnaryExpr:
		return &ast.UnaryExpr{
			Operator: t.Operator,
			Expr:     c.resolve_expr(t.Expr),
		}
	default:
		util.Exit_with_printf("unknown expr %v (%T)", expr, expr)
		return nil
	}
}

func make_new_namer(prefix string) func(string) string {
	counter := 0
	return func(var_name string) string {
		name := fmt.Sprintf("%s%s_%d", prefix, var_name, counter)
		counter++
		return name
	}
}
