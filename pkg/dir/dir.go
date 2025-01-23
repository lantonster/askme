package dir

import "os"

func CreateDirIfNotExist(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
