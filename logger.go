package easylog

import (
	"path"
	"runtime"
	"sync"
	"time"
)

type Logger struct {
	manager        *manager
	isPlaceholder  bool
	placeholderMap map[*Logger]interface{}

	name   string
	parent *Logger

	propagate bool

	level Level
	stack map[Level]bool
	Filters

	Handlers

	mu            sync.Mutex
	cached        bool
	cachedRecords []*Record
}

func (l *Logger) SetPropagate(propagate bool) {
	l.propagate = propagate
}

func (l *Logger) SetLevel(level Level) {
	if IsLevel(level) {
		l.level = level
	}
}

func (l *Logger) SetLevelByString(level string) {
	switch level {
	case "DEBUG":
		l.level = DEBUG
	case "INFO":
		l.level = INFO
	case "WARN":
		l.level = WARN
	case "WARNING":
		l.level = WARNING
	case "ERROR":
		l.level = ERROR
	case "FATAL":
		l.level = FATAL
	default:
		return
	}
}

func (l *Logger) SetStack(level Level, recordStack bool) {
	if IsLevel(level) {
		l.stack[level] = recordStack
	}
}

func (l *Logger) needStack(level Level) bool {
	need, ok := l.stack[level]
	if ok {
		return need
	}
	return false
}

func (l *Logger) SetCached(cached bool) {
	l.cached = cached
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	l.log(WARNING, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(FATAL, msg, args...)
}

func (l *Logger) log(level Level, msg string, args ...interface{}) {
	if level < l.level {
		return
	}

	record := &Record{
		Time:  time.Now(),
		Level: level,
		Msg:   msg,
		Args:  args,
	}
	if l.needStack(level) {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		} else {
			file = path.Base(file)
		}
		record.File = file
		record.Line = line
	}

	l.handle(record)
}

func (l *Logger) handle(record *Record) {
	if record.Level < l.level {
		return
	}

	if !l.Filters.Filter(record) {
		return
	}

	if l.cached {
		// 缓存 record
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.cachedRecords == nil {
			l.cachedRecords = make([]*Record, 0)
		}
		l.cachedRecords = append(l.cachedRecords, record)
		return
	} else {
		l.Handlers.Handle(record)
		if l.propagate && l.parent != nil {
			l.parent.handle(record)
		}
		return
	}
}

func (l *Logger) Flush() {
	if l.cached {
		l.mu.Lock()
		for _, record := range l.cachedRecords {
			l.Handlers.Handle(record)
			if l.propagate && l.parent != nil {
				l.parent.handle(record)
			}
		}
		l.cachedRecords = nil
		l.mu.Unlock()
	}

	l.Handlers.Flush()
}
