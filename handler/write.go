package handler

import (
	"git.qutoutiao.net/govine/easylog"
)

type IWriter interface {
	Write(p []byte) (n int, err error)
	Flush() error
	Close()
}

type WriteHandler struct {
	format easylog.Formatter
	writer IWriter
}

func (w *WriteHandler) Handle(record *easylog.Record) {
	if w.writer != nil {
		s := w.format(record)
		w.writer.Write([]byte(s + "\n"))
	}
}

func (w *WriteHandler) Flush() {
	if w.writer != nil {
		w.writer.Flush()
	}
}

func NewWriteHandler(format easylog.Formatter, writer IWriter) easylog.IEasyLogHandler {
	return easylog.NewEasyLogHandler(&WriteHandler{
		format: format,
		writer: writer,
	})
}
