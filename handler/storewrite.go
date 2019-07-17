package handler

import (
	"fmt"
	"sync"

	"git.qutoutiao.net/govine/easylog"
)

type StoreHandler struct {
	format easylog.Formatter
	writer IWriter
	logs   []string
	mu     sync.Mutex
}

func (s *StoreHandler) Handle(record easylog.Record) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var str string
	if s.format != nil {
		str = s.format(record)
	} else {
		str = fmt.Sprintf(record.Msg, record.Args)
	}

	s.logs = append(s.logs, str)
}

func (s *StoreHandler) Flush() {
	if s.writer == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, log := range s.logs {
		s.writer.Write([]byte(log + "\n"))
	}
	s.logs = nil
}

func (s *StoreHandler) Close() {
}

func NewStoreWriteHandler(format easylog.Formatter, writer IWriter) easylog.IEasyLogHandler {
	return easylog.NewEasyLogHandler(&StoreHandler{
		format: format,
		writer: writer,
	})
}
