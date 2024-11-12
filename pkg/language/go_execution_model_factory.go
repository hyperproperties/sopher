package language

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type GoExecutionModelFactory struct{}

func NewGoExecutionModelFactory() GoExecutionModelFactory {
	return GoExecutionModelFactory{}
}

// Constructs an execution model declaration for the function declaration.
func (factory GoExecutionModelFactory) Create(name string, function *ast.FuncDecl) (string, *ast.GenDecl) {
	// TODO: Handle generic functions.

	var fields []*ast.Field
	fields = append(fields, function.Type.Params.List...)

	for idx, output := range function.Type.Results.List {
		if len(output.Names) > 0 {
			field := &ast.Field{
				Names: output.Names,
				Type:  output.Type,
			}
			fields = append(fields, field)
		} else {
			ident := ast.NewIdent(fmt.Sprintf("ret%v", idx))
			field := &ast.Field{
				Names: []*ast.Ident{ident},
				Type:  output.Type,
			}
			fields = append(fields, field)
		}
	}

	model := &ast.StructType{
		Fields: &ast.FieldList{
			List: fields,
		},
	}

	modelName := strings.Join([]string{name, "ExecutionModel"}, "_")

	return modelName, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(modelName),
				Type: model,
			},
		},
	}
}
