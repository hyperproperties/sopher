package language

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"iter"
	"os"

	"github.com/hyperproperties/sopher/pkg/filesx"
	"golang.org/x/tools/go/ast/astutil"
)

type Injector struct{}

func NewGoInjector() Injector {
	return Injector{}
}

func (injector Injector) Imports(file *ast.File, fset *token.FileSet, imports map[string]string) {
	for name, path := range imports {
		astutil.AddNamedImport(fset, file, name, path)
	}
}

func (injector Injector) Model(function *ast.FuncDecl) (string, *ast.GenDecl) {
	name := function.Name.Name
	
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

	modelName := name+"_ExecutionModel"

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

func (injector Injector) Contract(model string, function *ast.FuncDecl) (string, *ast.GenDecl) {
	name := function.Name.Name

	parser := NewParser(LexGo(function.Doc))
	contract := parser.Parse()
	
	assumptionList := make([]ast.Expr, len(contract.regions[0].assumptions))
	guaranteeList := make([]ast.Expr, len(contract.regions[0].guarantees))
	
	monitors := NewGoMonitorFactory("sopher", model)
	for idx, assumption := range contract.regions[0].assumptions {
		assumptionList[idx] = monitors.Create(assumption)
	}
	for idx, guarantee := range contract.regions[0].guarantees {
		guaranteeList[idx] = monitors.Create(guarantee)
	}

	constructor := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			Sel: ast.NewIdent("NewAGContract"),
			X:   ast.NewIdent("sopher"),
		},
		Args: []ast.Expr{
			&ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: &ast.IndexExpr{
						X: &ast.SelectorExpr{
							Sel: ast.NewIdent("IncrementalMonitor"),
							X:   ast.NewIdent("sopher"),
						},
						Index: ast.NewIdent(model),
					},
				},
				Elts: assumptionList,
			},
			&ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: &ast.IndexExpr{
						X: &ast.SelectorExpr{
							Sel: ast.NewIdent("IncrementalMonitor"),
							X:   ast.NewIdent("sopher"),
						},
						Index: ast.NewIdent(model),
					},
				},
				Elts: guaranteeList,
			},
		},
	}

	contractName := name+"_Contract"

	return contractName, &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					ast.NewIdent(contractName),
				},
				Type: &ast.IndexExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("sopher"),
						Sel: ast.NewIdent("AGContract"),
					},
					Index: ast.NewIdent(model),
				},
				Values: []ast.Expr{ constructor },
			},
		},
	}
}

func (injector Injector) Inject(file *token.File, fset *token.FileSet, root *ast.File) {
	declarations := make([]ast.Decl, 0)

	astutil.Apply(root, nil, func(cursor *astutil.Cursor) bool {
		switch cast := cursor.Node().(type) {
		case *ast.FuncDecl:
			if len(cast.Doc.List) == 0 {
				return true
			}
			
			modelName, model := injector.Model(cast)
			_, contract := injector.Contract(modelName, cast)
			declarations = append(declarations, []ast.Decl{model, contract}...)

			wrap := &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("wrap"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.FuncLit{
						Type: &ast.FuncType{
							TypeParams: cast.Type.TypeParams,
							Params:     cast.Type.Params,
							Results:    cast.Type.Results,
						},
						Body: cast.Body,
					},
				},
			}

			body := make([]ast.Stmt, 0)
			body = append(body, wrap)

			declaration := &ast.FuncDecl{
				Doc:  cast.Doc,
				Recv: cast.Recv,
				Name: cast.Name,
				Type: cast.Type,
				Body: &ast.BlockStmt{
					List: body,
				},
			}

			cursor.Replace(declaration)
		}

		return true
	})

	if len(declarations) > 0 {
		injector.Imports(root, fset, map[string]string{
			"sopher": "github.com/hyperproperties/sopher/pkg/language",
		})
	}

	root.Decls = append(root.Decls, declarations...)

	path := file.Name()
	filesx.Move(path, path+"-sopher")

	newFile, err := filesx.Create(path)
	if err != nil {
		panic(err)
	}

	printer.Fprint(newFile, fset, root)
	newFile.Close()
}

func (injector Injector) Files(files iter.Seq[string]) {
	fset := token.NewFileSet()

	for path := range files {
		content, _ := os.ReadFile(path)
		root, _ := parser.ParseFile(fset, path, content, parser.ParseComments)
		file := fset.File(root.Pos())
		injector.Inject(file, fset, root)
	}
}

func (injector Injector) Restore(files iter.Seq[string]) {
	for path := range files {
		filesx.Move(path+"-sopher", path)
	}
}