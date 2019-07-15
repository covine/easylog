package easylog

import "time"

type Logger struct {
	Filters
	Handlers

	name      string
	manager   *manager
	level     Level
	disabled  bool
	parent    *Logger
	propagate bool

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
		disabled:       false,
		isPlaceholder:  false,
		placeholderMap: make(map[*Logger]interface{}),
	}
}

func newSparkLogger() *Logger {
	return &Logger{
		name:           "",
		manager:        nil,
		level:          WARNING,
		parent:         nil,
		propagate:      false,
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

func (l *Logger) setParent(p *Logger) {
	l.parent = p
}

func (l *Logger) SetPropagate(propagate bool) {
	l.propagate = propagate
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
		Time:  time.Now(),
		Level: level,
		Msg:   msg,
		Args:  args,
	}
	if l.Filters.Filter(record) {
		l.handle(record)
	}
}

func (l *Logger) handle(record Record) {
	if !l.disabled && l.Filters.Filter(record) {
		l.callHandlers(record)
	}
}

func (l *Logger) callHandlers(record Record) {
	logger := l
	for logger != nil {
		l.Handlers.Handle(record)
		if !l.propagate {
			logger = nil
		} else {
			logger = l.parent
		}
	}
}

func (l *Logger) Flush() {
	l.Handlers.Flush()
}

func (l *Logger) Close() {
	l.Handlers.Close()
}
