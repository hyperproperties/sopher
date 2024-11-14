package language

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestInjectV2(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	path := path.Join(filepath.Dir(thisFile), "/examples/monotone.go")

	files := NewFiles()
	files.Add(path)
	injector := NewGoInjector()
	injector.Files(files.Iterator())
}
