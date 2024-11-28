package language

import (
	goparser "go/parser"
	"go/token"
	"iter"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
)

type FileParser struct{}

func NewFileParser() FileParser {
	return FileParser{}
}

func (parser FileParser) Apply(file *dst.File, application func(function FunctionContract, cursor *dstutil.Cursor) bool) {
	var index uint = 0
	dstutil.Apply(file, nil, func(cursor *dstutil.Cursor) bool {
		if cast, ok := cursor.Node().(*dst.FuncDecl); ok {
			// If there are no comments then there are no contract.
			comments := cast.Decs.NodeDecs.Start
			if len(comments) == 0 {
				return true
			}

			// Parse contract.
			tokens := LexComments(comments)
			parser := NewParser(tokens)
			contract := parser.Parse()

			// Create function contract.
			name := cast.Name.Name
			function := NewFunctionContract(name, index, contract)
			index++

			return application(function, cursor)
		}
		return true
	})
}

func (parser FileParser) ApplyTo(files iter.Seq[string], application func(ffunction FunctionContract, cursor *dstutil.Cursor) bool) {
	for _, file := range parser.DstFiles(files) {
		parser.Apply(file, application)
	}
}

func (parser FileParser) File(file *dst.File) (functions []FunctionContract) {
	parser.Apply(file, func(function FunctionContract, _ *dstutil.Cursor) bool {
		functions = append(functions, function)
		return true
	})
	return
}

func (parser FileParser) DstFiles(files iter.Seq[string]) iter.Seq2[string, *dst.File] {
	fset := token.NewFileSet()
	decorator := decorator.NewDecorator(fset)

	return func(yield func(string, *dst.File) bool) {
		for path := range files {
			// Read the file.
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			// Parse the dst file.
			node, err := decorator.ParseFile(
				path, content, goparser.ParseComments,
			)
			if err != nil {
				continue
			}

			if !yield(path, node) {
				continue
			}
		}
	}
}

func (parser FileParser) Files(files iter.Seq[string]) iter.Seq2[string, File] {
	return func(yield func(string, File) bool) {
		for path, node := range parser.DstFiles(files) {
			functions := parser.File(node)
			file := NewFile(node.Name.Name, path, functions)
			if !yield(path, file) {
				return
			}
		}
	}
}
