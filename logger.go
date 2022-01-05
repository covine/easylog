package easylog

import (
	"runtime"
	"sync"
)

type logger struct {
	manager     *manager
	parent      *logger
	placeholder bool
	children    map[*logger]struct{}

	name      string
	propagate bool
	level     Level

	debugFrame bool
	infoFrame  bool
	warnFrame  bool
	errorFrame bool
	fatalFrame bool

	tags *sync.Map
	kvs  *sync.Map

	filters  *[]IFilter
	fMu      sync.Mutex
	handlers *[]IHandler
	hMu      sync.Mutex
}

func newLogger() *logger {
	filters := make([]IFilter, 0)
	handlers := make([]IHandler, 0)

	return &logger{
		children: make(map[*logger]struct{}),
		tags:     new(sync.Map),
		kvs:      new(sync.Map),
		filters:  &filters,
		handlers: &handlers,
	}
}

func (l *logger) Name() string {
	return l.name
}

func (l *logger) SetPropagate(propagate bool) {
	l.propagate = propagate
}

func (l *logger) GetPropagate() bool {
	return l.propagate
}

func (l *logger) SetLevel(level Level) {
	l.level = level
}

func (l *logger) GetLevel() Level {
	return l.level
}

func (l *logger) EnableFrame(level Level) {
	switch level {
	case DEBUG:
		l.debugFrame = true
	case INFO:
		l.infoFrame = true
	case WARN:
		l.warnFrame = true
	case ERROR:
		l.errorFrame = true
	case FATAL:
		l.fatalFrame = true
	}
}

func (l *logger) DisableFrame(level Level) {
	switch level {
	case DEBUG:
		l.debugFrame = false
	case INFO:
		l.infoFrame = false
	case WARN:
		l.warnFrame = false
	case ERROR:
		l.errorFrame = false
	case FATAL:
		l.fatalFrame = false
	}
}

func (l *logger) AddFilter(f IFilter) {
	if f == nil {
		return
	}

	l.fMu.Lock()
	defer l.fMu.Unlock()

	for _, filter := range *(l.filters) {
		if filter == f {
			return
		}
	}

	fs := make([]IFilter, len(*(l.filters)))
	copy(fs, *(l.filters))
	fs = append(fs, f)

	l.filters = &fs
}

func (l *logger) RemoveFilter(f IFilter) {
	if f == nil {
		return
	}

	l.fMu.Lock()
	defer l.fMu.Unlock()

	find := false
	for _, filter := range *(l.filters) {
		if filter == f {
			find = true
		}
	}
	if !find {
		return
	}

	fs := make([]IFilter, 0, len(*l.filters))
	for _, filter := range *(l.filters) {
		if filter == f {
			continue
		}
		fs = append(fs, filter)
	}

	l.filters = &fs
}

func (l *logger) AddHandler(h IHandler) {
	if h == nil {
		return
	}

	l.hMu.Lock()
	defer l.hMu.Unlock()

	for _, handler := range *(l.handlers) {
		if handler == h {
			return
		}
	}

	hs := make([]IHandler, len(*(l.handlers)))
	copy(hs, *(l.handlers))
	hs = append(hs, h)

	l.handlers = &hs
}

func (l *logger) RemoveHandler(h IHandler) {
	if h == nil {
		return
	}

	l.hMu.Lock()
	defer l.hMu.Unlock()

	find := false
	for _, handler := range *(l.handlers) {
		if handler == h {
			find = true
		}
	}
	if !find {
		return
	}

	hs := make([]IHandler, 0, len(*l.handlers))
	for _, handler := range *(l.handlers) {
		if handler == h {
			continue
		}
		hs = append(hs, handler)
	}

	l.handlers = &hs
}

func (l *logger) SetTag(k, v interface{}) {
	l.tags.Store(k, v)
}

func (l *logger) GetTags() *sync.Map {
	return l.tags
}

func (l *logger) SetKv(k, v interface{}) {
	l.kvs.Store(k, v)
}

func (l *logger) GetKvs() *sync.Map {
	return l.kvs
}

func (l *logger) Debug() *Event {
	return l.log(DEBUG, 2)
}

func (l *logger) Info() *Event {
	return l.log(INFO, 2)
}

func (l *logger) Warn() *Event {
	return l.log(WARN, 2)
}

func (l *logger) Error() *Event {
	return l.log(ERROR, 2)
}

func (l *logger) Fatal() *Event {
	return l.log(FATAL, 2)
}

func (l *logger) Flush() {
	for _, handler := range *(l.handlers) {
		handler.Flush()
	}
}

func (l *logger) Close() {
	for _, handler := range *(l.handlers) {
		handler.Flush()
		handler.Close()
	}
}

func (l *logger) needFrame(level Level) bool {
	switch level {
	case DEBUG:
		return l.debugFrame
	case INFO:
		return l.infoFrame
	case WARN:
		return l.warnFrame
	case ERROR:
		return l.errorFrame
	case FATAL:
		return l.fatalFrame
	default:
		return false
	}
}

func (l *logger) log(level Level, skip int) *Event {
	if level < l.level {
		return nil
	}

	event := newEvent()
	event.logger = l
	event.level = level
	if l.needFrame(level) {
		event.pc, event.file, event.line, event.ok = runtime.Caller(skip)
	}

	return event
}

func (l *logger) filter(record *Event) bool {
	for _, filter := range *(l.filters) {
		if filter.Filter(record) == false {
			return false
		}
	}

	return true
}

func (l *logger) handle(event *Event) {
	for _, handler := range *(l.handlers) {
		handler.Handle(event)
	}
}

func (l *logger) handleEvent(event *Event) {
	if !l.filter(event) {
		putEvent(event)
		return
	}

	l.handle(event)

	if l.propagate && l.parent != nil {
		l.parent.handleEvent(event)
	} else {
		putEvent(event)
	}

	return
}

/*
func (l *logger) string() string {
	var names []string
	for k, _ := range l.children {
		names = append(names, k.name)
	}

	s := fmt.Sprintf("%s:%s", l.name, strings.Join(names, ","))

	for k, _ := range l.children {
		s = fmt.Sprintf("%s\n%s", s, k.string())
	}

	return s
}
type CachedLogger struct {
	logger

	mu           sync.Mutex
	cached       bool
	cachedEvents []*Event
}

func (c *CachedLogger) SetCached(cached bool) {
	c.cached = cached
}

func (c *CachedLogger) handleEvent(record *Event) {
	if record.level < l.level {
		putEvent(record)
		return
	}

	if !l.filter(record) {
		putEvent(record)
		return
	}

	if l.cached {
		l.mu.Lock()
		defer l.mu.Unlock()

		if l.cachedEvents == nil {
			l.cachedEvents = make([]*Event, 0)
		}
		l.cachedEvents = append(l.cachedEvents, record)

		return
	} else {
		l.Handlers.Handle(record)

		if l.propagate && l.parent != nil {
			l.parent.handleEvent(record)
		} else {
			putEvent(record)
		}

		return
	}
}

func (c *CachedLogger) Flush() {
	if l.cached {
		l.mu.Lock()
		defer l.mu.Unlock()

		for _, record := range l.cachedEvents {
			l.Handlers.Handle(record)
			if l.propagate && l.parent != nil {
				l.parent.handleEvent(record)
			} else {
				putEvent(record)
			}
		}

		l.cachedEvents = nil
	}

	l.Handlers.Flush()
}
*/
