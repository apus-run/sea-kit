//来源: https://github.com/go-kratos/kratos/blob/main/config/file/file.go

package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/apus-run/sea-kit/conf"
)

var _ conf.Source = (*file)(nil)

type file struct {
	path string
}

// NewSource new a file source.
func NewSource(path string) conf.Source {
	return &file{path: path}
}

func (f *file) loadFile(path string) (*conf.KeyValue, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &conf.KeyValue{
		Key:    info.Name(),
		Format: Format(info.Name()),
		Value:  data,
		Path:   path,
	}, nil
}

func (f *file) loadDir(path string) (kvs []*conf.KeyValue, err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		// ignore hidden files
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		kv, err := f.loadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, kv)
	}
	return
}

func (f *file) Load() (kvs []*conf.KeyValue, err error) {
	fi, err := os.Stat(f.path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return f.loadDir(f.path)
	}
	kv, err := f.loadFile(f.path)
	if err != nil {
		return nil, err
	}
	return []*conf.KeyValue{kv}, nil
}
