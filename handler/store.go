package handler

import (
	"container/list"
	"fmt"
	"sync"

	"git.qutoutiao.net/govine/easylog"
)

type StoreHandler struct {
	level      easylog.Level
	fileWriter IWriter
	formatter  easylog.IFormatter
	filters    *list.List
	fMu        sync.RWMutex
	logs       []string
	mu         sync.RWMutex
	flushed    bool
}

func NewStoreHandler(level easylog.Level, fileWriter IWriter) (*StoreHandler, error) {
	return &StoreHandler{
		level:      level,
		fileWriter: fileWriter,
		filters:    list.New(),
		logs:       make([]string, 0),
		flushed:    false,
	}, nil
}

func (f *StoreHandler) AddFilter(ef easylog.IFilter) {
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

func (f *StoreHandler) RemoveFilter(ef easylog.IFilter) {
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

func (f *StoreHandler) filter(record easylog.Record) bool {
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

func (f *StoreHandler) SetLevel(level easylog.Level) {
	if easylog.IsLevel(level) {
		f.level = level
	}
}

func (f *StoreHandler) GetLevel() easylog.Level {
	return f.level
}

func (f *StoreHandler) SetFormatter(formatter easylog.IFormatter) {
	if formatter != nil {
		f.formatter = formatter
	}
}

func (f *StoreHandler) Handle(record easylog.Record) {
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

	f.mu.Lock()
	f.logs = append(f.logs, s)
	f.mu.Unlock()
}

func (f *StoreHandler) Flush() {
	if f.flushed {
		return
	}

	if f.fileWriter != nil {
		f.mu.RLock()
		defer f.mu.RUnlock()

		for _, lo := range f.logs {
			if lo == "" {
				continue
			} else {
				f.fileWriter.Write([]byte(lo + "\n"))
			}
		}
	}

	f.flushed = true
}

func (f *StoreHandler) Close() {
}
