package easylog

import (
	"container/list"
	"sync"
)

type Logger struct {
	name      string
	manager   *manager
	level     Level
	disabled  bool
	parent    *Logger
	propagate bool

	filters *list.List
	fMu     sync.RWMutex

	handlers *list.List
	hMu      sync.RWMutex

	isPlaceholder  bool
	placeholderMap map[*Logger]interface{}
}

func newRootLogger() *Logger {
	return &Logger{
		name:           "root",
		manager:        nil,
		level:          WARNING,
		parent:         nil,
		propagate:      true,
		filters:        list.New(),
		handlers:       list.New(),
		disabled:       false,
		isPlaceholder:  false,
		placeholderMap: make(map[*Logger]interface{}),
	}
}

func newPlaceholder() *Logger {
	return &Logger{
		name:           "",
		manager:        nil,
		level:          WARNING,
		parent:         nil,
		propagate:      true,
		filters:        list.New(),
		handlers:       list.New(),
		disabled:       false,
		isPlaceholder:  true,
		placeholderMap: make(map[*Logger]interface{}),
	}
}

func newLogger(name string) *Logger {
	return &Logger{
		name:           name,
		manager:        nil,
		level:          WARNING,
		parent:         nil,
		propagate:      true,
		filters:        list.New(),
		handlers:       list.New(),
		disabled:       false,
		isPlaceholder:  false,
		placeholderMap: make(map[*Logger]interface{}),
	}
}

func (l *Logger) setManager(manager *manager) {
	l.manager = manager
}

func (l *Logger) Name() string {
	return l.name
}

func (l *Logger) SetLevel(level Level) {
	if IsLevel(level) {
		l.level = level
	}
}

func (l *Logger) SetPropagate(propagate bool) {
	l.propagate = propagate
}

func (l *Logger) AddFilter(f IFilter) {
	if f == nil {
		return
	}

	l.fMu.Lock()
	defer l.fMu.Unlock()

	find := false
	for ele := l.filters.Front(); ele != nil; ele = ele.Next() {
		filter, ok := ele.Value.(IFilter)
		if ok && filter == f {
			find = true
			break
		}
	}

	if !find {
		l.filters.PushBack(f)
	}
}

func (l *Logger) RemoveFilter(f IFilter) {
	if f == nil {
		return
	}

	l.fMu.Lock()
	defer l.fMu.Unlock()

	var next *list.Element
	for ele := l.filters.Front(); ele != nil; ele = next {
		filter, ok := ele.Value.(IFilter)
		if ok && filter == f {
			next = ele.Next()
			l.filters.Remove(ele)
		}
	}
}

func (l *Logger) filter(record Record) bool {
	for ele := l.filters.Front(); ele != nil; ele = ele.Next() {
		filter, ok := ele.Value.(IFilter)
		if ok && filter != nil {
			if filter.Filter(record) == false {
				return false
			}
		}
	}
	return true
}

func (l *Logger) AddHandler(h IHandler) {
	if h == nil {
		return
	}

	l.hMu.Lock()
	defer l.hMu.Unlock()

	find := false
	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IHandler)
		if ok && handler == h {
			find = true
			break
		}
	}

	if !find {
		l.handlers.PushBack(h)
	}
}

func (l *Logger) RemoveHandler(h IHandler) {
	if h == nil {
		return
	}

	l.hMu.Lock()
	defer l.hMu.Unlock()

	var next *list.Element
	for ele := l.handlers.Front(); ele != nil; ele = next {
		handler, ok := ele.Value.(IHandler)
		if ok && handler == h {
			next = ele.Next()
			l.handlers.Remove(ele)
		}
	}
}

func (l *Logger) hasHandlers() bool {
	pl := l
	rv := false

	for pl != nil {
		if pl.handlers.Len() > 0 {
			rv = true
			break
		}
		if !pl.propagate {
			break
		}

		pl = pl.parent
	}
	return rv
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.isEnableFor(DEBUG) {
		l.log(DEBUG, msg, args...)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.isEnableFor(INFO) {
		l.log(INFO, msg, args...)
	}
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	if l.isEnableFor(WARNING) {
		l.log(WARNING, msg, args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.isEnableFor(WARN) {
		l.log(WARN, msg, args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if l.isEnableFor(ERROR) {
		l.log(ERROR, msg, args...)
	}
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	if l.isEnableFor(FATAL) {
		l.log(FATAL, msg, args...)
	}
}

func (l *Logger) getEffectiveLevel() Level {
	logger := l
	for logger != nil {
		if logger.level != NOTSET {
			return logger.level
		}
		logger = logger.parent
	}
	return NOTSET
}

func (l *Logger) isEnableFor(level Level) bool {
	if l.manager.disable >= level {
		return false
	}
	return level >= l.getEffectiveLevel()
}

func (l *Logger) log(level Level, msg string, args ...interface{}) {
	record := Record{
		Level: level,
		Msg:   msg,
		Args:  args,
	}
	if l.filter(record) {
		l.handle(record)
	}
}

func (l *Logger) handle(record Record) {
	if !l.disabled && l.filter(record) {
		l.callHandlers(record)
	}
}

func (l *Logger) callHandlers(record Record) {
	logger := l
	found := 0
	for logger != nil {
		for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
			handler, ok := ele.Value.(IHandler)
			if ok && handler != nil {
				found += 1
				if record.Level >= handler.GetLevel() {
					handler.Handle(record)
				}
			}
		}
		if !l.propagate {
			logger = nil
		} else {
			logger = l.parent
		}
	}

	if found == 0 {
		// TODO 兜底
	}
}

func (l *Logger) Flush() {
	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IHandler)
		if ok && handler != nil {
			handler.Flush()
		}
	}
}

func (l *Logger) Close() {
	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IHandler)
		if ok && handler != nil {
			handler.Close()
		}
	}
}
