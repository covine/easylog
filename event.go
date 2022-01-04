package easylog

import (
	"fmt"
	"sync"
	"time"
)

// Event not concurrency-safe
type Event struct {
	Logger *logger // which Logger Event belongs to

	Time      time.Time
	Level     Level
	Message   string
	Tags      map[string]interface{}
	Kvs       map[string]interface{}
	ExtraData interface{}

	// runtime stack
	PC   uintptr
	File string
	Line int
	OK   bool
}

var recordPool = &sync.Pool{
	New: func() interface{} {
		return &Event{}
	},
}

func newEvent() *Event {
	r := recordPool.Get().(*Event)

	r.Logger = nil
	r.Level = INFO

	r.Time = time.Time{}
	r.Level = INFO
	r.Message = ""
	r.Tags = nil
	r.Kvs = nil
	r.ExtraData = nil

	r.PC = 0
	r.File = ""
	r.Line = 0
	r.OK = false

	return r
}

func putEvent(r *Event) {
	if r.Level >= FATAL {
		panic(r.Message)
	}

	recordPool.Put(r)
}

func (r *Event) Tag(tags map[string]interface{}) *Event {
	if r == nil {
		return r
	}

	if r.Tags == nil {
		r.Tags = make(map[string]interface{})
	}

	for k, v := range tags {
		r.Tags[k] = v
	}
	return r
}

func (r *Event) Fields(kvs map[string]interface{}) *Event {
	if r == nil {
		return r
	}

	if r.Kvs == nil {
		r.Kvs = make(map[string]interface{})
	}

	for k, v := range kvs {
		r.Kvs[k] = v
	}
	return r
}

func (r *Event) Extra(extra interface{}) *Event {
	if r == nil {
		return r
	}

	r.ExtraData = extra
	return r
}

func (r *Event) Log() {
	if r == nil {
		return
	}

	r.Time = time.Now()

	r.Logger.handleEvent(r)
}

func (r *Event) Logf(msg string, args ...interface{}) {
	if r == nil {
		return
	}

	r.Time = time.Now()

	if len(args) > 0 {
		r.Message = fmt.Sprintf(msg, args...)
	} else {
		r.Message = msg
	}

	r.Logger.handleEvent(r)
}
