package language

import (
	"fmt"
	"go/ast"
	"go/token"
)

type GoAGContractFactory struct {
	packageName string
	modelName   string
	monitors    GoMonitorFactory
}

func NewGoAGContractFactory(packageName, modelName string) GoAGContractFactory {
	return GoAGContractFactory{
		packageName: packageName,
		modelName:   modelName,
		monitors:    NewGoMonitorFactory(packageName, modelName),
	}
}

func (factory *GoAGContractFactory) Declaration(packageName, modelName, name string, assumptions, guarantees []Node) (string, ast.Decl) {
	contractName := fmt.Sprintf("%v_Contract", name)

	return contractName, &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					ast.NewIdent(contractName),
				},
				Type: &ast.IndexExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(packageName),
						Sel: ast.NewIdent("AGContract"),
					},
					Index: ast.NewIdent(modelName),
				},
				Values: []ast.Expr{
					factory.Constructor(assumptions, guarantees),
				},
			},
		},
	}
}

func (factory *GoAGContractFactory) Constructor(assumptions, guarantees []Node) *ast.CallExpr {
	assumptionList := make([]ast.Expr, len(assumptions))
	guaranteeList := make([]ast.Expr, len(guarantees))

	for idx, assumption := range assumptions {
		assumptionList[idx] = factory.monitors.Create(assumption)
	}

	for idx, guarantee := range guarantees {
		guaranteeList[idx] = factory.monitors.Create(guarantee)
	}

	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			Sel: ast.NewIdent("NewAGContract"),
			X:   ast.NewIdent(factory.packageName),
		},
		Args: []ast.Expr{
			&ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: &ast.IndexExpr{
						X: &ast.SelectorExpr{
							Sel: ast.NewIdent("IncrementalMonitor"),
							X:   ast.NewIdent(factory.packageName),
						},
						Index: ast.NewIdent(factory.modelName),
					},
				},
				Elts: assumptionList,
			},
			&ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: &ast.IndexExpr{
						X: &ast.SelectorExpr{
							Sel: ast.NewIdent("IncrementalMonitor"),
							X:   ast.NewIdent(factory.packageName),
						},
						Index: ast.NewIdent(factory.modelName),
					},
				},
				Elts: guaranteeList,
			},
		},
	}
}
