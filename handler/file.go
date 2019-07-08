package handler

import (
	"container/list"
	"fmt"
	"sync"

	"git.qutoutiao.net/govine/easylog"
)

type IWriter interface {
	Write(p []byte) (n int, err error)
	Flush() error
	Close()
}

type FileHandler struct {
	level      easylog.Level
	fileWriter IWriter
	formatter  easylog.IFormatter
	filters    *list.List
	fMu        sync.RWMutex
}

func NewFileHandler(level easylog.Level, fileWriter IWriter) (*FileHandler, error) {
	return &FileHandler{
		level:      level,
		fileWriter: fileWriter,
		filters:    list.New(),
	}, nil
}

func (f *FileHandler) AddFilter(ef easylog.IFilter) {
	if ef == nil {
		return
	}

	f.fMu.Lock()
	defer f.fMu.Unlock()

	find := false
	for ele := f.filters.Front(); ele != nil; ele = ele.Next() {
		filter, ok := ele.Value.(easylog.IFilter)
		if ok && filter == ef {
			find = true
			break
		}
	}

	if !find {
		f.filters.PushBack(ef)
	}
}

func (f *FileHandler) RemoveFilter(ef easylog.IFilter) {
	if ef == nil {
		return
	}

	f.fMu.Lock()
	defer f.fMu.Unlock()

	var next *list.Element
	for ele := f.filters.Front(); ele != nil; ele = next {
		filter, ok := ele.Value.(easylog.IFilter)
		if ok && filter == ef {
			next = ele.Next()
			f.filters.Remove(ele)
		}
	}
}

func (f *FileHandler) filter(record easylog.Record) bool {
	for ele := f.filters.Front(); ele != nil; ele = ele.Next() {
		filter, ok := ele.Value.(easylog.IFilter)
		if ok && filter != nil {
			if filter.Filter(record) == false {
				return false
			}
		}
	}
	return true
}

func (f *FileHandler) SetLevel(level easylog.Level) {
	if easylog.IsLevel(level) {
		f.level = level
	}
}

func (f *FileHandler) GetLevel() easylog.Level {
	return f.level
}

func (f *FileHandler) SetFormatter(formatter easylog.IFormatter) {
	if formatter != nil {
		f.formatter = formatter
	}
}

func (f *FileHandler) Handle(record easylog.Record) {
	if !f.filter(record) {
		return
	}

	s := record.Msg
	if f.formatter != nil {
		s = f.formatter.Format(record)
	} else {
		if record.Args != nil && len(record.Args) > 0 {
			s = fmt.Sprintf(record.Msg, record.Args...)
		}
	}

	f.fileWriter.Write([]byte(s + "\n"))
}

func (f *FileHandler) Flush() {
	if f.fileWriter != nil {
		f.fileWriter.Flush()
	}
}

func (f *FileHandler) Close() {
}
