package handler

import (
	"os"

	"github.com/covine/easylog"
)

type StderrHandler struct {
	format easylog.Formatter
}

func (s *StderrHandler) Handle(record *easylog.Record) {
	var str string
	if s.format != nil {
		str = s.format(record)
	} else {
		str = record.Message
	}

	_, _ = os.Stderr.Write([]byte(str + "\n"))
}

func (s *StderrHandler) Flush() {
}

func (s *StderrHandler) Close() {
}

func NewStderrHandler(format easylog.Formatter) easylog.IEasyLogHandler {
	return easylog.NewEasyLogHandler(&StderrHandler{
		format: format,
	})
}
