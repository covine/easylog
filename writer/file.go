package writer

import (
	"os"
	"path/filepath"
)

type FileWriter struct {
	f *os.File
}

func (f *FileWriter) Write(b []byte) (n int, err error) {
	return f.f.Write(b)
}

func (f *FileWriter) Flush() error {
	return f.f.Sync()
}

func (f *FileWriter) Close() error {
	return f.f.Close()
}

func NewFileWriter(path string) (*FileWriter, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return &FileWriter{f: f}, nil
}
