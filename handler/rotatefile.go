package handler

import (
	"github.com/govine/easylog"
)

type IWriter interface {
	Write(p []byte) (n int, err error)
	Flush() error
	Close()
}

type RotateFileHandler struct {
	format easylog.Formatter
	writer IWriter
}

func (r *RotateFileHandler) Handle(record *easylog.Record) {
	if r.writer != nil {
		s := r.format(record)
		_, _ = r.writer.Write([]byte(s + "\n"))
	}
}

func (r *RotateFileHandler) Flush() {
	if r.writer != nil {
		_ = r.writer.Flush()
	}
}

func (r *RotateFileHandler) Close() {
	if r.writer != nil {
		r.writer.Close()
	}
}

func NewRotateFileHandler(format easylog.Formatter, writer IWriter) easylog.IEasyLogHandler {
	return easylog.NewEasyLogHandler(&RotateFileHandler{
		format: format,
		writer: writer,
	})
}
