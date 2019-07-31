package handler

import (
	"fmt"
	"os"

	"github.com/govine/easylog"
)

type StdoutHandler struct {
	format easylog.Formatter
}

func (s *StdoutHandler) Handle(record *easylog.Record) {
	var str string
	if s.format != nil {
		str = s.format(record)
	} else {
		str = fmt.Sprintf(record.Message, record.Args...)
	}

	os.Stdout.Write([]byte(str + "\n"))
}

func (f *StdoutHandler) Flush() {
}

func (f *StdoutHandler) Close() {
}

func NewStdoutHandler(format easylog.Formatter) easylog.IEasyLogHandler {
	return easylog.NewEasyLogHandler(&StdoutHandler{
		format: format,
	})
}
