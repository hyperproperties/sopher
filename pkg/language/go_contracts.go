package language

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strings"
)

type GoContracts struct {
	fset  *token.FileSet
	files map[token.Pos]*ast.File
}

func NewGoContracts() GoContracts {
	return GoContracts{
		fset:  token.NewFileSet(),
		files: map[token.Pos]*ast.File{},
	}
}

func (contracts GoContracts) Length() int {
	return len(contracts.files)
}

func (contracts GoContracts) File(pos token.Pos) (*ast.File, *token.File) {
	ast, ok := contracts.files[pos]
	if !ok {
		return nil, nil
	}
	return ast, contracts.fset.File(pos)
}

func (contracts GoContracts) Add(path string) ([]token.Pos, error) {
	info, err := os.Stat(path)
	if err != nil {
		return []token.Pos{}, err
	}

	if info.IsDir() {
		return contracts.AddDirectory(path)
	}

	pos, err := contracts.AddFile(path)
	return []token.Pos{pos}, err
}

func (contracts GoContracts) AddFile(filename string) (token.Pos, error) {
	if filepath.Ext(filename) == "" {
		filename = strings.Join([]string{filename, ".go"}, "")
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return token.NoPos, err
	}

	ast, err := parser.ParseFile(contracts.fset, filename, file, parser.ParseComments)
	if err != nil {
		return token.NoPos, err
	}

	position := ast.Pos()
	contracts.files[position] = ast

	return position, nil
}

func (contracts *GoContracts) AddDirectory(directory string) (positions []token.Pos, err error) {
	trimed, includeSubdirs := strings.CutSuffix(directory, "...")

	if includeSubdirs {
		err = filepath.WalkDir(trimed, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if entry.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".go" {
				return nil
			}

			position, err := contracts.AddFile(path)
			if err != nil {
				return err
			}
			positions = append(positions, position)

			return nil
		})
	} else {
		files, err := os.ReadDir(directory)
		if err != nil {
			return positions, err
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".go" {
				path := filepath.Join(directory, file.Name())
				position, err := contracts.AddFile(path)
				if err != nil {
					return positions, err
				}
				positions = append(positions, position)
			}
		}
	}

	return positions, err
}

func (contracts GoContracts) Iterator() iter.Seq2[token.Pos, iter.Seq2[*ast.FuncDecl, Contract]] {
	return func(yield func(token.Pos, iter.Seq2[*ast.FuncDecl, Contract]) bool) {
		for position, file := range contracts.files {
			if !yield(position, func(yield func(*ast.FuncDecl, Contract) bool) {
				for _, declaration := range file.Decls {
					if functionDeclaration, ok := declaration.(*ast.FuncDecl); ok {
						lexer := LexGo(functionDeclaration.Doc)
						parser := NewParser(lexer)

						if !yield(functionDeclaration, parser.Parse()) {
							return
						}
					}
				}
			}) {
				return
			}
		}
	}
}
