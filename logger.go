package easylog

import (
	"runtime"
	"sync"
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
	if l == nil {
		return
	}
	l.propagate = propagate
}

func (l *Logger) SetLevel(level Level) {
	if l == nil {
		return
	}
	if IsLevel(level) {
		l.level = level
	}
}

func (l *Logger) GetLevel() Level {
	if l == nil {
		return NOTSET
	}
	return l.level
}

func (l *Logger) SetLevelByString(level string) {
	if l == nil {
		return
	}
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

func (l *Logger) EnableFrame(level Level) {
	if l == nil {
		return
	}
	if IsLevel(level) {
		if l.stack == nil {
			l.stack = make(map[Level]bool)
		}
		l.stack[level] = true
	}
}

func (l *Logger) DisableFrame(level Level) {
	if l == nil {
		return
	}
	if IsLevel(level) {
		if l.stack == nil {
			l.stack = make(map[Level]bool)
		}
		l.stack[level] = false
	}
}

func (l *Logger) needRecordFrame(level Level) bool {
	if l.stack == nil {
		return false
	}
	need, ok := l.stack[level]
	if ok {
		return need
	}
	return false
}

func (l *Logger) SetCached(cached bool) {
	if l == nil {
		return
	}

	l.cached = cached
}

func (l *Logger) Debug() *Record {
	if l == nil {
		return nil
	}

	return l.log(DEBUG)
}

func (l *Logger) Info() *Record {
	if l == nil {
		return nil
	}

	return l.log(INFO)
}

func (l *Logger) Warning() *Record {
	if l == nil {
		return nil
	}

	return l.log(WARNING)
}

func (l *Logger) Warn() *Record {
	if l == nil {
		return nil
	}

	return l.log(WARN)
}

func (l *Logger) Error() *Record {
	if l == nil {
		return nil
	}

	return l.log(ERROR)
}

func (l *Logger) Fatal() *Record {
	if l == nil {
		return nil
	}
	return l.log(FATAL)
}

func (l *Logger) log(level Level) *Record {
	if level < l.level {
		return nil
	}

	record := newRecord()
	record.Logger = l
	record.Level = level
	if l.needRecordFrame(level) {
		record.PC, record.File, record.Line, record.OK = runtime.Caller(2)
	}

	return record
}

func (l *Logger) handleRecord(record *Record) {
	if record.Level < l.level {
		putRecord(record)
		return
	}

	if !l.Filters.Filter(record) {
		putRecord(record)
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
			l.parent.handleRecord(record)
		} else {
			putRecord(record)
		}
		return
	}
}

func (l *Logger) Flush() {
	if l == nil {
		return
	}

	if l.cached {
		l.mu.Lock()
		defer l.mu.Unlock()
		for _, record := range l.cachedRecords {
			l.Handlers.Handle(record)
			if l.propagate && l.parent != nil {
				l.parent.handleRecord(record)
			} else {
				putRecord(record)
			}
		}
		l.cachedRecords = nil
	}

	l.Handlers.Flush()
}

func (l *Logger) Close() {
	if l == nil {
		return
	}

	l.Flush()
	l.Handlers.Close()
}
