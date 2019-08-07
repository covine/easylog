package easylog

import (
	"strings"
	"sync"
)

type manager struct {
	mu        sync.RWMutex
	root      *Logger
	loggerMap map[string]*Logger
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
			l = &Logger{
				name:           name,
				manager:        m,
				level:          NOTSET,
				parent:         nil,
				propagate:      true,
				isPlaceholder:  false,
				placeholderMap: make(map[*Logger]interface{}),
			}
			m.loggerMap[name] = l
			m.fixUpChildren(ph, l)
			m.fixUpParents(l)
		}
		return l
	} else {
		l := &Logger{
			name:           name,
			manager:        m,
			level:          NOTSET,
			parent:         nil,
			propagate:      true,
			isPlaceholder:  false,
			placeholderMap: make(map[*Logger]interface{}),
		}
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
			placeHolder := &Logger{
				name:           "",
				manager:        nil,
				level:          NOTSET,
				parent:         nil,
				propagate:      true,
				isPlaceholder:  true,
				placeholderMap: make(map[*Logger]interface{}),
			}
			placeHolder.placeholderMap[l] = nil
			m.loggerMap[subStr] = placeHolder
		} else {
			if !m.loggerMap[subStr].isPlaceholder {
				rv = m.loggerMap[subStr]
			} else {
				m.loggerMap[subStr].placeholderMap[l] = nil
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
