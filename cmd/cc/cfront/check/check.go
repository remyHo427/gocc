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
	make_temporary func(string) string
}

// make_temporary needs to be reset between calls to Check()
// make a Checker struct and attach methods to it
var make_temporary = namer()

func New() Checker {
	return Checker{
		make_temporary: namer(),
	}
}

func (c *Checker) Check(tree ast.Program) ast.Program {
	result := []ast.BlockItem{}
	store := map[string]Symbol{}

	for _, item := range tree.FuncDef.Block.Blocks {
		result = append(result, *c.resolve_block_item(item, store))
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

func (c *Checker) resolve_block_item(item ast.BlockItem, s Store) *ast.BlockItem {
	result := ast.BlockItem{
		Type: item.Type,
	}

	switch item.Type {
	case ast.DECL:
		result.Decl = c.resolve_decl(item.Decl, s)
	case ast.STMT:
		result.Stmt = c.resolve_stmt(item.Stmt, s)
	default:
		util.Exit_with_printf("unknown block item type\n")
		return nil
	}

	return &result
}

func (c *Checker) resolve_decl(decl ast.Decl, s Store) ast.Decl {
	switch t := decl.(type) {
	case *ast.Declaration:
		return c.resolve_declaration(*t, s)
	default:
		util.Exit_with_printf("unknown decl %v (%T)\n", t, t)
		return nil
	}
}

func (c *Checker) resolve_declaration(decl ast.Declaration, s Store) *ast.Declaration {
	sym, ok := s[decl.Name]
	if ok && sym.fromCurrentBlock {
		util.Exit_with_printf("variable %s is already declared!\n", decl.Name)
		return nil
	}

	result := ast.Declaration{}
	unique_name := c.make_temporary(decl.Name)
	s[decl.Name] = Symbol{
		name:             unique_name,
		fromCurrentBlock: true,
	}
	result.Name = unique_name

	if decl.Init == nil {
		return &result
	}
	result.Init = c.resolve_expr(decl.Init, s)

	return &result
}

func (c *Checker) resolve_stmt(stmt ast.Stmt, s Store) ast.Stmt {
	switch t := stmt.(type) {
	case *ast.ExprStmt:
		return &ast.ExprStmt{Expr: c.resolve_expr(t.Expr, s)}
	case *ast.ReturnStmt:
		return &ast.ReturnStmt{Expr: c.resolve_expr(t.Expr, s)}
	case *ast.IfStmt:
		return &ast.IfStmt{
			Condition: c.resolve_expr(t.Condition, s),
			Then:      c.resolve_stmt(t.Then, s),
			Else:      c.resolve_stmt(t.Else, s),
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

func (c *Checker) resolve_expr(expr ast.Expr, s Store) ast.Expr {
	switch t := expr.(type) {
	case *ast.ConstantExpr:
		return t
	case *ast.AssignmentExpr:
		if _, ok := t.Left.(*ast.VarExpr); !ok {
			util.Exit_with_printf("Invalid lvalue!, got %v (%T)\n", t, t)
			return nil
		} else {
			return &ast.AssignmentExpr{
				Left:  c.resolve_expr(t.Left, s),
				Right: c.resolve_expr(t.Right, s),
			}
		}
	case *ast.VarExpr:
		sym, ok := s[t.Name]
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
			Left:     c.resolve_expr(t.Left, s),
			Right:    c.resolve_expr(t.Right, s),
		}
	case *ast.UnaryExpr:
		return &ast.UnaryExpr{
			Operator: t.Operator,
			Expr:     c.resolve_expr(t.Expr, s),
		}
	case *ast.TernaryExpr:
		return &ast.TernaryExpr{
			Condition: c.resolve_expr(t.Condition, s),
			Then:      c.resolve_expr(t.Then, s),
			Else:      c.resolve_expr(t.Else, s),
		}
	case nil:
		// do nothing
		return nil
	default:
		util.Exit_with_printf("unknown expr %v (%T)\n", expr, expr)
		return nil
	}
}

func namer() func(string) string {
	counter := 0
	return func(var_name string) string {
		name := fmt.Sprintf("%s.%d", var_name, counter)
		counter++
		return name
	}
}
func copy_store(s Store) Store {
	result := Store{}

	for k, v := range s {
		result[k] = Symbol{
			fromCurrentBlock: false,
			name:             v.name,
		}
	}

	return result
}
