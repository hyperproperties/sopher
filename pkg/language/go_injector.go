package language

import (
	"go/ast"
	"go/printer"

	"github.com/hyperproperties/sopher/pkg/filesx"
	"golang.org/x/tools/go/ast/astutil"
)

type GoInjector struct {
	contracts GoContracts
}

func NewGoInjector(contracts GoContracts) GoInjector {
	return GoInjector{
		contracts: contracts,
	}
}

func (injector GoInjector) Imports(file *ast.File, imports map[string]string) {
	for name, path := range imports {
		astutil.AddNamedImport(injector.contracts.fset, file, name, path)
	}
}

func (injector GoInjector) Contract(file *ast.File, declaration *ast.FuncDecl) {
	// Create execution model.
	models := NewGoExecutionModelFactory()
	name := declaration.Name.Name
	modelName, model := models.Create(name, declaration)
	file.Decls = append(file.Decls, model)

	// Parse the contract.
	parser := NewParser(LexGo(declaration.Doc))
	contract := parser.Parse()
	for _, region := range contract.regions {
		// Declare AG contract.
		contracts := NewGoAGContractFactory("sopher", modelName)
		_, contract := contracts.Declaration("sopher", modelName, name, region.assumptions, region.guarantees)
		file.Decls = append(file.Decls, contract)
	}

	// Wrap the function body.
	/*wrap := &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("wrap"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.FuncLit{
				Type: &ast.FuncType{
					TypeParams: declaration.Type.TypeParams,
					Params: declaration.Type.Params,
					Results: declaration.Type.Results,
				},
				Body: declaration.Body,
			},
		},
	}

	var buffer bytes.Buffer
	printer.Fprint(&buffer, token.NewFileSet(), wrap)

	log.Println(buffer.String())
	
	// Clear the function delcaration body.
	declaration.Body.List = append([]ast.Stmt{wrap}, declaration.Body.List...)*/


	// Construct the pre-execution model.

	// Check the assumption.

	// Store the return values.

	// Assign the return values to the execution model.

	// Check the guarantee.

	// Return original outputs.
}

func (injector GoInjector) Inject() {
	for position, contracts := range injector.contracts.Iterator() {
		astFile, tokenFile := injector.contracts.File(position)

		injector.Imports(astFile, map[string]string{
			"sopher": "github.com/hyperproperties/sopher/pkg/language",
		})

		for declaration, _ := range contracts {
			injector.Contract(astFile, declaration)
		}

		// Move file to temp and write a new one with generated code.
		path := tokenFile.Name()
		filesx.Move(path, path+"-sopher")

		newFile, err := filesx.Create(path)
		if err != nil {
			panic(err)
		}
		printer.Fprint(newFile, injector.contracts.fset, astFile)
		newFile.Close()
	}
}

func (injector GoInjector) Restore() {
	for position := range injector.contracts.Iterator() {
		_, tokenFile := injector.contracts.File(position)
		name := tokenFile.Name()
		filesx.Delete(name)
		filesx.Move(name+"-sopher", name)
	}
}
