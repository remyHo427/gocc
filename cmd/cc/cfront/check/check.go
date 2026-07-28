package check

import (
	"cc260717/cmd/cc/cfront/ast"
	"cc260717/cmd/cc/util"
	"fmt"
)

type Store map[string]Symbol
type Symbol struct {
	name             string
	fromCurrentBlock bool
}
type Checker struct {
	store     Store
	make_name func(string) string
}

func New() Checker {
	return Checker{
		store:     map[string]Symbol{},
		make_name: make_namer(),
	}
}

func (c *Checker) Check(tree ast.Program) ast.Program {
	result := []ast.BlockItem{}

	for _, item := range tree.FuncDef.Block.Blocks {
		result = append(result, *c.resolve_block_item(item))
	}

	return ast.Program{
		FuncDef: ast.FunctionDefinition{
			Name: tree.FuncDef.Name,
			Block: ast.Block{
				Blocks: result,
			},
		},
	}
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
	case *ast.IfStmt:
		return &ast.IfStmt{
			Condition: c.resolve_expr(t.Condition),
			Then:      c.resolve_stmt(t.Then),
			Else:      c.resolve_stmt(t.Else),
		}
	case *ast.CompoundStmt:
		// new_store = copy_store(c.store)
		// return &ast.CompoundStmt{
		// 	Block: c.resolve_block_item(),
		// }
		return nil
	case *ast.NullStmt:
		return &ast.NullStmt{}
	case nil:
		// do nothing
		return nil
	default:
		util.Exit_with_printf("unknown stmt %v (%T)\n", t, t)
		return nil
	}
}

func (c *Checker) resolve_declaration(decl ast.Declaration) *ast.Declaration {
	sym, ok := c.store[decl.Name]
	if ok && sym.fromCurrentBlock {
		util.Exit_with_printf("variable %s is already declared!\n", decl.Name)
		return nil
	}

	result := ast.Declaration{}
	unique_name := c.make_name(decl.Name)
	c.store[decl.Name] = Symbol{
		name:             unique_name,
		fromCurrentBlock: true,
	}
	result.Name = unique_name

	if decl.Init == nil {
		return &result
	}
	result.Init = c.resolve_expr(decl.Init)

	return &result
}

func (c *Checker) resolve_expr(expr ast.Expr) ast.Expr {
	switch t := expr.(type) {
	case *ast.ConstantExpr:
		return t
	case *ast.AssignmentExpr:
		if _, ok := t.Left.(*ast.VarExpr); !ok {
			util.Exit_with_printf("Invalid lvalue!, got %v (%T)\n",
				t.Left, t.Left)
			return nil
		} else {
			return &ast.AssignmentExpr{
				Left:  c.resolve_expr(t.Left),
				Right: c.resolve_expr(t.Right),
			}
		}
	case *ast.VarExpr:
		sym, ok := c.store[t.Name]
		if !ok {
			util.Exit_with_printf("Undeclared variable! got %v (%T)\n",
				t.Name, t.Name)
			return nil
		} else {
			return &ast.VarExpr{Name: sym.name}
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
	case *ast.TernaryExpr:
		return &ast.TernaryExpr{
			Condition: c.resolve_expr(t.Condition),
			Then:      c.resolve_expr(t.Then),
			Else:      c.resolve_expr(t.Else),
		}
	case nil:
		// do nothing
		return nil
	default:
		util.Exit_with_printf("unknown expr %v (%T)\n", expr, expr)
		return nil
	}
}

func make_namer() func(string) string {
	counter := 0
	return func(var_name string) string {
		name := fmt.Sprintf("%s.%d", var_name, counter)
		counter++
		return name
	}
}

func copy_store(store Store) Store {
	result := Store{}

	for k, v := range store {
		result[k] = Symbol{
			fromCurrentBlock: false,
			name:             v.name,
		}
	}

	return result
}
