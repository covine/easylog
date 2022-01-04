package handler

import (
	"os"
	"path/filepath"

	"github.com/covine/easylog"
)

type FileHandler struct {
	format easylog.Formatter
	f      *os.File
}

func (f *FileHandler) Handle(record *easylog.Event) {
	var str string
	if f.format != nil {
		str = f.format(record)
	} else {
		str = record.Message
	}

	_, _ = f.f.Write([]byte(str + "\n"))
}

func (f *FileHandler) Flush() {
}

func (f *FileHandler) Close() {
	_ = f.f.Close()
}

func NewFileHandler(filePath string, format easylog.Formatter) (easylog.IEasyLogHandler, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return easylog.NewEasyLogHandler(&FileHandler{
		format: format,
		f:      f,
	}), nil
}
