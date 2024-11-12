package language

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestInject(t *testing.T) {
	contracts := NewGoContracts()
	_, file, _, _ := runtime.Caller(0)
	path := path.Join(filepath.Dir(file), "/examples/non-interference")
	contracts.AddFile(path)
	injector := NewGoInjector(contracts)
	injector.Inject()
	t.Fail()
}
