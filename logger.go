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

	tags *map[string]interface{}
	kvs  *map[string]interface{}
	sync.Map

	handlers *[]Handler
	hMu      sync.Mutex

	errorHandler ErrorHandler
}

func newLogger() *logger {
	handlers := make([]Handler, 0)

	return &logger{
		children:     make(map[*logger]struct{}),
		tags:         new(sync.Map),
		kvs:          new(sync.Map),
		handlers:     &handlers,
		errorHandler: NewNopErrorHandler(),
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

func (l *logger) AddHandler(h Handler) {
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

	hs := make([]Handler, len(*(l.handlers)))
	copy(hs, *(l.handlers))
	hs = append(hs, h)

	l.handlers = &hs
}

func (l *logger) RemoveHandler(h Handler) {
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

	hs := make([]Handler, 0, len(*l.handlers))
	for _, handler := range *(l.handlers) {
		if handler == h {
			continue
		}
		hs = append(hs, handler)
	}

	l.handlers = &hs
}

func (l *logger) SetErrorHandler(w ErrorHandler) {
	l.errorHandler = w
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

func (l *logger) Panic() *Event {
	return l.log(PANIC)
}

func (l *logger) Fatal() *Event {
	return l.log(FATAL)
}

func (l *logger) Flush() {
	for _, handler := range *(l.handlers) {
		if err := handler.Flush(); err != nil {
			// ignore error produced by errorHandler
			_ = l.errorHandler.Handle(err)
		}
	}
}

func (l *logger) Close() {
	for _, handler := range *(l.handlers) {
		if err := handler.Close(); err != nil {
			// ignore error produced by errorHandler
			_ = l.errorHandler.Handle(err)
		}
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

// couldEnd could end the Logger with panic or os.exit().
func (l *logger) couldEnd(level Level, v interface{}) {
	// Note: If there is any Level bigger than PANIC added, the logic here should be updated.
	switch level {
	case PANIC:
		l.Flush()
		// ignore error produced by errorHandler
		_ = l.errorHandler.Flush()
		panic(v)
	case FATAL:
		l.Close()
		// ignore error produced by errorHandler
		_ = l.errorHandler.Close()
		os.Exit(1)
	}
}

func (l *logger) log(level Level) *Event {
	if level < l.level {
		l.couldEnd(level, "")
		// No need to generate an Event for and then be handled.
		return nil
	}

	return newEvent(l, level)
}

func (l *logger) handle(event *Event) {
	defer putEvent(event)

	for _, handler := range *(l.handlers) {
		next, err := handler.Handle(event)
		if err != nil {
			// ignore error produced by errorHandler
			_ = l.errorHandler.Handle(err)
		}
		if !next {
			return
		}
	}

	if l.propagate && l.parent != nil {
		l.parent.handle(event)
	}
}
