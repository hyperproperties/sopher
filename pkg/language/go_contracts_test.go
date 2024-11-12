package language

import (
	"go/token"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hyperproperties/sopher/pkg/iterx"
	"github.com/stretchr/testify/assert"
)

func TestGoContractsAddFile(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	path := path.Join(filepath.Dir(file), "/examples/non-interference")
	contracts := NewGoContracts()
	position, err := contracts.AddFile(path)
	assert.NotEqual(t, token.NoPos, position)
	assert.NoError(t, err)
}

func TestGoContractsAddDirectory(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	directory := path.Join(filepath.Dir(file), "/examples/")
	contracts := NewGoContracts()
	positions, err := contracts.AddDirectory(directory)
	assert.Len(t, positions, 1)
	assert.NoError(t, err)
}

func TestGoContractsAddDirectoryRecursive(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	directory := path.Join(filepath.Dir(file), "/examples/...")
	contracts := NewGoContracts()
	positions, err := contracts.AddDirectory(directory)
	assert.Len(t, positions, 2)
	assert.NoError(t, err)
}

func TestGoContracts(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	path := path.Join(filepath.Dir(file), "/examples/non-interference.go")
	contracts := NewGoContracts()
	contracts.AddFile(path)
	all := iterx.CollectMap(contracts.Iterator())
	assert.Len(t, all, 2)
}
