package easylog

import (
	"strings"
	"sync"
)

type manager struct {
	mu        sync.RWMutex
	root      *logger
	loggerMap map[string]*logger
}

func (m *manager) getLogger(name string) *logger {
	m.mu.RLock()
	if l, ok := m.loggerMap[name]; ok && l != nil && !l.placeholder {
		m.mu.RUnlock()
		return l
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if l, ok := m.loggerMap[name]; ok && l != nil {
		if l.placeholder {
			ph := l
			l = newLogger()
			l.name = name
			l.manager = m
			m.loggerMap[name] = l
			m.fixUpChildren(ph, l)
			m.fixUpParents(l)
		}
		return l
	} else {
		l := newLogger()
		l.name = name
		l.manager = m
		m.loggerMap[name] = l
		m.fixUpParents(l)
		return l
	}
}

func (m *manager) fixUpParents(l *logger) {
	name := l.name
	i := strings.LastIndexByte(name, '.')
	var rv *logger = nil

	for {
		if i < 0 || rv != nil {
			break
		}

		subStr := name[:i]
		if _, ok := m.loggerMap[subStr]; !ok {
			placeHolder := newLogger()
			placeHolder.placeholder = true
			placeHolder.children[l] = struct{}{}
			m.loggerMap[subStr] = placeHolder
		} else {
			if !m.loggerMap[subStr].placeholder {
				rv = m.loggerMap[subStr]
			} else {
				m.loggerMap[subStr].children[l] = struct{}{}
			}
		}

		i = strings.LastIndexByte(subStr, '.')
	}

	if rv == nil {
		rv = m.root
	}

	l.parent = rv
}

func (m *manager) fixUpChildren(placeholder *logger, l *logger) {
	for c := range placeholder.children {
		if !strings.HasPrefix(c.parent.name, l.name) {
			l.parent = c.parent
			c.parent = l
		}
	}
}
