package filesx

import "os"

func Move(source, destination string) error {
	return os.Rename(source, destination)
}

func Create(path string) (*os.File, error) {
	return os.Create(path)
}

func Clear(path string) error {
	return os.Truncate(path, 0)
}

func Delete(path string) error {
	return os.Remove(path)
}
