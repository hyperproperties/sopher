package language

import (
	"errors"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/hyperproperties/sopher/pkg/filesx"
)

var ErrFileNotFound = errors.New("file could not be found")

type Files struct {
	paths []string
}

func NewFiles() Files {
	return Files{}
}

func (files *Files) Add(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return files.AddDirectory(path)
	}
	return files.AddFile(path)
}

func (files *Files) AddFile(path string) error {
	if filepath.Ext(path) == "" {
		path = strings.Join([]string{path, ".go"}, "")
	}
	
	if !filesx.Exists(path) {
		return ErrFileNotFound
	}

	files.paths = append(files.paths, path)

	return nil
}

func (files *Files) AddDirectory(path string) (err error) {
	trimed, includeSubdirs := strings.CutSuffix(path, "...")

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

			files.paths = append(files.paths, path)

			return nil
		})
	} else {
		in, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, file := range in {
			if filepath.Ext(file.Name()) == ".go" {
				path := filepath.Join(path, file.Name())
				if err := files.AddFile(path); err != nil {
					return err
				}
			}
		}
	}

	return err
}

func (files *Files) Iterator() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, path := range files.paths {
			if !yield(path) {
				return
			}
		}
	}
}