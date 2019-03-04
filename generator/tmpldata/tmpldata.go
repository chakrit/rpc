package tmpldata

import (
	"io"
	"net/http"

	"github.com/rakyll/statik/fs"
)

var fileSystem http.FileSystem

func init() {
	fs_, err := fs.New()
	if err != nil {
		panic(err)
	}

	fileSystem = fs_
}

func Open(path string) (io.ReadCloser, error) {
	return fileSystem.Open(path)
}

func Read(path string) (string, error) {
	if bytes, err := fs.ReadFile(fileSystem, path); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}
