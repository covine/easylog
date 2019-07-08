package easylog

import (
	"strings"
	"sync"
)

type manager struct {
	mu        sync.RWMutex
	root      *Logger
	loggerMap map[string]*Logger
	disable   Level
}

func newManager() *manager {
	return &manager{
		loggerMap: make(map[string]*Logger),
		disable:   NOTSET,
	}
}

func (m *manager) setDisable(level Level) {
	if IsLevel(level) {
		m.disable = level
	}
}

func (m *manager) setRoot(root *Logger) {
	m.root = root
}

func (m *manager) setLoggerClass() {

}

func (m *manager) setLogRecordFactory() {

}

func (m *manager) clearCache() {

}

func (m *manager) getLogger(name string) *Logger {
	m.mu.RLock()
	if l, ok := m.loggerMap[name]; ok && l != nil && !l.isPlaceholder {
		m.mu.RUnlock()
		return l
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if l, ok := m.loggerMap[name]; ok && l != nil {
		if l.isPlaceholder {
			ph := l
			l = newLogger(name)
			l.setManager(m)
			m.loggerMap[name] = l
			m.fixUpChildren(ph, l)
			m.fixUpParents(l)
		}
		return l
	} else {
		l := newLogger(name)
		l.setManager(m)
		m.loggerMap[name] = l
		m.fixUpParents(l)
		return l
	}
}

func (m *manager) fixUpParents(l *Logger) {
	name := l.name
	i := strings.LastIndexByte(name, '.')
	var rv *Logger = nil

	for {
		if i < 0 || rv != nil {
			break
		}
		subStr := name[:i]
		if _, ok := m.loggerMap[subStr]; !ok {
			placeHolder := newPlaceholder()
			placeHolder.placeholderMap[l] = nil
			m.loggerMap[subStr] = placeHolder
		} else {
			tl, o := m.loggerMap[subStr]
			if !o {
				// should not be here
				rv = m.root
			}

			if !tl.isPlaceholder {
				rv = tl
			} else {
				tl.placeholderMap[l] = nil
			}
		}
		i = strings.LastIndexByte(subStr, '.')
	}

	if rv == nil {
		rv = m.root
	}

	l.parent = rv
}

func (m *manager) fixUpChildren(ph *Logger, l *Logger) {
	name := l.name
	nameLen := len(name)
	for c := range ph.placeholderMap {
		if len(c.parent.name) < nameLen {
			l.parent = c.parent
			c.parent = l
		} else if len(c.parent.name) >= nameLen && c.parent.name[:nameLen] != name {
			l.parent = c.parent
			c.parent = l
		}
	}
}
