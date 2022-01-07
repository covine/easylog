package easylog

import (
	"os"
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

	debugCaller bool
	infoCaller  bool
	warnCaller  bool
	errorCaller bool
	fatalCaller bool

	tags *sync.Map
	kvs  *sync.Map

	filters  *[]IFilter
	fMu      sync.Mutex
	handlers *[]IHandler
	hMu      sync.Mutex

	errWriter BufWriter
}

func newLogger() *logger {
	filters := make([]IFilter, 0)
	handlers := make([]IHandler, 0)

	return &logger{
		children:  make(map[*logger]struct{}),
		tags:      new(sync.Map),
		kvs:       new(sync.Map),
		filters:   &filters,
		handlers:  &handlers,
		errWriter: NewSerialBufWriter(os.Stderr, -1),
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

func (l *logger) EnableCaller(level Level) {
	switch level {
	case DEBUG:
		l.debugCaller = true
	case INFO:
		l.infoCaller = true
	case WARN:
		l.warnCaller = true
	case ERROR:
		l.errorCaller = true
	case FATAL:
		l.fatalCaller = true
	}
}

func (l *logger) DisableCaller(level Level) {
	switch level {
	case DEBUG:
		l.debugCaller = false
	case INFO:
		l.infoCaller = false
	case WARN:
		l.warnCaller = false
	case ERROR:
		l.errorCaller = false
	case FATAL:
		l.fatalCaller = false
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

func (l *logger) SetErrWriter(w BufWriter) {
	l.errWriter = w
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
	return l.log(DEBUG)
}

func (l *logger) Info() *Event {
	return l.log(INFO)
}

func (l *logger) Warn() *Event {
	return l.log(WARN)
}

func (l *logger) Error() *Event {
	return l.log(ERROR)
}

func (l *logger) Fatal() *Event {
	return l.log(FATAL)
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

func (l *logger) needCaller(level Level) bool {
	switch level {
	case DEBUG:
		return l.debugCaller
	case INFO:
		return l.infoCaller
	case WARN:
		return l.warnCaller
	case ERROR:
		return l.errorCaller
	case FATAL:
		return l.fatalCaller
	default:
		return false
	}
}

func (l *logger) log(level Level) *Event {
	if level < l.level {
		return nil
	}

	return newEvent(l, level)
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
		// TODO should flush and close here?
	}
}

func (l *logger) handleEvent(event *Event) {
	if !l.filter(event) {
		if event.Level >= FATAL {
			panic(event.Msg)
		}
		putEvent(event)
		return
	}

	l.handle(event)

	if l.propagate && l.parent != nil {
		l.parent.handleEvent(event)
	} else {
		if event.Level >= FATAL {
			panic(event.Msg)
		}
		putEvent(event)
	}

	return
}

/*
func (l *Logger) string() string {
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
	Logger

	mu           sync.Mutex
	cached       bool
	cachedEvents []*Event
}

func (c *CachedLogger) SetCached(cached bool) {
	c.cached = cached
}

func (c *CachedLogger) handleEvent(record *Event) {
	if record.Level < l.Level {
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
