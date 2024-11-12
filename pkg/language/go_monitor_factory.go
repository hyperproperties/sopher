package language

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type GoMonitorFactory struct {
	packageName string
	modelName   string
	offset      int
	variables   []string
}

func NewGoMonitorFactory(packageName, modelName string) GoMonitorFactory {
	return GoMonitorFactory{
		packageName: packageName,
		modelName:   modelName,
		offset:      0,
		variables:   make([]string, 0),
	}
}

func (factory *GoMonitorFactory) Create(node Node) *ast.CallExpr {
	switch cast := node.(type) {
	case GoExpresion:
		return factory.NewPredicateMonitorCall(cast)
	case Universal:
		return factory.NewUniversalMonitorCall(cast)
	case Existential:
		return factory.NewExistentialMonitorCall(cast)
	case Guarantee:
		factory.variables = nil
		factory.offset = 0
		return factory.Create(cast.assertion)
	case Assumption:
		factory.variables = nil
		factory.offset = 0
		return factory.Create(cast.assertion)
	}
	factory.offset = 0
	panic(fmt.Sprintf("unknown node type %t", node))
}

func (factory *GoMonitorFactory) NewPredicateMonitorCall(expression GoExpresion) *ast.CallExpr {
	// FIXME: Can accidentally define variables not in use. First we have to see what variables are in use and only define those.

	var body []ast.Stmt

	if len(factory.variables) > 0 {
		// e0, e1, e2 := assignments[0], assignments[1], assignments[2]
		var lhs, rhs []ast.Expr
		for idx, identifier := range factory.variables {
			lhs = append(lhs, ast.NewIdent(identifier))
			rhs = append(rhs, &ast.IndexExpr{
				X:     ast.NewIdent("assignments"),
				Index: &ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", idx)},
			})
		}

		executionAssignment := &ast.AssignStmt{
			Lhs: lhs, Tok: token.DEFINE, Rhs: rhs,
		}

		body = append(body, executionAssignment)

		// anonymous assignments for each execution variable.
		anonymousAssignment := &ast.AssignStmt{
			Tok: token.ASSIGN, Rhs: lhs,
		}
		for idx := 0; idx < len(lhs); idx++ {
			anonymousAssignment.Lhs = append(anonymousAssignment.Lhs, ast.NewIdent("_"))
		}
		body = append(body, anonymousAssignment)
	}

	code, err := parser.ParseExpr(expression.code)
	if err != nil {
		panic(err)
	}

	body = append(body, &ast.ReturnStmt{
		Results: []ast.Expr{
			code,
		},
	})

	predicate := &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("assignments")},
						Type:  &ast.ArrayType{Elt: ast.NewIdent(factory.modelName)},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: ast.NewIdent("bool")},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}

	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent(factory.packageName),
			Sel: ast.NewIdent("NewPredicateMonitor"),
		},
		Args: []ast.Expr{predicate},
	}
}

func (factory *GoMonitorFactory) NewUniversalMonitorCall(universal Universal) *ast.CallExpr {
	offset := factory.offset
	factory.offset += len(universal.variables)
	factory.variables = append(factory.variables, universal.variables...)
	call := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent(factory.packageName),
			Sel: ast.NewIdent("NewUniversalMonitor"),
		},
		Args: []ast.Expr{
			&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", offset)},
			&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", len(universal.variables))},
			factory.Create(universal.assertion),
		},
	}
	return call
}

func (factory *GoMonitorFactory) NewExistentialMonitorCall(existential Existential) *ast.CallExpr {
	offset := factory.offset
	factory.offset += len(existential.variables)
	factory.variables = append(factory.variables, existential.variables...)
	call := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent(factory.packageName),
			Sel: ast.NewIdent("NewExistentialMonitor"),
		},
		Args: []ast.Expr{
			&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", offset)},
			&ast.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", len(existential.variables))},
			factory.Create(existential.assertion),
		},
	}
	return call
}
