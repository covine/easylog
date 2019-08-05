package handler

import (
	"fmt"
	"os"

	"github.com/govine/easylog"
)

type StderrHandler struct {
	format easylog.Formatter
}

func (s *StderrHandler) Handle(record *easylog.Record) {
	var str string
	if s.format != nil {
		str = s.format(record)
	} else {
		str = fmt.Sprintf(record.Message, record.Args...)
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
