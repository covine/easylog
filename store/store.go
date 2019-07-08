package store

import (
	"container/list"
	"sync"

	"git.qutoutiao.net/govine/easylog"
)

type StoreLogger struct {
	level easylog.Level

	filters *list.List
	fMu     sync.RWMutex

	handlers *list.List
	hMu      sync.RWMutex
}

func NewStoreLogger() (*StoreLogger, error) {
	return &StoreLogger{
		level:    easylog.WARNING,
		filters:  list.New(),
		handlers: list.New(),
	}, nil
}

func (l *StoreLogger) SetLevel(level easylog.Level) {
	if easylog.IsLevel(level) {
		l.level = level
	}
}

func (l *StoreLogger) AddFilter(f easylog.IFilter) {
	if f == nil {
		return
	}

	l.fMu.Lock()
	defer l.fMu.Unlock()

	find := false
	for ele := l.filters.Front(); ele != nil; ele = ele.Next() {
		filter, ok := ele.Value.(easylog.IFilter)
		if ok && filter == f {
			find = true
			break
		}
	}

	if !find {
		l.filters.PushBack(f)
	}
}

func (l *StoreLogger) RemoveFilter(f easylog.IFilter) {
	if f == nil {
		return
	}

	l.fMu.Lock()
	defer l.fMu.Unlock()

	var next *list.Element
	for ele := l.filters.Front(); ele != nil; ele = next {
		filter, ok := ele.Value.(easylog.IFilter)
		if ok && filter == f {
			next = ele.Next()
			l.filters.Remove(ele)
		}
	}
}

func (l *StoreLogger) filter(record easylog.Record) bool {
	for ele := l.filters.Front(); ele != nil; ele = ele.Next() {
		filter, ok := ele.Value.(easylog.IFilter)
		if ok && filter != nil {
			if filter.Filter(record) == false {
				return false
			}
		}
	}
	return true
}

func (l *StoreLogger) AddHandler(h easylog.IHandler) {
	if h == nil {
		return
	}

	l.hMu.Lock()
	defer l.hMu.Unlock()

	find := false
	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(easylog.IHandler)
		if ok && handler == h {
			find = true
			break
		}
	}

	if !find {
		l.handlers.PushBack(h)
	}
}

func (l *StoreLogger) RemoveHandler(h easylog.IHandler) {
	if h == nil {
		return
	}

	l.hMu.Lock()
	defer l.hMu.Unlock()

	var next *list.Element
	for ele := l.handlers.Front(); ele != nil; ele = next {
		handler, ok := ele.Value.(easylog.IHandler)
		if ok && handler == h {
			next = ele.Next()
			l.handlers.Remove(ele)
		}
	}
}

func (l *StoreLogger) Debug(msg string, args ...interface{}) {
	if l.isEnableFor(easylog.DEBUG) {
		l.log(easylog.DEBUG, msg, args...)
	}
}

func (l *StoreLogger) Info(msg string, args ...interface{}) {
	if l.isEnableFor(easylog.INFO) {
		l.log(easylog.INFO, msg, args...)
	}
}

func (l *StoreLogger) Warning(msg string, args ...interface{}) {
	if l.isEnableFor(easylog.WARNING) {
		l.log(easylog.WARNING, msg, args...)
	}
}

func (l *StoreLogger) Warn(msg string, args ...interface{}) {
	if l.isEnableFor(easylog.WARN) {
		l.log(easylog.WARN, msg, args...)
	}
}

func (l *StoreLogger) Error(msg string, args ...interface{}) {
	if l.isEnableFor(easylog.ERROR) {
		l.log(easylog.ERROR, msg, args...)
	}
}

func (l *StoreLogger) Fatal(msg string, args ...interface{}) {
	if l.isEnableFor(easylog.FATAL) {
		l.log(easylog.FATAL, msg, args...)
	}
}

func (l *StoreLogger) isEnableFor(level easylog.Level) bool {
	return level >= l.level
}

func (l *StoreLogger) log(level easylog.Level, msg string, args ...interface{}) {
	record := easylog.Record{
		Level: level,
		Msg:   msg,
		Args:  args,
	}
	if l.filter(record) {
		l.handle(record)
	}
}

func (l *StoreLogger) handle(record easylog.Record) {
	if l.filter(record) {
		l.callHandlers(record)
	}
}

func (l *StoreLogger) callHandlers(record easylog.Record) {
	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(easylog.IHandler)
		if ok && handler != nil {
			if record.Level >= handler.GetLevel() {
				handler.Handle(record)
			}
		}
	}
}

func (l *StoreLogger) Flush() {
	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(easylog.IHandler)
		if ok && handler != nil {
			handler.Flush()
		}
	}
}

func (l *StoreLogger) Close() {
	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(easylog.IHandler)
		if ok && handler != nil {
			handler.Close()
		}
	}
}
