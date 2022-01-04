package easylog

import (
	"container/list"
	"runtime"
)

type logger struct {
	manager *manager

	parent   *logger
	children map[*logger]struct{}

	name        string
	placeholder bool
	propagate   bool
	level       Level
	stack       map[Level]bool

	tags map[string]interface{}
	kvs  map[string]interface{}

	filters  *list.List
	handlers *list.List
}

func newLogger() *logger {
	return &logger{
		manager:     nil,
		parent:      nil,
		children:    make(map[*logger]struct{}),
		name:        "",
		placeholder: false,
		propagate:   false,
		level:       INFO,
		stack:       make(map[Level]bool),
		tags:        make(map[string]interface{}),
		kvs:         make(map[string]interface{}),
		filters:     list.New(),
		handlers:    list.New(),
	}
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
	l.stack[level] = true
}

func (l *logger) DisableFrame(level Level) {
	l.stack[level] = false
}

func (l *logger) AddFilter(f IFilter) {
	find := false
	for ele := l.filters.Front(); ele != nil; ele = ele.Next() {
		filter, ok := ele.Value.(IFilter)
		if ok && filter == f {
			find = true
			break
		}
	}

	if find {
		return
	} else {
		l.filters.PushBack(f)
	}
}

func (l *logger) RemoveFilter(f IFilter) {
	var next *list.Element
	for ele := l.filters.Front(); ele != nil; ele = next {
		next = ele.Next()
		filter, ok := ele.Value.(IFilter)
		if ok && filter == f {
			l.filters.Remove(ele)
		}
	}
}

func (l *logger) AddHandler(h IHandler) {
	find := false
	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IHandler)
		if ok && handler == h {
			find = true
			break
		}
	}

	if find {
		return
	} else {
		l.handlers.PushBack(h)
	}
}

func (l *logger) RemoveHandler(h IHandler) {
	var next *list.Element
	for ele := l.handlers.Front(); ele != nil; ele = next {
		next = ele.Next()
		handler, ok := ele.Value.(IHandler)
		if ok && handler == h {
			l.handlers.Remove(ele)
		}
	}
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
	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IHandler)
		if ok && handler != nil {
			handler.Flush()
		}
	}
}

func (l *logger) Close() {
	l.Flush()

	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IHandler)
		if ok && handler != nil {
			handler.Close()
		}
	}
}

func (l *logger) needFrame(level Level) bool {
	if need, ok := l.stack[level]; ok {
		return need
	}

	return false
}

func (l *logger) log(level Level, skip int) *Event {
	if level < l.level {
		return nil
	}

	event := newEvent()
	event.Logger = l
	event.Level = level
	if l.needFrame(level) {
		event.PC, event.File, event.Line, event.OK = runtime.Caller(skip)
	}

	return event
}

func (l *logger) filter(record *Event) bool {
	for ele := l.filters.Front(); ele != nil; ele = ele.Next() {
		filter, ok := ele.Value.(IFilter)
		if ok && filter != nil {
			if filter.Filter(record) == false {
				return false
			}
		}
	}

	return true
}

func (l *logger) handle(event *Event) {
	for ele := l.handlers.Front(); ele != nil; ele = ele.Next() {
		handler, ok := ele.Value.(IHandler)
		if ok && handler != nil {
			handler.Handle(event)
		}
	}
}

func (l *logger) handleEvent(record *Event) {
	if record.Level < l.level {
		putEvent(record)
		return
	}

	if !l.filter(record) {
		putEvent(record)
		return
	}

	l.handle(record)

	if l.propagate && l.parent != nil {
		l.parent.handleEvent(record)
	} else {
		putEvent(record)
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
	if record.Level < l.level {
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
