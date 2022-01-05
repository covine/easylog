package easylog

import (
	"fmt"
	"sync"
	"time"
)

type Event struct {
	logger *logger

	time  time.Time
	level Level
	msg   string
	tags  *sync.Map
	kvs   *sync.Map
	extra interface{}

	pc   uintptr
	file string
	line int
	ok   bool
}

var eventPool = &sync.Pool{
	New: func() interface{} {
		return &Event{}
	},
}

func newEvent() *Event {
	r := eventPool.Get().(*Event)

	r.logger = nil

	r.time = time.Time{}
	r.level = INFO
	r.msg = ""
	r.tags = new(sync.Map)
	r.kvs = new(sync.Map)
	r.extra = nil

	r.pc = 0
	r.file = ""
	r.line = 0
	r.ok = false

	return r
}

func putEvent(r *Event) {
	if r.level >= FATAL {
		panic(r.msg)
	}

	eventPool.Put(r)
}

func (e *Event) Tag(k, v interface{}) *Event {
	if e == nil {
		return e
	}

	e.tags.Store(k, v)

	return e
}

func (e *Event) Kv(k, v interface{}) *Event {
	if e == nil {
		return e
	}

	e.kvs.Store(k, v)

	return e
}

func (e *Event) Extra(extra interface{}) *Event {
	if e == nil {
		return e
	}

	e.extra = extra
	return e
}

func (e *Event) Log() {
	if e == nil {
		return
	}

	e.time = time.Now()

	e.logger.handleEvent(e)
}

func (e *Event) Logf(msg string, args ...interface{}) {
	if e == nil {
		return
	}

	e.time = time.Now()
	e.msg = fmt.Sprintf(msg, args...)

	e.logger.handleEvent(e)
}

func (e *Event) GetLogger() *logger {
	return e.logger
}

func (e *Event) GetTime() time.Time {
	return e.time
}

func (e *Event) GetLevel() Level {
	return e.level
}

func (e *Event) GetMsg() string {
	return e.msg
}

func (e *Event) GetTags() *sync.Map {
	return e.tags
}

func (e *Event) GetKvs() *sync.Map {
	return e.kvs
}

func (e *Event) GetExtra() interface{} {
	return e.extra
}

func (e *Event) GetPC() uintptr {
	return e.pc
}

func (e *Event) GetFile() string {
	return e.file
}

func (e *Event) GetLine() int {
	return e.line
}

func (e *Event) GetOK() bool {
	return e.ok
}
