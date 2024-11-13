package language

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

type MonitorFactory struct {
	packageName string
	modelName   string
	offset      int
	variables   []string
}

func NewGoMonitorFactory(packageName, modelName string) MonitorFactory {
	return MonitorFactory{
		packageName: packageName,
		modelName:   modelName,
		offset:      0,
		variables:   make([]string, 0),
	}
}

func (factory *MonitorFactory) Create(node Node) *dst.CallExpr {
	switch cdst := node.(type) {
	case GoExpresion:
		return factory.NewPredicateMonitorCall(cdst)
	case Universal:
		return factory.NewUniversalMonitorCall(cdst)
	case Existential:
		return factory.NewExistentialMonitorCall(cdst)
	case Guarantee:
		factory.variables = nil
		factory.offset = 0
		return factory.Create(cdst.assertion)
	case Assumption:
		factory.variables = nil
		factory.offset = 0
		return factory.Create(cdst.assertion)
	}
	factory.offset = 0
	panic(fmt.Sprintf("unknown node type %t", node))
}

func (factory *MonitorFactory) NewPredicateMonitorCall(expression GoExpresion) *dst.CallExpr {
	// FIXME: Can accidentally define variables not in use. First we have to see what variables are in use and only define those.

	var body []dst.Stmt

	if len(factory.variables) > 0 {
		// e0, e1, e2 := assignments[0], assignments[1], assignments[2]
		var lhs, anon, rhs []dst.Expr
		for idx, identifier := range factory.variables {
			anon = append(anon, dst.NewIdent(identifier))
			lhs = append(lhs, dst.NewIdent(identifier))
			rhs = append(rhs, &dst.IndexExpr{
				X:     dst.NewIdent("assignments"),
				Index: &dst.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", idx)},
			})
		}

		executionAssignment := &dst.AssignStmt{
			Lhs: lhs, Tok: token.DEFINE, Rhs: rhs,
		}
		body = append(body, executionAssignment)

		// anonymous assignments for each execution variable.
		anonymousAssignment := &dst.AssignStmt{
			Tok: token.ASSIGN, Rhs: anon,
		}
		for idx := 0; idx < len(lhs); idx++ {
			anonymousAssignment.Lhs = append(anonymousAssignment.Lhs, dst.NewIdent("_"))
		}
		body = append(body, anonymousAssignment)
	}

	// Inject the expression.
	expr, err := parser.ParseExpr(expression.code)
	if err != nil {
		panic(err)
	}
	decor := decorator.NewDecorator(token.NewFileSet())
	code, _ := decor.DecorateNode(expr)
	body = append(body, &dst.ReturnStmt{
		Results: []dst.Expr{
			code.(dst.Expr),
		},
	})

	predicate := &dst.FuncLit{
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("assignments")},
						Type:  &dst.ArrayType{Elt: dst.NewIdent(factory.modelName)},
					},
				},
			},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{Type: dst.NewIdent("bool")},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: body,
		},
	}

	return &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			X:   dst.NewIdent(factory.packageName),
			Sel: dst.NewIdent("NewPredicateMonitor"),
		},
		Args: []dst.Expr{predicate},
	}
}

func (factory *MonitorFactory) NewUniversalMonitorCall(universal Universal) *dst.CallExpr {
	offset := factory.offset
	factory.offset += len(universal.variables)
	factory.variables = append(factory.variables, universal.variables...)
	call := &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			X:   dst.NewIdent(factory.packageName),
			Sel: dst.NewIdent("NewUniversalMonitor"),
		},
		Args: []dst.Expr{
			&dst.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", offset)},
			&dst.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", len(universal.variables))},
			factory.Create(universal.assertion),
		},
	}
	return call
}

func (factory *MonitorFactory) NewExistentialMonitorCall(existential Existential) *dst.CallExpr {
	offset := factory.offset
	factory.offset += len(existential.variables)
	factory.variables = append(factory.variables, existential.variables...)
	call := &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			X:   dst.NewIdent(factory.packageName),
			Sel: dst.NewIdent("NewExistentialMonitor"),
		},
		Args: []dst.Expr{
			&dst.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", offset)},
			&dst.BasicLit{Kind: token.INT, Value: fmt.Sprintf("%v", len(existential.variables))},
			factory.Create(existential.assertion),
		},
	}
	return call
}
