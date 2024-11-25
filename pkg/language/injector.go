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
	"github.com/hyperproperties/sopher/pkg/dstx"
	"github.com/hyperproperties/sopher/pkg/filesx"
)

type Injector struct{}

func NewGoInjector() Injector {
	return Injector{}
}

func (injector Injector) Imports(file *dst.File, imports map[string]string) {
	specs := make([]dst.Spec, 0)
	for name, path := range imports {
		importSpec := dstx.ImportS(name, path)
		file.Imports = append(file.Imports, importSpec)
		specs = append(specs, importSpec)
	}

	if declaration := dstx.FindImportDeclaration(file.Decls); declaration != nil {
		declaration.Specs = append(declaration.Specs, specs...)
	} else {
		declaration := dstx.DeclareImports(specs...)
		file.Decls = append([]dst.Decl{declaration}, file.Decls...)
	}
}

func (injector Injector) InputFields(function *dst.FuncDecl) (fields []*dst.Field) {
	return dstx.Clones(function.Type.Params.List)
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

func (injector Injector) DeclareModel(function *dst.FuncDecl) (string, *dst.GenDecl) {
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
	parser := NewParser(LexComments(comments))
	contract := parser.Parse()

	// TODO: Support multiple regions.
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
			Sel: dst.NewIdent("NewAGHyperContract"),
			X:   dst.NewIdent("sopher"),
		},
		Args: []dst.Expr{
			&dst.CallExpr{
				Fun: &dst.IndexExpr{
					X: &dst.SelectorExpr{
						Sel: dst.NewIdent("NewAllAssertion"),
						X:   dst.NewIdent("sopher"),
					},
					Index: dst.NewIdent(model),
				},
				Args: assumptionList,
			},
			&dst.CallExpr{
				Fun: &dst.IndexExpr{
					X: &dst.SelectorExpr{
						Sel: dst.NewIdent("NewAllAssertion"),
						X:   dst.NewIdent("sopher"),
					},
					Index: dst.NewIdent(model),
				},
				Args: guaranteeList,
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
						Sel: dst.NewIdent("AGHyperContract"),
					},
					Index: dst.NewIdent(model),
				},
				Values: []dst.Expr{constructor},
			},
		},
	}
}

func (injector Injector) DeclareWrap(function *dst.FuncDecl) *dst.AssignStmt {
	return dstx.DefineS("wrap").As(dstx.FunctionOf(function))
}

func (injector Injector) DeclareCall(model string, function *dst.FuncDecl) *dst.AssignStmt {
	block := dstx.Sequence(
		injector.DeclareWrap(function),
		injector.CallWrap(function),
	).Append(dstx.Statements(injector.Updates(function))...).
		Append(dstx.ReturnS("execution"))

	return dstx.DefineS("call").
		As(
			dstx.Function(
				dstx.TakingN(dstx.FieldS("execution").TypeS(model)).
					ResultsN(dstx.Field().TypeS(model)),
				block.Terminate(),
			),
		)
}

func (injector Injector) ConstructModel(model string, function *dst.FuncDecl) *dst.AssignStmt {
	composite := dstx.ComposeS("execution")
	for _, field := range function.Type.Params.List {
		for _, name := range field.Names {
			composite.Elements(dstx.KeyS(name.Name).ValueS(name.Name))
		}
	}
	return dstx.DefineS("execution").As(composite.Final())
}

func (injector Injector) Check(name string) *dst.IfStmt {
	return dstx.
		If(dstx.Call(dstx.SelectS("IsFalse").FromS(name)).Pass()).
		ThenN(dstx.ExprStmt(dstx.CallS("panic").Pass(dstx.BasicString("\"\"")))).
		End()
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
			inputs = append(inputs, dstx.SelectS(name.Name).FromS("execution"))
		}
	}

	if dstx.HasNamedOutputs(function) {
		return dstx.Assign(outputs...).To(dstx.CallS("wrap").Pass(inputs...))
	}
	return dstx.Define(outputs...).As(dstx.CallS("wrap").Pass(inputs...))
}

func (injector Injector) Updates(function *dst.FuncDecl) (updates []*dst.AssignStmt) {
	for _, output := range injector.OutputFields(function) {
		for _, name := range output.Names {
			assignment := dstx.Assign(
				dstx.SelectS(name.Name).FromS("execution"),
			).ToS(name.Name)
			updates = append(updates, assignment)
		}
	}

	return updates
}

func (injector Injector) Call(contract string) *dst.AssignStmt {
	return dstx.DefineS("assumption", "execution", "guarantee").As(
		dstx.Call(dstx.SelectS("Call").FromS(contract)).PassS("caller", "execution", "call"),
	)
}

func (injector Injector) Return(function *dst.FuncDecl) *dst.ReturnStmt {
	var results []dst.Expr
	for _, output := range injector.OutputFields(function) {
		for _, name := range output.Names {
			results = append(results, dstx.SelectS(name.Name).FromS("execution"))
		}
	}

	return dstx.Return(results...)
}

func (injector Injector) GetCaller() *dst.AssignStmt {
	return dstx.DefineS("caller").As(dstx.Call(dstx.SelectS("Caller").FromS("sopher")).Pass())
}

func (injector Injector) Inject(file *dst.File) {
	dstutil.Apply(file, nil, func(cursor *dstutil.Cursor) bool {
		switch cast := cursor.Node().(type) {
		case *dst.FuncDecl:
			comments := cast.Decs.NodeDecs.Start
			if len(comments) == 0 {
				return true
			}

			modelName, model := injector.DeclareModel(cast)
			cursor.InsertBefore(model)

			contractName, contract := injector.Contract(modelName, cast)
			cursor.InsertBefore(contract)

			declaration := &dst.FuncDecl{
				Recv: cast.Recv,
				Name: cast.Name,
				Type: cast.Type,
				Body: dstx.Block(
					injector.GetCaller(),
					injector.ConstructModel(modelName, cast),
					injector.DeclareCall(modelName, cast),
					injector.Call(contractName),
					injector.Check("assumption"),
					injector.Check("guarantee"),
					injector.Return(cast),
				),
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
			filesx.Move(path+"-sopher", path)
			// TODO: Error handling.
		}

		decorator.Fprint(newFile, dst)
	}
}

func (injector Injector) Restore(files iter.Seq[string]) {
	for path := range files {
		filesx.Move(path+"-sopher", path)
	}
}
