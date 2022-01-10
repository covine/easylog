package easylog

import (
	"os"
)

type logger struct {
	manager     *manager
	parent      *logger
	placeholder bool
	children    map[*logger]struct{}

	name      string
	propagate bool
	level     Level

	caller map[Level]bool
	stack  map[Level]bool

	tags map[string]interface{}
	kvs  map[string]interface{}

	handlers []Handler

	errorHandler ErrorHandler
}

func newLogger() *logger {
	handlers := make([]Handler, 0)

	return &logger{
		children:     make(map[*logger]struct{}),
		tags:         make(map[string]interface{}),
		kvs:          make(map[string]interface{}),
		handlers:     handlers,
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
	case PANIC:
		l.panicCaller = true
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
	case PANIC:
		l.panicCaller = false
	case FATAL:
		l.fatalCaller = false
	}
}

func (l *logger) AddHandler(h Handler) {
	if h == nil {
		return
	}

	for _, handler := range l.handlers {
		if handler == h {
			return
		}
	}

	l.handlers = append(l.handlers, h)
}

func (l *logger) RemoveHandler(h Handler) {
	if h == nil {
		return
	}

	for i, handler := range l.handlers {
		if handler == h {
			l.handlers = append(l.handlers[:i], l.handlers[i+1:]...)
			return
		}
	}
}

func (l *logger) SetErrorHandler(w ErrorHandler) {
	l.errorHandler = w
}

func (l *logger) SetTag(k string, v interface{}) {
	l.tags[k] = v
}

func (l *logger) GetTags() map[string]interface{} {
	return l.tags
}

func (l *logger) SetKv(k string, v interface{}) {
	l.kvs[k] = v
}

func (l *logger) GetKvs() map[string]interface{} {
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
	for _, handler := range l.handlers {
		if err := handler.Flush(); err != nil {
			// ignore error produced by errorHandler
			_ = l.errorHandler.Handle(err)
		}
	}
}

func (l *logger) Close() {
	for _, handler := range l.handlers {
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
	case PANIC:
		return l.panicCaller
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

	for _, handler := range l.handlers {
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
