package language

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelFor(t *testing.T) {
	source := `
	package main

	func Foo(a, b int) (int, int) { }
	`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	assert.NoError(t, err)

	function := file.Decls[0].(*ast.FuncDecl)

	factory := NewGoExecutionModelFactory()
	_, model := factory.Create("Foo", function)

	log.Println(model)

	var buffer bytes.Buffer
	printer.Fprint(&buffer, fset, model)

	t.Log(buffer.String())

	assert.Equal(t,
		`type Foo_ExecutionModel struct {
	a, b	int
	ret0	int
	ret1	int
}`, buffer.String())
}
