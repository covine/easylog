package handler

import (
	"os"

	"github.com/covine/easylog"
)

type StdoutHandler struct {
	format Formatter
}

func (s *StdoutHandler) Handle(e *easylog.Event) (bool, error) {
	b, err := s.format(e)
	if err != nil {
		return true, err
	}

	if _, err := os.Stdout.Write(b); err != nil {
		return true, err
	}

	if _, err := os.Stdout.WriteString("\n"); err != nil {
		return true, err
	}

	return true, nil
}

func (s *StdoutHandler) Flush() error {
	return os.Stdout.Sync()
}

func (s *StdoutHandler) Close() error {
	if err := os.Stdout.Sync(); err != nil {
		return err
	}

	if err := os.Stdout.Close(); err != nil {
		return err
	}

	return nil
}

func NewStdoutHandler(format Formatter) *StdoutHandler {
	return &StdoutHandler{
		format: format,
	}
}
