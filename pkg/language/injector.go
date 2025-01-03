package language

import (
	"fmt"
	"go/parser"
	"go/token"
	"iter"
	"os"
	"path/filepath"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"github.com/hyperproperties/sopher/pkg/filesx"
)

type Injector struct{}

func NewGoInjector() Injector {
	return Injector{}
}

func (injector Injector) Imports(file *dst.File, imports map[string]string) {
	specs := make([]dst.Spec, 0)

	for name, path := range imports {
		importSpec := &dst.ImportSpec{
			Name: dst.NewIdent(name),
			Path: &dst.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("\"%v\"", path),
			},
		}

		file.Imports = append(file.Imports, importSpec)
		specs = append(specs, importSpec)
	}

	exists := false

	// import declaration already exist.
	for _, decl := range file.Decls {
		if cast, ok := decl.(*dst.GenDecl); ok && cast.Tok == token.IMPORT {
			cast.Specs = append(cast.Specs, specs...)
			exists = true
			break
		}
	}

	// No import so we create one.
	if !exists {
		genDecl := &dst.GenDecl{
			Tok:   token.IMPORT,
			Specs: specs,
		}

		file.Decls = append([]dst.Decl{genDecl}, file.Decls...)
	}
}

func (injector Injector) InputFields(function *dst.FuncDecl) (fields []*dst.Field) {
	for _, input := range function.Type.Params.List {
		fields = append(fields, dst.Clone(input).(*dst.Field))
	}
	return fields
}

func (injector Injector) HasNamedOutputs(function *dst.FuncDecl) bool {
	for _, output := range function.Type.Results.List {
		if len(output.Names) > 0 {
			return true
		}
	}
	return false
}

func (injector Injector) OutputFields(function *dst.FuncDecl) (fields []*dst.Field) {
	for idx, output := range function.Type.Results.List {
		if len(output.Names) > 0 {
			fields = append(fields, dst.Clone(output).(*dst.Field))
		} else {
			field := &dst.Field{
				Names: []*dst.Ident{
					dst.NewIdent(fmt.Sprintf("ret%v", idx)),
				},
				Type: dst.Clone(output.Type).(dst.Expr),
			}
			fields = append(fields, field)
		}
	}

	return fields
}

func (injector Injector) Model(function *dst.FuncDecl) (string, *dst.GenDecl) {
	name := function.Name.Name

	fields := make([]*dst.Field, 0)
	fields = append(fields, injector.InputFields(function)...)
	fields = append(fields, injector.OutputFields(function)...)

	model := &dst.StructType{
		Fields: &dst.FieldList{
			List: fields,
		},
	}

	modelName := name + "_ExecutionModel"

	return modelName, &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(modelName),
				Type: model,
			},
		},
	}
}

func (injector Injector) Contract(model string, function *dst.FuncDecl) (string, *dst.GenDecl) {
	name := function.Name.Name

	comments := function.Decs.NodeDecs.Start
	parser := NewParser(LexDocStrings(comments))
	contract := parser.Parse()

	assumptionList := make([]dst.Expr, len(contract.regions[0].assumptions))
	guaranteeList := make([]dst.Expr, len(contract.regions[0].guarantees))

	monitors := NewGoMonitorFactory("sopher", model)
	for idx, assumption := range contract.regions[0].assumptions {
		assumptionList[idx] = monitors.Create(assumption)
	}
	for idx, guarantee := range contract.regions[0].guarantees {
		guaranteeList[idx] = monitors.Create(guarantee)
	}

	constructor := &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			Sel: dst.NewIdent("NewAGContract"),
			X:   dst.NewIdent("sopher"),
		},
		Args: []dst.Expr{
			&dst.CompositeLit{
				Type: &dst.ArrayType{
					Elt: &dst.IndexExpr{
						X: &dst.SelectorExpr{
							Sel: dst.NewIdent("IncrementalMonitor"),
							X:   dst.NewIdent("sopher"),
						},
						Index: dst.NewIdent(model),
					},
				},
				Elts: assumptionList,
			},
			&dst.CompositeLit{
				Type: &dst.ArrayType{
					Elt: &dst.IndexExpr{
						X: &dst.SelectorExpr{
							Sel: dst.NewIdent("IncrementalMonitor"),
							X:   dst.NewIdent("sopher"),
						},
						Index: dst.NewIdent(model),
					},
				},
				Elts: guaranteeList,
			},
		},
	}

	contractName := name + "_Contract"

	return contractName, &dst.GenDecl{
		Tok: token.VAR,
		Specs: []dst.Spec{
			&dst.ValueSpec{
				Names: []*dst.Ident{
					dst.NewIdent(contractName),
				},
				Type: &dst.IndexExpr{
					X: &dst.SelectorExpr{
						X:   dst.NewIdent("sopher"),
						Sel: dst.NewIdent("AGContract"),
					},
					Index: dst.NewIdent(model),
				},
				Values: []dst.Expr{constructor},
			},
		},
	}
}

func (injector Injector) Wrap(function *dst.FuncDecl) *dst.AssignStmt {
	var TypeParams *dst.FieldList = nil
	if function.Type.TypeParams != nil {
		TypeParams = dst.Clone(function.Type.TypeParams).(*dst.FieldList)
	}

	var Params *dst.FieldList = nil
	if function.Type.Params != nil {
		Params = dst.Clone(function.Type.Params).(*dst.FieldList)
	}

	var Results *dst.FieldList = nil
	if function.Type.Results != nil {
		Results = dst.Clone(function.Type.Results).(*dst.FieldList)
	}

	return &dst.AssignStmt{
		Lhs: []dst.Expr{
			dst.NewIdent("wrap"),
		},
		Tok: token.DEFINE,
		Rhs: []dst.Expr{
			&dst.FuncLit{
				Type: &dst.FuncType{
					TypeParams: TypeParams,
					Params:     Params,
					Results:    Results,
				},
				Body: function.Body,
			},
		},
	}
}

func (injector Injector) ConstructModel(model string, function *dst.FuncDecl) *dst.AssignStmt {
	fields := make([]dst.Expr, 0)
	for _, field := range function.Type.Params.List {
		for _, name := range field.Names {
			fields = append(fields, &dst.KeyValueExpr{
				Key:   dst.NewIdent(name.Name),
				Value: dst.NewIdent(name.Name),
			})
		}
	}

	return &dst.AssignStmt{
		Lhs: []dst.Expr{
			dst.NewIdent("execution"),
		},
		Tok: token.DEFINE,
		Rhs: []dst.Expr{
			&dst.CompositeLit{
				Type: dst.NewIdent(model),
				Elts: fields,
			},
		},
	}
}

func (injector Injector) Check(name string, contractName string) *dst.IfStmt {
	return &dst.IfStmt{
		Cond: &dst.CallExpr{
			Fun: &dst.SelectorExpr{
				X: &dst.CallExpr{
					Fun: &dst.SelectorExpr{
						X:   dst.NewIdent(contractName),
						Sel: dst.NewIdent(name),
					},
					Args: []dst.Expr{
						dst.NewIdent("execution"),
					},
				},
				Sel: dst.NewIdent("IsFalse"),
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ExprStmt{
					X: &dst.CallExpr{
						Fun: dst.NewIdent("panic"),
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind:  token.STRING,
								Value: "\"\"",
							},
						},
					},
				},
			},
		},
	}
}

func (injector Injector) CallWrap(function *dst.FuncDecl) *dst.AssignStmt {
	var outputs []dst.Expr
	for _, output := range injector.OutputFields(function) {
		for _, name := range output.Names {
			outputs = append(outputs, dst.NewIdent(name.Name))
		}
	}

	var inputs []dst.Expr
	for _, input := range injector.InputFields(function) {
		for _, name := range input.Names {
			inputs = append(inputs, dst.NewIdent(name.Name))
		}
	}

	operator := token.DEFINE
	if injector.HasNamedOutputs(function) {
		operator = token.ASSIGN
	}

	return &dst.AssignStmt{
		Lhs: outputs,
		Tok: operator,
		Rhs: []dst.Expr{
			&dst.CallExpr{
				Fun:  dst.NewIdent("wrap"),
				Args: inputs,
			},
		},
	}
}

func (injector Injector) Updates(function *dst.FuncDecl) (updates []*dst.AssignStmt) {
	for _, output := range injector.OutputFields(function) {
		for _, name := range output.Names {
			updates = append(updates, &dst.AssignStmt{
				Lhs: []dst.Expr{
					&dst.SelectorExpr{
						X:   dst.NewIdent("execution"),
						Sel: dst.NewIdent(name.Name),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []dst.Expr{
					dst.NewIdent(name.Name),
				},
			})
		}
	}

	return updates
}

func (injector Injector) Return(function *dst.FuncDecl) *dst.ReturnStmt {
	var results []dst.Expr
	for _, output := range injector.OutputFields(function) {
		for _, name := range output.Names {
			results = append(results, dst.NewIdent(name.Name))
		}
	}

	return &dst.ReturnStmt{
		Results: results,
	}
}

func (injector Injector) Inject(file *dst.File) {
	dstutil.Apply(file, nil, func(cursor *dstutil.Cursor) bool {
		switch cast := cursor.Node().(type) {
		case *dst.FuncDecl:
			comments := cast.Decs.NodeDecs.Start
			if len(comments) == 0 {
				return true
			}

			modelName, model := injector.Model(cast)
			cursor.InsertBefore(model)

			contractName, contract := injector.Contract(modelName, cast)
			cursor.InsertBefore(contract)

			body := make([]dst.Stmt, 0)

			wrap := injector.Wrap(cast)
			body = append(body, wrap)

			modelConstruction := injector.ConstructModel(modelName, cast)
			body = append(body, modelConstruction)

			assumptionCheck := injector.Check("Assume", contractName)
			body = append(body, assumptionCheck)

			wrapCall := injector.CallWrap(cast)
			body = append(body, wrapCall)

			for _, update := range injector.Updates(cast) {
				body = append(body, update)
			}

			guaranteeCheck := injector.Check("Guarantee", contractName)
			body = append(body, guaranteeCheck)

			returnStmt := injector.Return(cast)
			body = append(body, returnStmt)

			declaration := &dst.FuncDecl{
				Recv: cast.Recv,
				Name: cast.Name,
				Type: cast.Type,
				Body: &dst.BlockStmt{
					List: body,
				},
				Decs: cast.Decs,
			}

			cursor.Replace(declaration)
		}
		return true
	})

	injector.Imports(file, map[string]string{
		"sopher": "github.com/hyperproperties/sopher/pkg/language",
	})
}

func (injector Injector) Files(files iter.Seq[string]) {
	fset := token.NewFileSet()
	decor := decorator.NewDecorator(fset)

	for path := range files {
		// Read the file
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		// Parse the decorated syntax tree.
		dst, err := decor.ParseFile(filepath.Base(path), content, parser.ParseComments)
		if err != nil {
			continue
		}

		injector.Inject(dst)

		// Move original file to keep it.
		filesx.Move(path, path+"-sopher")

		// Create file for instrumented version.
		newFile, err := filesx.Create(path)
		if err != nil {
			// TODO: What to do here? Revert to the original file and report it?
			panic(err)
		}

		decorator.Fprint(newFile, dst)
	}
}

func (injector Injector) Restore(files iter.Seq[string]) {
	for path := range files {
		filesx.Move(path, path+"sopher")
	}
}
