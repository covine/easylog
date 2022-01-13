package handler

import (
	"os"

	"github.com/covine/easylog"
)

type StderrHandler struct {
	format Formatter
}

func (s *StderrHandler) Handle(e *easylog.Event) (bool, error) {
	b, err := s.format(e)
	if err != nil {
		return true, err
	}

	if _, err := os.Stderr.Write(b); err != nil {
		return true, err
	}

	if _, err := os.Stderr.WriteString("\n"); err != nil {
		return true, err
	}

	return true, nil
}

func (s *StderrHandler) Flush() error {
	return os.Stderr.Sync()
}

func (s *StderrHandler) Close() error {
	if err := os.Stderr.Sync(); err != nil {
		return err
	}

	if err := os.Stderr.Close(); err != nil {
		return err
	}

	return nil
}

func NewStderrHandler(format Formatter) *StderrHandler {
	return &StderrHandler{
		format: format,
	}
}
