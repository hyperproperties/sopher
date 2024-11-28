package language

import (
	"fmt"
	"iter"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
	"github.com/hyperproperties/sopher/pkg/dstx"
	"github.com/hyperproperties/sopher/pkg/filesx"
)

type Instrumentor struct{}

func NewInstrumentor() Instrumentor {
	return Instrumentor{}
}

func (instrumentor Instrumentor) OutputFields(function *dst.FuncDecl) (fields []*dst.Field) {
	if dstx.HasNamedOutputs(function) {
		return dstx.Clones(function.Type.Results.List)
	}

	for idx, output := range function.Type.Results.List {
		field := dstx.
			FieldS(fmt.Sprintf("ret%v", idx)).
			Type(dstx.Clone(output.Type))
		fields = append(fields, field)
	}

	return fields
}

func (instrumentor Instrumentor) ModelDeclaration(function FunctionContract, declaration *dst.FuncDecl) *dst.GenDecl {
	fields := make([]*dst.Field, 0)
	fields = append(fields, dstx.Clones(declaration.Type.Params.List)...)
	fields = append(fields, instrumentor.OutputFields(declaration)...)

	if dstx.HasReceiver(declaration) {
		fields = append(fields, dstx.Clones(declaration.Recv.List)...)
	}

	return dstx.DeclareStructTypeSN(function.ModelName(), fields...)
}

func (instrumentor Instrumentor) ContractDeclaration(file *dst.File, function FunctionContract, declaration *dst.FuncDecl) *dst.GenDecl {
	region := function.contract.regions[0]
	assumptions := make([]dst.Expr, len(region.assumptions))
	guarantees := make([]dst.Expr, len(region.guarantees))

	assertions := NewAssertionFactory("sopher", function.ModelName())
	for idx, assumption := range region.assumptions {
		assumptions[idx] = assertions.Create(assumption)
	}
	for idx, guarantee := range region.guarantees {
		guarantees[idx] = assertions.Create(guarantee)
	}

	constructor := dstx.
		Call(dstx.SelectS("NewAGHyperContract").FromS("sopher")).
		Pass(
			dstx.NewLineAround(
				dstx.Call(
					dstx.IndexS(function.ModelName()).Of(dstx.SelectS("NewAllAssertion").FromS("sopher")),
				).Pass(assumptions...),
			),
			dstx.NewLineAround(
				dstx.Call(
					dstx.IndexS(function.ModelName()).Of(dstx.SelectS("NewAllAssertion").FromS("sopher")),
				).Pass(guarantees...),
			),
		)

	return dstx.
		DeclareVariableS(function.ContractName()).
		Type(dstx.IndexS(function.ModelName()).Of(dstx.SelectS("AGHyperContract").FromS("sopher"))).
		Values(constructor)
}

func (instrumentor Instrumentor) DefineWrap(declaration *dst.FuncDecl) *dst.AssignStmt {
	function := dstx.Function(
		dstx.
			Taking(dstx.Clone(declaration.Recv)).
			Taking(dstx.Clone(declaration.Type.Params)).
			Results(dstx.Clone(declaration.Type.Results)),
		dstx.Clone(declaration.Body),
	)
	return dstx.DefineS("wrap").As(function)
}

func (instrumentor Instrumentor) CallWrap(declaration *dst.FuncDecl) *dst.AssignStmt {
	var outputs []dst.Expr
	for _, output := range instrumentor.OutputFields(declaration) {
		for _, name := range output.Names {
			outputs = append(outputs, dst.NewIdent(name.Name))
		}
	}

	var inputs []dst.Expr
	if dstx.HasReceiver(declaration) {
		for _, field := range declaration.Recv.List {
			for _, name := range field.Names {
				inputs = append(inputs, dstx.SelectS(name.Name).FromS("execution"))
			}
		}
	}

	for _, input := range declaration.Type.Params.List {
		for _, name := range input.Names {
			inputs = append(inputs, dstx.SelectS(name.Name).FromS("execution"))
		}
	}

	if dstx.HasNamedOutputs(declaration) {
		return dstx.Assign(outputs...).To(dstx.CallS("wrap").Pass(inputs...))
	}
	return dstx.Define(outputs...).As(dstx.CallS("wrap").Pass(inputs...))
}

func (instrumentor Instrumentor) ModelOutputAssignments(declaration *dst.FuncDecl) (updates []*dst.AssignStmt) {
	// TODO: Refactor
	for _, output := range instrumentor.OutputFields(declaration) {
		for _, name := range output.Names {
			assignment := dstx.Assign(
				dstx.SelectS(name.Name).FromS("execution"),
			).ToS(name.Name)
			updates = append(updates, assignment)
		}
	}

	return updates
}

func (instrumentor Instrumentor) CallDeclaration(function FunctionContract, declaration *dst.FuncDecl) *dst.GenDecl {
	block := dstx.Sequence(
		instrumentor.DefineWrap(declaration),
		instrumentor.CallWrap(declaration),
	).Append(dstx.Statements(instrumentor.ModelOutputAssignments(declaration))...).
		Append(dstx.ReturnS("execution"))

	return dstx.
		DeclareVariableS(function.CallName()).
		Values(
			dstx.Function(
				dstx.TakingN(dstx.FieldS("execution").TypeS(function.ModelName())).
					ResultsN(dstx.Field().TypeS(function.ModelName())),
				block.Terminate(),
			),
		)
}

func (instrumentor Instrumentor) CallerIdDefinition() *dst.AssignStmt {
	return dstx.DefineS("callerID").As(dstx.Call(dstx.SelectS("Caller").FromS("sopher")).Pass())
}

func (instrumentor Instrumentor) ModelDefinition(function FunctionContract, declaration *dst.FuncDecl) *dst.AssignStmt {
	composite := dstx.ComposeS(function.ModelName())
	for _, field := range declaration.Type.Params.List {
		for _, name := range field.Names {
			composite.Elements(dstx.KeyS(name.Name).ValueS(name.Name))
		}
	}

	if dstx.HasReceiver(declaration) {
		for _, field := range declaration.Recv.List {
			for _, name := range field.Names {
				composite.Elements(dstx.KeyS(name.Name).ValueS(name.Name))
			}
		}
	}

	return dstx.DefineS("execution").As(composite.Final())
}

func (instrumentor Instrumentor) TripleDefinition(function FunctionContract) *dst.AssignStmt {
	return dstx.DefineS("assumption", "execution", "guarantee").As(
		dstx.Call(dstx.SelectS("Call").FromS(function.ContractName())).PassS("callerID", "execution", function.CallName()),
	)
}

func (instrumentor Instrumentor) IsFalseThenPanic(name string) *dst.IfStmt {
	return dstx.
		If(dstx.Call(dstx.SelectS("IsFalse").FromS(name)).Pass()).
		ThenN(dstx.ExprStmt(dstx.CallS("panic").Pass(dstx.BasicString("\"\"")))).
		End()
}

func (instrumentor Instrumentor) ModelReturn(declaration *dst.FuncDecl) *dst.ReturnStmt {
	var results []dst.Expr
	for _, output := range instrumentor.OutputFields(declaration) {
		for _, name := range output.Names {
			results = append(results, dstx.SelectS(name.Name).FromS("execution"))
		}
	}

	return dstx.Return(results...)
}

func (instrumentor Instrumentor) Function(file *dst.File, function FunctionContract, declaration *dst.FuncDecl, cursor *dstutil.Cursor) {
	cursor.InsertBefore(instrumentor.ModelDeclaration(function, declaration))
	cursor.InsertBefore(instrumentor.ContractDeclaration(file, function, declaration))
	cursor.InsertBefore(instrumentor.CallDeclaration(function, declaration))

	declaration = &dst.FuncDecl{
		Recv: declaration.Recv,
		Name: declaration.Name,
		Type: declaration.Type,
		Body: dstx.Block(
			instrumentor.CallerIdDefinition(),
			instrumentor.ModelDefinition(function, declaration),
			instrumentor.TripleDefinition(function),
			instrumentor.IsFalseThenPanic("assumption"),
			instrumentor.IsFalseThenPanic("guarantee"),
			instrumentor.ModelReturn(declaration),
		),
		Decs: declaration.Decs,
	}

	cursor.Replace(declaration)

	dstx.PrependStart(file, "// Code generated by sopher (https://github.com/hyperproperties/sopher). DO NOT EDIT.")
}

func (instrumentor Instrumentor) Files(files iter.Seq[string]) {
	parser := NewFileParser()
	for path, file := range parser.DstFiles(files) {
		dstx.ImportIntoS(file, map[string]string{
			"sopher": "github.com/hyperproperties/sopher/pkg/language",
		})

		parser.Apply(file, func(contract FunctionContract, cursor *dstutil.Cursor) bool {
			declaration, ok := cursor.Node().(*dst.FuncDecl)
			if !ok {
				return true
			}

			// Instrument function declaration.
			instrumentor.Function(file, contract, declaration, cursor)

			// Move original file to keep it.
			filesx.Move(path, path+"-sopher")

			// Create file for instrumented version.
			newFile, err := filesx.Create(path)
			if err != nil {
				filesx.Move(path+"-sopher", path)
				// TODO: Error handling.
			}

			decorator.Fprint(newFile, file)

			return true
		})
	}
}

func (instrumentor Instrumentor) Restore(files iter.Seq[string]) {
	for file := range files {
		if filesx.Exists(file + "-sopher") {
			filesx.Move(file+"-sopher", file)
		}
	}
}
